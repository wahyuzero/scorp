# TestGetRepoMap

> 5 nodes

## Key Concepts

- **TestGetRepoMap()** (4 connections) — `agent/repomap_test.go`
- **GetRepoMap()** (3 connections) — `agent/repomap.go`
- **repomap.go** (2 connections) — `agent/repomap.go`
- **InvalidateRepoMap()** (2 connections) — `agent/repomap.go`
- **repomap_test.go** (1 connections) — `agent/repomap_test.go`

## Relationships

- [runPlanningTurns](runPlanningTurns.md) (1 shared connections)
- [testing.T](testing.T.md) (1 shared connections)

## Source Files

- `agent/repomap.go`
- `agent/repomap_test.go`

## Audit Trail

- EXTRACTED: 4 (57%)
- INFERRED: 3 (43%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*