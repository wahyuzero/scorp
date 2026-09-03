package models

import (
	"context"
	"errors"
	"testing"

	"scorp-agent/registry"
)

func TestIsRateLimitError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil error", nil, false},
		{"non-matching error", errors.New("something else"), false},
		{"HTTP 429", errors.New("HTTP 429 Too Many Requests"), true},
		{"HTTP 402", errors.New("HTTP 402 Payment Required"), true},
		{"rate_limit string", errors.New("rate_limit exceeded"), true},
		{"rate limit string", errors.New("rate limit exceeded"), true},
		{"quota string", errors.New("quota exceeded"), true},
		{"suspicious activity", errors.New("suspicious activity detected"), true},
		{"wrapped 429", errors.New("upstream error: HTTP 429"), true},
		{"case sensitive rate", errors.New("RATE_LIMIT"), false},
		{"similar but not exact", errors.New("rated as good"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsRateLimitError(tt.err)
			if got != tt.want {
				t.Errorf("IsRateLimitError(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

func TestCallModelWithToolsNilModel(t *testing.T) {
	ctx := context.Background()
	_, _, err := CallModelWithTools(ctx, nil, nil)
	if err == nil {
		t.Error("CallModelWithTools with nil model expected error, got nil")
	}
}

func TestGenerateNativeToolsSchema(t *testing.T) {
	registry.RegisterTool(registry.ToolDef{
		Name:        "test_schema_tool",
		Description: "A test tool",
		Native:      true,
	})
	registry.ResetNativeToolCache()
	schema := registry.GenerateNativeToolsSchema()
	if len(schema) == 0 {
		t.Error("GenerateNativeToolsSchema returned empty schema")
	}
	registry.UnregisterTool("test_schema_tool")
	registry.ResetNativeToolCache()
}
