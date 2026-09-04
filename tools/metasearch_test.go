package tools

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

// MockSearchEngine implements SearchEngine for testing
type MockSearchEngine struct {
	name    string
	results []SearchResult
	err     error
	delay   time.Duration
}

func (m *MockSearchEngine) Name() string { return m.name }

func (m *MockSearchEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	if m.delay > 0 {
		select {
		case <-time.After(m.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	if m.err != nil {
		return nil, m.err
	}
	if len(m.results) > limit {
		return m.results[:limit], nil
	}
	return m.results, nil
}

func TestNormalizeSearchURL(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{
			input:    "https://www.example.com/path/?utm_source=twitter&utm_medium=cpc",
			expected: "example.com/path",
		},
		{
			input:    "http://github.com/sipeed/picoclaw/",
			expected: "github.com/sipeed/picoclaw",
		},
		{
			input:    "https://en.wikipedia.org/wiki/Go_(programming_language)?fbclid=12345",
			expected: "en.wikipedia.org/wiki/Go_(programming_language)",
		},
		{
			input:    "https://example.com/search?q=test&utm_campaign=winter",
			expected: "example.com/search?q=test",
		},
	}

	for _, c := range cases {
		got := normalizeSearchURL(c.input)
		if got != c.expected {
			t.Errorf("normalizeSearchURL(%q) = %q, want %q", c.input, got, c.expected)
		}
	}
}

func TestDeduplicateAndRank(t *testing.T) {
	items := []SearchResult{
		{
			Title:   "Go Language - Wikipedia",
			URL:     "https://en.wikipedia.org/wiki/Go_(programming_language)?utm_source=duckduckgo",
			Snippet: "Go is a language developed at Google.",
			Engine:  "DuckDuckGo",
			Score:   1.0,
		},
		{
			Title:   "Golang Official Site",
			URL:     "https://golang.org",
			Snippet: "Build simple, secure, scalable systems.",
			Engine:  "DuckDuckGo",
			Score:   1.0,
		},
		{
			Title:   "Go (programming language) - Full Wikipedia Article",
			URL:     "https://en.wikipedia.org/wiki/Go_(programming_language)",
			Snippet: "Go is a statically typed, compiled high-level programming language designed at Google by Robert Griesemer, Rob Pike, and Ken Thompson.",
			Engine:  "Wikipedia",
			Score:   1.1,
		},
	}

	deduped := deduplicateAndRank(items)

	// Should deduplicate the Wikipedia article from 2 occurrences into 1
	if len(deduped) != 2 {
		t.Fatalf("expected 2 unique results, got %d", len(deduped))
	}

	// First item should be the consensus-boosted Wikipedia article
	topResult := deduped[0]
	if !strings.Contains(topResult.Engine, "DuckDuckGo") || !strings.Contains(topResult.Engine, "Wikipedia") {
		t.Errorf("expected combined engines in top result, got: %s", topResult.Engine)
	}

	// Should keep the richer/longer snippet
	if !strings.Contains(topResult.Snippet, "Robert Griesemer") {
		t.Errorf("expected richer snippet from consensus merge, got: %s", topResult.Snippet)
	}

	// Score should be boosted (1.0 or 1.1 + 0.5 = > 1.5)
	if topResult.Score < 1.5 {
		t.Errorf("expected boosted consensus score >= 1.5, got %f", topResult.Score)
	}
}

func TestMetaSearchAggregator_ConcurrentSuccess(t *testing.T) {
	engine1 := &MockSearchEngine{
		name: "MockGoogle",
		results: []SearchResult{
			{Title: "Result A", URL: "https://a.com", Snippet: "From Google", Engine: "MockGoogle", Score: 1.0},
			{Title: "Result B", URL: "https://b.com", Snippet: "Shared result", Engine: "MockGoogle", Score: 1.0},
		},
	}
	engine2 := &MockSearchEngine{
		name: "MockBing",
		results: []SearchResult{
			{Title: "Result B - Bing", URL: "https://b.com", Snippet: "Shared result from Bing", Engine: "MockBing", Score: 1.0},
			{Title: "Result C", URL: "https://c.com", Snippet: "From Bing", Engine: "MockBing", Score: 1.0},
		},
	}

	agg := &MetaSearchAggregator{
		engines: []SearchEngine{engine1, engine2},
		timeout: 2 * time.Second,
	}

	results, engines, err := agg.Search(context.Background(), "test query", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(engines) != 2 {
		t.Errorf("expected 2 successful engines, got %d: %v", len(engines), engines)
	}

	// 3 unique URLs: a.com, b.com, c.com
	if len(results) != 3 {
		t.Fatalf("expected 3 deduplicated results, got %d", len(results))
	}

	// Top result should be b.com due to consensus
	if results[0].URL != "https://b.com" {
		t.Errorf("expected b.com as top consensus result, got: %s", results[0].URL)
	}
}

func TestMetaSearchAggregator_FaultTolerance(t *testing.T) {
	// One engine fails, one engine times out, one engine succeeds
	failingEngine := &MockSearchEngine{
		name: "FailingEngine",
		err:  errors.New("HTTP 403 Forbidden"),
	}
	slowEngine := &MockSearchEngine{
		name:  "SlowEngine",
		delay: 500 * time.Millisecond,
		results: []SearchResult{
			{Title: "Late Result", URL: "https://late.com", Engine: "SlowEngine"},
		},
	}
	healthyEngine := &MockSearchEngine{
		name: "HealthyEngine",
		results: []SearchResult{
			{Title: "Healthy Result", URL: "https://healthy.com", Snippet: "Works fine", Engine: "HealthyEngine", Score: 1.0},
		},
	}

	agg := &MetaSearchAggregator{
		engines: []SearchEngine{failingEngine, slowEngine, healthyEngine},
		timeout: 100 * time.Millisecond, // Will cause slowEngine to timeout
	}

	results, engines, err := agg.Search(context.Background(), "fault test", 5)
	if err != nil {
		t.Fatalf("expected successful search from remaining healthy engine, got err: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 result from healthy engine, got %d", len(results))
	}

	if results[0].URL != "https://healthy.com" {
		t.Errorf("expected result from healthy engine, got: %s", results[0].URL)
	}

	if len(engines) != 1 || engines[0] != "HealthyEngine" {
		t.Errorf("expected [HealthyEngine], got: %v", engines)
	}
}

func TestLiveWebSearch(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping live network test in short mode")
	}
	out, ok := ExecuteWebSearch(map[string]interface{}{
		"query":       "eBPF Linux observability tools",
		"num_results": 5,
	})
	if !ok {
		t.Logf("Search failed or offline: %s", out)
		return
	}
	t.Logf("Live Search Output:\n%s", out)
}
