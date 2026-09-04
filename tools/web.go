package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"scorp-agent/internal/helpers"
)

// ── Web Fetch ──

var htmlTagRe = regexp.MustCompile(`<[^>]*>`)
var scriptRe = regexp.MustCompile(`(?is)<script[^>]*>.*?</script>`)
var styleRe = regexp.MustCompile(`(?is)<style[^>]*>.*?</style>`)
var whitespaceRe = regexp.MustCompile(`\s+`)
var ddgResultRe = regexp.MustCompile(`<a rel="nofollow" class="result__a" href="([^"]*)"[^>]*>(.*?)</a>`)
var ddgSnippetRe = regexp.MustCompile(`<a class="result__snippet"[^>]*>(.*?)</a>`)

func ExecuteWebFetch(args map[string]interface{}) (string, bool) {
	rawURL := helpers.GetStringArg(args, "url", "")
	if rawURL == "" {
		return "Error: 'url' argument is required", false
	}

	maxLength := helpers.GetIntArg(args, "max_length", 3000)

	req, err := http.NewRequest("GET", rawURL, nil)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ScorpAgent/1.0)")

	client := HttpShort
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error fetching URL: %v", err), false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 500*1024))
	if err != nil {
		return fmt.Sprintf("Error reading response: %v", err), false
	}

	content := string(body)

	// Strip scripts and styles
	content = scriptRe.ReplaceAllString(content, "")
	content = styleRe.ReplaceAllString(content, "")

	// Strip HTML tags
	content = htmlTagRe.ReplaceAllString(content, " ")

	// Clean whitespace
	content = whitespaceRe.ReplaceAllString(content, " ")
	content = strings.TrimSpace(content)

	if len(content) > maxLength {
		content = content[:maxLength] + "..."
	}

	return fmt.Sprintf("URL: %s\nStatus: %d\n\n%s", rawURL, resp.StatusCode, content), true
}

// ── Web Search (Native MetaSearch Aggregator) ──

func ExecuteWebSearch(args map[string]interface{}) (string, bool) {
	query := helpers.GetStringArg(args, "query", "")
	if query == "" {
		return "Error: 'query' argument is required", false
	}

	numResults := helpers.GetIntArg(args, "num_results", 5)
	if numResults > 15 {
		numResults = 15
	}

	agg := GetMetaSearchAggregator()
	results, engines, err := agg.Search(context.Background(), query, numResults)
	if err != nil || len(results) == 0 {
		return fmt.Sprintf("No results found for: %s (error: %v)", query, err), false
	}

	var sb strings.Builder
	enginesStr := strings.Join(engines, ", ")
	sb.WriteString(fmt.Sprintf("🔍 Metasearch Results for: \"%s\" (via %s)\n\n", query, enginesStr))

	for i, r := range results {
		engineLabel := ""
		if r.Engine != "" {
			engineLabel = fmt.Sprintf(" [%s]", r.Engine)
		}
		sb.WriteString(fmt.Sprintf("%d. %s%s\n   %s\n", i+1, r.Title, engineLabel, r.URL))
		if r.Snippet != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", r.Snippet))
		}
		sb.WriteString("\n")
	}

	return sb.String(), true
}

func decodeURL(s string) string {
	return strings.NewReplacer("%3A", ":", "%2F", "/", "%3F", "?", "%3D", "=", "%26", "&").Replace(s)
}

// ── Memory (persistent key-value store) ──
