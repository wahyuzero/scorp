package tools

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Native MetaSearch Engine (Embedded Multi-Engine Aggregator)
// Orchestrates concurrent queries across multiple engines, consensus ranking,
// URL normalization, and smart deduplication.
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

	// Google Custom Search API hook
	googleKey := os.Getenv("GOOGLE_SEARCH_API_KEY")
	if googleKey == "" {
		googleKey = os.Getenv("GOOGLE_API_KEY")
	}
	googleCX := os.Getenv("GOOGLE_SEARCH_CX")
	if googleCX == "" {
		googleCX = os.Getenv("GOOGLE_CX")
	}
	if googleKey != "" && googleCX != "" {
		agg.RegisterEngine(NewGoogleCSEEngine(googleKey, googleCX))
	}

	// 2. Built-in zero-dependency free engines (always active)
	agg.RegisterEngine(NewBingEngine())
	agg.RegisterEngine(NewDuckDuckGoHTMLEngine())
	agg.RegisterEngine(NewWikipediaEngine())
	agg.RegisterEngine(NewGitHubEngine())

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
		return raw
	}

	host := strings.TrimPrefix(parsed.Host, "www.")
	path := strings.TrimRight(parsed.Path, "/")

	// Strip common tracking query params
	query := parsed.Query()
	for k := range query {
		lowerK := strings.ToLower(k)
		if strings.HasPrefix(lowerK, "utm_") || lowerK == "ref" || lowerK == "source" || lowerK == "fbclid" || lowerK == "gclid" {
			query.Del(k)
		}
	}

	rawQuery := query.Encode()
	if rawQuery != "" {
		return fmt.Sprintf("%s%s?%s", host, path, rawQuery)
	}
	return fmt.Sprintf("%s%s", host, path)
}
