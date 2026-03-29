## Context

`tmd type validate` currently runs all validation phases (schema, object, relation, wiki-link, uniqueness) once and exits. The TUI already has a file watcher (`tui/watcher.go`) that uses `fsnotify` with debouncing. The CLI needs a similar continuous mode.

The existing watcher in `tui/watcher.go` is tightly coupled to Bubble Tea's `tea.Cmd`/`tea.Msg` pattern and cannot be directly reused. However, the fsnotify + debounce approach can be extracted or reimplemented for CLI use.

## Goals / Non-Goals

**Goals:**
- Add `--watch` / `-w` flag to `tmd type validate`
- Continuously monitor vault files and re-run validation on changes
- Provide a clean, `tsc --watch`-like terminal experience (clear + redisplay)
- Graceful exit on Ctrl+C

**Non-Goals:**
- Watch mode for other commands
- Incremental validation (always runs full validation)
- Extracting a shared watcher package from TUI (can be done later if needed)
- Fancy TUI output (spinners, color) — plain text is sufficient

## Decisions

### 1. Standalone watcher in `cmd/validate.go` rather than shared package

The TUI watcher uses Bubble Tea's `tea.Cmd` return pattern, which doesn't apply to a CLI loop. Rather than prematurely abstracting a shared watcher, implement the fsnotify + debounce loop directly in the validate command. The pattern is small (~40 lines) and the two use cases (TUI message-driven vs CLI blocking loop) have different control flow needs.

**Alternative considered:** Extract a shared `internal/watcher` package. Rejected because the TUI and CLI patterns differ enough that a shared abstraction would be forced. Revisit if a third consumer appears.

### 2. Watch paths: types dir, properties file, objects dir

Watch these paths:
- `.typemd/types/` (recursive) — schema changes
- `.typemd/properties.yaml` — shared property changes
- `objects/` (recursive) — object file changes

Skip missing directories gracefully (e.g., no `objects/` yet in a fresh vault).

### 3. Re-index before each validation run

After detecting changes, call `vault.Sync()` to re-index before running validation. This ensures the SQLite index reflects the latest file state. The vault is opened once at startup and reused across watch cycles.

### 4. Terminal clear via ANSI escape

Use `\033[H\033[2J` (cursor home + clear screen) for cross-platform terminal clearing. This avoids a dependency on `os/exec` to call `clear`/`cls`.

### 5. Signal handling via `os/signal` + context

Use `signal.NotifyContext` with `os.Interrupt` and `syscall.SIGTERM` to create a cancellable context. The watch loop checks context cancellation to exit cleanly.

## Risks / Trade-offs

- **Risk**: Rapid file changes (e.g., `git checkout` touching many files) could trigger excessive validation runs → **Mitigation**: 200ms debounce window collects rapid changes into a single run.
- **Risk**: Vault re-sync on each change adds latency → **Mitigation**: Sync is already optimized for incremental updates; full validation is the bottleneck, not sync.
- **Risk**: fsnotify doesn't detect new subdirectories created after watcher starts → **Mitigation**: For `objects/` and `types/`, new type directories are uncommon during a watch session. Document as known limitation; user can restart the watch.
