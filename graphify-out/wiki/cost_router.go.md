# cost_router.go

> 19 nodes

## Key Concepts

- **cost_router.go** (20 connections) — `models/cost_router.go`
- **RouteModelCostAware()** (10 connections) — `models/cost_router.go`
- **GetModelByName()** (8 connections) — `models/model_router.go`
- **handleCostCommand()** (6 connections) — `models/cost_router.go`
- **formatCostReport()** (5 connections) — `models/cost_router.go`
- **LoadCostConfig()** (5 connections) — `models/cost_router.go`
- **LoadCostTracker()** (5 connections) — `models/cost_router.go`
- **defaultCostConfig()** (4 connections) — `models/cost_router.go`
- **getCheapestModel()** (4 connections) — `models/cost_router.go`
- **FormatDailyCostSummary()** (3 connections) — `models/cost_router.go`
- **isOffPeak()** (3 connections) — `models/cost_router.go`
- **saveCostConfig()** (3 connections) — `models/cost_router.go`
- **CostConfig** (3 connections) — `models/cost_router.go`
- **CostTracker** (3 connections) — `models/cost_router.go`
- **.getTotal()** (3 connections) — `models/cost_router.go`
- **init()** (2 connections) — `models/cost_router.go`
- **isBudgetExceeded()** (2 connections) — `models/cost_router.go`
- **makeBudgetBar()** (2 connections) — `models/cost_router.go`
- **ModelCost** (2 connections) — `models/cost_router.go`

## Relationships

- [ConfigMgr](ConfigMgr.md) (4 shared connections)
- [TruncateStr](TruncateStr.md) (4 shared connections)
- [context.Context](context.Context.md) (4 shared connections)
- [startCLI](startCLI.md) (3 shared connections)
- [ChatMessage](ChatMessage.md) (2 shared connections)
- [StartDaemon](StartDaemon.md) (2 shared connections)
- [ToolCall](ToolCall.md) (2 shared connections)
- [runSubagent](runSubagent.md) (2 shared connections)
- [readInteractiveInput](readInteractiveInput.md) (1 shared connections)
- [GetStringArg](GetStringArg.md) (1 shared connections)
- [testing.T](testing.T.md) (1 shared connections)
- [SaveModelConfig](SaveModelConfig.md) (1 shared connections)

## Source Files

- `models/cost_router.go`
- `models/model_router.go`

## Audit Trail

- EXTRACTED: 54 (90%)
- INFERRED: 6 (10%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*