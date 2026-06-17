package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"math/bits"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// ──────────────────────────────────────────────
// Vector RAG — SimHash-Based Semantic Search
// ──────────────────────────────────────────────

// SimHashBits is the fingerprint size (64-bit simhash).
const SimHashBits = 64

// SimHashThreshold is the minimum similarity (0.0-1.0) for a match.
const SimHashThreshold = 0.75 // ~48/64 bits matching

// VecChunk represents a chunk with its simhash fingerprint.
type VecChunk struct {
	ID      string    `json:"id"`
	Source  string    `json:"source"`
	Content string    `json:"content"`
	Simhash uint64    `json:"simhash"`
	Created time.Time `json:"created"`
}

// VectorIndex is the in-memory simhash-based index.
type VectorIndex struct {
	mu     sync.RWMutex
	Chunks map[string]*VecChunk `json:"chunks"`
	DBPath string               `json:"-"`
	dirty  bool
}

var vecIndex *VectorIndex

func initVectorRAG() {
	vecIndex = &VectorIndex{
		Chunks: make(map[string]*VecChunk),
		DBPath: ragVectorDBPath(),
	}
	vecIndex.load()
	log.Printf("[ragvec] Loaded %d simhash chunks from disk", len(vecIndex.Chunks))
}

// ──────────────────────────────────────────────
// SimHash — Pure Go, Zero Dependencies
// ──────────────────────────────────────────────

// tokenHash hashes a token to a 64-bit value using FNV-1a.
func tokenHash(token string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(token))
	return h.Sum64()
}

// computeSimhash generates a 64-bit simhash fingerprint for text.
// Each token's hash bits are weighted by its frequency, then summed.
// The final bit is 1 if the weighted sum > 0, else 0.
func computeSimhash(text string) uint64 {
	tokens := tokenize(text)
	if len(tokens) == 0 {
		return 0
	}

	// Count term frequency
	tf := make(map[string]int)
	for _, t := range tokens {
		tf[t]++
	}

	// Weighted bit accumulator
	var weights [SimHashBits]int64

	for term, count := range tf {
		h := tokenHash(term)
		w := int64(count) // weight = term frequency
		for i := 0; i < SimHashBits; i++ {
			if h&(1<<uint(i)) != 0 {
				weights[i] += w
			} else {
				weights[i] -= w
			}
		}
	}

	// Threshold: bit = 1 if sum > 0
	var fp uint64
	for i := 0; i < SimHashBits; i++ {
		if weights[i] > 0 {
			fp |= 1 << uint(i)
		}
	}
	return fp
}

// hammingDistance counts differing bits between two simhashes.
func hammingDistance(a, b uint64) int {
	return bits.OnesCount64(a ^ b)
}

// simhashSimilarity returns similarity score (0.0-1.0).
func simhashSimilarity(a, b uint64) float64 {
	return 1.0 - float64(hammingDistance(a, b))/SimHashBits
}

// ──────────────────────────────────────────────
// Index Operations
// ──────────────────────────────────────────────

// addVecChunk adds a chunk with simhash to the index.
// Re-indexes if source already exists.
func (idx *VectorIndex) addVecChunk(source, content string) string {
	fp := computeSimhash(content)
	return idx.addVecChunkWithHash(source, content, fp)
}

// addVecChunkWithHash adds a pre-computed simhash chunk.
func (idx *VectorIndex) addVecChunkWithHash(source, content string, fp uint64) string {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	id := fmt.Sprintf("sh_%d", time.Now().UnixNano())

	// Remove old chunks from same source
	for cid, c := range idx.Chunks {
		if c.Source == source {
			delete(idx.Chunks, cid)
		}
	}

	chunk := &VecChunk{
		ID:      id,
		Source:  source,
		Content: content,
		Simhash: fp,
		Created: time.Now(),
	}
	idx.Chunks[id] = chunk
	idx.dirty = true
	return id
}

