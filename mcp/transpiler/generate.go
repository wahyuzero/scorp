package transpiler

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"scorp-agent/models"
)

// codegenSystemPrompt is the contract-driven codegen contract for mark3labs/mcp-go.
// It deliberately forbids everything the Layer 2 security audit rejects
// (raw shell exec, secret exfiltration patterns) so generated code passes CI.
const codegenSystemPrompt = `You are the Scorp AI Transpiler, an expert Go engineer.
You transpile Model Context Protocol (MCP) servers into idiomatic Go using the
official github.com/mark3labs/mcp-go SDK (v1.0.0 API).

RESPONSE MODE (most important):
- You have NO tools, no shell, no file access, and cannot browse anything.
- Do NOT announce intentions ("I'll examine...", "Let me check...").
- Your FIRST token must be the Go source. Produce the COMPLETE final file in
  this single response.

OUTPUT CONTRACT (non-negotiable):
- Emit exactly ONE complete Go source file: package main, compile-ready.
- No markdown fences, no commentary outside the code, no build tags.
- Every tool from the provided benchmark MUST be implemented with the SAME
  tool name, parameter names, JSON types, and required fields.
- Mirror the upstream's command-line argument handling (benchmark.args shows
  how it was booted) — e.g. scoped directory allow-lists must be accepted the
  same way so the binary boots identically.
- Implement real logic where the benchmark makes intent clear (HTTP calls via
  net/http, file access via os/filepath, SQLite via modernc.org/sqlite, etc).
  Do NOT stub tools with placeholder returns unless the upstream behavior is
  unknowable from the contract; prefer a faithful best-effort implementation.

SDK CONTRACT (mcp-go v1.0.0) — AUTHORITATIVE API CHEAT-SHEET.
Use ONLY these symbols; inventing any other function fails compilation.

  server.NewMCPServer(name, version string, opts ...server.ServerOption) *server.MCPServer
  server.ServeStdio(s *server.MCPServer) error
  mcp.NewTool(name string, opts ...mcp.ToolOption) mcp.Tool
  s.AddTool(tool, handler) with handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error)
  Results: mcp.NewToolResultText(s string), mcp.NewToolResultError(s string), mcp.NewToolResultErrorf(format, args...)

  Param TYPE options (pick by benchmark JSON type — TYPE FIDELITY REQUIRED):
    "string"  → mcp.WithString(name, opts...)
    "number"  → mcp.WithNumber(name, opts...)
    "integer" → mcp.WithInteger(name, opts...)
    "boolean" → mcp.WithBoolean(name, opts...)
    "array"   → mcp.WithArray(name, opts..., mcp.WithStringItems())  // also WithNumberItems/WithStringEnumItems
    "object"  → mcp.WithObject(name, opts...)
  Param PROPERTY options (exact names, no others exist):
    mcp.Description(s), mcp.Required(), mcp.Title(s)
    mcp.DefaultString(v), mcp.DefaultBool(v)   // NO generic mcp.Default, NO numeric default
    mcp.Enum(vals ...string), mcp.MaxLength(n), mcp.MinLength(n), mcp.Pattern(re)
    (there are NO numeric Min/Max constraint options — enforce ranges in code)

  Reading params in the handler (EXACT signatures — copy these shapes):
    v, err := req.RequireString("k")          // returns (string, error)
    s, err := req.RequireStringSlice("k")     // returns ([]string, error)
    n, err := req.RequireInt("k")             // returns (int, error)
    x := req.GetString("k", "default")        // 2 args, ONE return, NO error
    i := req.GetInt("k", 0)                   // 2 args, ONE return, NO error
    f := req.GetFloat("k", 0.0)               // 2 args, ONE return, NO error
    b := req.GetBool("k", false)              // 2 args, ONE return, NO error
    sl := req.GetStringSlice("k", nil)        // 2 args, ONE return, NO error
  NEVER write "x, err := req.GetX(...)" — Get* NEVER returns an error.
  ONLY Require* returns (value, error), and Require* takes just the key.

  WIRING EXAMPLE:
    s := server.NewMCPServer("name", "1.0.0", server.WithToolCapabilities(false), server.WithRecovery())
    if err := server.ServeStdio(s); err != nil { fmt.Fprintf(os.Stderr, "%v\n", err); os.Exit(1) }
  Imports: "github.com/mark3labs/mcp-go/mcp" and "github.com/mark3labs/mcp-go/server".

  COMMON HALLUCINATIONS (these symbols DO NOT exist — never emit them):
    req.GetParamString / req.GetParamInt / req.GetParamNumber → use req.GetString / req.GetInt / req.GetFloat
    mcp.Default / mcp.DefaultInt / mcp.DefaultNumber → use mcp.DefaultString(v) or mcp.DefaultBool(v)
    mcp.NewServer / mcp.NewMCPServer → use server.NewMCPServer(...)
    (s *MCPServer).ServeStdio() method → use the package function server.ServeStdio(s)
    mcp.Minimum / mcp.Maximum / mcp.ExclusiveMinimum / mcp.ExclusiveMaximum → do not exist; enforce ranges in handler code
    mcp.WithRequired / mcp.Optional → do not exist; required-ness is mcp.Required(), optional is the absence of it

SECURITY CONTRACT (Layer 2 zero-trust audit — violations fail CI):
- NEVER use os/exec, syscall.StartProcess, or shell invocation of any kind.
- No embedding of credentials, tokens, or secret Material of any kind.
- All network access via net/http with context timeouts and bounded response
  readers (io.LimitReader). Filesystem access must validate paths against the
  configured scope. Never write diagnostics to stdout (it is the JSON-RPC
  channel); use os.Stderr.
- Prefer the Go standard library; third-party imports are limited to
  github.com/mark3labs/mcp-go plus, when the tool contract requires it:
  modernc.org/sqlite, golang.org/x/net/html, github.com/PuerkitoBio/goquery.

STYLE: small focused functions, doc comments on exported behavior, errors
wrapped with %w, context propagated everywhere, no globals except the server.`

