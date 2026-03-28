## Context

typemd currently has no structured logging. The TUI has a dead `logf()` function in `tui/log.go` that writes to `.typemd/logs/{date}.log` but is only initialized — never actually called for diagnostic output. Errors are surfaced via `fmt.Println` (CLI), toast notifications (TUI), and returned errors (MCP), with no way to inspect internal operations.

Go 1.21+ includes `log/slog` in the standard library, providing structured, leveled logging with JSON output — no external dependencies needed.

## Goals / Non-Goals

**Goals:**
- Single `InitLogging` function that all entry points (CLI, TUI, MCP) call at startup
- `--debug` persistent flag on root CLI command for debug-level JSON output to stderr
- TUI always logs to `.typemd/logs/{date}.log` (preserving existing path convention)
- MCP server logs to stderr (stdout is reserved for MCP protocol)
- Instrument key operations: sync, object CRUD, queries, AI calls, file watcher events
- Normal mode is effectively silent (WARN+ level); debug mode enables DEBUG level

**Non-Goals:**
- Log rotation or cleanup (daily files are sufficient for now)
- Remote log shipping or aggregation
- Configurable log levels beyond debug on/off
- `--verbose` as a separate level
- `context.Context` propagation (slog attributes are sufficient for now)

## Decisions

### D1: Use `log/slog` with `slog.JSONHandler`

**Choice:** Go standard library `log/slog` with JSON output format.

**Rationale:** Zero external dependencies, standard Go ecosystem, structured key-value attributes, configurable handlers and levels. JSON format is machine-parseable for future tooling.

**Alternatives considered:**
- `zerolog` / `zap` — faster but add external dependencies for marginal benefit in a CLI tool
- Plain `log` package — no structure, no levels, no JSON

### D2: Single `InitLogging` function in `core/logging.go`

**Choice:** `core.InitLogging(level slog.Level, output io.Writer)` sets `slog.SetDefault()` once at startup.

**Rationale:** Every entry point (CLI, TUI, MCP) needs logging configured before any work starts. A single function in `core/` keeps it DRY and ensures consistent configuration. Using `slog.Level` directly (not a custom enum) keeps the API standard.

**The function:**
- Creates a `slog.JSONHandler` with the given level and output
- Calls `slog.SetDefault()` to configure the global logger
- Must be called once at startup, before goroutines log

### D3: `--debug` as a persistent flag on root command

**Choice:** Add `--debug` to `cmd/root.go` as a `PersistentPreRunE` hook.

**Rationale:** Persistent flags are inherited by all subcommands. `PersistentPreRunE` runs before any command body, ensuring logging is initialized before any operation. When `--debug` is set, `InitLogging(slog.LevelDebug, os.Stderr)` is called; otherwise `InitLogging(slog.LevelWarn, io.Discard)` keeps output silent.

### D4: TUI logging to file only, never stderr

**Choice:** TUI calls `InitLogging(slog.LevelDebug, logFile)` where `logFile` is `.typemd/logs/{date}.log`.

**Rationale:** TUI owns the terminal — any stderr output would corrupt the Bubble Tea rendering. Always logging at DEBUG level to file is cheap and provides diagnostic data when users report issues. The existing `.typemd/logs/` path convention is preserved from the current dead code.

### D5: Lightweight instrumentation with `slog.Default()`

**Choice:** Use `slog.Default().Info(...)` / `slog.Default().Debug(...)` calls at key points. No logger injection or dependency inversion.

**Rationale:** `slog.Default()` is the idiomatic Go approach for application-level logging. Logger injection (passing `*slog.Logger` through constructors) would require touching every service constructor and every caller — high churn for minimal benefit. The global logger is configured once at startup, which is sufficient for a CLI/TUI tool.

**Instrumentation points:**
- `Projector.Sync` / `Projector.SyncFiles` — sync start/complete, objects processed, errors
- `ObjectService` — create/save/delete operations
- `QueryService` — query execution, search terms
- `SQLiteObjectIndex` — index operations, query timing
- `ai/Service` — provider calls, response times, errors
- TUI file watcher — file change events, debounce, schema reload

## Risks / Trade-offs

- **[Global state]** `slog.SetDefault()` is process-global → Mitigated by calling `InitLogging` once at startup before goroutines. Tests that need logging isolation can create their own `slog.Logger` instances.
- **[Log volume]** DEBUG level in TUI always writes to file → Mitigated by daily rotation (one file per day). Files are small for a local-first tool. No auto-cleanup for now.
- **[Performance]** Structured logging adds overhead to hot paths → Mitigated by using `slog.LevelDebug` which is skipped when level is WARN+. The JSON handler short-circuits when the log level is below threshold.
