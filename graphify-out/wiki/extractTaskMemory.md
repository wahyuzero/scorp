# extractTaskMemory

> 18 nodes

## Key Concepts

- **extractTaskMemory()** (10 connections) — `agent/memory_md.go`
- **ReadMemoryMD()** (7 connections) — `agent/memory_md.go`
- **memory_md.go** (6 connections) — `agent/memory_md.go`
- **AppendMemoryMD()** (6 connections) — `agent/memory_md.go`
- **memory_md_test.go** (6 connections) — `agent/memory_md_test.go`
- **TestAppendMemoryMD_QuotaDropsOldest()** (6 connections) — `agent/memory_md_test.go`
- **memoryMDTestFile()** (5 connections) — `agent/memory_md_test.go`
- **TestAppendMemoryMD_DedupQuotaAndHeader()** (5 connections) — `agent/memory_md_test.go`
- **TestReadMemoryMD_Bounded()** (4 connections) — `agent/memory_md_test.go`
- **caseMemoryMD()** (4 connections) — `eval/core.go`
- **parseMemoryEntries()** (3 connections) — `agent/memory_md.go`
- **readMemoryMDLocked()** (3 connections) — `agent/memory_md.go`
- **SetMemoryMDFile()** (3 connections) — `agent/memory_md.go`
- **TestParseMemoryEntries()** (3 connections) — `agent/memory_md_test.go`
- **nonEmptyLines()** (2 connections) — `agent/memory_md_test.go`
- **Durable Memory MEMORY.md (P1.7)** (2 connections) — `docs/IMPLEMENTATION_PLAN_SCORP.md`
- **Claude Code 2026 Reference Architecture** (2 connections) — `docs/RESEARCH_AI_AGENT_FEATURES_2026.md`
- **AgentMessage** (1 connections)

## Relationships

- [testing.T](testing.T.md) (5 shared connections)
- [context.Context](context.Context.md) (2 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (2 shared connections)
- [TruncateStr](TruncateStr.md) (1 shared connections)
- [runPlanningTurns](runPlanningTurns.md) (1 shared connections)
- [ScorpPath](ScorpPath.md) (1 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (1 shared connections)
- [CheckDenyRules](CheckDenyRules.md) (1 shared connections)

## Source Files

- `agent/memory_md.go`
- `agent/memory_md_test.go`
- `docs/IMPLEMENTATION_PLAN_SCORP.md`
- `docs/RESEARCH_AI_AGENT_FEATURES_2026.md`
- `eval/core.go`

## Audit Trail

- EXTRACTED: 37 (80%)
- INFERRED: 9 (20%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*