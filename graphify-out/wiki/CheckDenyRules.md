# CheckDenyRules

> 15 nodes

## Key Concepts

- **CheckDenyRules()** (10 connections) — `config/deny.go`
- **ReloadDenyRules()** (8 connections) — `config/deny.go`
- **resetDenyRules()** (6 connections) — `config/deny_test.go`
- **TestCheckDenyRulesHoldInYOLO()** (6 connections) — `config/deny_test.go`
- **deny.go** (5 connections) — `config/deny.go`
- **config/deny_test.go** (5 connections) — `config/deny_test.go`
- **loadDenyRules()** (4 connections) — `config/deny.go`
- **ParseDenyRule()** (4 connections) — `config/deny.go`
- **TestCheckDenyRulesInvalidSpecsSkipped()** (4 connections) — `config/deny_test.go`
- **TestCheckDenyRulesMatches()** (4 connections) — `config/deny_test.go`
- **TestParseDenyRule()** (3 connections) — `config/deny_test.go`
- **DenyRule** (3 connections) — `config/deny.go`
- **Auto-Mode Classifier (P3.13)** (3 connections) — `docs/IMPLEMENTATION_PLAN_SCORP.md`
- **Deny-Rule Engine (P0.2)** (2 connections) — `docs/IMPLEMENTATION_PLAN_SCORP.md`
- **regexp.Regexp** (1 connections)

## Relationships

- [GetAutonomyLevel](GetAutonomyLevel.md) (7 shared connections)
- [testing.T](testing.T.md) (5 shared connections)
- [ExecuteTool](ExecuteTool.md) (2 shared connections)
- [time.Time](time.Time.md) (2 shared connections)
- [extractTaskMemory](extractTaskMemory.md) (1 shared connections)
- [PermissionDecision](PermissionDecision.md) (1 shared connections)

## Source Files

- `config/deny.go`
- `config/deny_test.go`
- `docs/IMPLEMENTATION_PLAN_SCORP.md`

## Audit Trail

- EXTRACTED: 35 (81%)
- INFERRED: 8 (19%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*