# ConfigMgr

> 48 nodes

## Key Concepts

- **ConfigMgr()** (17 connections) — `config/config_manager.go`
- **agent/autonomous.go** (16 connections) — `agent/autonomous.go`
- **executeAutonomousAction()** (11 connections) — `agent/autonomous.go`
- **RunAutonomousCycle()** (11 connections) — `agent/autonomous.go`
- **phase7_test.go** (10 connections) — `agent/phase7_test.go`
- **tools/monitor.go** (10 connections) — `tools/monitor.go`
- **ExecuteAutonomous()** (8 connections) — `tools/autonomous.go`
- **gatherContext()** (7 connections) — `agent/autonomous.go`
- **makeDecision()** (7 connections) — `agent/autonomous.go`
- **ExecuteMonitor()** (7 connections) — `tools/monitor.go`
- **SaveAutonomousConfig()** (6 connections) — `agent/autonomous.go`
- **setupTestPaths()** (6 connections) — `agent/phase7_test.go`
- **TestPhase7_ConfigPersistence()** (6 connections) — `agent/phase7_test.go`
- **monitorCheckOne()** (6 connections) — `tools/monitor.go`
- **monitorLoop()** (6 connections) — `tools/monitor.go`
- **AppendAutoLog()** (5 connections) — `agent/autonomous.go`
- **TestPhase7_AuditLog()** (5 connections) — `agent/phase7_test.go`
- **TestPhase7_KillSwitch()** (5 connections) — `agent/phase7_test.go`
- **tools/autonomous.go** (5 connections) — `tools/autonomous.go`
- **AutonomousLoop()** (4 connections) — `agent/autonomous.go`
- **LoadAutonomousConfig()** (4 connections) — `agent/autonomous.go`
- **saveAutonomousConfigLocked()** (4 connections) — `agent/autonomous.go`
- **SetKillSwitch()** (4 connections) — `agent/autonomous.go`
- **AutonomousContext** (4 connections) — `agent/autonomous.go`
- **saveMonitorTargets()** (4 connections) — `tools/monitor.go`
- *... and 23 more nodes in this community*

## Relationships

- [testing.T](testing.T.md) (10 shared connections)
- [time.Time](time.Time.md) (7 shared connections)
- [TruncateStr](TruncateStr.md) (4 shared connections)
- [cost_router.go](cost_router.go.md) (4 shared connections)
- [startCLI](startCLI.md) (2 shared connections)
- [SaveModelConfig](SaveModelConfig.md) (2 shared connections)
- [GetIntArg](GetIntArg.md) (2 shared connections)
- [StartDaemon](StartDaemon.md) (1 shared connections)
- [client.go](client.go.md) (1 shared connections)
- [IsDangerousCommand](IsDangerousCommand.md) (1 shared connections)
- [collector_system_native.go](collector_system_native.go.md) (1 shared connections)
- [collector_docker.go](collector_docker.go.md) (1 shared connections)

## Source Files

- `agent/autonomous.go`
- `agent/phase7_test.go`
- `config/config_manager.go`
- `models/cost_router.go`
- `tools/autonomous.go`
- `tools/monitor.go`

## Audit Trail

- EXTRACTED: 124 (87%)
- INFERRED: 18 (13%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*