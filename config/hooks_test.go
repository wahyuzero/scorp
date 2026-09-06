package config

import (
	"os"
	"testing"
)

func TestParseHookSpec(t *testing.T) {
	t.Run("valid spec with colons in command", func(t *testing.T) {
		h, err := ParseHookSpec("shell:logger.sh --event pre:tool")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if h.Pattern != "shell" || h.Command != "logger.sh --event pre:tool" {
			t.Errorf("got pattern=%q command=%q", h.Pattern, h.Command)
		}
	})

	t.Run("empty pattern defaults to wildcard", func(t *testing.T) {
		h, err := ParseHookSpec(":audit.sh")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if h.Pattern != "*" {
			t.Errorf("expected pattern '*', got %q", h.Pattern)
		}
	})

	t.Run("no colon is invalid", func(t *testing.T) {
		if _, err := ParseHookSpec("audit.sh"); err == nil {
			t.Error("expected error for spec without colon")
		}
	})

	t.Run("empty command is invalid", func(t *testing.T) {
		if _, err := ParseHookSpec("shell:   "); err == nil {
			t.Error("expected error for empty command")
		}
	})
}

func TestHookMatches(t *testing.T) {
	cases := []struct {
		pattern, tool string
		want          bool
	}{
		{"*", "shell", true},
		{"", "read_file", true},
		{"shell", "shell", true},
		{"shell", "read_file", false},
		{"write_*", "write_file", true},
		{"write_*", "read_file", false},
	}
	for _, c := range cases {
		got := HookMatches(HookEntry{Pattern: c.pattern}, c.tool)
		if got != c.want {
			t.Errorf("HookMatches(pattern=%q, tool=%q) = %v, want %v", c.pattern, c.tool, got, c.want)
		}
	}
}

func TestLoadHooksFromEnv(t *testing.T) {
	t.Setenv("SCORP_HOOKS_PRE", "shell:audit.sh;not-a-valid-spec;write_file:policy.sh")
	t.Setenv("SCORP_HOOKS_POST", "*:log-result.sh")
	ReloadHooks()
	defer ReloadHooks()

	pre := PreToolHooks()
	if len(pre) != 2 {
		t.Fatalf("expected 2 valid pre hooks (invalid spec skipped), got %d: %+v", len(pre), pre)
	}
	if pre[0].Pattern != "shell" || pre[0].Command != "audit.sh" {
		t.Errorf("pre[0] = %+v", pre[0])
	}
	if pre[1].Pattern != "write_file" {
		t.Errorf("pre[1] = %+v", pre[1])
	}

	post := PostToolHooks()
	if len(post) != 1 || post[0].Pattern != "*" {
		t.Fatalf("post hooks = %+v", post)
	}
}

func TestReloadHooksClearsState(t *testing.T) {
	t.Setenv("SCORP_HOOKS_PRE", "shell:audit.sh")
	defer unsetHookEnv(t, "SCORP_HOOKS_PRE")
	ReloadHooks()
	if got := len(PreToolHooks()); got != 1 {
		t.Fatalf("expected 1 hook, got %d", got)
	}
	os.Unsetenv("SCORP_HOOKS_PRE")
	ReloadHooks()
	if got := len(PreToolHooks()); got != 0 {
		t.Errorf("expected 0 hooks after unset+reload, got %d", got)
	}
}

// unsetHookEnv removes an env var on test cleanup so other tests never see it.
func unsetHookEnv(t *testing.T, keys ...string) {
	t.Cleanup(func() {
		for _, k := range keys {
			os.Unsetenv(k)
		}
		ReloadHooks()
	})
}
