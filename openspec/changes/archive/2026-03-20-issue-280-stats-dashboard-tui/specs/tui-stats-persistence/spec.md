## ADDED Requirements

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

#### Scenario: Fallback when saved type no longer exists
- **WHEN** the state file contains `stats_type_name: deleted-type`
- **AND** "deleted-type" no longer exists in the vault
- **THEN** the TUI SHALL enter stats mode at the Vault Overview (ignoring the invalid type)

#### Scenario: Stats state takes precedence over view state
- **WHEN** the state file contains both `stats_cursor` and `view_type_name` fields
- **THEN** stats mode state SHALL take precedence (stats mode was the last active mode)
