# runSelfReview

> 19 nodes · cohesion 0.17

## Key Concepts

- **runSelfReview()** (9 connections) — `agent/self_improve.go`
- **memory.go** (8 connections) — `tools/memory.go`
- **TestSelfReviewIntegration()** (7 connections) — `agent/self_improve_test.go`
- **ExecuteMemory()** (7 connections) — `tools/memory.go`
- **SetMemory()** (5 connections) — `tools/memory.go`
- **InitMemoryCache()** (4 connections) — `tools/memory.go`
- **ListMemory()** (4 connections) — `tools/memory.go`
- **persistMemory()** (4 connections) — `tools/memory.go`
- **saveToMemory()** (3 connections) — `agent/chat.go`
- **self_improve.go** (3 connections) — `agent/self_improve.go`
- **config_helper.go** (3 connections) — `config/config_helper.go`
- **SaveJSON()** (3 connections) — `config/config_helper.go`
- **SaveJSONPerm()** (3 connections) — `config/config_helper.go`
- **deleteMemory()** (3 connections) — `tools/memory.go`
- **LoadJSON()** (2 connections) — `config/config_helper.go`
- **os.FileMode** (2 connections)
- **getMemory()** (2 connections) — `tools/memory.go`
- **memoryFact** (1 connections) — `agent/self_improve.go`
- **AgentMessage** (1 connections)

## Relationships

- [chat.go](chat.go.md) (6 shared connections)
- [TruncateStr](TruncateStr.md) (4 shared connections)
- [testing.T](testing.T.md) (2 shared connections)
- [GetStringArg](GetStringArg.md) (2 shared connections)
- [ConfigMgr](ConfigMgr.md) (1 shared connections)
- [handleAction](handleAction.md) (1 shared connections)

## Source Files

- `agent/chat.go`
- `agent/self_improve.go`
- `agent/self_improve_test.go`
- `config/config_helper.go`
- `tools/memory.go`

## Audit Trail

- EXTRACTED: 43 (96%)
- INFERRED: 2 (4%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*