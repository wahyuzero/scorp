# Benchmark

> 41 nodes

## Key Concepts

- **Benchmark** (14 connections) — `mcp/transpiler/probe.go`
- **Rebuild()** (11 connections) — `mcp/transpiler/transpiler.go`
- **VerifyContract()** (10 connections) — `mcp/transpiler/verify.go`
- **Probe()** (9 connections) — `mcp/transpiler/probe.go`
- **buildAndVerify()** (9 connections) — `mcp/transpiler/transpiler.go`
- **verify.go** (9 connections) — `mcp/transpiler/verify.go`
- **probe.go** (7 connections) — `mcp/transpiler/probe.go`
- **transpiler.go** (7 connections) — `mcp/transpiler/transpiler.go`
- **probe_test.go** (6 connections) — `mcp/transpiler/probe_test.go`
- **TestProbeAndVerifyGoldenRoundTrip()** (6 connections) — `mcp/transpiler/probe_test.go`
- **probeFromManifest()** (6 connections) — `mcp/transpiler/transpiler.go`
- **contractsFromBenchmark()** (6 connections) — `mcp/transpiler/verify.go`
- **contractsFromTools()** (6 connections) — `mcp/transpiler/verify.go`
- **Sandbox** (6 connections) — `mcp/transpiler/build.go`
- **build.go** (5 connections) — `mcp/transpiler/build.go`
- **DumpBenchmark()** (5 connections) — `mcp/transpiler/probe.go`
- **MarshalBenchmark()** (5 connections) — `mcp/transpiler/probe.go`
- **keepSandboxForInspection()** (5 connections) — `mcp/transpiler/transpiler.go`
- **saveContribution()** (5 connections) — `mcp/transpiler/transpiler.go`
- **NewSandbox()** (4 connections) — `mcp/transpiler/build.go`
- **detectRuntime()** (4 connections) — `mcp/transpiler/probe.go`
- **TestDetectRuntimeMessages()** (4 connections) — `mcp/transpiler/probe_test.go`
- **diffContract()** (4 connections) — `mcp/transpiler/verify.go`
- **Result** (4 connections) — `mcp/transpiler/transpiler.go`
- **.Build()** (4 connections) — `mcp/transpiler/build.go`
- *... and 16 more nodes in this community*

## Relationships

- [context.Context](context.Context.md) (12 shared connections)
- [testing.T](testing.T.md) (5 shared connections)
- [MCPServer](MCPServer.md) (4 shared connections)
- [Manifest](Manifest.md) (4 shared connections)
- [metasearch_engines.go](metasearch_engines.go.md) (2 shared connections)
- [LoadMCPConfig](LoadMCPConfig.md) (1 shared connections)
- [ScorpPath](ScorpPath.md) (1 shared connections)
- [ShareToMarketplace](ShareToMarketplace.md) (1 shared connections)

## Source Files

- `mcp/transpiler/build.go`
- `mcp/transpiler/probe.go`
- `mcp/transpiler/probe_test.go`
- `mcp/transpiler/transpiler.go`
- `mcp/transpiler/verify.go`

## Audit Trail

- EXTRACTED: 102 (87%)
- INFERRED: 15 (13%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*