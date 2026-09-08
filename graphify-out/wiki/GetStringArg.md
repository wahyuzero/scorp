# GetStringArg

> 17 nodes

## Key Concepts

- **GetStringArg()** (47 connections) — `internal/helpers/helpers.go`
- **init()** (9 connections) — `bootstrap/core.go`
- **exec.go** (9 connections) — `tools/exec.go`
- **ExecuteListDir()** (5 connections) — `tools/exec.go`
- **ExecuteReadFile()** (5 connections) — `tools/exec.go`
- **isPathAllowed()** (5 connections) — `tools/exec.go`
- **ExecuteProcess()** (5 connections) — `tools/process.go`
- **ExecuteSearchCode()** (5 connections) — `tools/search.go`
- **ExecuteSendFile()** (4 connections) — `tools/exec.go`
- **ExecuteSystemInfo()** (4 connections) — `tools/exec.go`
- **ExecuteWriteFile()** (4 connections) — `tools/exec.go`
- **SendDocumentBytes()** (2 connections) — `telegram/files.go`
- **bootstrap/core.go** (1 connections) — `bootstrap/core.go`
- **needsShellExecution()** (1 connections) — `tools/exec.go`
- **shellQuote()** (1 connections) — `tools/exec.go`
- **process.go** (1 connections) — `tools/process.go`
- **tools/search.go** (1 connections) — `tools/search.go`

## Relationships

- [GetIntArg](GetIntArg.md) (9 shared connections)
- [TruncOutput](TruncOutput.md) (5 shared connections)
- [patch.go](patch.go.md) (4 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (4 shared connections)
- [LoadMCPConfig](LoadMCPConfig.md) (3 shared connections)
- [registry/registry.go](registry-registry.go.md) (3 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (3 shared connections)
- [PermissionDecision](PermissionDecision.md) (2 shared connections)
- [startCLI](startCLI.md) (2 shared connections)
- [RegisterTool](RegisterTool.md) (1 shared connections)
- [testing.T](testing.T.md) (1 shared connections)
- [ClearTaskPlan](ClearTaskPlan.md) (1 shared connections)

## Source Files

- `bootstrap/core.go`
- `internal/helpers/helpers.go`
- `telegram/files.go`
- `tools/exec.go`
- `tools/process.go`
- `tools/search.go`

## Audit Trail

- EXTRACTED: 79 (98%)
- INFERRED: 2 (2%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*