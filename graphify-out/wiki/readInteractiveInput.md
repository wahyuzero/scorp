# readInteractiveInput

> 14 nodes

## Key Concepts

- **readInteractiveInput()** (8 connections) — `tui_input.go`
- **renderStatusFooter()** (8 connections) — `tui_status.go`
- **tui_input.go** (4 connections) — `tui_input.go`
- **tui_status.go** (4 connections) — `tui_status.go`
- **getContextPill()** (4 connections) — `tui_status.go`
- **SlashCommand** (3 connections) — `tui_input.go`
- **filterCommands()** (3 connections) — `tui_input.go`
- **renderPopupBox()** (3 connections) — `tui_input.go`
- **GetDailyTotalUSD()** (2 connections) — `models/cost_router.go`
- **tui_paste.go** (2 connections) — `tui_paste.go`
- **disableBracketedPaste()** (2 connections) — `tui_paste.go`
- **enableBracketedPaste()** (2 connections) — `tui_paste.go`
- **getGitStatus()** (2 connections) — `tui_status.go`
- **getShortCwd()** (2 connections) — `tui_status.go`

## Relationships

- [startCLI](startCLI.md) (2 shared connections)
- [context.Context](context.Context.md) (2 shared connections)
- [cost_router.go](cost_router.go.md) (1 shared connections)
- [chat.go](chat.go.md) (1 shared connections)
- [GetAutonomyLevel](GetAutonomyLevel.md) (1 shared connections)

## Source Files

- `models/cost_router.go`
- `tui_input.go`
- `tui_paste.go`
- `tui_status.go`

## Audit Trail

- EXTRACTED: 23 (82%)
- INFERRED: 5 (18%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*