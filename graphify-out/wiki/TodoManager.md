# TodoManager

> 11 nodes

## Key Concepts

- **TodoManager** (8 connections) — `tools/todo.go`
- **todo.go** (7 connections) — `tools/todo.go`
- **.Execute()** (6 connections) — `tools/todo.go`
- **ExecuteTodo()** (3 connections) — `tools/todo.go`
- **.formatLocked()** (3 connections) — `tools/todo.go`
- **formatTodoList()** (2 connections) — `tools/todo.go`
- **GetDefaultTodoManager()** (2 connections) — `tools/todo.go`
- **getStringArgFromMap()** (2 connections) — `tools/todo.go`
- **NewTodoManager()** (2 connections) — `tools/todo.go`
- **TodoItem** (2 connections) — `tools/todo.go`
- **.ResolveKey()** (2 connections) — `tools/todo.go`

## Relationships

- [init](init.md) (1 shared connections)
- [runSubagent](runSubagent.md) (1 shared connections)
- [GetIntArg](GetIntArg.md) (1 shared connections)

## Source Files

- `tools/todo.go`

## Audit Trail

- EXTRACTED: 21 (100%)
- INFERRED: 0 (0%)
- AMBIGUOUS: 0 (0%)

---

*Part of the graphify knowledge wiki. See [index](index.md) to navigate.*