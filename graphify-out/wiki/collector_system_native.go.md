# collector_system_native.go

> 50 nodes · cohesion 0.07

## Key Concepts

- **collector_system_native.go** (18 connections) — `collectors/collector_system_native.go`
- **CollectSystem()** (16 connections) — `collectors/collector_system_native.go`
- **ReadProcessList()** (10 connections) — `collectors/collector_system_native.go`
- **uptime.go** (10 connections) — `tools/uptime.go`
- **time.Duration** (8 connections)
- **ExecuteUptime()** (8 connections) — `tools/uptime.go`
- **GetTopProcesses()** (6 connections) — `collectors/collector_system_native.go`
- **transports.go** (6 connections) — `models/transports.go`
- **getClient()** (6 connections) — `models/transports.go`
- **tools/callbacks.go** (6 connections) — `tools/callbacks.go`
- **sync.RWMutex** (5 connections)
- **UptimeMonitor** (5 connections) — `tools/uptime.go`
- **UptimeResult** (5 connections) — `tools/uptime.go`
- **TopProcess** (4 connections) — `collectors/collector_system.go`
- **getMemInfo()** (4 connections) — `collectors/collector_system_native.go`
- **StartCPUSampler()** (4 connections) — `collectors/collector_system_native.go`
- **nativeTopProcess** (4 connections) — `collectors/collector_system_native.go`
- **getTransport()** (4 connections) — `models/transports.go`
- **checkTarget()** (4 connections) — `tools/uptime.go`
- **runUptimeCheck()** (4 connections) — `tools/uptime.go`
- **UptimeTarget** (4 connections) — `tools/uptime.go`
- **TestCollectorNative_GetTopProcesses()** (3 connections) — `collectors/collector_native_test.go`
- **TestCollectorNative_GetTopProcesses_Limit()** (3 connections) — `collectors/collector_native_test.go`
- **FormatDuration()** (3 connections) — `collectors/collector_system.go`
- **getProcesses()** (3 connections) — `collectors/collector_system_native.go`
- *... and 25 more nodes in this community*

## Relationships

- [testing.T](testing.T.md) (9 shared connections)
- [collector_system.go](collector_system.go.md) (6 shared connections)
- [handleAction](handleAction.md) (5 shared connections)
- [ConfigMgr](ConfigMgr.md) (3 shared connections)
- [TruncateStr](TruncateStr.md) (3 shared connections)
- [time.Time](time.Time.md) (3 shared connections)
- [GetStringArg](GetStringArg.md) (3 shared connections)
- [rag_vector.go](rag_vector.go.md) (2 shared connections)
- [runSubagent](runSubagent.md) (1 shared connections)
- [chat.go](chat.go.md) (1 shared connections)

## Source Files

- `collectors/collector_native_test.go`
- `collectors/collector_system.go`
- `collectors/collector_system_native.go`
- `models/transports.go`
- `tools/callbacks.go`
- `tools/uptime.go`

## Audit Trail

- EXTRACTED: 120 (96%)
- INFERRED: 5 (4%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*