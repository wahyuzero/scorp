package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Native MetaSearch Engine (Embedded Multi-Engine Aggregator)
// Provides zero-RAM, zero-Docker metasearch directly in Scorp:
// - Fan-out concurrent queries across multiple engines (DDG HTML, DDG Lite, Wikipedia, GitHub)
// - Optional pluggable providers (SearXNG API, Brave API, Tavily API via ENV)
// - Consensus ranking, URL normalization, and smart deduplication
// ──────────────────────────────────────────────

// SearchResult represents a normalized item found by any search engine.
type SearchResult struct {
	Title   string  `json:"title"`
	URL     string  `json:"url"`
	Snippet string  `json:"snippet"`
	Engine  string  `json:"engine"`
	Score   float64 `json:"score,omitempty"`
}

// SearchEngine defines the interface for all pluggable search backends.
type SearchEngine interface {
	Name() string
	Search(ctx context.Context, query string, limit int) ([]SearchResult, error)
}

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

var (
	ddgResultRegex  = regexp.MustCompile(`<a rel="nofollow" class="result__a" href="([^"]*)"[^>]*>(.*?)</a>`)
	ddgSnippetRegex = regexp.MustCompile(`<a class="result__snippet"[^>]*>(.*?)</a>`)
	htmlCleaner     = regexp.MustCompile(`<[^>]*>`)
)

func (e *DuckDuckGoHTMLEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	data := url.Values{}
	data.Set("q", query)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://html.duckduckgo.com/html/", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return nil, err
	}

	htmlStr := string(body)
	resultMatches := ddgResultRegex.FindAllStringSubmatch(htmlStr, limit*2)
	snippetMatches := ddgSnippetRegex.FindAllStringSubmatch(htmlStr, limit*2)

	var results []SearchResult
	for i, match := range resultMatches {
		if len(results) >= limit {
			break
		}
		rawLink := match[1]
		cleanTitle := html.UnescapeString(strings.TrimSpace(htmlCleaner.ReplaceAllString(match[2], "")))

		// Decode redirect URL: /l/?uddg=...
		actualURL := rawLink
		if strings.Contains(rawLink, "uddg=") {
			parts := strings.Split(rawLink, "uddg=")
			if len(parts) > 1 {
				decoded := parts[1]
				if idx := strings.Index(decoded, "&"); idx > 0 {
					decoded = decoded[:idx]
				}
				if unescaped, err := url.QueryUnescape(decoded); err == nil {
					actualURL = unescaped
				}
			}
		}

		snippet := ""
		if i < len(snippetMatches) {
			snippet = html.UnescapeString(strings.TrimSpace(htmlCleaner.ReplaceAllString(snippetMatches[i][1], "")))
		}

		if actualURL != "" && cleanTitle != "" {
			results = append(results, SearchResult{
				Title:   cleanTitle,
				URL:     actualURL,
				Snippet: snippet,
				Engine:  "DuckDuckGo",
				Score:   1.0,
			})
		}
	}

	return results, nil
}

// ── 2. DuckDuckGo Lite Engine (Ultra-compact fallback) ──

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
	ddgLiteLinkRegex    = regexp.MustCompile(`<a class="result-link" href="([^"]*)"[^>]*>(.*?)</a>`)
	ddgLiteSnippetRegex = regexp.MustCompile(`<td class="result-snippet">(.*?)</td>`)
)

func (e *DuckDuckGoLiteEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	data := url.Values{}
	data.Set("q", query)

	req, err := http.NewRequestWithContext(ctx, "POST", "https://lite.duckduckgo.com/lite/", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024))
	if err != nil {
		return nil, err
	}

	htmlStr := string(body)
	linkMatches := ddgLiteLinkRegex.FindAllStringSubmatch(htmlStr, limit*2)
	snippetMatches := ddgLiteSnippetRegex.FindAllStringSubmatch(htmlStr, limit*2)

	var results []SearchResult
	for i, match := range linkMatches {
		if len(results) >= limit {
			break
		}
		cleanTitle := html.UnescapeString(strings.TrimSpace(htmlCleaner.ReplaceAllString(match[2], "")))
		rawLink := match[1]

		actualURL := rawLink
		if strings.Contains(rawLink, "uddg=") {
			parts := strings.Split(rawLink, "uddg=")
			if len(parts) > 1 {
				decoded := parts[1]
				if idx := strings.Index(decoded, "&"); idx > 0 {
					decoded = decoded[:idx]
				}
				if unescaped, err := url.QueryUnescape(decoded); err == nil {
					actualURL = unescaped
				}
			}
		}

		snippet := ""
		if i < len(snippetMatches) {
			snippet = html.UnescapeString(strings.TrimSpace(htmlCleaner.ReplaceAllString(snippetMatches[i][1], "")))
		}

		if actualURL != "" && cleanTitle != "" {
			results = append(results, SearchResult{
				Title:   cleanTitle,
				URL:     actualURL,
				Snippet: snippet,
				Engine:  "DuckDuckGo-Lite",
				Score:   0.9,
			})
		}
	}

	return results, nil
}

