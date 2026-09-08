# ExecuteTermuxAPI

> 8 nodes

## Key Concepts

- **ExecuteTermuxAPI()** (8 connections) — `tools/termux.go`
- **IsTermux()** (6 connections) — `tools/termux.go`
- **termux.go** (5 connections) — `tools/termux.go`
- **AcquireTermuxWakeLock()** (4 connections) — `tools/termux.go`
- **ReleaseTermuxWakeLock()** (4 connections) — `tools/termux.go`
- **SendTermuxNotification()** (3 connections) — `tools/termux.go`
- **TestExecuteTermuxAPI_Simulation()** (3 connections) — `tools/termux_test.go`
- **termux_test.go** (1 connections) — `tools/termux_test.go`

## Relationships

- [RunAgentSessionLoop](RunAgentSessionLoop.md) (4 shared connections)
- [init](init.md) (1 shared connections)
- [GetIntArg](GetIntArg.md) (1 shared connections)
- [GetStringArg](GetStringArg.md) (1 shared connections)
- [testing.T](testing.T.md) (1 shared connections)

## Source Files

- `tools/termux.go`
- `tools/termux_test.go`

## Audit Trail

- EXTRACTED: 20 (95%)
- INFERRED: 1 (5%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*