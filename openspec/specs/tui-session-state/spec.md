## Requirements

### Requirement: TUI saves session state on exit
The TUI SHALL save session state to `.typemd/tui-state.yaml` when the user quits (q or ctrl+c).

#### Scenario: State file created on first exit
- **WHEN** the TUI exits and no `.typemd/tui-state.yaml` exists
- **THEN** the file SHALL be created with the current session state

#### Scenario: State file updated on subsequent exits
- **WHEN** the TUI exits and `.typemd/tui-state.yaml` already exists
- **THEN** the file SHALL be overwritten with the current session state

### Requirement: Session state includes navigation context
The saved state SHALL include the selected object ID, expanded type groups, scroll offset, focus panel, left panel width, properties panel width, and properties panel visibility. Note: focus panel is saved for completeness but is NOT restored on startup — the TUI always starts with focus on the sidebar for consistent UX.

#### Scenario: Full state persisted
- **WHEN** the user has selected object `book/clean-code-01jqr...`, expanded groups `book` and `person`, focus on body panel, left panel width 35, properties panel width 30, properties visible, and scroll offset 5
- **THEN** the state file SHALL contain all of these values

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

#### Scenario: Restore expanded groups
- **WHEN** the state file contains `expandedGroups: ["book", "person"]`
- **THEN** only the `book` and `person` type groups SHALL be expanded on startup

#### Scenario: Restore panel dimensions
- **WHEN** the state file contains panel width values
- **THEN** the TUI SHALL apply those widths, subject to terminal size constraints (existing clamp logic)

#### Scenario: Focus always resets to sidebar
- **WHEN** the state file contains `focus: "body"`
- **THEN** the TUI SHALL start with focus on the sidebar (focus is not restored; it always resets to the left panel for consistent UX)

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

### Requirement: TUI saves view mode state on exit
The TUI SHALL save view mode context to `.typemd/tui-state.yaml` when the user quits while in view mode.

#### Scenario: View mode state captured on quit
- **WHEN** the user quits the TUI while in view mode for type `book` with view `by-rating`
- **AND** the cursor is at position 3 with scroll offset 2
- **AND** groups "5 stars" and "4 stars" are expanded
- **THEN** the state file SHALL contain `view_type_name: book`, `view_name: by-rating`, `view_cursor: 3`, `view_scroll: 2`, and `view_expanded_groups: ["5 stars", "4 stars"]`

#### Scenario: View mode fields omitted when not in view mode
- **WHEN** the user quits the TUI while in normal sidebar mode
- **THEN** the state file SHALL NOT contain `view_type_name`, `view_name`, `view_cursor`, `view_scroll`, or `view_expanded_groups` fields

#### Scenario: Full state persisted in view mode
- **WHEN** the user is in view mode for type `book` with view `default`, cursor at 2, scroll at 1, with group "fiction" expanded
- **THEN** the state file SHALL contain `view_type_name`, `view_name`, `view_cursor`, `view_scroll`, and `view_expanded_groups` in addition to the sidebar state fields

### Requirement: TUI restores view mode on launch
The TUI SHALL enter view mode on startup when the state file contains valid view mode fields.

#### Scenario: Restore view mode with full state
- **WHEN** the state file contains `view_type_name: book`, `view_name: by-rating`, `view_cursor: 3`, `view_scroll: 2`, and `view_expanded_groups: ["5 stars", "4 stars"]`
- **AND** type `book` exists with view `by-rating`
- **THEN** the TUI SHALL start in view mode showing the `by-rating` view for `book`
- **AND** the cursor SHALL be at position 3 with scroll offset 2
- **AND** groups "5 stars" and "4 stars" SHALL be expanded

#### Scenario: Restore view mode without cursor/scroll
- **WHEN** the state file contains `view_type_name: book` and `view_name: default` but no cursor, scroll, or expanded groups fields
- **AND** type `book` exists with view `default`
- **THEN** the TUI SHALL start in view mode with cursor at 0, scroll at 0, and default group expansion

#### Scenario: Restore view mode takes precedence
- **WHEN** the state file contains both `selected_object_id` and `view_type_name`/`view_name`
- **THEN** the TUI SHALL enter view mode (view state takes precedence over sidebar selection)

