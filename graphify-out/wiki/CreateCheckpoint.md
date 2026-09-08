# CreateCheckpoint

> 21 nodes

## Key Concepts

- **CreateCheckpoint()** (13 connections) — `tools/checkpoint.go`
- **checkpoint.go** (11 connections) — `tools/checkpoint.go`
- **ckptGit()** (9 connections) — `tools/checkpoint.go`
- **RestoreCheckpoint()** (9 connections) — `tools/checkpoint.go`
- **checkpointRepoRoot()** (8 connections) — `tools/checkpoint.go`
- **ListCheckpoints()** (8 connections) — `tools/checkpoint.go`
- **TestCheckpointLifecycle()** (8 connections) — `tools/checkpoint_test.go`
- **DeleteCheckpoint()** (7 connections) — `tools/checkpoint.go`
- **TestCheckpointPruneCap()** (7 connections) — `tools/checkpoint_test.go`
- **checkpoint_test.go** (6 connections) — `tools/checkpoint_test.go`
- **CheckpointDiffStat()** (5 connections) — `tools/checkpoint.go`
- **setupCkptRepo()** (5 connections) — `tools/checkpoint_test.go`
- **TestCheckpointRejectsForeignRef()** (5 connections) — `tools/checkpoint_test.go`
- **caseCheckpoint()** (4 connections) — `eval/core.go`
- **listCheckpointRefs()** (4 connections) — `tools/checkpoint.go`
- **pruneCheckpoints()** (4 connections) — `tools/checkpoint.go`
- **writeFile()** (4 connections) — `tools/checkpoint_test.go`
- **ckptRefFor()** (3 connections) — `tools/checkpoint.go`
- **TestCheckpointNonRepoNoop()** (3 connections) — `tools/checkpoint_test.go`
- **CheckpointInfo** (3 connections) — `tools/checkpoint.go`
- **Checkpoint and Rewind (P1.6)** (3 connections) — `docs/IMPLEMENTATION_PLAN_SCORP.md`

## Relationships

- [testing.T](testing.T.md) (6 shared connections)
- [startCLI](startCLI.md) (4 shared connections)
- [HandleTelegramAction](HandleTelegramAction.md) (4 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (2 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (1 shared connections)
- [time.Time](time.Time.md) (1 shared connections)
- [RunPreToolUseHooks](RunPreToolUseHooks.md) (1 shared connections)

## Source Files

- `docs/IMPLEMENTATION_PLAN_SCORP.md`
- `eval/core.go`
- `tools/checkpoint.go`
- `tools/checkpoint_test.go`

## Audit Trail

- EXTRACTED: 64 (86%)
- INFERRED: 10 (14%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*