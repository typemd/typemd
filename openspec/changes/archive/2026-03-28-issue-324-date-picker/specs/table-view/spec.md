## MODIFIED Requirements

### Requirement: Date cell editing
The table view SHALL allow inline editing of date property cells using the dedicated date picker with segmented input and inline calendar modes.

#### Scenario: Activate date cell editing
- **WHEN** the cursor is on a date property cell and the user presses Enter
- **THEN** the date picker SHALL open in segmented input mode (not a generic textinput)

#### Scenario: Reject invalid date in cell
- **WHEN** the date picker produces an invalid date in a cell edit
- **THEN** the edit SHALL be rejected and a toast notification SHALL display the validation error
