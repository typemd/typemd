### Requirement: Watcher collects changed file paths during debounce

The file watcher SHALL collect all changed file paths during the debounce window and deliver them in `fileChangedMsg`. Duplicate paths within the same debounce window SHALL be deduplicated.

#### Scenario: Single file changed

- **WHEN** a single object file is written
- **THEN** `fileChangedMsg` SHALL contain that file's path

#### Scenario: Multiple files changed within debounce window

- **WHEN** three object files are written within the debounce window
- **THEN** `fileChangedMsg` SHALL contain all three file paths

#### Scenario: Same file changed multiple times within debounce window

- **WHEN** the same file is written twice within the debounce window
- **THEN** `fileChangedMsg` SHALL contain that file path only once

### Requirement: Watcher debounce interval is configurable

The watcher SHALL read the debounce interval from `tui.debounce_ms` in `.typemd/config.yaml`. If not configured or zero, the default SHALL be 200 milliseconds.

#### Scenario: Custom debounce interval

- **WHEN** `config.yaml` contains `tui:\n  debounce_ms: 500`
- **THEN** the watcher SHALL use a 500ms debounce window

#### Scenario: Default debounce interval

- **WHEN** `config.yaml` does not contain `tui.debounce_ms`
- **THEN** the watcher SHALL use a 200ms debounce window

### Requirement: Projector supports incremental sync by file paths

The Projector SHALL provide a `SyncFiles(paths []string)` method that synchronizes only the specified files to the index, rather than walking all objects.

#### Scenario: File created or updated

- **WHEN** `SyncFiles` receives a path to an existing object file
- **THEN** the Projector SHALL read that object via `ObjectRepository.Get()`, filter its properties against the type schema, and upsert it into the index

#### Scenario: File deleted

- **WHEN** `SyncFiles` receives a path to a file that no longer exists on disk
- **THEN** the Projector SHALL remove that object's entry from the index

#### Scenario: Wikilinks and tags rebuilt after incremental sync

- **WHEN** `SyncFiles` completes incremental object sync
- **THEN** the Projector SHALL perform a full wikilink sync and full tag relation sync

### Requirement: Incremental FTS update on Upsert and Remove

`ObjectIndex.Upsert()` SHALL update the corresponding FTS entry atomically. `ObjectIndex.Remove()` SHALL delete the corresponding FTS entry. When incremental FTS is used, the full `Rebuild()` SHALL NOT be called.

#### Scenario: Upsert updates FTS entry

- **WHEN** an object is upserted with updated body text containing "new keyword"
- **THEN** a subsequent `Search("new keyword")` SHALL return that object

#### Scenario: Remove deletes FTS entry

- **WHEN** an object is removed from the index
- **THEN** a subsequent `Search()` for terms unique to that object SHALL NOT return it

#### Scenario: Full Rebuild still works for full sync

- **WHEN** `Rebuild()` is called after a full sync
- **THEN** the FTS index SHALL be fully rebuilt from the objects table

### Requirement: Fallback to full sync

The TUI SHALL fall back to full `Projector.Sync()` when incremental sync is not possible.

#### Scenario: Empty file paths triggers full sync

- **WHEN** `fileChangedMsg` has an empty paths list
- **THEN** the TUI SHALL call `Projector.Sync()` (full sync)

#### Scenario: Incremental sync error triggers full sync

- **WHEN** `Projector.SyncFiles()` returns an error
- **THEN** the TUI SHALL fall back to `Projector.Sync()` (full sync)

#### Scenario: Initial startup uses full sync

- **WHEN** the TUI starts up and performs its first data load
- **THEN** the TUI SHALL use full `Projector.Sync()`

### Requirement: Watcher monitors types directory for schema changes

The watcher SHALL additionally monitor the `.typemd/types/` directory. Schema file changes SHALL produce a distinct message that triggers schema cache invalidation and full data refresh.

#### Scenario: Schema file edited externally

- **WHEN** a file in `.typemd/types/` is modified
- **THEN** the watcher SHALL emit a schema change message
- **AND** the TUI SHALL invalidate the schema cache and perform a full refresh
