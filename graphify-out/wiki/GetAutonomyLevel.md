# GetAutonomyLevel

> 29 nodes

## Key Concepts

- **GetAutonomyLevel()** (32 connections) — `config/autonomy.go`
- **ExecuteShell()** (22 connections) — `tools/exec.go`
- **SetAutonomyLevel()** (21 connections) — `config/autonomy.go`
- **eval/core.go** (16 connections) — `eval/core.go`
- **autonomy.go** (9 connections) — `config/autonomy.go`
- **ConfirmationRequired()** (9 connections) — `config/autonomy.go`
- **caseAutoClassifier()** (9 connections) — `eval/core.go`
- **IsToolAllowed()** (7 connections) — `config/autonomy.go`
- **IsPathRestricted()** (6 connections) — `config/autonomy.go`
- **TestAutonomyLevels()** (6 connections) — `config/autonomy_test.go`
- **TestIsToolAllowedPlanMode()** (6 connections) — `config/planmode_test.go`
- **caseDenyShellYOLO()** (6 connections) — `eval/core.go`
- **TestShellDenyRuleBlocksAllModesAndConfirmation()** (6 connections) — `tools/deny_integration_test.go`
- **SetPlanningMode()** (5 connections) — `config/autonomy.go`
- **casePlanModeGate()** (5 connections) — `eval/core.go`
- **withEnv()** (5 connections) — `eval/core.go`
- **TestShellDangerGateRespectsAutonomy()** (5 connections) — `tools/exec_gate_test.go`
- **TestShellSensitivePathSandboxAllModes()** (5 connections) — `tools/exec_gate_test.go`
- **TestConfirmationRequired()** (4 connections) — `config/autonomy_test.go`
- **caseDangerGate()** (4 connections) — `eval/core.go`
- **caseDenyInvalidSkipped()** (4 connections) — `eval/core.go`
- **caseSensitivePath()** (4 connections) — `eval/core.go`
- **PlanningModeActive()** (3 connections) — `config/autonomy.go`
- **autonomy_test.go** (2 connections) — `config/autonomy_test.go`
- **AutonomyLevel** (2 connections) — `config/autonomy.go`
- *... and 4 more nodes in this community*

## Relationships

- [PermissionDecision](PermissionDecision.md) (10 shared connections)
- [startCLI](startCLI.md) (9 shared connections)
- [CheckDenyRules](CheckDenyRules.md) (7 shared connections)
- [SandboxActive](SandboxActive.md) (6 shared connections)
- [testing.T](testing.T.md) (6 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (4 shared connections)
- [ExecuteTool](ExecuteTool.md) (4 shared connections)
- [RegisterTool](RegisterTool.md) (4 shared connections)
- [GetIntArg](GetIntArg.md) (3 shared connections)
- [GetStringArg](GetStringArg.md) (3 shared connections)
- [HandleTelegramAction](HandleTelegramAction.md) (3 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (3 shared connections)

## Source Files

- `config/autonomy.go`
- `config/autonomy_test.go`
- `config/planmode_test.go`
- `docs/IMPLEMENTATION_PLAN_SCORP.md`
- `eval/core.go`
- `tools/deny_integration_test.go`
- `tools/exec.go`
- `tools/exec_gate_test.go`

## Audit Trail

- EXTRACTED: 122 (85%)
- INFERRED: 21 (15%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*