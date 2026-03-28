## 1. Core: Logging initialization

- [ ] 1.1 Write BDD scenarios for InitLogging (debug level, warn level, filtering, JSON format)
- [ ] 1.2 Implement BDD step definitions for InitLogging scenarios
- [ ] 1.3 Add `core/logging.go` with `InitLogging(level slog.Level, output io.Writer)` to make scenarios pass
- [ ] 1.4 Add unit tests for edge cases (nil writer, multiple calls)

## 2. CLI: Debug flag

- [ ] 2.1 Write BDD scenarios for `--debug` flag (enables logging, normal mode silent, subcommand inheritance)
- [ ] 2.2 Implement BDD step definitions for debug flag scenarios
- [ ] 2.3 Add `--debug` persistent flag to `cmd/root.go` with `PersistentPreRunE` hook to make scenarios pass

## 3. TUI: File logging

- [ ] 3.1 Refactor `tui/log.go`: remove dead `logf()` and `tuiLog`, call `core.InitLogging` with `.typemd/logs/{date}.log` file
- [ ] 3.2 Add unit test verifying TUI log file creation and DEBUG level output

## 4. MCP: Server logging

- [ ] 4.1 Add `core.InitLogging` call to MCP server startup in `mcp/server.go`
- [ ] 4.2 Add unit test verifying MCP logs to stderr, not stdout

## 5. Instrumentation: Sync operations

- [ ] 5.1 Add DEBUG-level logging to `Projector.Sync()` (start/complete with object count)
- [ ] 5.2 Add DEBUG-level logging to `Projector.SyncFiles()` (file count)
- [ ] 5.3 Add WARN-level logging for per-object sync errors

## 6. Instrumentation: Object and query services

- [ ] 6.1 Add DEBUG-level logging to `ObjectService.CreateObject()` and `SaveObject()`
- [ ] 6.2 Add DEBUG-level logging to `QueryService.QueryObjects()` and `SearchObjects()`

## 7. Instrumentation: AI service and file watcher

- [ ] 7.1 Add DEBUG-level logging to AI provider calls (provider name, operation type)
- [ ] 7.2 Add WARN-level logging for AI provider errors
- [ ] 7.3 Add DEBUG-level logging to TUI file watcher (file path, event type, schema reload)
