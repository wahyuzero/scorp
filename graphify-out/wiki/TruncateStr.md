# TruncateStr

> 95 nodes · cohesion 0.07

## Key Concepts

- **TruncateStr()** (36 connections) — `internal/helpers/helpers.go`
- **ModelConfig** (28 connections) — `models/model_router.go`
- **model_router.go** (26 connections) — `models/model_router.go`
- **context.Context** (21 connections)
- **ChatMessage** (21 connections) — `models/model_router.go`
- **CallCommandCodeWithTools()** (17 connections) — `models/api_commandcode.go`
- **CallModel()** (16 connections) — `models/model_router.go`
- **KeySourceLabel()** (16 connections) — `models/providers.go`
- **ResolveAPIKey()** (16 connections) — `models/providers.go`
- **CallModelStream()** (14 connections) — `models/model_router.go`
- **ToolCall** (14 connections) — `models/callback.go`
- **CallModelWithToolsAndFallback()** (14 connections) — `models/tools.go`
- **CallAnthropicWithTools()** (13 connections) — `models/api_anthropic.go`
- **api_commandcode.go** (13 connections) — `models/api_commandcode.go`
- **api_gemini.go** (13 connections) — `models/api_gemini.go`
- **CallCommandCodeStream()** (12 connections) — `models/api_commandcode.go`
- **CallModelWithFallback()** (12 connections) — `models/model_router.go`
- **CallModelWithTools()** (12 connections) — `models/tools.go`
- **GetAIClient()** (12 connections) — `models/transports.go`
- **callAnthropic()** (11 connections) — `models/api_anthropic.go`
- **geminiDoRequest()** (11 connections) — `models/api_gemini.go`
- **CallOpenAI()** (11 connections) — `models/model_router.go`
- **CallOpenAIWithTools()** (11 connections) — `models/tools.go`
- **CallGeminiWithTools()** (10 connections) — `models/api_gemini.go`
- **RecordCost()** (10 connections) — `models/cost_router.go`
- *... and 70 more nodes in this community*

## Relationships

- [chat.go](chat.go.md) (17 shared connections)
- [startCLI](startCLI.md) (15 shared connections)
- [time.Time](time.Time.md) (9 shared connections)
- [runSubagent](runSubagent.md) (8 shared connections)
- [client.go](client.go.md) (7 shared connections)
- [testing.T](testing.T.md) (7 shared connections)
- [GetStringArg](GetStringArg.md) (6 shared connections)
- [ConfigMgr](ConfigMgr.md) (6 shared connections)
- [wizard.go](wizard.go.md) (6 shared connections)
- [config_paths.go](config_paths.go.md) (5 shared connections)
- [runSelfReview](runSelfReview.md) (4 shared connections)
- [collector_system_native.go](collector_system_native.go.md) (3 shared connections)

## Source Files

- `internal/helpers/helpers.go`
- `models/api_anthropic.go`
- `models/api_commandcode.go`
- `models/api_gemini.go`
- `models/callback.go`
- `models/cost_router.go`
- `models/model_router.go`
- `models/providers.go`
- `models/tools.go`
- `models/transports.go`
- `registry/registry.go`
- `tools/native.go`
- `tools/provider.go`
- `tools/vision.go`

## Audit Trail

- EXTRACTED: 338 (84%)
- INFERRED: 63 (16%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*