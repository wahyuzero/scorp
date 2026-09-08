# TruncateStr

> 26 nodes

## Key Concepts

- **TruncateStr()** (39 connections) — `internal/helpers/helpers.go`
- **ModelConfig** (36 connections) — `models/config_io.go`
- **ResolveAPIKey()** (16 connections) — `models/providers.go`
- **KeySourceLabel()** (15 connections) — `models/providers.go`
- **CallModelStream()** (14 connections) — `models/model_router.go`
- **CallAnthropicWithTools()** (13 connections) — `models/api_anthropic.go`
- **CallCommandCodeStream()** (12 connections) — `models/api_commandcode.go`
- **callAnthropic()** (11 connections) — `models/api_anthropic.go`
- **GetAIClient()** (11 connections) — `models/transports.go`
- **providers.go** (10 connections) — `models/providers.go`
- **telemetry.go** (8 connections) — `models/telemetry.go`
- **RecordCost()** (7 connections) — `models/cost_router.go`
- **CheckModelHealth()** (7 connections) — `models/telemetry.go`
- **TrackModelUsage()** (7 connections) — `models/telemetry.go`
- **ResolveAPIFormat()** (6 connections) — `models/providers.go`
- **InitModelUsage()** (5 connections) — `models/telemetry.go`
- **hasAPIKey()** (4 connections) — `models/providers.go`
- **StreamChunk** (4 connections) — `models/model_router.go`
- **applyProviderDefaults()** (3 connections) — `models/providers.go`
- **resolveCommandCodeKeyFromDisk()** (3 connections) — `models/providers.go`
- **resolveOpenCodeKeyFromDisk()** (3 connections) — `models/providers.go`
- **FormatModelList()** (3 connections) — `models/telemetry.go`
- **FormatModelListWithHealth()** (3 connections) — `models/telemetry.go`
- **FormatUsageStats()** (3 connections) — `models/telemetry.go`
- **SwitchModel()** (3 connections) — `models/telemetry.go`
- *... and 1 more nodes in this community*

## Relationships

- [ChatMessage](ChatMessage.md) (25 shared connections)
- [api_gemini.go](api_gemini.go.md) (14 shared connections)
- [CallCommandCodeWithTools](CallCommandCodeWithTools.md) (12 shared connections)
- [context.Context](context.Context.md) (12 shared connections)
- [SaveModelConfig](SaveModelConfig.md) (10 shared connections)
- [time.Time](time.Time.md) (7 shared connections)
- [startCLI](startCLI.md) (6 shared connections)
- [ConfigMgr](ConfigMgr.md) (4 shared connections)
- [cost_router.go](cost_router.go.md) (4 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (3 shared connections)
- [chat.go](chat.go.md) (3 shared connections)
- [runSubagent](runSubagent.md) (3 shared connections)

## Source Files

- `internal/helpers/helpers.go`
- `models/api_anthropic.go`
- `models/api_commandcode.go`
- `models/config_io.go`
- `models/cost_router.go`
- `models/model_router.go`
- `models/providers.go`
- `models/telemetry.go`
- `models/transports.go`

## Audit Trail

- EXTRACTED: 143 (76%)
- INFERRED: 44 (24%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*