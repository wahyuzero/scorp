package models

// ToolCall represents a parsed tool call from LLM response
type ToolCall struct {
	Name string                 `json:"name"`
	Args map[string]interface{} `json:"args"`
	// AutoDecision (P3.13): preset by the agent loop's auto-mode gate so
	// ExecuteTool trusts the classification instead of re-grading. Never
	// marshaled into transcripts.
	AutoDecision string `json:"-"`
}

var UpdateEnvFile func(key, value string)

// CustomProvider stored in models.json under "custom_providers"
type CustomProvider struct {
	BaseURL     string   `json:"base_url"`
	API         string   `json:"api"`
	KeyEnvs     []string `json:"key_envs,omitempty"`
	DisplayName string   `json:"display_name,omitempty"`
	NoAuth      bool     `json:"no_auth,omitempty"`
}
