# RecordToolReceipt

> 12 nodes

## Key Concepts

- **RecordToolReceipt()** (12 connections) — `tools/receipts.go`
- **GetRecentReceipts()** (10 connections) — `tools/receipts.go`
- **ToolReceipt** (8 connections) — `tools/receipts.go`
- **RedactSecrets()** (6 connections) — `tools/redact.go`
- **receipts.go** (5 connections) — `tools/receipts.go`
- **loadReceiptsLocked()** (5 connections) — `tools/receipts.go`
- **TestRecordToolReceipt()** (4 connections) — `tools/receipts_test.go`
- **saveReceiptsLocked()** (3 connections) — `tools/receipts.go`
- **TestRedactSecrets()** (3 connections) — `tools/redact_test.go`
- **receipts_test.go** (1 connections) — `tools/receipts_test.go`
- **redact.go** (1 connections) — `tools/redact.go`
- **redact_test.go** (1 connections) — `tools/redact_test.go`

## Relationships

- [testgate.go](testgate.go.md) (4 shared connections)
- [HasGreenTestRun](HasGreenTestRun.md) (3 shared connections)
- [TestIntegrityStatus](TestIntegrityStatus.md) (3 shared connections)
- [testing.T](testing.T.md) (3 shared connections)
- [startCLI](startCLI.md) (2 shared connections)
- [ScorpPath](ScorpPath.md) (2 shared connections)
- [ExecuteTool](ExecuteTool.md) (2 shared connections)
- [time.Time](time.Time.md) (2 shared connections)
- [RegisterTool](RegisterTool.md) (1 shared connections)
- [TestGatewayEndpoints](TestGatewayEndpoints.md) (1 shared connections)
- [registerMCPToolsAsNative](registerMCPToolsAsNative.md) (1 shared connections)
- [RunPreToolUseHooks](RunPreToolUseHooks.md) (1 shared connections)

## Source Files

- `tools/receipts.go`
- `tools/receipts_test.go`
- `tools/redact.go`
- `tools/redact_test.go`

## Audit Trail

- EXTRACTED: 31 (74%)
- INFERRED: 11 (26%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*