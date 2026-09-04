package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// ──────────────────────────────────────────────
// Pluggable Search Engines Implementation
// ──────────────────────────────────────────────

var (
	ddgResultRegex  = regexp.MustCompile(`<a rel="nofollow" class="result__a" href="([^"]*)"[^>]*>(.*?)</a>`)
	ddgSnippetRegex = regexp.MustCompile(`<a class="result__snippet"[^>]*>(.*?)</a>`)
	htmlCleaner     = regexp.MustCompile(`<[^>]*>`)
)

// ── 1. DuckDuckGo HTML Engine ──

type DuckDuckGoHTMLEngine struct {
	client *http.Client
}

func NewDuckDuckGoHTMLEngine() *DuckDuckGoHTMLEngine {
	return &DuckDuckGoHTMLEngine{
		client: &http.Client{Timeout: 4 * time.Second},
	}
}

func (e *DuckDuckGoHTMLEngine) Name() string { return "DuckDuckGo" }

func (e *DuckDuckGoHTMLEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	data := url.Values{}
	data.Set("q", query)
	data.Set("kl", "us-en")

	req, err := http.NewRequestWithContext(ctx, "POST", "https://html.duckduckgo.com/html/", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("duckduckgo html returned status: %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return nil, err
	}
	body := string(bodyBytes)

	titlesAndURLs := ddgResultRegex.FindAllStringSubmatch(body, -1)
	snippets := ddgSnippetRegex.FindAllStringSubmatch(body, -1)

	var results []SearchResult
	for i, match := range titlesAndURLs {
		if len(results) >= limit {
			break
		}
		if len(match) < 3 {
			continue
		}

		rawHref := match[1]
		rawTitle := match[2]

		targetURL := rawHref
		if u, err := url.Parse(rawHref); err == nil && strings.Contains(u.Path, "/l/") {
			if realURL := u.Query().Get("uddg"); realURL != "" {
				targetURL = realURL
			}
		}

		cleanTitle := html.UnescapeString(htmlCleaner.ReplaceAllString(rawTitle, ""))
		cleanTitle = strings.TrimSpace(cleanTitle)

		snippet := ""
		if i < len(snippets) && len(snippets[i]) > 1 {
			snippet = html.UnescapeString(htmlCleaner.ReplaceAllString(snippets[i][1], ""))
			snippet = strings.TrimSpace(snippet)
		}

		if targetURL != "" && cleanTitle != "" {
			results = append(results, SearchResult{
				Title:   cleanTitle,
				URL:     targetURL,
				Snippet: snippet,
				Engine:  e.Name(),
				Score:   1.0,
			})
		}
	}

	return results, nil
}

// ── 2. Bing HTML Search Engine ──

type BingEngine struct {
	client *http.Client
}

func NewBingEngine() *BingEngine {
	return &BingEngine{
		client: &http.Client{Timeout: 4 * time.Second},
	}
}

func (e *BingEngine) Name() string { return "Bing" }

var (
	bingItemRegex    = regexp.MustCompile(`(?s)<li class="b_algo".*?<h2><a\s+href="([^"]+)".*?>(.*?)</a></h2>.*?<div class="b_caption"><p>(.*?)</p>`)
	bingAltItemRegex = regexp.MustCompile(`(?s)<li class="b_algo".*?<h2><a\s+href="([^"]+)".*?>(.*?)</a></h2>`)
)

func (e *BingEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("https://www.bing.com/search?q=%s&setlang=en-US", url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:125.0) Gecko/20100101 Firefox/125.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bing search status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return nil, err
	}
	body := string(bodyBytes)

	matches := bingItemRegex.FindAllStringSubmatch(body, -1)
	var results []SearchResult

	if len(matches) > 0 {
		for _, m := range matches {
			if len(results) >= limit {
				break
			}
			if len(m) < 4 {
				continue
			}
			targetURL := m[1]
			title := html.UnescapeString(htmlCleaner.ReplaceAllString(m[2], ""))
			snippet := html.UnescapeString(htmlCleaner.ReplaceAllString(m[3], ""))

			if strings.HasPrefix(targetURL, "http") && title != "" {
				results = append(results, SearchResult{
					Title:   strings.TrimSpace(title),
					URL:     targetURL,
					Snippet: strings.TrimSpace(snippet),
					Engine:  e.Name(),
					Score:   1.0,
				})
			}
		}
	} else {
		altMatches := bingAltItemRegex.FindAllStringSubmatch(body, -1)
		for _, m := range altMatches {
			if len(results) >= limit {
				break
			}
			if len(m) < 3 {
				continue
			}
			targetURL := m[1]
			title := html.UnescapeString(htmlCleaner.ReplaceAllString(m[2], ""))
			if strings.HasPrefix(targetURL, "http") && title != "" {
				results = append(results, SearchResult{
					Title:   strings.TrimSpace(title),
					URL:     targetURL,
					Snippet: "",
					Engine:  e.Name(),
					Score:   0.8,
				})
			}
		}
	}

	return results, nil
}

