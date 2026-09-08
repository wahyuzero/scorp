# .listenSSEStream

> 6 nodes

## Key Concepts

- **.listenSSEStream()** (4 connections) — `mcp/sse.go`
- **MCPServer** (3 connections) — `mcp/sse.go`
- **io.ReadCloser** (2 connections)
- **.sendRemoteRequest()** (2 connections) — `mcp/sse.go`
- **net/url.URL** (1 connections)
- **.sendRemoteNotification()** (1 connections) — `mcp/sse.go`

## Relationships

- [runSubagent](runSubagent.md) (1 shared connections)
- [MCPServer](MCPServer.md) (1 shared connections)
- [client.go](client.go.md) (1 shared connections)

## Source Files

- `mcp/sse.go`

## Audit Trail

- EXTRACTED: 8 (100%)
- INFERRED: 0 (0%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*