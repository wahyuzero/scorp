# LoadMCPConfig

> 14 nodes

## Key Concepts

- **LoadMCPConfig()** (12 connections) — `mcp/client.go`
- **AddServerEntry()** (12 connections) — `mcp/manage.go`
- **StartMCPServers()** (10 connections) — `mcp/client.go`
- **ReloadMCPServers()** (7 connections) — `mcp/client.go`
- **mcp/manage.go** (7 connections) — `mcp/manage.go`
- **ExecuteMCPManage()** (7 connections) — `mcp/manage.go`
- **MCPConfigFilePath()** (5 connections) — `config/config_paths.go`
- **sanitizeMCPName()** (5 connections) — `mcp/client.go`
- **mcpManageAdd()** (5 connections) — `mcp/manage.go`
- **RemoveServerEntry()** (5 connections) — `mcp/manage.go`
- **mcpManageList()** (4 connections) — `mcp/manage.go`
- **mcpManageReload()** (4 connections) — `mcp/manage.go`
- **mcpManageRemove()** (4 connections) — `mcp/manage.go`
- **rebuildMCPToolList()** (2 connections) — `mcp/client.go`

## Relationships

- [client.go](client.go.md) (8 shared connections)
- [registerMCPToolsAsNative](registerMCPToolsAsNative.md) (6 shared connections)
- [ScorpPath](ScorpPath.md) (3 shared connections)
- [Manifest](Manifest.md) (3 shared connections)
- [GetStringArg](GetStringArg.md) (3 shared connections)
- [StopMCPServers](StopMCPServers.md) (2 shared connections)
- [MCPServer](MCPServer.md) (2 shared connections)
- [registry/registry.go](registry-registry.go.md) (1 shared connections)
- [StartDaemon](StartDaemon.md) (1 shared connections)
- [CheckServerContracts](CheckServerContracts.md) (1 shared connections)
- [Benchmark](Benchmark.md) (1 shared connections)
- [init](init.md) (1 shared connections)

## Source Files

- `config/config_paths.go`
- `mcp/client.go`
- `mcp/manage.go`

## Audit Trail

- EXTRACTED: 47 (77%)
- INFERRED: 14 (23%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*