// ── 3. DuckDuckGo Lite Engine ──

type DuckDuckGoLiteEngine struct {
	client *http.Client
}

func NewDuckDuckGoLiteEngine() *DuckDuckGoLiteEngine {
	return &DuckDuckGoLiteEngine{
		client: &http.Client{Timeout: 4 * time.Second},
	}
}

func (e *DuckDuckGoLiteEngine) Name() string { return "DuckDuckGo-Lite" }

var (
	liteLinkRegex    = regexp.MustCompile(`<a class='result-link' href="([^"]*)">(.*?)</a>`)
	liteSnippetRegex = regexp.MustCompile(`<td class='result-snippet'>(.*?)</td>`)
)

func (e *DuckDuckGoLiteEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	data := url.Values{}
	data.Set("q", query)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://lite.duckduckgo.com/lite/", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Lynx/2.8.9rel.1 libwww-FM/2.14 SSL-MM/1.4.1 OpenSSL/1.0.2k")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ddg lite status %d", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return nil, err
	}
	body := string(bodyBytes)

	links := liteLinkRegex.FindAllStringSubmatch(body, -1)
	snippets := liteSnippetRegex.FindAllStringSubmatch(body, -1)

	var results []SearchResult
	for i, m := range links {
		if len(results) >= limit {
			break
		}
		if len(m) < 3 {
			continue
		}
		targetURL := m[1]
		if u, err := url.Parse(targetURL); err == nil && strings.Contains(u.Path, "/l/") {
			if realURL := u.Query().Get("uddg"); realURL != "" {
				targetURL = realURL
			}
		}

		title := html.UnescapeString(htmlCleaner.ReplaceAllString(m[2], ""))
		snippet := ""
		if i < len(snippets) && len(snippets[i]) > 1 {
			snippet = html.UnescapeString(htmlCleaner.ReplaceAllString(snippets[i][1], ""))
			snippet = strings.TrimSpace(snippet)
		}

		if targetURL != "" && title != "" {
			results = append(results, SearchResult{
				Title:   strings.TrimSpace(title),
				URL:     targetURL,
				Snippet: snippet,
				Engine:  e.Name(),
				Score:   0.9,
			})
		}
	}

	return results, nil
}

// ── 4. Wikipedia OpenSearch Engine ──

type WikipediaEngine struct {
	client *http.Client
}

func NewWikipediaEngine() *WikipediaEngine {
	return &WikipediaEngine{
		client: &http.Client{Timeout: 4 * time.Second},
	}
}

func (e *WikipediaEngine) Name() string { return "Wikipedia" }

func (e *WikipediaEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	apiURL := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=opensearch&search=%s&limit=%d&namespace=0&format=json",
		url.QueryEscape(query), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ScorpMetaSearchBot/2.0 (Open Source AI Agent)")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wikipedia api status %d", resp.StatusCode)
	}

	var data []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	if len(data) < 4 {
		return nil, nil
	}

	titles, _ := data[1].([]interface{})
	snippets, _ := data[2].([]interface{})
	urls, _ := data[3].([]interface{})

	var results []SearchResult
	for i := 0; i < len(titles) && len(results) < limit; i++ {
		t, _ := titles[i].(string)
		s, _ := snippets[i].(string)
		u, _ := urls[i].(string)

		if u != "" && t != "" {
			results = append(results, SearchResult{
				Title:   t,
				URL:     u,
				Snippet: s,
				Engine:  e.Name(),
				Score:   1.1,
			})
		}
	}

	return results, nil
}

// ── 5. GitHub Code/Repo Search Engine ──

type GitHubEngine struct {
	client *http.Client
}

func NewGitHubEngine() *GitHubEngine {
	return &GitHubEngine{
		client: &http.Client{Timeout: 4 * time.Second},
	}
}

func (e *GitHubEngine) Name() string { return "GitHub" }

type ghSearchResp struct {
	Items []struct {
		FullName    string `json:"full_name"`
		HTMLURL     string `json:"html_url"`
		Description string `json:"description"`
		Stars       int    `json:"stargazers_count"`
	} `json:"items"`
}

