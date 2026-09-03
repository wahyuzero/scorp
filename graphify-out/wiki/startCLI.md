# startCLI

> 40 nodes · cohesion 0.10

## Key Concepts

- **startCLI()** (28 connections) — `cli.go`
- **cost_router.go** (18 connections) — `models/cost_router.go`
- **cli.go** (17 connections) — `cli.go`
- **executeOneShot()** (11 connections) — `cli.go`
- **printStatus()** (7 connections) — `cli.go`
- **wireCLICallbacks()** (7 connections) — `cli.go`
- **printCostUsage()** (6 connections) — `cli.go`
- **printCurrentModel()** (6 connections) — `cli.go`
- **handleCostCommand()** (6 connections) — `models/cost_router.go`
- **executeTurn()** (5 connections) — `cli.go`
- **setupCLILogging()** (5 connections) — `cli.go`
- **formatCostReport()** (5 connections) — `models/cost_router.go`
- **LoadCostConfig()** (5 connections) — `models/cost_router.go`
- **LoadCostTracker()** (5 connections) — `models/cost_router.go`
- **SwitchActiveModel()** (5 connections) — `models/model_router.go`
- **HasPendingConfirmation()** (4 connections) — `agent/loop.go`
- **formatFinalResponse()** (4 connections) — `cli.go`
- **formatTerminalText()** (4 connections) — `cli.go`
- **printModelList()** (4 connections) — `cli.go`
- **printToolList()** (4 connections) — `cli.go`
- **stripHTML()** (4 connections) — `cli.go`
- **defaultCostConfig()** (4 connections) — `models/cost_router.go`
- **hasDebugFlag()** (3 connections) — `cli.go`
- **isTerminal()** (3 connections) — `cli.go`
- **printBanner()** (3 connections) — `cli.go`
- *... and 15 more nodes in this community*

## Relationships

- [TruncateStr](TruncateStr.md) (15 shared connections)
- [chat.go](chat.go.md) (13 shared connections)
- [ConfigMgr](ConfigMgr.md) (8 shared connections)
- [handleAction](handleAction.md) (5 shared connections)
- [testing.T](testing.T.md) (4 shared connections)
- [config_paths.go](config_paths.go.md) (3 shared connections)
- [GetStringArg](GetStringArg.md) (2 shared connections)
- [skills.go](skills.go.md) (1 shared connections)
- [client.go](client.go.md) (1 shared connections)
- [RunAgentLoop](RunAgentLoop.md) (1 shared connections)

## Source Files

- `agent/loop.go`
- `cli.go`
- `cli_test.go`
- `models/cost_router.go`
- `models/model_router.go`

## Audit Trail

- EXTRACTED: 127 (96%)
- INFERRED: 5 (4%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*