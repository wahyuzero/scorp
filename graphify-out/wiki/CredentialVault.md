# CredentialVault

> 11 nodes

## Key Concepts

- **CredentialVault** (10 connections) — `tools/vault.go`
- **ExecuteVault()** (8 connections) — `tools/vault.go`
- **CredentialEntry** (3 connections) — `tools/vault.go`
- **.Add()** (3 connections) — `tools/vault.go`
- **.Decrypt()** (3 connections) — `tools/vault.go`
- **.Get()** (3 connections) — `tools/vault.go`
- **tools/vault.go** (3 connections) — `tools/vault.go`
- **.Encrypt()** (2 connections) — `tools/vault.go`
- **.LoadMasterKey()** (2 connections) — `tools/vault.go`
- **.Persist()** (2 connections) — `tools/vault.go`
- **.Load()** (1 connections) — `tools/vault.go`

## Relationships

- [TestPhase6_AllTools](TestPhase6_AllTools.md) (2 shared connections)
- [time.Time](time.Time.md) (1 shared connections)
- [runSubagent](runSubagent.md) (1 shared connections)
- [ScorpPath](ScorpPath.md) (1 shared connections)
- [GetStringArg](GetStringArg.md) (1 shared connections)

## Source Files

- `tools/vault.go`

## Audit Trail

- EXTRACTED: 23 (100%)
- INFERRED: 0 (0%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*