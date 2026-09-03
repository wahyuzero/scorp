# ConfigMgr

> 52 nodes · cohesion 0.08

## Key Concepts

- **ConfigMgr()** (18 connections) — `config/config_manager.go`
- **agent/autonomous.go** (16 connections) — `agent/autonomous.go`
- **ConfigManager** (15 connections) — `config/config_manager.go`
- **executeAutonomousAction()** (11 connections) — `agent/autonomous.go`
- **RunAutonomousCycle()** (11 connections) — `agent/autonomous.go`
- **phase7_test.go** (10 connections) — `agent/phase7_test.go`
- **ExecuteAutonomous()** (8 connections) — `tools/autonomous.go`
- **gatherContext()** (7 connections) — `agent/autonomous.go`
- **makeDecision()** (7 connections) — `agent/autonomous.go`
- **SaveAutonomousConfig()** (6 connections) — `agent/autonomous.go`
- **setupTestPaths()** (6 connections) — `agent/phase7_test.go`
- **TestPhase7_ConfigPersistence()** (6 connections) — `agent/phase7_test.go`
- **InitConfigManager()** (6 connections) — `config/config_manager.go`
- **.Path()** (6 connections) — `config/config_manager.go`
- **AppendAutoLog()** (5 connections) — `agent/autonomous.go`
- **LoadAutonomousConfig()** (5 connections) — `agent/autonomous.go`
- **TestPhase7_AuditLog()** (5 connections) — `agent/phase7_test.go`
- **TestPhase7_KillSwitch()** (5 connections) — `agent/phase7_test.go`
- **config_manager.go** (5 connections) — `config/config_manager.go`
- **ExecuteToolByName()** (5 connections) — `registry/registry.go`
- **tools/autonomous.go** (5 connections) — `tools/autonomous.go`
- **AutonomousLoop()** (4 connections) — `agent/autonomous.go`
- **LoadAutoLog()** (4 connections) — `agent/autonomous.go`
- **saveAutonomousConfigLocked()** (4 connections) — `agent/autonomous.go`
- **SetKillSwitch()** (4 connections) — `agent/autonomous.go`
- *... and 27 more nodes in this community*

## Relationships

- [testing.T](testing.T.md) (10 shared connections)
- [startCLI](startCLI.md) (8 shared connections)
- [TruncateStr](TruncateStr.md) (6 shared connections)
- [handleAction](handleAction.md) (4 shared connections)
- [collector_system_native.go](collector_system_native.go.md) (3 shared connections)
- [chat.go](chat.go.md) (3 shared connections)
- [time.Time](time.Time.md) (3 shared connections)
- [config_paths.go](config_paths.go.md) (3 shared connections)
- [client.go](client.go.md) (2 shared connections)
- [collector_system.go](collector_system.go.md) (1 shared connections)
- [collector_security.go](collector_security.go.md) (1 shared connections)
- [wizard.go](wizard.go.md) (1 shared connections)

## Source Files

- `agent/autonomous.go`
- `agent/phase7_test.go`
- `config/config_manager.go`
- `registry/registry.go`
- `tools/autonomous.go`

## Audit Trail

- EXTRACTED: 131 (87%)
- INFERRED: 19 (13%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*