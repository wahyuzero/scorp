# HasGreenTestRun

> 8 nodes

## Key Concepts

- **HasGreenTestRun()** (8 connections) — `tools/testgate.go`
- **caseClaimGate()** (6 connections) — `eval/core.go`
- **TestHasGreenTestRun()** (6 connections) — `tools/claimgate_test.go`
- **LooksLikeTestPassClaim()** (6 connections) — `tools/testgate.go`
- **MarkTaskBoundary()** (5 connections) — `tools/testgate.go`
- **Evidence-Based Claim Gate (P4.16)** (4 connections) — `docs/IMPLEMENTATION_PLAN_SCORP.md`
- **TestLooksLikeTestPassClaim()** (3 connections) — `tools/claimgate_test.go`
- **claimgate_test.go** (2 connections) — `tools/claimgate_test.go`

## Relationships

- [RunAgentSessionLoop](RunAgentSessionLoop.md) (5 shared connections)
- [TestIntegrityStatus](TestIntegrityStatus.md) (4 shared connections)
- [RecordToolReceipt](RecordToolReceipt.md) (3 shared connections)
- [testgate.go](testgate.go.md) (3 shared connections)
- [testing.T](testing.T.md) (2 shared connections)
- [ScorpPath](ScorpPath.md) (1 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (1 shared connections)
- [CheckServerContracts](CheckServerContracts.md) (1 shared connections)

## Source Files

- `docs/IMPLEMENTATION_PLAN_SCORP.md`
- `eval/core.go`
- `tools/claimgate_test.go`
- `tools/testgate.go`

## Audit Trail

- EXTRACTED: 22 (73%)
- INFERRED: 8 (27%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*