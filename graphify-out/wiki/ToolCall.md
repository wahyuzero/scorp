# ToolCall

> 14 nodes

## Key Concepts

- **ToolCall** (20 connections) — `models/callback.go`
- **CallModelWithToolsAndFallback()** (16 connections) — `models/tools.go`
- **CallModelWithTools()** (9 connections) — `models/tools.go`
- **tools.go** (6 connections) — `models/tools.go`
- **ParseAllToolCalls()** (6 connections) — `models/tools.go`
- **ParseToolCalls()** (4 connections) — `models/tools.go`
- **TestParseToolCalls()** (3 connections) — `agent/prompt_test.go`
- **IsRateLimitError()** (3 connections) — `models/tools.go`
- **ParseCodeBlockFallback()** (3 connections) — `models/tools.go`
- **tools_test.go** (3 connections) — `models/tools_test.go`
- **TestCallModelWithToolsNilModel()** (3 connections) — `models/tools_test.go`
- **TestIsRateLimitError()** (3 connections) — `models/tools_test.go`
- **models/callback.go** (2 connections) — `models/callback.go`
- **CustomProvider** (2 connections) — `models/callback.go`

## Relationships

- [RunAgentSessionLoop](RunAgentSessionLoop.md) (6 shared connections)
- [ChatMessage](ChatMessage.md) (5 shared connections)
- [context.Context](context.Context.md) (5 shared connections)
- [testing.T](testing.T.md) (4 shared connections)
- [TruncateStr](TruncateStr.md) (3 shared connections)
- [CallCommandCodeWithTools](CallCommandCodeWithTools.md) (3 shared connections)
- [ExecuteTool](ExecuteTool.md) (2 shared connections)
- [api_gemini.go](api_gemini.go.md) (2 shared connections)
- [runSubagent](runSubagent.md) (2 shared connections)
- [cost_router.go](cost_router.go.md) (2 shared connections)
- [SaveModelConfig](SaveModelConfig.md) (1 shared connections)
- [GetProvider](GetProvider.md) (1 shared connections)

## Source Files

- `agent/prompt_test.go`
- `models/callback.go`
- `models/tools.go`
- `models/tools_test.go`

## Audit Trail

- EXTRACTED: 52 (85%)
- INFERRED: 9 (15%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*