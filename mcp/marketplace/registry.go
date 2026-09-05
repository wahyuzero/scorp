package marketplace

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"scorp-agent/config"
)

// Artifact describes a prebuilt release artifact for one platform.
type Artifact struct {
	URL   string `json:"url"`
	SHA256 string `json:"sha256"`
}

// IndexEntry is one server record inside registry.json.
type IndexEntry struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Version     string              `json:"version"`
	Description string              `json:"description"`
	Origin      string              `json:"origin"` // port | native
	SourceDir   string              `json:"source_dir"`
	Health      string              `json:"health"` // full | partial
	Tools       []string            `json:"tools"`
	Drift       string              `json:"drift"` // synchronized | outdated | unknown
	ManifestURL string              `json:"manifest_url"`
	Artifacts   map[string]Artifact `json:"artifacts"`
	Variant     *Variant            `json:"variant,omitempty"` // flavors carry their canonical base
}

// RegistryIndex is the top-level registry.json document.
type RegistryIndex struct {
	SchemaVersion int          `json:"schema_version"`
	UpdatedAt     string       `json:"updated_at"`
	Description   string       `json:"description"`
	Servers       []IndexEntry `json:"servers"`
}

var (
	cacheMu       sync.Mutex
	httpClient    = &http.Client{Timeout: HTTPTimeout}
)

// IndexURL resolves the registry index source. SCORP_MCP_REGISTRY_URL wins so
// local checkouts (file://) or forks can be used for development.
func IndexURL() string {
	if u := os.Getenv("SCORP_MCP_REGISTRY_URL"); u != "" {
		return u
	}
	return DefaultIndexURL
}

func cachePath() string   { return config.ScorpPath("marketplace", "registry.json") }
func etagPath() string    { return cachePath() + ".etag" }

// GetIndex returns the registry index, using this freshness order:
//  1. Local cache younger than CacheTTL
//  2. Conditional GET (If-None-Match) against the registry source
//  3. Full fetch
//  4. Stale cache fallback when the network is unreachable
func GetIndex(ctx context.Context) (*RegistryIndex, error) {
	cacheMu.Lock()
	defer cacheMu.Unlock()

	if info, err := os.Stat(cachePath()); err == nil && time.Since(info.ModTime()) < CacheTTL {
		if idx, err := readIndexFile(cachePath()); err == nil {
			return idx, nil
		}
	}

	data, err := fetchURL(ctx, IndexURL(), readEtag())
	if err != nil {
		if err == errNotModified {
			// ETag match: server content unchanged — treat cache as fresh.
			if idx, cerr := readIndexFile(cachePath()); cerr == nil {
				now := time.Now()
				_ = os.Chtimes(cachePath(), now, now)
				return idx, nil
			}
		}
		// Offline fallback: serve the stale cache if we have one.
		if idx, cerr := readIndexFile(cachePath()); cerr == nil {
			log.Printf("[marketplace] registry fetch failed (%v), using stale cache", err)
			return idx, nil
		}
		return nil, fmt.Errorf("fetch registry index: %w", err)
	}

	var idx RegistryIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, fmt.Errorf("parse registry index: %w", err)
	}
	if len(idx.Servers) == 0 {
		return nil, fmt.Errorf("registry index has no servers")
	}

	if err := os.MkdirAll(config.ScorpPath("marketplace"), 0o755); err != nil {
		log.Printf("[marketplace] cache dir: %v", err)
		return &idx, nil
	}
	if werr := os.WriteFile(cachePath(), data, 0o644); werr != nil {
		log.Printf("[marketplace] cache write: %v", werr)
	}
	return &idx, nil
}

// fetchURL performs a GET with optional ETag revalidation. It returns the body
// bytes and remembers the ETag for the next conditional request. file:// URLs
// are read directly so local checkouts of scorp-mcp-registry can be used.
func fetchURL(ctx context.Context, url string, etag string) ([]byte, error) {
	if strings.HasPrefix(url, "file://") {
		return os.ReadFile(strings.TrimPrefix(url, "file://"))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}
	req.Header.Set("User-Agent", "scorp-agent-marketplace")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch {
	case resp.StatusCode == http.StatusNotModified:
		// Recompute from cache — the caller falls through to the stale path.
		return nil, errNotModified
	case resp.StatusCode != http.StatusOK:
		return nil, fmt.Errorf("HTTP %s from %s", resp.Status, url)
	}

	data, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return nil, err
	}

	if newEtag := resp.Header.Get("ETag"); newEtag != "" && url == IndexURL() {
		writeEtag(newEtag)
	}
	return data, nil
}

var errNotModified = fmt.Errorf("not modified")

func readIndexFile(path string) (*RegistryIndex, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var idx RegistryIndex
	if err := json.Unmarshal(data, &idx); err != nil {
		return nil, err
	}
	return &idx, nil
}

func readEtag() string {
	data, err := os.ReadFile(etagPath())
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

func writeEtag(etag string) {
	_ = os.MkdirAll(config.ScorpPath("marketplace"), 0o755)
	_ = os.WriteFile(etagPath(), []byte(etag), 0o644)
}
