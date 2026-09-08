# runSelfReview

> 21 nodes

## Key Concepts

- **runSelfReview()** (9 connections) — `agent/self_improve.go`
- **memory.go** (8 connections) — `tools/memory.go`
- **TestSelfReviewIntegration()** (7 connections) — `agent/self_improve_test.go`
- **ExecuteMemory()** (7 connections) — `tools/memory.go`
- **maybeRunSelfReview()** (5 connections) — `agent/self_improve.go`
- **SetMemory()** (5 connections) — `tools/memory.go`
- **InitMemoryCache()** (4 connections) — `tools/memory.go`
- **ListMemory()** (4 connections) — `tools/memory.go`
- **persistMemory()** (4 connections) — `tools/memory.go`
- **self_improve.go** (3 connections) — `agent/self_improve.go`
- **config_helper.go** (3 connections) — `config/config_helper.go`
- **SaveJSON()** (3 connections) — `config/config_helper.go`
- **SaveJSONPerm()** (3 connections) — `config/config_helper.go`
- **deleteMemory()** (3 connections) — `tools/memory.go`
- **GetMemorySummary()** (3 connections) — `tools/memory.go`
- **getSharedMemorySummary()** (2 connections) — `agent/chat.go`
- **LoadJSON()** (2 connections) — `config/config_helper.go`
- **os.FileMode** (2 connections)
- **getMemory()** (2 connections) — `tools/memory.go`
- **memoryFact** (1 connections) — `agent/self_improve.go`
- **AgentMessage** (1 connections)

## Relationships

- [chat.go](chat.go.md) (3 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (2 shared connections)
- [context.Context](context.Context.md) (2 shared connections)
- [testing.T](testing.T.md) (2 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (2 shared connections)
- [TruncateStr](TruncateStr.md) (1 shared connections)
- [SaveModelConfig](SaveModelConfig.md) (1 shared connections)
- [time.Time](time.Time.md) (1 shared connections)
- [init](init.md) (1 shared connections)
- [GetStringArg](GetStringArg.md) (1 shared connections)
- [StartDaemon](StartDaemon.md) (1 shared connections)

## Source Files

- `agent/chat.go`
- `agent/self_improve.go`
- `agent/self_improve_test.go`
- `config/config_helper.go`
- `tools/memory.go`

## Audit Trail

- EXTRACTED: 44 (90%)
- INFERRED: 5 (10%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*