package eval

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"scorp-agent/config"

	"github.com/joho/godotenv"
)

// ──────────────────────────────────────────────
// Live cases — real agent tasks, independently verified (verify26 method).
// Each case spawns a fresh one-shot agent subprocess, then a checker inspects
// the artifacts on disk. The agent's final message is NEVER the source of
// truth; tokens per task come from model_usage.json deltas.
// ──────────────────────────────────────────────

const (
	liveTaskTimeout   = 6 * time.Minute
	liveSessionPrefix = "scorp-eval"
)

type liveCase struct {
	name    string
	prompt  string
	checker func(dir string) error
}

func LiveCases() []Case {
	lcs := []liveCase{
		{
			name:   "write_artifact_exact_content",
			prompt: "Buat file hello.txt di direktori kerjamu sekarang dengan konten PERSIS satu baris: eval-ok. Tidak ada teks lain di file itu. Lalu selesai.",
			checker: func(dir string) error {
				data, err := os.ReadFile(filepath.Join(dir, "hello.txt"))
				if err != nil {
					return fmt.Errorf("artifact missing: %v", err)
				}
				if strings.TrimSpace(string(data)) != "eval-ok" {
					return fmt.Errorf("content mismatch: %q", string(data))
				}
				return nil
			},
		},
		{
			name:   "go_module_tests_green",
			prompt: "Di direktori kerjamu sekarang buat Go module kecil bernama evaldemo (go 1.21) dengan add.go (fungsi Add(a, b int) int) dan add_test.go yang menguji Add(1,2)==3. Jalankan go test ./... dan pastikan hijau selesai.",
			checker: func(dir string) error {
				for _, f := range []string{"go.mod", "add.go", "add_test.go"} {
					if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
						return fmt.Errorf("%s missing", f)
					}
				}
				ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
				defer cancel()
				cmd := exec.CommandContext(ctx, "go", "test", "./...")
				cmd.Dir = dir
				if out, err := cmd.CombinedOutput(); err != nil {
					return fmt.Errorf("independent go test failed: %v (%.200s)", err, out)
				}
				return nil
			},
		},
		{
			name:   "durable_memory_write",
			prompt: fmt.Sprintf("Gunakan tool memory dengan action=remember untuk menyimpan satu fakta persis seperti ini: eval-marker-%d tersimpan. Lalu selesai.", time.Now().Unix()),
			checker: func(dir string) error {
				data, err := os.ReadFile(config.ScorpPath("MEMORY.md"))
				if err != nil {
					return fmt.Errorf("MEMORY.md unreadable: %v", err)
				}
				marker := fmt.Sprintf("eval-marker-%d", time.Now().Unix()-60)
				if !strings.Contains(string(data), "eval-marker-") {
					return fmt.Errorf("no eval marker found in MEMORY.md (wanted ~%s)", marker)
				}
				return nil
			},
		},
	}

	cases := make([]Case, 0, len(lcs))
	for _, lc := range lcs {
		lc := lc
		cases = append(cases, Case{
			Name:     lc.name,
			Category: "live-agent",
			Run:      func() error { return runLiveCase(lc) },
			Tokens:   func() int { return liveTokenBudget },
		})
	}
	return cases
}

// deploymentEnv builds the subprocess environment: the parent env plus the
// deployment .env candidates (same order as config.LoadConfig). Without this,
// a one-shot spawned outside systemd lacks the API keys systemd injects for
// the daemon.
func deploymentEnv() ([]string, error) {
	candidates := []string{config.ScorpPath(".env"), filepath.Join(config.HomeDir(), "scorp", ".env")}
	existing := map[string]bool{}
	for _, kv := range os.Environ() {
		if i := strings.Index(kv, "="); i > 0 {
			existing[kv[:i]] = true
		}
	}
	var extra []string
	var readErrs []string
	for _, path := range candidates {
		kv, err := godotenv.Read(path)
		if err != nil {
			readErrs = append(readErrs, fmt.Sprintf("%s: %v", path, err))
			// Lenient fallback: godotenv rejects the whole file on one
			// malformed line; recover the simple KEY=VALUE lines instead of
			// losing every API key over one bad quoting.
			data, rerr := os.ReadFile(path)
			if rerr != nil {
				continue
			}
			kv = map[string]string{}
			for _, line := range strings.Split(string(data), "\n") {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				if i := strings.Index(line, "="); i > 0 {
					kv[strings.TrimSpace(line[:i])] = strings.Trim(strings.TrimSpace(line[i+1:]), "\"'")
				}
			}
		}
		for k, v := range kv {
			if !existing[k] {
				existing[k] = true
				extra = append(extra, k+"="+v)
			}
		}
	}
	if extra == nil {
		return nil, fmt.Errorf("no deployment .env found (tried %v; errors: %v)", candidates, readErrs)
	}
	return append(os.Environ(), extra...), nil
}

// runLiveCase spawns `scorp -p <task> -s scorp-eval-<name>` in a fresh sandbox
// dir, then runs the independent checker. Token usage is captured from the
// model_usage.json delta for the tokens/task metric.
func runLiveCase(lc liveCase) error {
	dir, err := evalSandboxDir(strings.ReplaceAll(lc.name, "_", "-"))
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	exe, err := os.Executable()
	if err != nil {
		return err
	}

	env, err := deploymentEnv()
	if err != nil {
		return fmt.Errorf("live case needs deployment env: %v", err)
	}

	before := totalTokens()

	ctx, cancel := context.WithTimeout(context.Background(), liveTaskTimeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, exe, "-p", lc.prompt, "-s", fmt.Sprintf("%s-%s-%d", liveSessionPrefix, lc.name, time.Now().Unix()))
	cmd.Dir = dir
	cmd.Env = env
	out, err := cmd.CombinedOutput()

	if ctx.Err() == context.DeadlineExceeded {
		return fmt.Errorf("task timed out after %s (output tail: %.200s)", liveTaskTimeout, tail(out))
	}
	if err != nil {
		return fmt.Errorf("agent subprocess failed: %v (tail: %.200s)", err, tail(out))
	}

	if cerr := lc.checker(dir); cerr != nil {
		return fmt.Errorf("INDEPENDENT CHECK FAILED: %v", cerr)
	}

	tokens := totalTokens() - before
	if tokens < 0 {
		tokens = 0
	}
	liveTokenBudget = tokens // reported via report-side helper
	return nil
}

// liveTokenBudget carries the last live case's token usage to the report.
// (Eval runs sequentially, so one slot is enough.)
var liveTokenBudget int

func tokensReported() int { return liveTokenBudget }

// totalTokens walks model_usage.json summing every numeric field whose name
// mentions "token" — tolerant to schema drift.
func totalTokens() int {
	data, err := os.ReadFile(config.ScorpPath("model_usage.json"))
	if err != nil {
		return 0
	}
	var v interface{}
	if json.Unmarshal(data, &v) != nil {
		return 0
	}
	return sumTokenFields(v)
}

func sumTokenFields(v interface{}) int {
	switch t := v.(type) {
	case map[string]interface{}:
		sum := 0
		for k, sub := range t {
			if n, ok := sub.(float64); ok && strings.Contains(strings.ToLower(k), "token") {
				sum += int(n)
			} else {
				sum += sumTokenFields(sub)
			}
		}
		return sum
	case []interface{}:
		sum := 0
		for _, sub := range t {
			sum += sumTokenFields(sub)
		}
		return sum
	}
	return 0
}

func tail(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 200 {
		return s[len(s)-200:]
	}
	return s
}