// vecSearch searches by simhash similarity.
func (idx *VectorIndex) vecSearch(queryFP uint64, topK int) []searchResult {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if len(idx.Chunks) == 0 {
		return nil
	}

	results := make([]searchResult, 0, len(idx.Chunks))
	for _, chunk := range idx.Chunks {
		score := simhashSimilarity(queryFP, chunk.Simhash)
		if score >= SimHashThreshold {
			results = append(results, searchResult{Chunk: chunk, Score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if topK > 0 && len(results) > topK {
		results = results[:topK]
	}
	return results
}

type searchResult struct {
	Chunk *VecChunk
	Score float64
}

// vecSearchTFIDF uses TF-IDF cosine similarity (existing implementation)
// for a fallback when simhash gives no results or hybrid mode.
func (idx *VectorIndex) vecSearchTFIDF(queryText string, topK int) []searchResult {
	if ragIndex == nil {
		return nil
	}
	tfidfResults := ragIndex.search(queryText, topK)
	if len(tfidfResults) == 0 {
		return nil
	}

	results := make([]searchResult, 0, len(tfidfResults))
	for _, r := range tfidfResults {
		results = append(results, searchResult{
			Chunk: &VecChunk{
				ID:      fmt.Sprintf("chunk_%d", time.Now().UnixNano()),
				Source:  r.Chunk.Source,
				Content: r.Chunk.Content,
				Simhash: computeSimhash(r.Chunk.Content),
			},
			Score: r.Score,
		})
	}
	return results
}

// hybridSearch combines simhash + TF-IDF with weighted scoring.
func (idx *VectorIndex) hybridSearch(queryFP uint64, queryText string, topK int, simhashWeight float64) []hybridResult {
	// Get simhash results
	shResults := idx.vecSearch(queryFP, topK*2)

	// Get TF-IDF results
	tfidfWeight := 1.0 - simhashWeight

	// Merge scores: collect all unique chunk sources
	results := make([]hybridResult, 0, len(shResults))
	sourceSeen := make(map[string]bool)

	// Add simhash results
	for _, sr := range shResults {
		sourceSeen[sr.Chunk.Source] = true
		results = append(results, hybridResult{
			Chunk:  sr.Chunk,
			SHSore: sr.Score,
			TScore: 0,
			Final:  sr.Score * simhashWeight,
		})
	}

	// Add TF-IDF results (only if not already in results)
	if ragIndex != nil {
		tfResults := ragIndex.search(queryText, topK*2)
		for _, tr := range tfResults {
			if sourceSeen[tr.Chunk.Source] {
				for i := range results {
					if results[i].Chunk.Source == tr.Chunk.Source {
						results[i].TScore = tr.Score
						results[i].Final = simhashWeight*results[i].SHSore + tfidfWeight*tr.Score
						break
					}
				}
			} else {
				results = append(results, hybridResult{
					Chunk: &VecChunk{
						ID:      fmt.Sprintf("hybrid_%d", time.Now().UnixNano()),
						Source:  tr.Chunk.Source,
						Content: tr.Chunk.Content,
						Simhash: computeSimhash(tr.Chunk.Content),
					},
					SHSore: 0,
					TScore: tr.Score,
					Final:  tr.Score * tfidfWeight,
				})
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Final > results[j].Final
	})

	if topK > 0 && len(results) > topK {
		results = results[:topK]
	}
	return results
}

// vecListSources returns unique sources with chunk counts.
func (idx *VectorIndex) vecListSources() []sourceEntry {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	sourceMap := make(map[string]int)
	for _, c := range idx.Chunks {
		sourceMap[c.Source]++
	}

	entries := make([]sourceEntry, 0, len(sourceMap))
	for s, n := range sourceMap {
		entries = append(entries, sourceEntry{Source: s, Chunks: n})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Source < entries[j].Source
	})
	return entries
}

type sourceEntry struct {
	Source string
	Chunks int
}

type hybridResult struct {
	Chunk  *VecChunk
	SHSore float64
	TScore float64
	Final  float64
}

// vecRemoveSource removes all chunks from a source.
func (idx *VectorIndex) vecRemoveSource(source string) int {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	removed := 0
	for cid, c := range idx.Chunks {
		if c.Source == source {
			delete(idx.Chunks, cid)
			removed++
		}
	}
	if removed > 0 {
		idx.dirty = true
	}
	return removed
}

// ──────────────────────────────────────────────
// Persistence
// ──────────────────────────────────────────────

func (idx *VectorIndex) persist() error {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return err
	}
	dir := filepath.Dir(idx.DBPath)
	os.MkdirAll(dir, 0755)
	return os.WriteFile(idx.DBPath, data, 0644)
}

func (idx *VectorIndex) load() error {
	data, err := os.ReadFile(idx.DBPath)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, idx)
}

// ──────────────────────────────────────────────
// Smart Chunking
// ──────────────────────────────────────────────

// SmartChunk splits text into semantically meaningful chunks.
// Uses paragraph boundaries and sentence-aware splitting.
func SmartChunk(text string, maxChunkSize int) []string {
	if len(text) <= maxChunkSize {
		return []string{text}
	}

	var chunks []string
	// Try paragraph breaks first
	paragraphs := splitParagraphs(text)

	current := ""
	for _, p := range paragraphs {
		if len(current)+len(p) <= maxChunkSize {
			if current != "" {
				current += "\n\n"
			}
			current += p
		} else {
			if current != "" {
				chunks = append(chunks, current)
			}
			// If single paragraph exceeds max size, split by sentences
			if len(p) > maxChunkSize {
				sentences := splitSentences(p)
				for _, s := range sentences {
					if len(current)+len(s) <= maxChunkSize {
						if current != "" {
							current += " "
						}
						current += s
					} else {
						if current != "" {
							chunks = append(chunks, current)
						}
						current = s
						// If single sentence exceeds max, hard-split
						for len(current) > maxChunkSize {
							chunks = append(chunks, current[:maxChunkSize])
							current = current[maxChunkSize:]
						}
					}
				}
			} else {
				current = p
			}
		}
	}
	if current != "" {
		chunks = append(chunks, current)
	}
	return chunks
}

// splitParagraphs splits text into paragraphs.
func splitParagraphs(text string) []string {
	// Return single element if already small
	var result []string
	start := 0
	for i := 0; i < len(text); i++ {
		if i+1 < len(text) && text[i] == '\n' && text[i+1] == '\n' {
			if i > start {
				result = append(result, text[start:i])
			}
			i++ // skip second newline
			start = i + 1
		}
	}
	if start < len(text) {
		result = append(result, text[start:])
	}
	if len(result) == 0 {
		result = []string{text}
	}
	return result
}

// splitSentences splits text into sentences on punctuation boundaries.
func splitSentences(text string) []string {
	var result []string
	start := 0
	for i := 0; i < len(text); i++ {
		if text[i] == '.' || text[i] == '!' || text[i] == '?' {
			if i+1 >= len(text) || text[i+1] == ' ' || text[i+1] == '\n' {
				result = append(result, text[start:i+1])
				start = i + 2
			}
		}
	}
	if start < len(text) {
		result = append(result, text[start:])
	}
	if len(result) == 0 {
		result = []string{text}
	}
	return result
}

// ──────────────────────────────────────────────
// Vector RAG Tools (exposed to agent)
// ──────────────────────────────────────────────

func ragVecIngest(args map[string]interface{}) (string, bool) {
	path := getStringArg(args, "path", "")
	if path == "" {
		return "Error: 'path' is required", false
	}

	// Expand ~ to home
	if len(path) > 0 && path[0] == '~' {
		home, _ := os.UserHomeDir()
		if len(path) > 1 && path[1] == '/' {
			path = filepath.Join(home, path[2:])
		} else {
			path = home
		}
	}

	info, err := os.Stat(path)
	if err != nil {
		return fmt.Sprintf("Error: %v", err), false
	}

	maxChunkSize := getIntArg(args, "chunk_size", 2000)
	if maxChunkSize < 100 {
		maxChunkSize = 2000
	}
	if maxChunkSize > 10000 {
		maxChunkSize = 10000
	}

	totalChunks := 0

	if info.IsDir() {
		files := []string{}
		filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			ext := stringsToLower(filepath.Ext(p))
			if ext == ".so" || ext == ".o" || ext == ".a" || ext == ".pyc" ||
				ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".gif" ||
				ext == ".zip" || ext == ".gz" || ext == ".tar" || ext == ".db" ||
				ext == ".sqlite" || ext == ".bin" || ext == ".exe" {
				return nil
			}
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
			text := string(content)
			chunks := SmartChunk(text, maxChunkSize)
			for i, chunk := range chunks {
				sourceKey := fmt.Sprintf("%s#chunk%d", f, i)
				vecIndex.addVecChunk(sourceKey, chunk)
				totalChunks++
			}
		}
	} else {
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Sprintf("Error: %v", err), false
		}
		text := string(content)
		chunks := SmartChunk(text, maxChunkSize)
		for i, chunk := range chunks {
			sourceKey := fmt.Sprintf("%s#chunk%d", path, i)
			vecIndex.addVecChunk(sourceKey, chunk)
			totalChunks++
		}
	}

	vecIndex.persist()
	return fmt.Sprintf("✅ Indexed %d chunks from %s (simhash)\nTotal vector chunks: %d",
		totalChunks, path, len(vecIndex.Chunks)), true
}

