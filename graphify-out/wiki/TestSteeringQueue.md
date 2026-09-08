# TestSteeringQueue

> 7 nodes

## Key Concepts

- **TestSteeringQueue()** (6 connections) — `agent/steering_test.go`
- **steering.go** (4 connections) — `agent/steering.go`
- **PopSteeringMessage()** (4 connections) — `agent/steering.go`
- **HasSteeringMessage()** (3 connections) — `agent/steering.go`
- **QueueSteeringMessage()** (3 connections) — `agent/steering.go`
- **ClearSteeringQueue()** (2 connections) — `agent/steering.go`
- **steering_test.go** (1 connections) — `agent/steering_test.go`

## Relationships

- [startCLI](startCLI.md) (1 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (1 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (1 shared connections)
- [HandleTelegramAction](HandleTelegramAction.md) (1 shared connections)
- [testing.T](testing.T.md) (1 shared connections)

## Source Files

- `agent/steering.go`
- `agent/steering_test.go`

## Audit Trail

- EXTRACTED: 8 (57%)
- INFERRED: 6 (43%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*