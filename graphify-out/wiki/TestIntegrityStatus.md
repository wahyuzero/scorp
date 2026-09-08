# TestIntegrityStatus

> 17 nodes

## Key Concepts

- **TestIntegrityStatus()** (16 connections) — `tools/testgate.go`
- **testgate_test.go** (11 connections) — `tools/testgate_test.go`
- **setTestReceipts()** (11 connections) — `tools/testgate_test.go`
- **mkReceipt()** (9 connections) — `tools/testgate_test.go`
- **TestTestIntegrityStatus_TaskBoundaryExcludesOldReceipts()** (6 connections) — `tools/testgate_test.go`
- **TestTestIntegrityStatus_FailingSuiteDoesNotCount()** (5 connections) — `tools/testgate_test.go`
- **TestTestIntegrityStatus_GreenRunMustBeAfterEdit()** (5 connections) — `tools/testgate_test.go`
- **TestTestIntegrityStatus_NoTouches()** (5 connections) — `tools/testgate_test.go`
- **TestTestIntegrityStatus_ShellEditDetected()** (5 connections) — `tools/testgate_test.go`
- **TestTestIntegrityStatus_TestRunWithRedirectIsNotATouch()** (5 connections) — `tools/testgate_test.go`
- **TestTestIntegrityStatus_TouchBlocksUntilGreenRun()** (5 connections) — `tools/testgate_test.go`
- **IsTestRelatedPath()** (4 connections) — `tools/testgate.go`
- **IsTestRunCommand()** (4 connections) — `tools/testgate.go`
- **shellTouchesTestFile()** (3 connections) — `tools/testgate.go`
- **TestIsTestRelatedPath()** (3 connections) — `tools/testgate_test.go`
- **TestIsTestRunCommand()** (3 connections) — `tools/testgate_test.go`
- **Test-Integrity Gate (P0.3)** (2 connections) — `docs/IMPLEMENTATION_PLAN_SCORP.md`

## Relationships

- [testing.T](testing.T.md) (10 shared connections)
- [testgate.go](testgate.go.md) (5 shared connections)
- [HasGreenTestRun](HasGreenTestRun.md) (4 shared connections)
- [RecordToolReceipt](RecordToolReceipt.md) (3 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (2 shared connections)

## Source Files

- `docs/IMPLEMENTATION_PLAN_SCORP.md`
- `tools/testgate.go`
- `tools/testgate_test.go`

## Audit Trail

- EXTRACTED: 50 (79%)
- INFERRED: 13 (21%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*