## MODIFIED Requirements

### Requirement: TUI startup initializes from restored state
The TUI `Start()` function SHALL attempt to load session state from `.typemd/tui-state.yaml` before applying default values. Restored state values take precedence over hardcoded defaults. If no state file exists or loading fails, the TUI SHALL use the current default behavior (first group expanded, first object selected).

#### Scenario: Startup with saved state
- **WHEN** the TUI starts and `.typemd/tui-state.yaml` contains valid state
- **THEN** the TUI SHALL initialize with the restored state instead of hardcoded defaults

#### Scenario: Startup without saved state
- **WHEN** the TUI starts and no `.typemd/tui-state.yaml` exists
- **THEN** the TUI SHALL initialize with current defaults (first group expanded, cursor at top, focus on left panel)
