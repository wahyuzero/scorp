# collector_system.go

> 65 nodes · cohesion 0.07

## Key Concepts

- **collector_system.go** (22 connections) — `collectors/collector_system.go`
- **FormatHourlyReport()** (14 connections) — `scheduler/reporter.go`
- **resourceAlertLoop()** (13 connections) — `main.go`
- **alerter_test.go** (11 connections) — `scheduler/alerter_test.go`
- **collector_docker.go** (10 connections) — `collectors/collector_docker.go`
- **StorageData** (10 connections) — `collectors/collector_system.go`
- **reporter.go** (10 connections) — `scheduler/reporter.go`
- **CollectDocker()** (9 connections) — `collectors/collector_docker.go`
- **DockerData** (9 connections) — `collectors/collector_docker.go`
- **CollectNetwork()** (9 connections) — `collectors/collector_system.go`
- **CollectStorage()** (9 connections) — `collectors/collector_system.go`
- **NetworkData** (8 connections) — `collectors/collector_system.go`
- **SystemData** (7 connections) — `collectors/collector_system.go`
- **CanFire()** (7 connections) — `scheduler/alerter.go`
- **contains()** (7 connections) — `scheduler/alerter_test.go`
- **resetAlertCooldowns()** (7 connections) — `scheduler/alerter_test.go`
- **setAlertThresholds()** (7 connections) — `scheduler/alerter_test.go`
- **StartDockerStatsSampler()** (6 connections) — `collectors/collector_docker.go`
- **alerter.go** (6 connections) — `scheduler/alerter.go`
- **CheckSystemAlerts()** (6 connections) — `scheduler/alerter.go`
- **TestCheckDockerAlerts()** (6 connections) — `scheduler/alerter_test.go`
- **TestCheckNetworkAlerts()** (6 connections) — `scheduler/alerter_test.go`
- **TestCheckStorageAlerts()** (6 connections) — `scheduler/alerter_test.go`
- **TestCheckSystemAlerts()** (6 connections) — `scheduler/alerter_test.go`
- **FormatStatusResponse()** (6 connections) — `scheduler/reporter.go`
- *... and 40 more nodes in this community*

## Relationships

- [handleAction](handleAction.md) (18 shared connections)
- [testing.T](testing.T.md) (9 shared connections)
- [collector_system_native.go](collector_system_native.go.md) (6 shared connections)
- [collector_security.go](collector_security.go.md) (4 shared connections)
- [collector_coolify.go](collector_coolify.go.md) (2 shared connections)
- [ConfigMgr](ConfigMgr.md) (1 shared connections)
- [config_paths.go](config_paths.go.md) (1 shared connections)
- [chat.go](chat.go.md) (1 shared connections)

## Source Files

- `collectors/collector_docker.go`
- `collectors/collector_system.go`
- `collectors/utils.go`
- `main.go`
- `scheduler/alerter.go`
- `scheduler/alerter_test.go`
- `scheduler/reporter.go`

## Audit Trail

- EXTRACTED: 187 (95%)
- INFERRED: 9 (5%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*