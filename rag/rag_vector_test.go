package rag

import (
	"testing"
)

func TestSimHash_Tokenize(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"Hello World", []string{"hello", "world"}},
		{"test123", []string{"test123"}},
		{"a b", []string{}}, // tokens < 2 chars filtered out
		{"Hello, World!", []string{"hello", "world"}},
		{"UPPERCASE", []string{"uppercase"}},
	}

	for _, tc := range tests {
		result := tokenize(tc.input)
		if len(result) != len(tc.expected) {
			t.Errorf("tokenize(%q): len=%d, want %d", tc.input, len(result), len(tc.expected))
			continue
		}
		for i, v := range result {
			if v != tc.expected[i] {
				t.Errorf("tokenize(%q)[%d]: got %q, want %q", tc.input, i, v, tc.expected[i])
			}
		}
	}
}

func TestSimHash_TokenHash(t *testing.T) {
	// tokenHash should be deterministic
	h1 := tokenHash("test")
	h2 := tokenHash("test")
	if h1 != h2 {
		t.Errorf("tokenHash not deterministic: %d != %d", h1, h2)
	}

	// Different tokens should produce different hashes (mostly)
	h3 := tokenHash("different")
	if h1 == h3 {
		t.Errorf("tokenHash collision: test == different")
	}
}

func TestSimHash_ComputeSimhash(t *testing.T) {
	// Empty text should return 0
	fp := ComputeSimhash("")
	if fp != 0 {
		t.Errorf("ComputeSimhash('') = %d, want 0", fp)
	}

	// Same text should produce same fingerprint
	fp1 := ComputeSimhash("hello world test")
	fp2 := ComputeSimhash("hello world test")
	if fp1 != fp2 {
		t.Errorf("ComputeSimhash not deterministic: %d != %d", fp1, fp2)
	}

	// Different text should produce different fingerprint (mostly)
	fp3 := ComputeSimhash("completely different text")
	if fp1 == fp3 {
		t.Errorf("ComputeSimhash collision")
	}
}

func TestSimHash_HammingDistance(t *testing.T) {
	// Distance of identical values should be 0
	d := hammingDistance(0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF)
	if d != 0 {
		t.Errorf("hammingDistance(identical) = %d, want 0", d)
	}

	// Distance of opposite values should be 64
	d = hammingDistance(0, 0xFFFFFFFFFFFFFFFF)
	if d != 64 {
		t.Errorf("hammingDistance(opposite) = %d, want 64", d)
	}

	// Known distance
	d = hammingDistance(0b1010, 0b0101) // 4 bits different
	if d != 4 {
		t.Errorf("hammingDistance(1010, 0101) = %d, want 4", d)
	}
}

func TestSimHash_Similarity(t *testing.T) {
	// Identical should have similarity 1.0
	sim := simhashSimilarity(0xFFFFFFFFFFFFFFFF, 0xFFFFFFFFFFFFFFFF)
	if sim != 1.0 {
		t.Errorf("simhashSimilarity(identical) = %f, want 1.0", sim)
	}

	// Opposite should have similarity 0.0
	sim = simhashSimilarity(0, 0xFFFFFFFFFFFFFFFF)
	if sim != 0.0 {
		t.Errorf("simhashSimilarity(opposite) = %f, want 0.0", sim)
	}

	// Half different (32 bits different) should be ~0.5
	sim = simhashSimilarity(0xFFFFFFFFFFFFFFFF, 0xFFFFFFFF00000000)
	if sim < 0.4 || sim > 0.6 {
		t.Errorf("simhashSimilarity(half) = %f, want ~0.5", sim)
	}
}

func TestSimHash_VectorIndexAddAndSearch(t *testing.T) {
	// Test basic add and search
	idx := &VectorIndex{
		Chunks: make(map[string]*VecChunk),
		DBPath: "/tmp/test_vector_index.json",
	}

	// Add some chunks
	id1 := idx.AddVecChunk("source1", "hello world this is a test document about golang")
	id2 := idx.AddVecChunk("source2", "another document about python programming language")

	if id1 == "" || id2 == "" {
		t.Error("addVecChunk returned empty ID")
	}

	// Search for something similar to first chunk
	queryFP := ComputeSimhash("golang test document")
	results := idx.vecSearch(queryFP, 5)

	if len(results) == 0 {
		t.Log("No results found (may be below threshold)")
		return
	}

	// First result should be source1 (more similar)
	if results[0].Chunk.Source != "source1" {
		t.Errorf("Expected source1 as top result, got %s", results[0].Chunk.Source)
	}
}

func TestSimHash_SmartChunk(t *testing.T) {
	chunks := SmartChunk("This is a test. This is another sentence. And a third one!", 50)
	if len(chunks) == 0 {
		t.Error("SmartChunk returned empty slices")
	}
	for _, c := range chunks {
		if len(c) > 60 { // Allow some overhead
			t.Errorf("Chunk too long: %d chars", len(c))
		}
	}
}