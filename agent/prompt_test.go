package agent

import (
	"scorp-agent/models"
	"scorp-agent/internal/helpers"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func TestGetStringArg(t *testing.T) {
	args := map[string]interface{}{
		"name":  "test",
		"count": 42,
		"flag":  true,
	}

	tests := []struct {
		name       string
		key        string
		defaultVal string
		want       string
	}{
		{"existing string key", "name", "default", "test"},
		{"missing key returns default", "missing", "default", "default"},
		{"non-string value returns default", "count", "default", "default"},
		{"bool value returns default", "flag", "default", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.GetStringArg(args, tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("helpers.GetStringArg(%q, %q) = %q, want %q", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestGetIntArg(t *testing.T) {
	args := map[string]interface{}{
		"float":  3.14,
		"int":    42,
		"string": "not-a-number",
	}

	tests := []struct {
		name       string
		key        string
		defaultVal int
		want       int
	}{
		{"float64 value", "float", 0, 3},
		{"int value", "int", 0, 42},
		{"string value returns default", "string", 99, 99},
		{"missing key returns default", "missing", 99, 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.GetIntArg(args, tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("helpers.GetIntArg(%q, %d) = %d, want %d", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestGetBoolArg(t *testing.T) {
	args := map[string]interface{}{
		"true":  true,
		"false": false,
		"str":   "yes",
	}

	tests := []struct {
		name       string
		key        string
		defaultVal bool
		want       bool
	}{
		{"bool true", "true", false, true},
		{"bool false", "false", true, false},
		{"string value returns default", "str", true, true},
		{"missing key returns default", "missing", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.GetBoolArg(args, tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("helpers.GetBoolArg(%q, %v) = %v, want %v", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestGetFloatArg(t *testing.T) {
	args := map[string]interface{}{
		"float64": 3.14,
		"float32": float32(2.5),
		"int":     42,
		"int64":   int64(100),
		"string":  "not-a-number",
	}

	tests := []struct {
		name       string
		key        string
		defaultVal float64
		want       float64
	}{
		{"float64 value", "float64", 0, 3.14},
		{"float32 value", "float32", 0, 2.5},
		{"int value", "int", 0, 42},
		{"int64 value", "int64", 0, 100},
		{"string value returns default", "string", 99, 99},
		{"missing key returns default", "missing", 99, 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.GetFloatArg(args, tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("helpers.GetFloatArg(%q, %v) = %v, want %v", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestGetStringSliceArg(t *testing.T) {
	args := map[string]interface{}{
		"slice": []interface{}{"a", "b", "c"},
		"mixed": []interface{}{"a", 1, "b"},
		"str":   "not-a-slice",
	}

	tests := []struct {
		name       string
		key        string
		want       []string
		wantNil    bool
	}{
		{"string slice", "slice", []string{"a", "b", "c"}, false},
		{"mixed slice filters non-strings", "mixed", []string{"a", "b"}, false},
		{"non-slice returns nil", "str", nil, true},
		{"missing key returns nil", "missing", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.GetStringSliceArg(args, tt.key)
			if tt.wantNil {
				if got != nil {
					t.Errorf("helpers.GetStringSliceArg(%q) = %v, want nil", tt.key, got)
				}
			} else {
				if len(got) != len(tt.want) {
					t.Errorf("helpers.GetStringSliceArg(%q) = %v, want %v", tt.key, got, tt.want)
				} else {
					for i := range got {
						if got[i] != tt.want[i] {
							t.Errorf("helpers.GetStringSliceArg(%q)[%d] = %q, want %q", tt.key, i, got[i], tt.want[i])
						}
					}
				}
			}
		})
	}
}

func TestGetInt64Arg(t *testing.T) {
	args := map[string]interface{}{
		"float64": 3.14,
		"int":     42,
		"int64":   int64(100),
		"jsonNum": json.Number("12345"),
		"string":  "not-a-number",
	}

	tests := []struct {
		name       string
		key        string
		defaultVal int64
		want       int64
	}{
		{"float64 value", "float64", 0, 3},
		{"int value", "int", 0, 42},
		{"int64 value", "int64", 0, 100},
		{"json.Number value", "jsonNum", 0, 12345},
		{"string value returns default", "string", 99, 99},
		{"missing key returns default", "missing", 99, 99},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.GetInt64Arg(args, tt.key, tt.defaultVal)
			if got != tt.want {
				t.Errorf("helpers.GetInt64Arg(%q, %d) = %d, want %d", tt.key, tt.defaultVal, got, tt.want)
			}
		})
	}
}

func TestIsDangerousCommand(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		want     bool
	}{
		{"rm -rf /", "rm -rf /", true},
		{"rm -rf /*", "rm -rf /*", true},
		{"mkfs", "mkfs.ext4 /dev/sda1", true},
		{"dd if=", "dd if=/dev/zero of=/dev/sda", true},
		{"fork bomb", ":(){ :|:& };:", true},
		{"drop table", "DROP TABLE users", true},
		{"drop database", "drop database prod", true},
		{"delete from", "DELETE FROM users WHERE id=1", true},
		{"kill -9", "kill -9 1234", true},
		{"killall", "killall nginx", true},
		{"pkill", "pkill -f chrome", true},
		{"systemctl stop", "systemctl stop nginx", true},
		{"systemctl disable", "systemctl disable docker", true},
		{"apt remove", "apt remove nginx", true},
		{"apt purge", "apt purge nginx", true},
		{"pip uninstall", "pip uninstall requests", true},
		{"docker rm", "docker rm container1", true},
		{"docker rmi", "docker rmi image1", true},
		{"docker prune is NOT in dangerous list", "docker system prune", false},
		{"docker-compose down", "docker-compose down", true},
		{"docker compose down", "docker compose down", true},
		{"> /dev/", "echo test > /dev/sda", true},
		{"chmod 777", "chmod 777 /etc/passwd", true},
		{"safe command", "ls -la", false},
		{"safe echo", "echo hello world", false},
		{"safe docker ps", "docker ps", false},
		{"safe systemctl status", "systemctl status nginx", false},
		{"case insensitive rm -RF /", "RM -RF /", true},
		{"case insensitive systemctl stop", "Systemctl Stop nginx", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsDangerousCommand(tt.cmd)
			if got != tt.want {
				t.Errorf("IsDangerousCommand(%q) = %v, want %v", tt.cmd, got, tt.want)
			}
		})
	}
}

func TestTruncOutput(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"short string unchanged", "hello", 10, "hello"},
		{"exact length unchanged", "hello", 5, "hello"},
		{"truncated with ellipsis", "hello world", 8, "hello wo\n... (truncated)"},
		{"empty string", "", 10, ""},
		{"maxLen zero", "hello", 0, "\n... (truncated)"},
		{"unicode truncated at byte boundary", "😀😁😂😃😄", 6, "😀\xf0\x9f\n... (truncated)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helpers.TruncOutput(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("helpers.TruncOutput(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func TestParseToolCalls(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantCalls  int
		wantClean  string
	}{
		{"no tool calls", "hello world", 0, "hello world"},
		{"single tool call", `start <tool_call>{"name": "shell", "args": {"command": "ls"}}</tool_call> end`, 1, "start  end"},
		{"multiple tool calls",
			`a<tool_call>{"name": "shell", "args": {"command": "ls"}}</tool_call>` +
				`b<tool_call>{"name": "read_file", "args": {"path": "/tmp/x"}}</tool_call>c`,
			2, "abc"},
		{"malformed json skipped", `x<tool_call>not-json</tool_call>y`, 0, "xy"},
		{"whitespace trimmed", " <tool_call>{\"name\": \"shell\", \"args\": {}}</tool_call> ", 1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls, clean := models.ParseToolCalls(tt.input)
			if len(calls) != tt.wantCalls {
				t.Errorf("models.ParseToolCalls calls = %d, want %d; input=%q", len(calls), tt.wantCalls, tt.input)
			}
			if clean != tt.wantClean {
				t.Errorf("models.ParseToolCalls clean = %q, want %q; input=%q", clean, tt.wantClean, tt.input)
			}
		})
	}
}

func TestToolDescription(t *testing.T) {
	tests := []struct {
		name string
		tc   ToolCall
		desc string
	}{
		{"shell", ToolCall{Name: "shell", Args: map[string]interface{}{"command": "ls -la"}}, "🖥 shell: ls -la"},
		{"shell long cmd truncated", ToolCall{Name: "shell", Args: map[string]interface{}{"command": "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"}},
			"🖥 shell: abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqr..."},
		{"read_file", ToolCall{Name: "read_file", Args: map[string]interface{}{"path": "/tmp/x"}}, "📖 read: /tmp/x"},
		{"write_file", ToolCall{Name: "write_file", Args: map[string]interface{}{"path": "/tmp/x"}}, "✏️ write: /tmp/x"},
		{"web_fetch", ToolCall{Name: "web_fetch", Args: map[string]interface{}{"url": "https://example.com"}}, "🌐 fetch: https://example.com"},
		{"web_search", ToolCall{Name: "web_search", Args: map[string]interface{}{"query": "test query"}}, "🔍 search: test query"},
		{"memory", ToolCall{Name: "memory", Args: map[string]interface{}{"action": "get", "key": "name"}}, "🧠 memory.get(name)"},
		{"browser goto", ToolCall{Name: "browser", Args: map[string]interface{}{"action": "goto", "url": "https://example.com"}}, "🌐 browser→https://example.com"},
		{"browser other", ToolCall{Name: "browser", Args: map[string]interface{}{"action": "snapshot"}}, "🌐 browser.snapshot"},
		{"mcp_tool", ToolCall{Name: "mcp_tool", Args: map[string]interface{}{"server": "fs", "tool": "list"}}, "🔌 mcp: fs.list"},
		{"delegate", ToolCall{Name: "delegate", Args: map[string]interface{}{"task": "do something"}}, "🤖 delegate: do something"},
		{"delegate long truncated", ToolCall{Name: "delegate", Args: map[string]interface{}{"task": "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz"}},
			"🤖 delegate: abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ01234567..."},
		{"unknown tool", ToolCall{Name: "unknown_tool", Args: nil}, "🔧 unknown_tool"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toolDescription(tt.tc)
			if got != tt.desc {
				t.Errorf("toolDescription(%v) = %q, want %q", tt.tc, got, tt.desc)
			}
		})
	}
}

func TestBuildThinkingMessage(t *testing.T) {
	lines := []string{"🔍 searching...", "📦 found 5 files"}
	msg := buildThinkingMessage(lines, 5*time.Second, false)
	if !strings.Contains(msg, "🔍 searching...") {
		t.Errorf("buildThinkingMessage missing line content: %s", msg)
	}
	if !strings.Contains(msg, "working...") {
		t.Errorf("buildThinkingMessage should show working indicator when not done: %s", msg)
	}

	msgDone := buildThinkingMessage(lines, 10*time.Second, true)
	if !strings.Contains(msgDone, "[10s]") {
		t.Errorf("buildThinkingMessage should include elapsed time: %s", msgDone)
	}
	if strings.Contains(msgDone, "working...") {
		t.Errorf("buildThinkingMessage should NOT show working indicator when done: %s", msgDone)
	}
}

func TestMaxIterations(t *testing.T) {
	t.Run("default value", func(t *testing.T) {
		// Unset env
		os.Unsetenv("SCORP_MAX_ITERATIONS")
		if got := maxIterations(); got != 20 {
			t.Errorf("maxIterations() = %d, want 20", got)
		}
	})

	t.Run("custom value", func(t *testing.T) {
		os.Setenv("SCORP_MAX_ITERATIONS", "5")
		defer os.Unsetenv("SCORP_MAX_ITERATIONS")
		if got := maxIterations(); got != 5 {
			t.Errorf("maxIterations() = %d, want 5", got)
		}
	})

	t.Run("invalid value returns default", func(t *testing.T) {
		os.Setenv("SCORP_MAX_ITERATIONS", "not-a-number")
		defer os.Unsetenv("SCORP_MAX_ITERATIONS")
		if got := maxIterations(); got != 20 {
			t.Errorf("maxIterations() = %d, want 20", got)
		}
	})

	t.Run("negative value returns default", func(t *testing.T) {
		os.Setenv("SCORP_MAX_ITERATIONS", "-5")
		defer os.Unsetenv("SCORP_MAX_ITERATIONS")
		if got := maxIterations(); got != 20 {
			t.Errorf("maxIterations() = %d, want 20", got)
		}
	})
}

func TestBase64Encode(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{"empty bytes", []byte{}, ""},
		{"hello", []byte("hello"), "aGVsbG8="},
		{"binary", []byte{0x00, 0x01, 0x02}, "AAEC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := base64Encode(tt.input)
			if got != tt.want {
				t.Errorf("base64Encode(%v) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
