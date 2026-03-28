## Why

typemd has no debug/diagnostic logging. The existing `tui/log.go` defines a `logf()` function that is never called (dead code). Errors surface via `fmt.Println` (CLI) and toast notifications (TUI), but there is no way to inspect internal operations like sync flow, query execution, or AI provider calls. This makes diagnosing issues in production vaults difficult and slows down development.

## What Changes

- Add `core/logging.go` with `InitLogging(debug bool, output io.Writer)` using Go's standard `log/slog` package
- Add `--debug` persistent flag on the root CLI command that enables debug-level JSON logging to stderr
- Refactor `tui/log.go`: remove dead `logf()`, use `InitLogging` with `.typemd/logs/{date}.log` file output
- MCP server: call `InitLogging(debug, os.Stderr)` on startup
- Instrument key operations: Projector sync, ObjectService writes, QueryService queries, SQLite index ops, AI provider calls, TUI file watcher events
- Normal mode logs at ERROR+ level (effectively silent); debug mode logs at DEBUG level

## Capabilities

### New Capabilities
- `structured-logging`: Centralized logging initialization, debug mode toggle, JSON output format, and instrumentation of key operations across all entry points (CLI, TUI, MCP)

### Modified Capabilities
<!-- No existing spec-level behavior changes — logging is additive -->

## Impact

- **core/**: New `logging.go` module; instrumentation added to `projector.go`, `object_service.go`, `query_service.go`, `sqlite_object_index.go`, `ai/service.go`
- **cmd/**: New `--debug` persistent flag on `root.go`
- **tui/**: Refactored `log.go` (remove dead code, use slog); instrumentation in file watcher
- **mcp/**: Logging initialization on server startup
- **Dependencies**: None — uses Go standard library `log/slog`
