# PermissionDecision

> 22 nodes

## Key Concepts

- **PermissionDecision()** (18 connections) — `agent/auto.go`
- **setAutoMode()** (12 connections) — `agent/auto_test.go`
- **auto.go** (11 connections) — `agent/auto.go`
- **TestExecuteToolAutoAllowlistReachesExec()** (9 connections) — `agent/auto_allowlist_test.go`
- **auto_test.go** (9 connections) — `agent/auto_test.go`
- **TestExecuteToolAutoDenyNotBypassedByConfirmedArgs()** (8 connections) — `agent/auto_allowlist_test.go`
- **ResetAutoAllowlist()** (6 connections) — `agent/auto.go`
- **ResetAutoStats()** (6 connections) — `agent/auto.go`
- **autoClassify()** (5 connections) — `agent/auto.go`
- **TestAutoFallbackAfterRepeatedUncertain()** (5 connections) — `agent/auto_test.go`
- **TestPermissionDecisionDestructiveDenyAndAllowlist()** (5 connections) — `agent/auto_test.go`
- **TestAutoAllowlistPrefixAnchoring()** (4 connections) — `agent/auto_allowlist_test.go`
- **IsReadOnlyShellCommand()** (4 connections) — `agent/auto.go`
- **TestExecuteToolAutoDeniesOnNoChannelPath()** (4 connections) — `agent/auto_test.go`
- **TestPermissionDecisionModelPaths()** (4 connections) — `agent/auto_test.go`
- **TestPermissionDecisionReadOnlyFastPath()** (4 connections) — `agent/auto_test.go`
- **auto_allowlist_test.go** (3 connections) — `agent/auto_allowlist_test.go`
- **autoAllowlisted()** (3 connections) — `agent/auto.go`
- **TestIsReadOnlyShellCommand()** (3 connections) — `agent/auto_test.go`
- **IsReadOnlyTool()** (3 connections) — `config/autonomy.go`
- **AutoStatsSnapshot()** (2 connections) — `agent/auto.go`
- **bumpAutoStat()** (2 connections) — `agent/auto.go`

## Relationships

- [testing.T](testing.T.md) (10 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (10 shared connections)
- [RegisterTool](RegisterTool.md) (6 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (4 shared connections)
- [ExecuteTool](ExecuteTool.md) (4 shared connections)
- [GetStringArg](GetStringArg.md) (2 shared connections)
- [context.Context](context.Context.md) (1 shared connections)
- [runSubagent](runSubagent.md) (1 shared connections)
- [TruncateStr](TruncateStr.md) (1 shared connections)
- [IsDangerousCommand](IsDangerousCommand.md) (1 shared connections)
- [CheckDenyRules](CheckDenyRules.md) (1 shared connections)
- [startCLI](startCLI.md) (1 shared connections)

## Source Files

- `agent/auto.go`
- `agent/auto_allowlist_test.go`
- `agent/auto_test.go`
- `config/autonomy.go`

## Audit Trail

- EXTRACTED: 64 (74%)
- INFERRED: 22 (26%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*