func (e *GitHubEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&per_page=%d&sort=stars",
		url.QueryEscape(query), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ScorpMetaSearch/2.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api status %d", resp.StatusCode)
	}

	var ghResp ghSearchResp
	if err := json.NewDecoder(resp.Body).Decode(&ghResp); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, item := range ghResp.Items {
		snippet := item.Description
		if item.Stars > 0 {
			snippet = fmt.Sprintf("★ %d | %s", item.Stars, item.Description)
		}
		results = append(results, SearchResult{
			Title:   item.FullName,
			URL:     item.HTMLURL,
			Snippet: snippet,
			Engine:  e.Name(),
			Score:   0.9,
		})
	}

	return results, nil
}

// ── 6. Optional: SearXNG Engine ──

type SearXNGEngine struct {
	baseURL string
	client  *http.Client
}

func NewSearXNGEngine(baseURL string) *SearXNGEngine {
	return &SearXNGEngine{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: 4 * time.Second},
	}
}

func (e *SearXNGEngine) Name() string { return "SearXNG" }

func (e *SearXNGEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("%s/search?q=%s&format=json&categories=general",
		e.baseURL, url.QueryEscape(query))

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ScorpMetaSearchBot/2.0")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("searxng status %d", resp.StatusCode)
	}

	var data struct {
		Results []struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Content string `json:"content"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, r := range data.Results {
		if len(results) >= limit {
			break
		}
		results = append(results, SearchResult{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Content,
			Engine:  e.Name(),
			Score:   1.3,
		})
	}
	return results, nil
}

// ── 7. Optional: Brave Search API ──

type BraveSearchEngine struct {
	apiKey string
	client *http.Client
}

func NewBraveSearchEngine(apiKey string) *BraveSearchEngine {
	return &BraveSearchEngine{
		apiKey: apiKey,
		client: &http.Client{Timeout: 4 * time.Second},
	}
}

func (e *BraveSearchEngine) Name() string { return "Brave" }

func (e *BraveSearchEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("https://api.search.brave.com/res/v1/web/search?q=%s&count=%d",
		url.QueryEscape(query), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Subscription-Token", e.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("brave status %d", resp.StatusCode)
	}

	var data struct {
		Web struct {
			Results []struct {
				Title       string `json:"title"`
				URL         string `json:"url"`
				Description string `json:"description"`
			} `json:"results"`
		} `json:"web"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, r := range data.Web.Results {
		results = append(results, SearchResult{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Description,
			Engine:  e.Name(),
			Score:   1.2,
		})
	}
	return results, nil
}

// ── 8. Optional: Tavily Search API ──

type TavilySearchEngine struct {
	apiKey string
	client *http.Client
}

func NewTavilySearchEngine(apiKey string) *TavilySearchEngine {
	return &TavilySearchEngine{
		apiKey: apiKey,
		client: &http.Client{Timeout: 4 * time.Second},
	}
}

func (e *TavilySearchEngine) Name() string { return "Tavily" }

func (e *TavilySearchEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	reqBody := map[string]interface{}{
		"api_key":             e.apiKey,
		"query":               query,
		"max_results":         limit,
		"search_depth":        "basic",
		"include_answer":      false,
		"include_raw_content": false,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.tavily.com/search", strings.NewReader(string(jsonBytes)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tavily status %d", resp.StatusCode)
	}

	var data struct {
		Results []struct {
			Title   string  `json:"title"`
			URL     string  `json:"url"`
			Content string  `json:"content"`
			Score   float64 `json:"score"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, r := range data.Results {
		score := 1.2
		if r.Score > 0 {
			score += r.Score
		}
		results = append(results, SearchResult{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Content,
			Engine:  e.Name(),
			Score:   score,
		})
	}
	return results, nil
}

// ── 9. Optional: Google Custom Search API ──

type GoogleCSEEngine struct {
	apiKey string
	cx     string
	client *http.Client
}

func NewGoogleCSEEngine(apiKey, cx string) *GoogleCSEEngine {
	return &GoogleCSEEngine{
		apiKey: apiKey,
		cx:     cx,
		client: &http.Client{Timeout: 4 * time.Second},
	}
}

func (e *GoogleCSEEngine) Name() string { return "Google" }

func (e *GoogleCSEEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	searchURL := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s&num=%d",
		e.apiKey, e.cx, url.QueryEscape(query), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google cse status %d", resp.StatusCode)
	}

	var data struct {
		Items []struct {
			Title   string `json:"title"`
			Link    string `json:"link"`
			Snippet string `json:"snippet"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, item := range data.Items {
		results = append(results, SearchResult{
			Title:   item.Title,
			URL:     item.Link,
			Snippet: item.Snippet,
			Engine:  e.Name(),
			Score:   1.2,
		})
	}
	return results, nil
}
