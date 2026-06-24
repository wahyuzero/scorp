package tools

import (
	"scorp-agent/internal/helpers"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
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

// ── Web Search (DuckDuckGo) ──

func ExecuteWebSearch(args map[string]interface{}) (string, bool) {
	query := helpers.GetStringArg(args, "query", "")
	if query == "" {
		return "Error: 'query' argument is required", false
	}

	numResults := helpers.GetIntArg(args, "num_results", 5)
	if numResults > 10 {
		numResults = 10
	}

	searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", strings.ReplaceAll(query, " ", "+"))

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), false
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; ScorpAgent/1.0)")

	client := HttpShort
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 500*1024))
	if err != nil {
		return fmt.Sprintf("Error: %v", err), false
	}

	html := string(body)

	// Extract search results
	resultMatches := ddgResultRe.FindAllStringSubmatch(html, numResults)
	snippetMatches := ddgSnippetRe.FindAllStringSubmatch(html, numResults)

	if len(resultMatches) == 0 {
		return fmt.Sprintf("No results found for: %s", query), true
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔍 Search results for: %s\n\n", query))

	for i, match := range resultMatches {
		title := htmlTagRe.ReplaceAllString(match[2], "")
		link := match[1]

		// Decode DuckDuckGo redirect URL
		if strings.Contains(link, "uddg=") {
			if parts := strings.Split(link, "uddg="); len(parts) > 1 {
				decoded := parts[1]
				if idx := strings.Index(decoded, "&"); idx > 0 {
					decoded = decoded[:idx]
				}
				link = decodeURL(decoded)
			}
		}

		snippet := ""
		if i < len(snippetMatches) {
			snippet = htmlTagRe.ReplaceAllString(snippetMatches[i][1], "")
			snippet = strings.TrimSpace(snippet)
		}

		sb.WriteString(fmt.Sprintf("%d. %s\n   %s\n", i+1, title, link))
		if snippet != "" {
			sb.WriteString(fmt.Sprintf("   %s\n", snippet))
		}
		sb.WriteString("\n")
	}

	return sb.String(), true
}

func decodeURL(s string) string {
	return strings.NewReplacer("%3A", ":", "%2F", "/", "%3F", "?", "%3D", "=", "%26", "&").Replace(s)
}

// ── Memory (persistent key-value store) ──
