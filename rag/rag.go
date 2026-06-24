package rag

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
)

// ──────────────────────────────────────────────
// RAG / Semantic Search — TF-IDF In-Memory Index
// ──────────────────────────────────────────────

// RAGChunk represents a piece of indexed text.
type RAGChunk struct {
	ID       string    `json:"id"`
	Source   string    `json:"source"`   // file path or URL
	Content  string    `json:"content"`  // text content
	 TF map[string]float64 `json:"tf"` // term frequencies
	Created  time.Time `json:"created"`
}

// RAGIndex is the in-memory TF-IDF index.
type RAGIndex struct {
	mu        sync.RWMutex
	Chunks    map[string]*RAGChunk `json:"chunks"`
	DF        map[string]int       `json:"df"`     // document frequency per term
	IndexPath string               `json:"-"`      // disk persistence path
	dirty     bool
}

var ragIndex *RAGIndex

func InitRAG() {
	ragIndex = &RAGIndex{
		Chunks:    make(map[string]*RAGChunk),
		DF:        make(map[string]int),
		IndexPath: ragIndexPath(),
	}
	ragIndex.load()
	log.Printf("[rag] Loaded %d chunks from disk", len(ragIndex.Chunks))
}

// tokenize splits text into lowercase terms, stripping punctuation.
func tokenize(text string) []string {
	text = strings.ToLower(text)
	words := strings.FieldsFunc(text, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})
	// Filter very short tokens
	result := make([]string, 0, len(words))
	for _, w := range words {
		if len(w) >= 2 {
			result = append(result, w)
		}
	}
	return result
}

// computeTF computes term frequency map for a piece of text.
func computeTF(text string) map[string]float64 {
	tokens := tokenize(text)
	if len(tokens) == 0 {
		return map[string]float64{}
	}
	counts := make(map[string]int)
	for _, t := range tokens {
		counts[t]++
	}
	tf := make(map[string]float64, len(counts))
	total := float64(len(tokens))
	for term, count := range counts {
		tf[term] = float64(count) / total
	}
	return tf
}

// computeIDF computes inverse document frequency for a term.
func (idx *RAGIndex) computeIDF(term string) float64 {
	n := len(idx.Chunks)
	if n == 0 {
		return 0
	}
	df := idx.DF[term]
	if df == 0 {
		return 0
	}
	return math.Log(float64(1+n) / float64(1+df))
}

// addChunk adds or updates a chunk in the index.
func (idx *RAGIndex) addChunk(source, content string) string {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	id := fmt.Sprintf("chunk_%d", time.Now().UnixNano())

	// Remove old chunks from same source (re-index)
	for cid, c := range idx.Chunks {
		if c.Source == source {
			// Update DF for old terms
			for term := range c.TF {
				idx.DF[term]--
				if idx.DF[term] <= 0 {
					delete(idx.DF, term)
				}
			}
			delete(idx.Chunks, cid)
		}
	}

	tf := computeTF(content)
	chunk := &RAGChunk{
		ID:      id,
		Source:  source,
		Content: content,
		TF:      tf,
		Created: time.Now(),
	}
	idx.Chunks[id] = chunk

	// Update document frequency
	seen := make(map[string]bool)
	for term := range tf {
		if !seen[term] {
			idx.DF[term]++
			seen[term] = true
		}
	}

	idx.dirty = true
	return id
}

// search performs TF-IDF cosine similarity search.
func (idx *RAGIndex) search(query string, topK int) []struct {
	Chunk   *RAGChunk
	Score   float64
} {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if len(idx.Chunks) == 0 {
		return nil
	}

	// Compute query TF
	queryTF := computeTF(query)
	if len(queryTF) == 0 {
		return nil
	}

	// Compute query TF-IDF vector
	queryVec := make(map[string]float64)
	for term, tf := range queryTF {
		idf := idx.computeIDF(term)
		if idf > 0 {
			queryVec[term] = tf * idf
		}
	}

	if len(queryVec) == 0 {
		// Fall back to raw term matching
		for term := range queryTF {
			queryVec[term] = queryTF[term]
		}
	}

	// Compute cosine similarity with each chunk
	type result struct {
		Chunk *RAGChunk
		Score float64
	}
	results := make([]result, 0, len(idx.Chunks))

	for _, chunk := range idx.Chunks {
		// Compute chunk TF-IDF vector (on the fly)
		var dotProduct, normQ, normC float64
		for term, qWeight := range queryVec {
			idf := idx.computeIDF(term)
			cWeight := chunk.TF[term] * idf
			dotProduct += qWeight * cWeight
			normQ += qWeight * qWeight
		}
		// Compute chunk norm
		for term, tf := range chunk.TF {
			idf := idx.computeIDF(term)
			cWeight := tf * idf
			normC += cWeight * cWeight
		}

		if normQ > 0 && normC > 0 {
			score := dotProduct / (math.Sqrt(normQ) * math.Sqrt(normC))
			if score > 0 {
				results = append(results, result{Chunk: chunk, Score: score})
			}
		}
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	// Limit to topK
	if topK > 0 && len(results) > topK {
		results = results[:topK]
	}

	// Convert to return type
	ret := make([]struct {
		Chunk *RAGChunk
		Score float64
	}, len(results))
	for i, r := range results {
		ret[i].Chunk = r.Chunk
		ret[i].Score = r.Score
	}
	return ret
}

// listSources returns unique indexed sources.
func (idx *RAGIndex) listSources() []struct {
	Source string
	Chunks int
} {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	sourceMap := make(map[string]int)
	for _, c := range idx.Chunks {
		sourceMap[c.Source]++
	}

	type src struct {
		Source string
		Chunks int
	}
	result := make([]src, 0, len(sourceMap))
	for s, n := range sourceMap {
		result = append(result, src{Source: s, Chunks: n})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Source < result[j].Source
	})

	ret := make([]struct {
		Source string
		Chunks int
	}, len(result))
	for i, r := range result {
		ret[i].Source = r.Source
		ret[i].Chunks = r.Chunks
	}
	return ret
}

