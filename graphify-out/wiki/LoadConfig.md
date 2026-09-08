# LoadConfig

> 10 nodes

## Key Concepts

- **LoadConfig()** (8 connections) — `config/config.go`
- **Config** (7 connections) — `config/config.go`
- **StartServer()** (4 connections) — `metrics/metrics.go`
- **EnvStr()** (3 connections) — `config/config.go`
- **metrics.go** (3 connections) — `metrics/metrics.go`
- **EnvBool()** (2 connections) — `config/config.go`
- **EnvFloat()** (2 connections) — `config/config.go`
- **EnvInt()** (2 connections) — `config/config.go`
- **Init()** (2 connections) — `metrics/metrics.go`
- **StopServer()** (2 connections) — `metrics/metrics.go`

## Relationships

- [StartDaemon](StartDaemon.md) (2 shared connections)
- [startCLI](startCLI.md) (1 shared connections)
- [ScorpPath](ScorpPath.md) (1 shared connections)
- [HandleTelegramAction](HandleTelegramAction.md) (1 shared connections)

## Source Files

- `config/config.go`
- `metrics/metrics.go`

## Audit Trail

- EXTRACTED: 18 (90%)
- INFERRED: 2 (10%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*