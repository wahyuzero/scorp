package tools

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReadURL_LocalMock(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintln(w, `<!DOCTYPE html>
<html>
<head><title>Scorp Super Agent Test Article</title></head>
<body>
<header><p>Banner and ad that should be removed</p></header>
<article>
<h1>Scorp Agent Modernization v2</h1>
<p>Scorp is a high performance AI coding agent written in Go, specifically optimized for Android Termux and low resource VPS.</p>
<p>By using Go-Readability, it parses article contents in less than 50 milliseconds using less than 5 megabytes of memory.</p>
</article>
<footer><p>Footer copyright</p></footer>
</body>
</html>`)
	}))
	defer ts.Close()

	res, ok := ExecuteReadURL(map[string]interface{}{
		"url": ts.URL,
	})
	if !ok {
		t.Fatalf("ExecuteReadURL failed: %s", res)
	}

	if !strings.Contains(res, "Scorp") {
		t.Errorf("Expected extracted text to contain 'Scorp', got: %s", res)
	}
	if !strings.Contains(res, "Android Termux") {
		t.Errorf("Expected extracted text to contain 'Android Termux', got: %s", res)
	}
}
