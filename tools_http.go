package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// HTTP/API Testing Tool — full request builder
// ──────────────────────────────────────────────

func executeHTTP(args map[string]interface{}) (string, bool) {
	method := strings.ToUpper(getStringArg(args, "method", "GET"))
	url := getStringArg(args, "url", "")
	if url == "" {
		return "Error: 'url' argument is required", false
	}

	bodyStr := getStringArg(args, "body", "")
	headersStr := getStringArg(args, "headers", "")
	authType := getStringArg(args, "auth_type", "")
	authValue := getStringArg(args, "auth_value", "")
	timeoutSec := getIntArg(args, "timeout", 15)
	followRedirects := getBoolArg(args, "follow_redirects", true)

	if timeoutSec > 60 {
		timeoutSec = 60
	}

	// Build request body
	var bodyReader io.Reader
	if bodyStr != "" {
		bodyReader = strings.NewReader(bodyStr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return fmt.Sprintf("Error creating request: %v", err), false
	}

	// Default Content-Type for POST/PUT/PATCH with body
	if bodyStr != "" && (method == "POST" || method == "PUT" || method == "PATCH") {
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}
	}

	// Parse headers (JSON object string)
	if headersStr != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(headersStr), &headers); err == nil {
			for k, v := range headers {
				req.Header.Set(k, v)
			}
		}
	}

	// Auth
	switch authType {
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+authValue)
	case "basic":
		req.SetBasicAuth(strings.SplitN(authValue, ":", 2)[0], func() string {
			parts := strings.SplitN(authValue, ":", 2)
			if len(parts) > 1 {
				return parts[1]
			}
			return ""
		}())
	case "api_key":
		headerName := getStringArg(args, "auth_header", "X-API-Key")
		req.Header.Set(headerName, authValue)
	}

	// HTTP client
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: getBoolArg(args, "insecure", false)},
	}
	client := &http.Client{
		Timeout:   time.Duration(timeoutSec) * time.Second,
		Transport: transport,
	}
	if !followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	start := time.Now()
	resp, err := client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		return fmt.Sprintf("Request failed after %s: %v", elapsed.Round(time.Millisecond), err), false
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	respStr := string(respBody)

	// Pretty-print JSON response
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") && len(respBody) > 0 {
		var pretty bytes.Buffer
		if json.Indent(&pretty, respBody, "", "  ") == nil {
			respStr = pretty.String()
		}
	}

	// Build response headers summary
	var hdrSB strings.Builder
	for k, v := range resp.Header {
		if len(v) > 0 {
			hdrSB.WriteString(fmt.Sprintf("  %s: %s\n", k, v[0]))
		}
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📡 %s %s\n", method, url))
	sb.WriteString(fmt.Sprintf("⏱ %s | Status: %d %s\n", elapsed.Round(time.Millisecond), resp.StatusCode, http.StatusText(resp.StatusCode)))
	sb.WriteString(fmt.Sprintf("📦 %d bytes\n\n", len(respBody)))
	if hdrSB.Len() > 0 {
		sb.WriteString("Response Headers:\n")
		sb.WriteString(truncOutput(hdrSB.String(), 500))
		sb.WriteString("\n")
	}
	sb.WriteString("Response Body:\n")
	sb.WriteString(truncOutput(respStr, maxToolOutput))

	return sb.String(), true
}