// codegenMaxTokens overrides the per-model output cap: reasoning models can
// exhaust small caps (e.g. 4096) before emitting any code, which comes back
// as an empty message. Codegen needs a generous budget.
const codegenMaxTokens = 16384

// maxGenerateAttempts bounds retries per candidate model (first shot plus
// corrective retries for prose/empty responses).
const maxGenerateAttempts = 3

// Generate produces the transpiled main.go for a captured benchmark via the
// internal LLM gateway. Routing: complex/premium first, chat model as
// fallback — bypassing cost/off-peak rerouting so codegen quality is never
// silently downgraded mid-pipeline.
func Generate(ctx context.Context, bench *Benchmark, extraHints string) (string, error) {
	benchJSON, err := MarshalBenchmark(bench)
	if err != nil {
		return "", fmt.Errorf("serialize benchmark: %w", err)
	}

	userPrompt := fmt.Sprintf(`Transpile this MCP server contract into ONE Go file.

## Captured contract benchmark (authoritative — tools/list output of the live upstream server)
%s

## Additional implementation hints
%s

Remember: single compile-ready main.go, no markdown fences, every benchmark
tool implemented faithfully, security contract respected.`, benchJSON, extraHints)

	messages := []models.ChatMessage{
		{Role: "system", Content: codegenSystemPrompt},
		{Role: "user", Content: userPrompt},
	}

	genCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	var candidates []*models.ModelConfig
	for _, taskType := range []string{"complex", "chat"} {
		if m := models.RouteModel(taskType); m != nil {
			clone := *m
			if clone.MaxTokens < codegenMaxTokens {
				clone.MaxTokens = codegenMaxTokens
			}
			// Skip duplicates (complex may resolve to the chat model).
			dup := false
			for _, c := range candidates {
				if c.Model == clone.Model && c.Provider == clone.Provider {
					dup = true
					break
				}
			}
			if !dup {
				candidates = append(candidates, &clone)
			}
		}
	}
	if len(candidates) == 0 {
		return "", fmt.Errorf("no models configured for codegen (check models.json)")
	}

// codegenRetryPrompt is appended for one corrective attempt when a model
// answers with prose instead of code (common with agentic-tuned models).
const codegenRetryPrompt = `Your previous response was prose, not code. You have NO tools and cannot examine or browse anything. Respond NOW with the complete Go source file only — your first characters must be "package main".`

	var lastErr error
	for _, model := range candidates {
		// Generation with this class of models is flaky: prose preambles,
		// empty completions, and refusals all happen transiently. Give each
		// candidate up to maxGenerateAttempts shots (first + corrective).
		for attempt := 0; attempt < maxGenerateAttempts; attempt++ {
			msgs := messages
			if attempt > 0 {
				msgs = append(append([]models.ChatMessage{}, messages...),
					models.ChatMessage{Role: "user", Content: codegenRetryPrompt})
			}
			source, err := models.CallModel(genCtx, model, msgs)
			if err != nil {
				lastErr = err
				continue
			}
			cleaned := stripCodeFences(source)
			if strings.Contains(cleaned, "func main(") {
				return cleaned, nil
			}
			lastErr = fmt.Errorf("model %s returned non-program output (%d bytes received): %.160q",
				model.Model, len(source), strings.TrimSpace(source))
		}
	}
	if lastErr != nil {
		return "", fmt.Errorf("codegen failed across %d candidate model(s): %w", len(candidates), lastErr)
	}
	return "", fmt.Errorf("codegen failed: no models available")
}

