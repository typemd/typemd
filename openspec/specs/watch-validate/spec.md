## ADDED Requirements

### Requirement: Watch flag activates continuous validation
The `tmd type validate` command SHALL accept a `--watch` / `-w` flag that enters continuous validation mode instead of running once and exiting.

#### Scenario: Watch flag starts continuous mode
- **WHEN** user runs `tmd type validate --watch`
- **THEN** the command runs a full validation and continues watching for file changes instead of exiting

#### Scenario: Short flag works
- **WHEN** user runs `tmd type validate -w`
- **THEN** the command enters continuous validation mode identical to `--watch`

#### Scenario: Without watch flag behaves as before
- **WHEN** user runs `tmd type validate` (no `--watch`)
- **THEN** the command runs validation once and exits, with no change in behavior

### Requirement: Watch monitors vault directories for changes
The watch mode SHALL monitor `types/` (recursive), `properties/` (recursive), and `objects/` (recursive) for file changes using fsnotify.

#### Scenario: Schema file change triggers re-validation
- **WHEN** watch mode is active
- **AND** a file in `types/` is modified
- **THEN** the system re-runs full validation after a debounce period

#### Scenario: Object file change triggers re-validation
- **WHEN** watch mode is active
- **AND** a file in `objects/` is modified
- **THEN** the system re-runs full validation after a debounce period

#### Scenario: Property file change triggers re-validation
- **WHEN** watch mode is active
- **AND** a file in `properties/` is modified
- **THEN** the system re-runs full validation after a debounce period

#### Scenario: Missing watched directory is skipped
- **WHEN** watch mode starts
- **AND** the `objects/` directory does not exist
- **THEN** the system skips watching that path without crashing and watches the remaining paths

### Requirement: Watch debounces rapid file changes
The watch mode SHALL debounce file change events with a 200ms window to prevent excessive re-validation during rapid changes.

#### Scenario: Rapid file changes produce single validation run
- **WHEN** watch mode is active
- **AND** multiple files change within 200ms
- **THEN** the system runs validation only once after the debounce window closes

### Requirement: Watch clears terminal and re-displays results
On each validation cycle, the watch mode SHALL clear the terminal and display the latest validation results cleanly.

#### Scenario: Terminal cleared on re-validation
- **WHEN** a file change triggers re-validation
- **THEN** the terminal is cleared before displaying the new validation results

#### Scenario: Timestamp shown with results
- **WHEN** validation results are displayed in watch mode
- **THEN** a timestamp is shown indicating when the validation ran

### Requirement: Watch re-indexes before validation
On each validation cycle, the watch mode SHALL re-sync the vault index before running validation phases, ensuring the index reflects the latest file state.

#### Scenario: Index updated before re-validation
- **WHEN** a file change triggers re-validation
- **THEN** the vault syncs its index before running the validation phases

### Requirement: Watch exits gracefully on interrupt
The watch mode SHALL exit cleanly when the user presses Ctrl+C (SIGINT) or the process receives SIGTERM.

#### Scenario: Ctrl+C stops watch mode
- **WHEN** watch mode is active
- **AND** user presses Ctrl+C
- **THEN** the command exits with code 0 and releases all file watchers

### Requirement: Watch shows initial validation on start
The watch mode SHALL run a full validation immediately on startup before entering the watch loop.

#### Scenario: Initial validation on start
- **WHEN** user runs `tmd type validate --watch`
- **THEN** the command immediately runs full validation and displays results before waiting for changes
