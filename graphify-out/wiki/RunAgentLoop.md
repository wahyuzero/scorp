# RunAgentLoop

> 21 nodes · cohesion 0.17

## Key Concepts

- **bg.go** (12 connections) — `tools/bg.go`
- **ExecuteBgProcess()** (8 connections) — `tools/bg.go`
- **sync.Mutex** (7 connections)
- **getBGProcess()** (7 connections) — `tools/bg.go`
- **BGProcess** (7 connections) — `tools/bg.go`
- **StartTestEndpoint()** (5 connections) — `testutil/endpoint.go`
- **bgPoll()** (5 connections) — `tools/bg.go`
- **bgWait()** (4 connections) — `tools/bg.go`
- **os/exec.Cmd** (3 connections)
- **CostTracker** (3 connections) — `models/cost_router.go`
- **getChatLock()** (3 connections) — `testutil/endpoint.go`
- **bgKill()** (3 connections) — `tools/bg.go`
- **bgSpawn()** (3 connections) — `tools/bg.go`
- **bgWrite()** (3 connections) — `tools/bg.go`
- **io.WriteCloser** (2 connections)
- **endpoint.go** (2 connections) — `testutil/endpoint.go`
- **bgList()** (2 connections) — `tools/bg.go`
- **closeStdin()** (2 connections) — `tools/bg.go`
- **nextBGProcessID()** (2 connections) — `tools/bg.go`
- **truncateString()** (2 connections) — `tools/bg.go`
- **bytes.Buffer** (1 connections)

## Relationships

- [runSubagent](runSubagent.md) (3 shared connections)
- [client.go](client.go.md) (2 shared connections)
- [chat.go](chat.go.md) (2 shared connections)
- [GetStringArg](GetStringArg.md) (2 shared connections)
- [config_paths.go](config_paths.go.md) (1 shared connections)
- [startCLI](startCLI.md) (1 shared connections)
- [TruncateStr](TruncateStr.md) (1 shared connections)
- [handleAction](handleAction.md) (1 shared connections)
- [time.Time](time.Time.md) (1 shared connections)

## Source Files

- `models/cost_router.go`
- `testutil/endpoint.go`
- `tools/bg.go`

## Audit Trail

- EXTRACTED: 50 (100%)
- INFERRED: 0 (0%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*