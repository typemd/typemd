## MODIFIED Requirements

### Requirement: TUI saves session state on exit
The TUI SHALL save session state to `.typemd/tui-state.yaml` when the user quits (q or ctrl+c).

#### Scenario: State file created on first exit
- **WHEN** the TUI exits and no `.typemd/tui-state.yaml` exists
- **THEN** the file SHALL be created with the current session state

#### Scenario: State file updated on subsequent exits
- **WHEN** the TUI exits and `.typemd/tui-state.yaml` already exists
- **THEN** the file SHALL be overwritten with the current session state
