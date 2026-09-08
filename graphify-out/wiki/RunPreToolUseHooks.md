# RunPreToolUseHooks

> 22 nodes

## Key Concepts

- **RunPreToolUseHooks()** (17 connections) — `tools/hooks.go`
- **setHookEnv()** (13 connections) — `tools/hooks_test.go`
- **RunPostToolUseHooks()** (12 connections) — `tools/hooks.go`
- **tools/hooks_test.go** (11 connections) — `tools/hooks_test.go`
- **tools/hooks.go** (6 connections) — `tools/hooks.go`
- **runHookCommand()** (5 connections) — `tools/hooks.go`
- **TestNoHooksConfiguredIsNoop()** (5 connections) — `tools/hooks_test.go`
- **TestHookMatcherFiltersTools()** (4 connections) — `tools/hooks_test.go`
- **TestHookPayloadFieldsAndRedaction()** (4 connections) — `tools/hooks_test.go`
- **TestPostToolUseHookAddsContextAndNeverBlocks()** (4 connections) — `tools/hooks_test.go`
- **TestPostToolUseHookExit2IsAdvisoryOnly()** (4 connections) — `tools/hooks_test.go`
- **TestPreToolUseHookExit0AddsContext()** (4 connections) — `tools/hooks_test.go`
- **TestPreToolUseHookExit2Blocks()** (4 connections) — `tools/hooks_test.go`
- **TestPreToolUseHookExit2DefaultReason()** (4 connections) — `tools/hooks_test.go`
- **TestPreToolUseHookNonZeroNonTwoIsNonBlocking()** (4 connections) — `tools/hooks_test.go`
- **TestPreToolUseHookTimeoutIsNonBlocking()** (4 connections) — `tools/hooks_test.go`
- **PreToolUse & PostToolUse Hooks (P3.12)** (4 connections) — `docs/IMPLEMENTATION_PLAN_SCORP.md`
- **appendHookContext()** (3 connections) — `tools/hooks.go`
- **truncateHookLog()** (3 connections) — `tools/hooks.go`
- **Community Praised Agent Patterns (2026)** (3 connections) — `docs/RESEARCH_AI_AGENT_FEATURES_2026.md`
- **hookPayload** (1 connections) — `tools/hooks.go`
- **2026 AI Coding Agent Competitor Landscape** (1 connections) — `docs/RESEARCH_AI_AGENT_FEATURES_2026.md`

## Relationships

- [testing.T](testing.T.md) (11 shared connections)
- [config/hooks.go](config-hooks.go.md) (7 shared connections)
- [startCLI](startCLI.md) (2 shared connections)
- [ExecuteTool](ExecuteTool.md) (2 shared connections)
- [RecordToolReceipt](RecordToolReceipt.md) (1 shared connections)
- [CreateCheckpoint](CreateCheckpoint.md) (1 shared connections)

## Source Files

- `docs/IMPLEMENTATION_PLAN_SCORP.md`
- `docs/RESEARCH_AI_AGENT_FEATURES_2026.md`
- `tools/hooks.go`
- `tools/hooks_test.go`

## Audit Trail

- EXTRACTED: 59 (82%)
- INFERRED: 13 (18%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*