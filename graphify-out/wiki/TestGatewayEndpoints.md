# TestGatewayEndpoints

> 15 nodes

## Key Concepts

- **TestGatewayEndpoints()** (9 connections) — `gateway/gateway_test.go`
- **gateway.go** (8 connections) — `gateway/gateway.go`
- **net/http.Request** (8 connections)
- **net/http.ResponseWriter** (7 connections)
- **WebhookHandler()** (7 connections) — `telegram/telegram.go`
- **handleChat()** (6 connections) — `gateway/gateway.go`
- **contextWithTimeout()** (5 connections) — `gateway/gateway.go`
- **handleSOPs()** (5 connections) — `gateway/gateway.go`
- **handleStatus()** (5 connections) — `gateway/gateway.go`
- **handleTools()** (5 connections) — `gateway/gateway.go`
- **handleDashboard()** (4 connections) — `gateway/gateway.go`
- **handleReceipts()** (4 connections) — `gateway/gateway.go`
- **TestWebhookHandlerMalformedJSON()** (3 connections) — `telegram/telegram_test.go`
- **telegram_test.go** (2 connections) — `telegram/telegram_test.go`
- **gateway_test.go** (1 connections) — `gateway/gateway_test.go`

## Relationships

- [startCLI](startCLI.md) (3 shared connections)
- [testing.T](testing.T.md) (2 shared connections)
- [metasearch_engines.go](metasearch_engines.go.md) (1 shared connections)
- [MCPServer](MCPServer.md) (1 shared connections)
- [context.Context](context.Context.md) (1 shared connections)
- [ToolCall](ToolCall.md) (1 shared connections)
- [ExecuteTool](ExecuteTool.md) (1 shared connections)
- [RecordToolReceipt](RecordToolReceipt.md) (1 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (1 shared connections)
- [registry/registry.go](registry-registry.go.md) (1 shared connections)
- [time.Time](time.Time.md) (1 shared connections)
- [SaveModelConfig](SaveModelConfig.md) (1 shared connections)

## Source Files

- `gateway/gateway.go`
- `gateway/gateway_test.go`
- `telegram/telegram.go`
- `telegram/telegram_test.go`

## Audit Trail

- EXTRACTED: 45 (90%)
- INFERRED: 5 (10%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*