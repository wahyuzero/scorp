package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync/atomic"
	"time"
)

// ──────────────────────────────────────────────
// MCP SSE / HTTP Transport Implementation (Anthropic Standard)
// ──────────────────────────────────────────────

func isRemoteMCP(cfg MCPServerConfig) bool {
	return cfg.URL != "" || cfg.Transport == "sse" || cfg.Transport == "http"
}

func startSSEServer(name string, cfg MCPServerConfig) (*MCPServer, error) {
	log.Printf("[mcp] Connecting to remote MCP server %s at %s", name, cfg.URL)

	parsedURL, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	srv := &MCPServer{
		Name:       name,
		Config:     cfg,
		alive:      true,
		isSSE:      true,
		ssePostURL: cfg.URL, // default post endpoint
		sseClient:  &http.Client{Timeout: 35 * time.Second},
		sseRespCh:  make(chan *jsonRPCResponse, 32),
		sseCancel:  cancel,
	}

	// If SSE transport is specified or URL contains sse, establish SSE listener
	if cfg.Transport == "sse" || strings.Contains(strings.ToLower(cfg.URL), "sse") {
		req, err := http.NewRequestWithContext(ctx, "GET", cfg.URL, nil)
		if err != nil {
			cancel()
			return nil, fmt.Errorf("create SSE request: %w", err)
		}
		req.Header.Set("Accept", "text/event-stream")
		req.Header.Set("Cache-Control", "no-cache")
		for k, v := range cfg.Headers {
			req.Header.Set(k, v)
		}

		// Connect to SSE stream
		sseClient := &http.Client{Timeout: 0}
		resp, err := sseClient.Do(req)
		if err != nil {
			log.Printf("[mcp] SSE connect error for %s (%v), will use direct HTTP POST", name, err)
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Stream established
			go srv.listenSSEStream(resp.Body, parsedURL)
			// Wait briefly for endpoint event
			time.Sleep(100 * time.Millisecond)
		} else {
			resp.Body.Close()
			log.Printf("[mcp] SSE endpoint returned status %d for %s, falling back to HTTP POST", resp.StatusCode, name)
		}
	}

	// Initialize connection
	if err := srv.initialize(); err != nil {
		srv.Close()
		return nil, fmt.Errorf("initialize remote MCP: %w", err)
	}

	// Discover tools
	tools, err := srv.listTools()
	if err != nil {
		log.Printf("[mcp] Warning: couldn't list tools for remote server %s: %v", name, err)
	}
	srv.tools = tools

	return srv, nil
}

func (s *MCPServer) listenSSEStream(body io.ReadCloser, baseURL *url.URL) {
	defer body.Close()
	scanner := bufio.NewScanner(body)
	var currentEvent string

	for scanner.Scan() {
		if !s.alive {
			return
		}
		line := scanner.Text()
		if line == "" {
			currentEvent = ""
			continue
		}

		if strings.HasPrefix(line, "event:") {
			currentEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			continue
		}

		if strings.HasPrefix(line, "data:") {
			data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
			if currentEvent == "endpoint" {
				// Server provided a relative or absolute message endpoint
				if endpointURL, err := url.Parse(data); err == nil {
					resolved := baseURL.ResolveReference(endpointURL).String()
					s.mu.Lock()
					s.ssePostURL = resolved
					s.mu.Unlock()
					log.Printf("[mcp] Discovered SSE message endpoint for %s: %s", s.Name, resolved)
				}
			} else {
				// Regular message or default event
				var rpcResp jsonRPCResponse
				if err := json.Unmarshal([]byte(data), &rpcResp); err == nil && (rpcResp.ID != 0 || rpcResp.Result != nil || rpcResp.Error != nil) {
					select {
					case s.sseRespCh <- &rpcResp:
					default:
					}
				}
			}
		}
	}
}

// sendRemoteRequest sends JSON-RPC 2.0 via HTTP POST to the MCP endpoint
func (s *MCPServer) sendRemoteRequest(method string, params interface{}) (*jsonRPCResponse, error) {
	s.mu.Lock()
	if !s.alive {
		s.mu.Unlock()
		return nil, fmt.Errorf("server %s is closed", s.Name)
	}
	postURL := s.ssePostURL
	client := s.sseClient
	headers := s.Config.Headers
	s.mu.Unlock()

	id := atomic.AddInt64(&s.reqID, 1)
	reqObj := jsonRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}

	payload, err := json.Marshal(reqObj)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", postURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create POST request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json, text/event-stream")
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("POST request failed: %w", err)
	}
	defer resp.Body.Close()

	// If response is returned directly in HTTP body (common for HTTP JSON-RPC)
	respBody, err := io.ReadAll(resp.Body)
	if err == nil && len(respBody) > 0 {
		var directResp jsonRPCResponse
		if err := json.Unmarshal(respBody, &directResp); err == nil && directResp.ID == id {
			if directResp.Error != nil {
				return nil, fmt.Errorf("MCP error %d: %s", directResp.Error.Code, directResp.Error.Message)
			}
			return &directResp, nil
		}
	}

	// Otherwise, wait for SSE event response matching id
	timeout := time.After(30 * time.Second)
	for {
		select {
		case sseResp := <-s.sseRespCh:
			if sseResp.ID == id {
				if sseResp.Error != nil {
					return nil, fmt.Errorf("MCP error %d: %s", sseResp.Error.Code, sseResp.Error.Message)
				}
				return sseResp, nil
			}
		case <-timeout:
			return nil, fmt.Errorf("timeout waiting for response from %s (id=%d)", s.Name, id)
		}
	}
}

// sendRemoteNotification sends a notification without expecting an ID response
func (s *MCPServer) sendRemoteNotification(method string, params interface{}) error {
	s.mu.Lock()
	postURL := s.ssePostURL
	client := s.sseClient
	headers := s.Config.Headers
	s.mu.Unlock()

	reqObj := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	payload, err := json.Marshal(reqObj)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequest("POST", postURL, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		httpReq.Header.Set(k, v)
	}

	resp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
