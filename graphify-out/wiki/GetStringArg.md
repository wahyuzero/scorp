# GetStringArg

> 80 nodes · cohesion 0.06

## Key Concepts

- **GetStringArg()** (40 connections) — `internal/helpers/helpers.go`
- **init()** (30 connections) — `bootstrap/extended.go`
- **GetIntArg()** (18 connections) — `internal/helpers/helpers.go`
- **TruncOutput()** (18 connections) — `internal/helpers/helpers.go`
- **browser_session.go** (14 connections) — `browser/browser_session.go`
- **GetOrCreateBrowserSession()** (14 connections) — `browser/browser_session.go`
- **ExecuteBrowser()** (13 connections) — `browser/browser.go`
- **GetBoolArg()** (12 connections) — `internal/helpers/helpers.go`
- **browserSessionNavigate()** (10 connections) — `browser/browser_session.go`
- **GetAllTools()** (9 connections) — `registry/registry.go`
- **exec.go** (9 connections) — `tools/exec.go`
- **ExecuteShell()** (9 connections) — `tools/exec.go`
- **init()** (7 connections) — `bootstrap/core.go`
- **ExecuteCompose()** (7 connections) — `tools/compose.go`
- **patch.go** (7 connections) — `tools/patch.go`
- **browserSessionScreenshot()** (6 connections) — `browser/browser_session.go`
- **ExecuteSQL()** (6 connections) — `tools/db.go`
- **ExecuteToolList()** (6 connections) — `tools/deferred.go`
- **ExecuteGit()** (6 connections) — `tools/git.go`
- **ExecuteHTTP()** (6 connections) — `tools/http.go`
- **ExecuteLog()** (6 connections) — `tools/log.go`
- **patchReplace()** (6 connections) — `tools/patch.go`
- **ExecuteTodo()** (6 connections) — `tools/todo.go`
- **browserSessionClick()** (5 connections) — `browser/browser_session.go`
- **browserSessionEvaluate()** (5 connections) — `browser/browser_session.go`
- *... and 55 more nodes in this community*

## Relationships

- [config_paths.go](config_paths.go.md) (13 shared connections)
- [client.go](client.go.md) (10 shared connections)
- [chat.go](chat.go.md) (10 shared connections)
- [rag_vector.go](rag_vector.go.md) (9 shared connections)
- [time.Time](time.Time.md) (9 shared connections)
- [runSubagent](runSubagent.md) (6 shared connections)
- [TruncateStr](TruncateStr.md) (6 shared connections)
- [testing.T](testing.T.md) (4 shared connections)
- [collector_system_native.go](collector_system_native.go.md) (3 shared connections)
- [runSelfReview](runSelfReview.md) (2 shared connections)
- [RunAgentLoop](RunAgentLoop.md) (2 shared connections)
- [handleAction](handleAction.md) (2 shared connections)

## Source Files

- `bootstrap/core.go`
- `bootstrap/extended.go`
- `browser/browser.go`
- `browser/browser_session.go`
- `internal/helpers/helpers.go`
- `registry/registry.go`
- `telegram/files.go`
- `tools/compose.go`
- `tools/db.go`
- `tools/deferred.go`
- `tools/exec.go`
- `tools/git.go`
- `tools/http.go`
- `tools/log.go`
- `tools/patch.go`
- `tools/process.go`
- `tools/search.go`
- `tools/todo.go`
- `tools/vision.go`
- `tools/web.go`

## Audit Trail

- EXTRACTED: 242 (93%)
- INFERRED: 18 (7%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*