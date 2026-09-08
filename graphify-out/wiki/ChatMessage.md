# ChatMessage

> 19 nodes

## Key Concepts

- **ChatMessage** (30 connections) — `models/model_router.go`
- **CallOpenAIWithTools()** (14 connections) — `models/api_openai.go`
- **CallOpenAI()** (12 connections) — `models/api_openai.go`
- **RecordCostWithCache()** (7 connections) — `models/cost_router.go`
- **.CallWithTools()** (6 connections) — `models/api_anthropic.go`
- **api_anthropic.go** (6 connections) — `models/api_anthropic.go`
- **.CallWithTools()** (6 connections) — `models/api_openai.go`
- **TrackModelUsageWithCache()** (6 connections) — `models/telemetry.go`
- **.Call()** (5 connections) — `models/api_anthropic.go`
- **.Call()** (5 connections) — `models/api_openai.go`
- **AnthropicProvider** (4 connections) — `models/api_anthropic.go`
- **api_openai.go** (4 connections) — `models/api_openai.go`
- **formatOpenAIMessages()** (4 connections) — `models/api_openai.go`
- **OpenAIProvider** (4 connections) — `models/api_openai.go`
- **anthropicRequest** (3 connections) — `models/api_anthropic.go`
- **anthropicResponse** (2 connections) — `models/api_anthropic.go`
- **anthropicTool** (2 connections) — `models/api_anthropic.go`
- **.Format()** (1 connections) — `models/api_anthropic.go`
- **.Format()** (1 connections) — `models/api_openai.go`

## Relationships

- [TruncateStr](TruncateStr.md) (25 shared connections)
- [context.Context](context.Context.md) (10 shared connections)
- [CallCommandCodeWithTools](CallCommandCodeWithTools.md) (7 shared connections)
- [api_gemini.go](api_gemini.go.md) (6 shared connections)
- [ToolCall](ToolCall.md) (5 shared connections)
- [registry/registry.go](registry-registry.go.md) (2 shared connections)
- [cost_router.go](cost_router.go.md) (2 shared connections)
- [client.go](client.go.md) (1 shared connections)
- [ConfigMgr](ConfigMgr.md) (1 shared connections)
- [ScorpPath](ScorpPath.md) (1 shared connections)

## Source Files

- `models/api_anthropic.go`
- `models/api_openai.go`
- `models/cost_router.go`
- `models/model_router.go`
- `models/telemetry.go`

## Audit Trail

- EXTRACTED: 79 (87%)
- INFERRED: 12 (13%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*