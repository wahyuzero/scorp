# IsDangerousCommand

> 7 nodes

## Key Concepts

- **IsDangerousCommand()** (13 connections) — `agent/safety.go`
- **safety.go** (3 connections) — `agent/safety.go`
- **TestDevNullRedirectionNotDangerous()** (3 connections) — `agent/safety_devsink_test.go`
- **TestDevTcpPseudoDeviceAllowed()** (3 connections) — `agent/safety_devsink_test.go`
- **devOverwriteTarget()** (2 connections) — `agent/safety.go`
- **safety_devsink_test.go** (2 connections) — `agent/safety_devsink_test.go`
- **isHarmlessDevSink()** (2 connections) — `agent/safety.go`

## Relationships

- [testing.T](testing.T.md) (3 shared connections)
- [RunAgentSessionLoop](RunAgentSessionLoop.md) (2 shared connections)
- [startCLI](startCLI.md) (1 shared connections)
- [StartDaemon](StartDaemon.md) (1 shared connections)
- [PermissionDecision](PermissionDecision.md) (1 shared connections)
- [ConfigMgr](ConfigMgr.md) (1 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (1 shared connections)

## Source Files

- `agent/safety.go`
- `agent/safety_devsink_test.go`

## Audit Trail

- EXTRACTED: 11 (58%)
- INFERRED: 8 (42%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*