// removeSource removes all chunks from a given source.
func (idx *RAGIndex) removeSource(source string) int {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	removed := 0
	for cid, c := range idx.Chunks {
		if c.Source == source {
			for term := range c.TF {
				idx.DF[term]--
				if idx.DF[term] <= 0 {
					delete(idx.DF, term)
				}
			}
			delete(idx.Chunks, cid)
			removed++
		}
	}
	if removed > 0 {
		idx.dirty = true
	}
	return removed
}

// persist saves the index to disk (JSON).
func (idx *RAGIndex) persist() error {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(idx.IndexPath)
	os.MkdirAll(dir, 0755)
	return os.WriteFile(idx.IndexPath, data, 0644)
}

// load reads the index from disk.
func (idx *RAGIndex) load() error {
	data, err := os.ReadFile(idx.IndexPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, idx)
}

// ──────────────────────────────────────────────
// RAG Tools (exposed to agent)
// ──────────────────────────────────────────────

func RagIndexAdd(args map[string]interface{}) (string, bool) {
	path := getStringArg(args, "path", "")
	if path == "" {
		return "Error: 'path' is required", false
	}

	// Expand ~ to home
	if strings.HasPrefix(path, "~/") {
		home := homeDir()
		path = filepath.Join(home, path[2:])
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), false
	}

	totalChunks := 0

	if info.IsDir() {
		// Index all files in directory (max 50 files, skip binary)
		files := []string{}
		filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			// Skip common binary/large files
			ext := strings.ToLower(filepath.Ext(p))
			if ext == ".so" || ext == ".o" || ext == ".a" || ext == ".pyc" ||
				ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" ||
				ext == ".zip" || ext == ".gz" || ext == ".tar" || ext == ".db" ||
				ext == ".sqlite" || ext == ".bin" || ext == ".exe" {
				return nil
			}
			// Skip files > 1MB
			if info.Size() > 1024*1024 {
				return nil
			}
			files = append(files, p)
			return nil
		})

		if len(files) > 50 {
			files = files[:50]
		}

		for _, f := range files {
			content, err := os.ReadFile(f)
			if err != nil {
				continue
			}
			// Chunk large files (2000 chars per chunk)
			text := string(content)
			if len(text) > 2000 {
				for i := 0; i < len(text); i += 2000 {
					end := i + 2000
					if end > len(text) {
						end = len(text)
					}
					chunkID := ragIndex.addChunk(f+"#chunk"+fmt.Sprintf("%d", i/2000), text[i:end])
					_ = chunkID
					totalChunks++
				}
			} else {
				ragIndex.addChunk(f, text)
				totalChunks++
			}
		}
	} else {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Sprintf("Error: %v", err), false
		}
		text := string(content)
		if len(text) > 2000 {
			for i := 0; i < len(text); i += 2000 {
				end := i + 2000
				if end > len(text) {
					end = len(text)
				}
				ragIndex.addChunk(path+"#chunk"+fmt.Sprintf("%d", i/2000), text[i:end])
				totalChunks++
			}
		} else {
			ragIndex.addChunk(path, text)
			totalChunks++
		}
	}

	// Persist to disk
	ragIndex.persist()

	return fmt.Sprintf("✅ Indexed %d chunks from %s\nTotal indexed chunks: %d",
		totalChunks, path, len(ragIndex.Chunks)), true
}

func RagIndexSearch(args map[string]interface{}) (string, bool) {
	query := getStringArg(args, "query", "")
	topK := getIntArg(args, "top_k", 5)

	if query == "" {
		return "Error: 'query' is required", false
	}

	results := ragIndex.search(query, topK)
	if len(results) == 0 {
		return "No matching results found. Try indexing more files with index_add.", true
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🔍 Found %d results for \"%s\":\n\n", len(results), query))
	for i, r := range results {
		preview := r.Chunk.Content
		if len(preview) > 300 {
			preview = preview[:300] + "..."
		}
		sb.WriteString(fmt.Sprintf("**%d.** %.2f | `%s`\n```\n%s\n```\n\n",
			i+1, r.Score, r.Chunk.Source, preview))
	}

	return sb.String(), true
}

func RagIndexList(args map[string]interface{}) (string, bool) {
	sources := ragIndex.listSources()
	if len(sources) == 0 {
		return "Index is empty. Use index_add to add files or directories.", true
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📚 Indexed sources (%d chunks total):\n\n", len(ragIndex.Chunks)))
	for i, s := range sources {
		sb.WriteString(fmt.Sprintf("%d. `%s` (%d chunks)\n", i+1, s.Source, s.Chunks))
	}

	return sb.String(), true
}

func RagIndexRemove(args map[string]interface{}) (string, bool) {
	source := getStringArg(args, "source", "")
	if source == "" {
		return "Error: 'source' is required", false
	}

	removed := ragIndex.removeSource(source)
	if removed == 0 {
		return fmt.Sprintf("No chunks found for source: %s", source), false
	}

	ragIndex.persist()
	return fmt.Sprintf("✅ Removed %d chunks from %s", removed, source), true
}
