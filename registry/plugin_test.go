package registry

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

type echoPlugin struct{}

func (e *echoPlugin) Name() string {
	return "test_echo"
}

func (e *echoPlugin) Description() string {
	return "Echo back received message"
}

func (e *echoPlugin) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		Msg string `json:"msg"`
	}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return "Echo: " + params.Msg, nil
}

func TestRegisterPlugin(t *testing.T) {
	plugin := &echoPlugin{}
	RegisterPlugin(plugin)

	def, ok := GetTool("test_echo")
	if !ok {
		t.Fatalf("Expected test_echo tool to be registered")
	}
	if def.Name != "test_echo" {
		t.Errorf("Expected name 'test_echo', got %s", def.Name)
	}

	res, ok := ExecuteToolByName("test_echo", map[string]interface{}{"msg": "hello world"}, 0)
	if !ok {
		t.Fatalf("Expected execution to succeed")
	}
	if !strings.Contains(res, "Echo: hello world") {
		t.Errorf("Unexpected result: %s", res)
	}

	UnregisterTool("test_echo")
}