var (
	// Fences are only recognized at the START OF A LINE. Generated Go code
	// legitimately contains inline triple-backticks inside string literals
	// (e.g. a WriteString emitting markdown code fences), which must not be
	// mistaken for a closing fence.
	fenceOpen  = regexp.MustCompile("(?m)^```[a-zA-Z]*[ \t]*\n")
	fenceClose = regexp.MustCompile("(?m)^```[ \t]*$")
)

// stripCodeFences removes markdown fencing models sometimes add despite instructions.
func stripCodeFences(s string) string {
	s = strings.TrimSpace(s)
	open := fenceOpen.FindStringIndex(s)
	if open == nil {
		return s // no fence: the model may have emitted bare Go already
	}
	contentStart := open[1]

	closes := fenceClose.FindAllStringIndex(s[contentStart:], -1)
	if len(closes) == 0 {
		// Unterminated fence: keep everything after the opening fence and
		// let the compiler catch any real truncation.
		return strings.TrimSpace(s[contentStart:])
	}
	return strings.TrimSpace(s[contentStart : contentStart+closes[len(closes)-1][0]])
}

// Repair runs one self-healing pass: the compiler/build errors and the
// offending source are fed back so the model emits a corrected file.
func Repair(ctx context.Context, bench *Benchmark, brokenSource, buildErr string) (string, error) {
	benchJSON, err := MarshalBenchmark(bench)
	if err != nil {
		return "", err
	}

	userPrompt := fmt.Sprintf(`The following Go program you generated for this MCP contract does not compile.

## Contract benchmark
%s

## Broken source
%s

## Compiler / build errors
%s

## Fix instructions
Every "undefined: X" error means you invented a symbol. Cross-check each one
against the API cheat-sheet in your instructions and replace it with the exact
listed equivalent (e.g. GetParamString→GetString, mcp.DefaultInt→remove or
mcp.DefaultString). Fix ONLY what the compiler flagged plus directly related
lines — keep everything else byte-identical.

Return the COMPLETE corrected Go file (package main, compile-ready, no
markdown fences). Do not explain anything — first characters must be
"package main".`, benchJSON, brokenSource, buildErr)

	genCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	m := models.RouteModel("complex")
	if m == nil {
		m = models.RouteModel("chat")
	}
	if m == nil {
		return "", fmt.Errorf("no models configured for repair")
	}
	clone := *m
	if clone.MaxTokens < codegenMaxTokens {
		clone.MaxTokens = codegenMaxTokens
	}

	source, err := models.CallModel(genCtx, &clone, []models.ChatMessage{
		{Role: "system", Content: codegenSystemPrompt},
		{Role: "user", Content: userPrompt},
	})
	if err != nil {
		return "", err
	}
	cleaned := stripCodeFences(source)
	if !strings.Contains(cleaned, "func main(") {
		return "", fmt.Errorf("repair returned non-program output (%d bytes)", len(source))
	}
	return cleaned, nil
}

