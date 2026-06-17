package main

import (
	"net/http"
	"sync"
	"time"
)

// TransportPool manages per-provider HTTP transports for optimal connection pooling
type TransportPool struct {
	mu          sync.RWMutex
	transports  map[string]*http.Transport // key = baseURL host
	defaultPool *http.Transport
}

var transportPool = &TransportPool{
	transports: make(map[string]*http.Transport),
	defaultPool: &http.Transport{
		MaxIdleConns:        50,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		MaxConnsPerHost:     50,
	},
}

// getTransport returns a transport optimized for the given base URL.
// Uses per-host connection pooling with tuned sizes.
func getTransport(baseURL string) *http.Transport {
	host := extractHost(baseURL)
	if host == "" {
		return transportPool.defaultPool
	}

	transportPool.mu.RLock()
	if t, ok := transportPool.transports[host]; ok {
		transportPool.mu.RUnlock()
		return t
	}
	transportPool.mu.RUnlock()

	// Create new transport with provider-specific tuning
	transportPool.mu.Lock()
	defer transportPool.mu.Unlock()

	// Double-check after acquiring write lock
	if t, ok := transportPool.transports[host]; ok {
		return t
	}

	var t *http.Transport

	switch host {
	case "127.0.0.1", "localhost":
		// Local 9router - high concurrency, fast
		t = &http.Transport{
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 10,
			MaxConnsPerHost:     20,
			IdleConnTimeout:     120 * time.Second,
		}
	case "api.openrouter.ai", "openrouter.ai":
		// OpenRouter - external, moderate pool
		t = &http.Transport{
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 5,
			MaxConnsPerHost:     20,
			IdleConnTimeout:     90 * time.Second,
		}
	case "api.moonshot.cn":
		// Moonshot — standard OpenAI-compatible
		t = &http.Transport{
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 5,
			MaxConnsPerHost:     20,
			IdleConnTimeout:     90 * time.Second,
		}
	default:
		// Generic provider - balanced
		t = &http.Transport{
			MaxIdleConns:        20,
			MaxIdleConnsPerHost: 5,
			MaxConnsPerHost:     20,
			IdleConnTimeout:     90 * time.Second,
		}
	}

	transportPool.transports[host] = t
	return t
}

// extractHost extracts host from URL (e.g., "https://api.openrouter.ai/v1" -> "api.openrouter.ai")
func extractHost(baseURL string) string {
	// Remove protocol
	url := baseURL
	if len(url) > 8 && url[:8] == "https://" {
		url = url[8:]
	} else if len(url) > 7 && url[:7] == "http://" {
		url = url[7:]
	}
	// Remove path
	for i := 0; i < len(url); i++ {
		if url[i] == '/' {
			return url[:i]
		}
	}
	return url
}

// getClient returns an HTTP client with the appropriate transport for the base URL
func getClient(baseURL string, timeout time.Duration) *http.Client {
	transport := getTransport(baseURL)
	return &http.Client{
		Timeout:   timeout,
		Transport: transport,
	}
}

// getAIClient returns a client configured for AI model calls (5 min timeout)
func getAIClient(baseURL string) *http.Client {
	return getClient(baseURL, 5*time.Minute)
}

// getShortClient returns a client for quick health checks (15s timeout)
func getShortClient(baseURL string) *http.Client {
	return getClient(baseURL, 15*time.Second)
}