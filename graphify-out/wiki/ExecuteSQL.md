# ExecuteSQL

> 4 nodes

## Key Concepts

- **ExecuteSQL()** (9 connections) — `tools/db.go`
- **db.go** (3 connections) — `tools/db.go`
- **loadDBConnections()** (3 connections) — `tools/db.go`
- **dbConnection** (2 connections) — `tools/db.go`

## Relationships

- [testing.T](testing.T.md) (2 shared connections)
- [TruncOutput](TruncOutput.md) (1 shared connections)
- [GetStringArg](GetStringArg.md) (1 shared connections)
- [GetIntArg](GetIntArg.md) (1 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (1 shared connections)
- [startCLI](startCLI.md) (1 shared connections)

## Source Files

- `tools/db.go`

## Audit Trail

- EXTRACTED: 9 (75%)
- INFERRED: 3 (25%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*