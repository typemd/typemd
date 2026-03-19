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

### Requirement: Graceful fallback when saved type is deleted
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

### Requirement: Cursor and scroll clamped to valid range
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
