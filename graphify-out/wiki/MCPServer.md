# MCPServer

> 22 nodes

## Key Concepts

- **MCPServer** (22 connections) — `mcp/client.go`
- **MCPTool** (11 connections) — `mcp/client.go`
- **startMCPServer()** (11 connections) — `mcp/client.go`
- **MCPServerConfig** (8 connections) — `mcp/client.go`
- **ProbeServer()** (7 connections) — `mcp/client.go`
- **.sendRequest()** (6 connections) — `mcp/client.go`
- **startSSEServer()** (6 connections) — `mcp/sse.go`
- **.initialize()** (4 connections) — `mcp/client.go`
- **.listTools()** (4 connections) — `mcp/client.go`
- **context.CancelFunc** (3 connections)
- **FindMCPTool()** (3 connections) — `mcp/client.go`
- **.CallTool()** (3 connections) — `mcp/client.go`
- **isRemoteMCP()** (3 connections) — `mcp/sse.go`
- **TestRemoteMCPServer_HTTP()** (3 connections) — `mcp/sse_test.go`
- **.sendNotification()** (2 connections) — `mcp/client.go`
- **sse.go** (2 connections) — `mcp/sse.go`
- **bufio.Scanner** (1 connections)
- **encoding/json.Encoder** (1 connections)
- **sync.Once** (1 connections)
- **.wasClosed()** (1 connections) — `mcp/client.go`
- **MCPServer** (1 connections)
- **sse_test.go** (1 connections) — `mcp/sse_test.go`

## Relationships

- [client.go](client.go.md) (11 shared connections)
- [Benchmark](Benchmark.md) (4 shared connections)
- [registerMCPToolsAsNative](registerMCPToolsAsNative.md) (3 shared connections)
- [runSubagent](runSubagent.md) (2 shared connections)
- [StopMCPServers](StopMCPServers.md) (2 shared connections)
- [CheckServerContracts](CheckServerContracts.md) (2 shared connections)
- [LoadMCPConfig](LoadMCPConfig.md) (2 shared connections)
- [TruncOutput](TruncOutput.md) (1 shared connections)
- [TestGatewayEndpoints](TestGatewayEndpoints.md) (1 shared connections)
- [metasearch_engines.go](metasearch_engines.go.md) (1 shared connections)
- [TruncateStr](TruncateStr.md) (1 shared connections)
- [.listenSSEStream](listenSSEStream.md) (1 shared connections)

## Source Files

- `mcp/client.go`
- `mcp/sse.go`
- `mcp/sse_test.go`

## Audit Trail

- EXTRACTED: 63 (93%)
- INFERRED: 5 (7%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*