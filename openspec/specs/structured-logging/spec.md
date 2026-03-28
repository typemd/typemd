## ADDED Requirements

### Requirement: Logging initialization
The system SHALL provide a `core.InitLogging(level slog.Level, output io.Writer)` function that configures the global `slog.Default()` logger with a JSON handler at the specified level writing to the specified output.

#### Scenario: Initialize with debug level
- **WHEN** `InitLogging(slog.LevelDebug, buffer)` is called
- **THEN** `slog.Default()` SHALL be configured with a JSON handler at DEBUG level writing to the buffer

#### Scenario: Initialize with warn level
- **WHEN** `InitLogging(slog.LevelWarn, buffer)` is called
- **THEN** `slog.Default()` SHALL be configured with a JSON handler at WARN level writing to the buffer

#### Scenario: Debug messages filtered at warn level
- **WHEN** `InitLogging(slog.LevelWarn, buffer)` is called
- **AND** a DEBUG-level message is logged
- **THEN** the buffer SHALL remain empty (message filtered)

#### Scenario: JSON output format
- **WHEN** `InitLogging(slog.LevelDebug, buffer)` is called
- **AND** a message is logged with attributes
- **THEN** the output SHALL be valid JSON with `time`, `level`, `msg` fields plus any attributes

### Requirement: CLI debug flag
The root CLI command SHALL accept a `--debug` persistent flag that enables debug-level JSON logging to stderr.

#### Scenario: Debug flag enables logging
- **WHEN** a CLI command is invoked with `--debug`
- **THEN** debug-level JSON log output SHALL appear on stderr

#### Scenario: Normal mode is silent
- **WHEN** a CLI command is invoked without `--debug`
- **THEN** no log output SHALL appear on stderr

#### Scenario: Debug flag inherited by subcommands
- **WHEN** a subcommand is invoked with `--debug` on the root command
- **THEN** the subcommand SHALL have debug logging enabled

### Requirement: TUI file logging
The TUI SHALL always log at DEBUG level to `.typemd/logs/{date}.log` where `{date}` is the current date in `YYYY-MM-DD` format.

#### Scenario: TUI creates log file on startup
- **WHEN** the TUI starts
- **THEN** a log file SHALL be created at `.typemd/logs/{YYYY-MM-DD}.log`
- **AND** the global logger SHALL write to that file at DEBUG level

#### Scenario: TUI never writes logs to stderr
- **WHEN** the TUI is running
- **THEN** no log output SHALL appear on stderr (to avoid corrupting terminal rendering)

### Requirement: MCP server logging
The MCP server SHALL log to stderr in JSON format. Stdout is reserved for MCP protocol messages.

#### Scenario: MCP server logs to stderr
- **WHEN** the MCP server starts
- **THEN** the global logger SHALL be configured to write to stderr

#### Scenario: MCP server preserves stdout for protocol
- **WHEN** the MCP server is running
- **THEN** no log output SHALL be written to stdout

### Requirement: Sync operation instrumentation
The Projector sync operations SHALL emit structured log messages at key points.

#### Scenario: Full sync logging
- **WHEN** `Projector.Sync()` is called
- **THEN** a DEBUG-level log SHALL be emitted at sync start with the vault path
- **AND** a DEBUG-level log SHALL be emitted at sync completion with the count of objects processed

#### Scenario: Incremental sync logging
- **WHEN** `Projector.SyncFiles()` is called with file paths
- **THEN** a DEBUG-level log SHALL be emitted with the number of files to sync

#### Scenario: Sync error logging
- **WHEN** an error occurs during sync for a specific object
- **THEN** a WARN-level log SHALL be emitted with the object path and error details

### Requirement: Object service instrumentation
The ObjectService SHALL emit structured log messages for write operations.

#### Scenario: Object creation logging
- **WHEN** an object is created via `ObjectService.CreateObject()`
- **THEN** a DEBUG-level log SHALL be emitted with the object ID and type

#### Scenario: Object save logging
- **WHEN** an object is saved via `ObjectService.SaveObject()`
- **THEN** a DEBUG-level log SHALL be emitted with the object ID

### Requirement: Query service instrumentation
The QueryService SHALL emit structured log messages for query operations.

#### Scenario: Query execution logging
- **WHEN** `QueryService.QueryObjects()` is called
- **THEN** a DEBUG-level log SHALL be emitted with the query parameters

#### Scenario: Search execution logging
- **WHEN** `QueryService.SearchObjects()` is called
- **THEN** a DEBUG-level log SHALL be emitted with the search term and result count

### Requirement: AI service instrumentation
The AI service SHALL emit structured log messages for provider calls.

#### Scenario: AI call logging
- **WHEN** an AI provider call is made
- **THEN** a DEBUG-level log SHALL be emitted with the provider name and operation type

#### Scenario: AI error logging
- **WHEN** an AI provider call fails
- **THEN** a WARN-level log SHALL be emitted with the provider name and error details

### Requirement: File watcher instrumentation
The TUI file watcher SHALL emit structured log messages for file system events.

#### Scenario: File change detected
- **WHEN** a file change event is received by the watcher
- **THEN** a DEBUG-level log SHALL be emitted with the file path and event type

#### Scenario: Schema change detected
- **WHEN** a schema file change event is received
- **THEN** a DEBUG-level log SHALL be emitted indicating schema reload
