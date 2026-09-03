package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	md "github.com/JohannesKaufmann/html-to-markdown"
	readability "github.com/go-shiori/go-readability"
	"scorp-agent/internal/helpers"
)

// ──────────────────────────────────────────────
// Ultra-Low-RAM Web Engine (< 5MB RAM)
// Implements Tiered Web Architecture:
// 1. Zero-RAM HTTP Fetch + Readability Parser (~2MB RAM)
// 2. Cloud Offload / Remote Scraper Fallback (Firecrawl / Tavily API, 0MB RAM)
// ──────────────────────────────────────────────

// ExecuteReadURL extracts clean article content and converts it to Markdown.
// Extremely fast (~50ms) and uses < 5MB RAM compared to heavy headless browsers (~500MB).
func ExecuteReadURL(args map[string]interface{}) (string, bool) {
	rawURL := helpers.GetStringArg(args, "url", "")
	if rawURL == "" {
		return "Error: 'url' parameter is required", false
	}

	maxLength := helpers.GetIntArg(args, "max_length", 8000)
	forceRemote := helpers.GetBoolArg(args, "remote", false)

	res, err := ReadURL(rawURL, forceRemote, maxLength)
	if err != nil {
		return fmt.Sprintf("Error reading URL '%s': %v", rawURL, err), false
	}

	return res, true
}

// ReadURL extracts reader-mode markdown from a URL
func ReadURL(targetURL string, forceRemote bool, maxLength int) (string, error) {
	if maxLength <= 0 {
		maxLength = 8000
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return "", fmt.Errorf("invalid URL: %w", err)
	}

	// If remote is forced, try remote scrapers directly
	if forceRemote {
		if remoteRes, err := tryRemoteScrape(targetURL); err == nil && len(strings.TrimSpace(remoteRes)) > 100 {
			return truncateURLOutput(remoteRes, maxLength), nil
		}
	}

	// 1. Local Zero-RAM HTTP Fetch
	req, err := http.NewRequestWithContext(context.Background(), "GET", targetURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,id;q=0.8")

	client := &http.Client{
		Timeout: 12 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		// If local request failed (e.g. timeout or blocked), attempt remote fallback
		if remoteRes, remoteErr := tryRemoteScrape(targetURL); remoteErr == nil && len(strings.TrimSpace(remoteRes)) > 100 {
			return truncateURLOutput(remoteRes, maxLength), nil
		}
		return "", fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		// Try remote fallback for HTTP errors (e.g. 403 Forbidden / Cloudflare challenge)
		if remoteRes, remoteErr := tryRemoteScrape(targetURL); remoteErr == nil && len(strings.TrimSpace(remoteRes)) > 100 {
			return truncateURLOutput(remoteRes, maxLength), nil
		}
		return "", fmt.Errorf("HTTP status %d: %s", resp.StatusCode, resp.Status)
	}

	// Limit reader to 2MB to ensure strict memory ceiling
	bodyLimit := io.LimitReader(resp.Body, 2*1024*1024)

	// 2. Readability Extraction (Firefox Reader engine port)
	article, err := readability.FromReader(bodyLimit, parsedURL)
	if err != nil || len(strings.TrimSpace(article.TextContent)) < 80 {
		// If readability returned very little (likely an SPA that needs JS rendering), try remote fallback
		if remoteRes, remoteErr := tryRemoteScrape(targetURL); remoteErr == nil && len(strings.TrimSpace(remoteRes)) > 100 {
			return truncateURLOutput(remoteRes, maxLength), nil
		}
		if err != nil {
			return "", fmt.Errorf("readability extraction failed: %w", err)
		}
	}

	// 3. Convert clean HTML to Markdown
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(article.Content)
	if err != nil || strings.TrimSpace(markdown) == "" {
		markdown = article.TextContent
	}

	var sb strings.Builder
	if article.Title != "" {
		sb.WriteString(fmt.Sprintf("# %s\n\n", article.Title))
	}
	if article.Byline != "" {
		sb.WriteString(fmt.Sprintf("*By %s*\n\n", article.Byline))
	}
	sb.WriteString(markdown)

	result := sb.String()
	return truncateURLOutput(result, maxLength), nil
}

func truncateURLOutput(s string, maxLength int) string {
	s = strings.TrimSpace(s)
	if len(s) > maxLength {
		return s[:maxLength] + fmt.Sprintf("\n\n... [Content truncated to %d chars]", maxLength)
	}
	return s
}

// ──────────────────────────────────────────────
// Remote Scraper Fallbacks (0MB Local RAM)
// Supports Firecrawl & Tavily
// ──────────────────────────────────────────────

func tryRemoteScrape(rawURL string) (string, error) {
	// 1. Try Firecrawl API if configured
	if apiKey := os.Getenv("FIRECRAWL_API_KEY"); apiKey != "" {
		res, err := scrapeFirecrawl(rawURL, apiKey)
		if err == nil && len(strings.TrimSpace(res)) > 0 {
			return fmt.Sprintf("🌐 [Offloaded via Firecrawl Cloud]\n\n%s", res), nil
		}
	}

	// 2. Try Tavily API if configured
	if apiKey := os.Getenv("TAVILY_API_KEY"); apiKey != "" {
		res, err := scrapeTavily(rawURL, apiKey)
		if err == nil && len(strings.TrimSpace(res)) > 0 {
			return fmt.Sprintf("🌐 [Offloaded via Tavily Cloud]\n\n%s", res), nil
		}
	}

	return "", fmt.Errorf("no remote scraper configured or available")
}

func scrapeFirecrawl(rawURL, apiKey string) (string, error) {
	payload := map[string]interface{}{
		"url":     rawURL,
		"formats": []string{"markdown"},
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.firecrawl.dev/v1/scrape", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Success bool `json:"success"`
		Data    struct {
			Markdown string `json:"markdown"`
			Metadata struct {
				Title string `json:"title"`
			} `json:"metadata"`
		} `json:"data"`
		Error string `json:"error"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if !result.Success {
		return "", fmt.Errorf("firecrawl error: %s", result.Error)
	}

	out := result.Data.Markdown
	if result.Data.Metadata.Title != "" && !strings.HasPrefix(out, "# ") {
		out = fmt.Sprintf("# %s\n\n%s", result.Data.Metadata.Title, out)
	}
	return out, nil
}

func scrapeTavily(rawURL, apiKey string) (string, error) {
	payload := map[string]interface{}{
		"urls": []string{rawURL},
	}
	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.tavily.com/extract", bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Results []struct {
			URL        string `json:"url"`
			RawContent string `json:"raw_content"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if len(result.Results) > 0 {
		return result.Results[0].RawContent, nil
	}
	return "", fmt.Errorf("no content returned by tavily")
}