### Requirement: Graceful fallback when saved view type is deleted
The TUI SHALL fall back to sidebar mode when the saved view type no longer exists.

#### Scenario: Type deleted between sessions
- **WHEN** the state file contains `view_type_name: recipe` and `view_name: default`
- **AND** type `recipe` no longer exists in the vault
- **THEN** the TUI SHALL start in normal sidebar mode (ignoring view state fields)

### Requirement: Graceful fallback when saved view is deleted
The TUI SHALL fall back gracefully when the saved view no longer exists for the type.

#### Scenario: Named view deleted, default exists
- **WHEN** the state file contains `view_type_name: book` and `view_name: by-rating`
- **AND** type `book` exists but view `by-rating` does not
- **AND** type `book` has a `default` view
- **THEN** the TUI SHALL start in view mode showing the `default` view for `book`
- **AND** cursor and scroll SHALL reset to 0

#### Scenario: Named view deleted, no views exist
- **WHEN** the state file contains `view_type_name: book` and `view_name: by-rating`
- **AND** type `book` exists but has no views at all
- **THEN** the TUI SHALL start in normal sidebar mode

### Requirement: View cursor and scroll clamped to valid range
The TUI SHALL clamp restored cursor and scroll values to the current view's item count.

#### Scenario: Cursor exceeds item count
- **WHEN** the state file contains `view_cursor: 50` but the view only has 10 items
- **THEN** the cursor SHALL be clamped to the last valid position (9)

#### Scenario: Scroll exceeds item count
- **WHEN** the state file contains `view_scroll: 50` but the view only has 10 items
- **THEN** the scroll SHALL be clamped to a valid offset

### Requirement: Unknown view expanded groups are silently ignored
The TUI SHALL silently ignore group labels in `view_expanded_groups` that do not match any current group in the view.

#### Scenario: Stale group label
- **WHEN** the state file contains `view_expanded_groups: ["5 stars", "deleted-group"]`
- **AND** the view only has groups "5 stars", "4 stars", "3 stars"
- **THEN** only "5 stars" SHALL be expanded; "deleted-group" SHALL be silently ignored

### Requirement: Stats mode state saved on exit
The TUI SHALL save stats mode state to `.typemd/tui-state.yaml` when the user quits while in stats mode.

#### Scenario: Save vault overview state
- **WHEN** the user quits while on the Vault Overview screen with cursor on type "book" at position 2
- **THEN** the state file SHALL contain `stats_cursor: 2` and `stats_scroll` fields

#### Scenario: Save type detail state
- **WHEN** the user quits while on the Type Detail screen for type "book"
- **THEN** the state file SHALL contain `stats_type_name: book`, `stats_cursor`, and `stats_scroll` fields

#### Scenario: Stats state cleared when not in stats mode
- **WHEN** the user quits while NOT in stats mode
- **THEN** the state file SHALL NOT contain any `stats_*` fields

### Requirement: Stats mode state restored on launch
The TUI SHALL restore stats mode from `.typemd/tui-state.yaml` on startup when stats state fields are present.

#### Scenario: Restore vault overview
- **WHEN** the state file contains `stats_cursor: 2` but no `stats_type_name`
- **AND** the vault has at least 3 types
- **THEN** the TUI SHALL enter stats mode with the cursor at position 2

#### Scenario: Restore type detail
- **WHEN** the state file contains `stats_type_name: book` and `stats_cursor: 1`
- **AND** the "book" type exists in the vault
- **THEN** the TUI SHALL enter stats mode showing the Type Detail for "book"

#### Scenario: Fallback when saved stats type no longer exists
- **WHEN** the state file contains `stats_type_name: deleted-type`
- **AND** "deleted-type" no longer exists in the vault
- **THEN** the TUI SHALL enter stats mode at the Vault Overview (ignoring the invalid type)

#### Scenario: Stats state takes precedence over view state
- **WHEN** the state file contains both `stats_cursor` and `view_type_name` fields
- **THEN** stats mode state SHALL take precedence (stats mode was the last active mode)
