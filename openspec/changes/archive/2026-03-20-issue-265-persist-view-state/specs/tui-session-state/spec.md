## MODIFIED Requirements

### Requirement: Session state includes navigation context
The saved state SHALL include the selected object ID, expanded type groups, scroll offset, focus panel, left panel width, properties panel width, properties panel visibility, and optionally the active view type name, view name, view cursor, view scroll, and view expanded groups.

#### Scenario: Full state persisted
- **WHEN** the user has selected object `book/clean-code-01jqr...`, expanded groups `book` and `person`, focus on body panel, left panel width 35, properties panel width 30, properties visible, and scroll offset 5
- **THEN** the state file SHALL contain all of these values

#### Scenario: Full state persisted in view mode
- **WHEN** the user is in view mode for type `book` with view `default`, cursor at 2, scroll at 1, with group "fiction" expanded
- **THEN** the state file SHALL contain `view_type_name`, `view_name`, `view_cursor`, `view_scroll`, and `view_expanded_groups` in addition to the sidebar state fields

#### Scenario: Search state excluded
- **WHEN** the user is in search mode with an active query and results
- **THEN** the state file SHALL NOT include search mode, search query, or search results

### Requirement: TUI restores session state on launch
The TUI SHALL read `.typemd/tui-state.yaml` on startup and restore the saved state, including view mode if view state fields are present.

#### Scenario: Restore selected object
- **WHEN** the state file contains `selectedObjectID: "book/clean-code-01jqr..."`
- **AND** that object exists in the vault
- **AND** no view state fields are present
- **THEN** the TUI SHALL start with that object selected and its type group expanded

#### Scenario: Restore view mode takes precedence
- **WHEN** the state file contains both `selected_object_id` and `view_type_name`/`view_name`
- **THEN** the TUI SHALL enter view mode (view state takes precedence over sidebar selection)

#### Scenario: Restore expanded groups
- **WHEN** the state file contains `expandedGroups: ["book", "person"]`
- **THEN** only the `book` and `person` type groups SHALL be expanded on startup

#### Scenario: Restore panel dimensions
- **WHEN** the state file contains panel width values
- **THEN** the TUI SHALL apply those widths, subject to terminal size constraints (existing clamp logic)

#### Scenario: Restore focus panel
- **WHEN** the state file contains `focus: "body"`
- **THEN** the TUI SHALL start with focus on the body panel

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