func ragVecSearch(args map[string]interface{}) (string, bool) {
	query := getStringArg(args, "query", "")
	topK := getIntArg(args, "top_k", 5)
	hybridMode := getBoolArg(args, "hybrid", true)
	vectorWeight := getFloatArg(args, "vector_weight", 0.7)

	if query == "" {
		return "Error: 'query' is required", false
	}

	queryFP := computeSimhash(query)

	var results []hybridResult

	if hybridMode {
		results = vecIndex.hybridSearch(queryFP, query, topK, vectorWeight)
	} else {
		simResults := vecIndex.vecSearch(queryFP, topK)
		for _, sr := range simResults {
			results = append(results, hybridResult{
				Chunk:  sr.Chunk,
				SHSore: sr.Score,
				Final:  sr.Score,
			})
		}
	}

	if len(results) == 0 {
		return "No matching results found. Try indexing more files first.", true
	}

	var sb stringsBuilder
	sb.WriteString(fmt.Sprintf("🔍 Found %d results for \"%s\" (%s):\n\n", len(results), query,
		map[bool]string{true: "hybrid", false: "simhash only"}[hybridMode]))
	for i, r := range results {
		preview := r.Chunk.Content
		if len(preview) > 300 {
			preview = preview[:300] + "..."
		}
		sourceLabel := r.Chunk.Source
		if len(sourceLabel) > 60 {
			sourceLabel = "..." + sourceLabel[len(sourceLabel)-57:]
		}
		sb.WriteString(fmt.Sprintf("**%d.** `%.3f` | `%s`\n```\n%s\n```\n\n",
			i+1, r.Final, sourceLabel, preview))
	}

	return sb.String(), true
}

