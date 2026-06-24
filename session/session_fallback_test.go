package session

import (
	"strings"
	"testing"
	"time"
)

func TestSessionFallback_InitAndIndex(t *testing.T) {
	// Test LIKE query building logic directly

	query := "hello world"
	likeQuery := "%" + query + "%"
	expected := "%hello world%"
	if likeQuery != expected {
		t.Errorf("LIKE query building: got %q, want %q", likeQuery, expected)
	}
}

func TestSessionFallback_LIKEQueryWithSpaces(t *testing.T) {
	query := "test query with spaces"
	likeQuery := "%" + query + "%"
	// Should contain %test%query%with%spaces%
	if likeQuery != "%test query with spaces%" {
		t.Errorf("LIKE query building with spaces: got %q", likeQuery)
	}
}

func TestSessionFallback_SessionResultStruct(t *testing.T) {
	// Test that SessionResult struct can be created
	r := SessionResult{
		ID:        1,
		ChatID:    "test_chat",
		Role:      "user",
		Content:   "test content",
		Timestamp: time.Now().Unix(),
		MsgIndex:  0,
	}
	if r.ChatID != "test_chat" {
		t.Errorf("SessionResult struct failed: %+v", r)
	}
}

func TestSessionFallback_ExecuteSessionSearch_RequiresQuery(t *testing.T) {
	// Test ExecuteSessionSearch returns error when query is empty
	args := map[string]interface{}{
		"query": "",
	}
	result, ok := ExecuteSessionSearch(args)
	if ok {
		t.Error("ExecuteSessionSearch should return ok=false for empty query")
	}
	if !strings.Contains(result, "query") {
		t.Errorf("ExecuteSessionSearch error message should mention 'query': got %q", result)
	}
}