// ── 3. Wikipedia OpenSearch Engine (Free Encyclopedic Knowledge) ──

type WikipediaEngine struct {
	client *http.Client
}

func NewWikipediaEngine() *WikipediaEngine {
	return &WikipediaEngine{
		client: &http.Client{Timeout: 3 * time.Second},
	}
}

func (e *WikipediaEngine) Name() string { return "Wikipedia" }

func (e *WikipediaEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	apiURL := fmt.Sprintf("https://en.wikipedia.org/w/api.php?action=opensearch&search=%s&limit=%d&format=json",
		url.QueryEscape(query), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ScorpAgent/2.0 (https://github.com/wahyuzero/scorp; bot@scorp-agent.local)")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 256*1024))
	if err != nil {
		return nil, err
	}

	var parsed []interface{}
	if err := json.Unmarshal(body, &parsed); err != nil || len(parsed) < 4 {
		return nil, fmt.Errorf("invalid opensearch format")
	}

	titles, _ := parsed[1].([]interface{})
	descriptions, _ := parsed[2].([]interface{})
	links, _ := parsed[3].([]interface{})

	var results []SearchResult
	for i := 0; i < len(titles) && i < len(links); i++ {
		t, ok1 := titles[i].(string)
		u, ok2 := links[i].(string)
		if !ok1 || !ok2 || u == "" || t == "" {
			continue
		}
		desc := ""
		if i < len(descriptions) {
			desc, _ = descriptions[i].(string)
		}
		results = append(results, SearchResult{
			Title:   t,
			URL:     u,
			Snippet: desc,
			Engine:  "Wikipedia",
			Score:   1.1,
		})
	}

	return results, nil
}

// ── 4. GitHub Repositories Engine (Code, Libraries & Tools) ──

type GitHubEngine struct {
	client *http.Client
}

func NewGitHubEngine() *GitHubEngine {
	return &GitHubEngine{
		client: &http.Client{Timeout: 3 * time.Second},
	}
}

func (e *GitHubEngine) Name() string { return "GitHub" }

func (e *GitHubEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	// Only query GitHub if query looks like a package, library, software, or repo search
	qLower := strings.ToLower(query)
	isTechQuery := strings.Contains(qLower, "github") ||
		strings.Contains(qLower, "golang") ||
		strings.Contains(qLower, "python") ||
		strings.Contains(qLower, "rust") ||
		strings.Contains(qLower, "library") ||
		strings.Contains(qLower, "package") ||
		strings.Contains(qLower, "tool") ||
		strings.Contains(qLower, "agent") ||
		strings.Contains(qLower, "sdk") ||
		strings.Contains(qLower, "repo") ||
		strings.Contains(qLower, "cli") ||
		len(strings.Fields(query)) <= 3

	if !isTechQuery {
		return nil, nil // Skip non-technical queries gracefully
	}

	cleanQuery := strings.ReplaceAll(query, "github", "")
	cleanQuery = strings.TrimSpace(cleanQuery)
	if cleanQuery == "" {
		cleanQuery = query
	}

	apiURL := fmt.Sprintf("https://api.github.com/search/repositories?q=%s&per_page=%d",
		url.QueryEscape(cleanQuery), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ScorpAgent/2.0")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var data struct {
		Items []struct {
			FullName    string `json:"full_name"`
			HTMLURL     string `json:"html_url"`
			Description string `json:"description"`
			Stars       int    `json:"stargazers_count"`
			Language    string `json:"language"`
		} `json:"items"`
	}

	if err := json.NewDecoder(io.LimitReader(resp.Body, 512*1024)).Decode(&data); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, item := range data.Items {
		snippet := item.Description
		if item.Language != "" {
			snippet = fmt.Sprintf("[%s ⭐%d] %s", item.Language, item.Stars, snippet)
		} else if item.Stars > 0 {
			snippet = fmt.Sprintf("[⭐%d] %s", item.Stars, snippet)
		}

		results = append(results, SearchResult{
			Title:   item.FullName,
			URL:     item.HTMLURL,
			Snippet: snippet,
			Engine:  "GitHub",
			Score:   1.0,
		})
	}

	return results, nil
}

// ── 5. SearXNG Engine (Pluggable JSON Metasearch API) ──

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
	if e.baseURL == "" {
		return nil, nil
	}

	searchURL := fmt.Sprintf("%s/search?q=%s&format=json", e.baseURL, url.QueryEscape(query))
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "ScorpAgent/2.0")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var data struct {
		Results []struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Content string `json:"content"`
			Engine  string `json:"engine"`
		} `json:"results"`
	}

	if err := json.NewDecoder(io.LimitReader(resp.Body, 1024*1024)).Decode(&data); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, r := range data.Results {
		if len(results) >= limit {
			break
		}
		eng := "SearXNG"
		if r.Engine != "" {
			eng = fmt.Sprintf("SearXNG (%s)", r.Engine)
		}
		results = append(results, SearchResult{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Content,
			Engine:  eng,
			Score:   1.3, // Prioritize SearXNG results
		})
	}

	return results, nil
}