func ragVecList(args map[string]interface{}) (string, bool) {
	sources := vecIndex.vecListSources()
	total := len(vecIndex.Chunks)
	if total == 0 {
		return "Vector index is empty. Use rag_index_add to add files.", true
	}

	var sb stringsBuilder
	sb.WriteString(fmt.Sprintf("📚 Vector index — %d chunks total:\n\n", total))
	for i, s := range sources {
		sb.WriteString(fmt.Sprintf("%d. `%s` (%d chunks)\n", i+1, s.Source, s.Chunks))
	}
	sb.WriteString(fmt.Sprintf("\n💡 _Index methods: simhash (semantic), tfidf (keyword), hybrid (combined)_"))
	return sb.String(), true
}

func ragVecRemove(args map[string]interface{}) (string, bool) {
	source := getStringArg(args, "source", "")
	if source == "" {
		return "Error: 'source' is required", false
	}
	removed := vecIndex.vecRemoveSource(source)
	if removed == 0 {
		return fmt.Sprintf("No chunks found for source: %s", source), false
	}
	vecIndex.persist()
	return fmt.Sprintf("✅ Removed %d chunks from %s", removed, source), true
}

func ragVecProvider(args map[string]interface{}) (string, bool) {
	return "Vector RAG provider: simhash (pure Go, zero dependencies). Embedding via 9router not yet configured.", true
}

// Helper: stringsToLower for a single string (avoid importing strings just for this)
func stringsToLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}

// stringsBuilder is a thin wrapper for efficient string concat.
type stringsBuilder struct {
	parts []string
}

func (b *stringsBuilder) WriteString(s string) {
	b.parts = append(b.parts, s)
}

func (b *stringsBuilder) String() string {
	total := 0
	for _, p := range b.parts {
		total += len(p)
	}
	result := make([]byte, 0, total)
	for _, p := range b.parts {
		result = append(result, []byte(p)...)
	}
	return string(result)
}
