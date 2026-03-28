## MODIFIED Requirements

### Requirement: Date property editing
The system SHALL allow inline editing of date properties using a dedicated date picker with segmented input and inline calendar modes.

#### Scenario: Activate date editing
- **WHEN** the cursor is on a date property and the user presses Enter
- **THEN** the date picker SHALL open in segmented input mode (not a generic textinput)

#### Scenario: Confirm valid date
- **WHEN** the user confirms a date in the date picker (Enter in either mode)
- **THEN** the property value SHALL be updated and saved

#### Scenario: Reject invalid date
- **WHEN** the date picker produces an invalid date (e.g., Feb 30)
- **THEN** the edit SHALL be rejected and a toast notification SHALL display the validation error