// ── 6. Brave Search API Engine (Optional Commercial Hook) ──

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
	if e.apiKey == "" {
		return nil, nil
	}

	searchURL := fmt.Sprintf("https://api.search.brave.com/res/v1/web/search?q=%s&count=%d",
		url.QueryEscape(query), limit)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Subscription-Token", e.apiKey)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
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

	if err := json.NewDecoder(io.LimitReader(resp.Body, 512*1024)).Decode(&data); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, r := range data.Web.Results {
		results = append(results, SearchResult{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Description,
			Engine:  "Brave",
			Score:   1.2,
		})
	}

	return results, nil
}

// ── 7. Tavily Search API Engine (Optional Commercial Hook) ──

type TavilySearchEngine struct {
	apiKey string
	client *http.Client
}

func NewTavilySearchEngine(apiKey string) *TavilySearchEngine {
	return &TavilySearchEngine{
		apiKey: apiKey,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (e *TavilySearchEngine) Name() string { return "Tavily" }

func (e *TavilySearchEngine) Search(ctx context.Context, query string, limit int) ([]SearchResult, error) {
	if e.apiKey == "" {
		return nil, nil
	}

	payload := map[string]interface{}{
		"api_key":      e.apiKey,
		"query":        query,
		"max_results":  limit,
		"search_depth": "basic",
	}

	jsonBytes, _ := json.Marshal(payload)
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
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	var data struct {
		Results []struct {
			Title   string `json:"title"`
			URL     string `json:"url"`
			Content string `json:"content"`
		} `json:"results"`
	}

	if err := json.NewDecoder(io.LimitReader(resp.Body, 512*1024)).Decode(&data); err != nil {
		return nil, err
	}

	var results []SearchResult
	for _, r := range data.Results {
		results = append(results, SearchResult{
			Title:   r.Title,
			URL:     r.URL,
			Snippet: r.Content,
			Engine:  "Tavily",
			Score:   1.4,
		})
	}

	return results, nil
}

// ──────────────────────────────────────────────
// MetaSearch Aggregator & Deduplicator
// ──────────────────────────────────────────────

type MetaSearchAggregator struct {
	engines []SearchEngine
	timeout time.Duration
}

var (
	metaSearchOnce       sync.Once
	globalMetaAggregator *MetaSearchAggregator
)

// GetMetaSearchAggregator returns the initialized singleton aggregator
func GetMetaSearchAggregator() *MetaSearchAggregator {
	metaSearchOnce.Do(func() {
		globalMetaAggregator = NewDefaultMetaSearchAggregator()
	})
	return globalMetaAggregator
}

// NewDefaultMetaSearchAggregator registers all free native engines and active hooks
func NewDefaultMetaSearchAggregator() *MetaSearchAggregator {
	agg := &MetaSearchAggregator{
		timeout: 4 * time.Second,
	}

	// 1. Pluggable engines via environment (highest priority if configured)
	if searxngURL := os.Getenv("SEARXNG_URL"); searxngURL != "" {
		agg.RegisterEngine(NewSearXNGEngine(searxngURL))
	}
	if tavilyKey := os.Getenv("TAVILY_API_KEY"); tavilyKey != "" {
		agg.RegisterEngine(NewTavilySearchEngine(tavilyKey))
	}
	if braveKey := os.Getenv("BRAVE_API_KEY"); braveKey != "" {
		agg.RegisterEngine(NewBraveSearchEngine(braveKey))
	}

	// 2. Built-in zero-dependency free engines (always active)
	agg.RegisterEngine(NewDuckDuckGoHTMLEngine())
	agg.RegisterEngine(NewWikipediaEngine())
	agg.RegisterEngine(NewGitHubEngine())
	agg.RegisterEngine(NewDuckDuckGoLiteEngine())

	return agg
}

func (m *MetaSearchAggregator) RegisterEngine(engine SearchEngine) {
	if engine != nil {
		m.engines = append(m.engines, engine)
	}
}

// Search executes concurrent queries across all engines and returns deduplicated results.
func (m *MetaSearchAggregator) Search(ctx context.Context, query string, limit int) ([]SearchResult, []string, error) {
	if limit <= 0 {
		limit = 5
	}
	if limit > 15 {
		limit = 15
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()

	type engineOutcome struct {
		engineName string
		results    []SearchResult
		err        error
	}

	resultCh := make(chan engineOutcome, len(m.engines))
	var wg sync.WaitGroup

	for _, eng := range m.engines {
		wg.Add(1)
		go func(e SearchEngine) {
			defer wg.Done()
			res, err := e.Search(ctxTimeout, query, limit)
			resultCh <- engineOutcome{
				engineName: e.Name(),
				results:    res,
				err:        err,
			}
		}(eng)
	}

	wg.Wait()
	close(resultCh)

	var allResults []SearchResult
	var successfulEngines []string

	for outcome := range resultCh {
		if outcome.err != nil {
			log.Printf("[metasearch] Engine '%s' error: %v", outcome.engineName, outcome.err)
			continue
		}
		if len(outcome.results) > 0 {
			successfulEngines = append(successfulEngines, outcome.engineName)
			allResults = append(allResults, outcome.results...)
		}
	}

	if len(allResults) == 0 {
		return nil, nil, fmt.Errorf("no search results returned from active engines")
	}

	// ── Deduplication & Consensus Ranking ──
	deduped := deduplicateAndRank(allResults)

	if len(deduped) > limit {
		deduped = deduped[:limit]
	}

	sort.Strings(successfulEngines)
	return deduped, successfulEngines, nil
}

// deduplicateAndRank groups identical URLs, merges metadata, and boosts consensus score
func deduplicateAndRank(items []SearchResult) []SearchResult {
	grouped := make(map[string]*SearchResult)
	var urlOrder []string

	for _, item := range items {
		cleanURL := normalizeSearchURL(item.URL)
		if cleanURL == "" {
			continue
		}

		if existing, found := grouped[cleanURL]; found {
			// Consensus boost: result found in multiple engines is more authoritative
			existing.Score += 0.5
			if !strings.Contains(existing.Engine, item.Engine) {
				existing.Engine = existing.Engine + ", " + item.Engine
			}
			// Keep longest snippet
			if len(item.Snippet) > len(existing.Snippet) {
				existing.Snippet = item.Snippet
			}
			// Keep cleaner title if longer
			if len(item.Title) > len(existing.Title) {
				existing.Title = item.Title
			}
		} else {
			grouped[cleanURL] = &SearchResult{
				Title:   item.Title,
				URL:     item.URL,
				Snippet: item.Snippet,
				Engine:  item.Engine,
				Score:   item.Score,
			}
			urlOrder = append(urlOrder, cleanURL)
		}
	}

	var results []SearchResult
	for _, u := range urlOrder {
		if r, ok := grouped[u]; ok {
			results = append(results, *r)
		}
	}

	// Sort by consensus score descending
	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// normalizeSearchURL cleans URLs for accurate deduplication
func normalizeSearchURL(raw string) string {
	raw = strings.TrimSpace(raw)
	parsed, err := url.Parse(raw)
	if err != nil {
		return strings.ToLower(raw)
	}

	host := strings.ToLower(parsed.Host)
	host = strings.TrimPrefix(host, "www.")

	path := strings.TrimRight(parsed.Path, "/")

	// Strip common tracking and referral params
	q := parsed.Query()
	trackingParams := []string{
		"utm_source", "utm_medium", "utm_campaign", "utm_term", "utm_content",
		"fbclid", "gclid", "ref", "source", "ref_src", "ved", "ei",
	}
	for _, p := range trackingParams {
		q.Del(p)
	}

	normalized := host + path
	if len(q) > 0 {
		normalized += "?" + q.Encode()
	}

	return normalized
}
