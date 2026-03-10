## ADDED Requirements

### Requirement: TUI saves session state on exit
The TUI SHALL save session state to `.typemd/tui-state.yaml` when the user quits (q or ctrl+c).

#### Scenario: State file created on first exit
- **WHEN** the TUI exits and no `.typemd/tui-state.yaml` exists
- **THEN** the file SHALL be created with the current session state

#### Scenario: State file updated on subsequent exits
- **WHEN** the TUI exits and `.typemd/tui-state.yaml` already exists
- **THEN** the file SHALL be overwritten with the current session state

### Requirement: Session state includes navigation context
The saved state SHALL include the selected object ID, expanded type groups, scroll offset, focus panel, left panel width, properties panel width, and properties panel visibility.

#### Scenario: Full state persisted
- **WHEN** the user has selected object `book/clean-code-01jqr...`, expanded groups `book` and `person`, focus on body panel, left panel width 35, properties panel width 30, properties visible, and scroll offset 5
- **THEN** the state file SHALL contain all of these values

#### Scenario: Search state excluded
- **WHEN** the user is in search mode with an active query and results
- **THEN** the state file SHALL NOT include search mode, search query, or search results

### Requirement: TUI restores session state on launch
The TUI SHALL read `.typemd/tui-state.yaml` on startup and restore the saved state.

#### Scenario: Restore selected object
- **WHEN** the state file contains `selectedObjectID: "book/clean-code-01jqr..."`
- **AND** that object exists in the vault
- **THEN** the TUI SHALL start with that object selected and its type group expanded

#### Scenario: Restore expanded groups
- **WHEN** the state file contains `expandedGroups: ["book", "person"]`
- **THEN** only the `book` and `person` type groups SHALL be expanded on startup

#### Scenario: Restore panel dimensions
- **WHEN** the state file contains panel width values
- **THEN** the TUI SHALL apply those widths, subject to terminal size constraints (existing clamp logic)

#### Scenario: Restore focus panel
- **WHEN** the state file contains `focus: "body"`
- **THEN** the TUI SHALL start with focus on the body panel

### Requirement: Graceful fallback when selected object is deleted
The TUI SHALL fall back gracefully when the previously selected object no longer exists.

#### Scenario: Object deleted, same type has other objects
- **WHEN** the state file references `book/clean-code-01jqr...` which no longer exists
- **AND** there are other objects of type `book`
- **THEN** the TUI SHALL select the first object in the `book` type group

#### Scenario: Object deleted, entire type removed
- **WHEN** the state file references an object whose type no longer exists
- **THEN** the TUI SHALL select the first object in the first type group (default behavior)

### Requirement: Silent failure on corrupt or missing state file
The TUI SHALL fall back to default startup behavior when the state file is missing, unreadable, or contains invalid data.

#### Scenario: No state file
- **WHEN** `.typemd/tui-state.yaml` does not exist
- **THEN** the TUI SHALL start with default behavior (first group expanded, cursor at top)

#### Scenario: Corrupt state file
- **WHEN** `.typemd/tui-state.yaml` contains invalid YAML
- **THEN** the TUI SHALL start with default behavior without displaying an error

#### Scenario: Partial state file
- **WHEN** the state file is valid YAML but missing some fields (e.g., no `expandedGroups`)
- **THEN** the TUI SHALL use default values for missing fields and restored values for present fields

### Requirement: Unknown type groups in expanded list are ignored
The TUI SHALL silently ignore type names in `expandedGroups` that do not correspond to any current type in the vault.

#### Scenario: Stale type group
- **WHEN** the state file contains `expandedGroups: ["book", "deleted-type"]`
- **AND** `deleted-type` no longer exists
- **THEN** only `book` SHALL be expanded; `deleted-type` SHALL be silently ignored
