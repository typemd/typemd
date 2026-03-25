## ADDED Requirements

### Requirement: String property editing
The system SHALL allow inline editing of string properties using a textinput component.

#### Scenario: Activate string editing
- **WHEN** the cursor is on a string property and the user presses Enter
- **THEN** a textinput SHALL appear pre-filled with the current value

#### Scenario: Confirm string edit
- **WHEN** the user edits the textinput and presses Enter
- **THEN** the property value SHALL be updated and saved

#### Scenario: Cancel string edit
- **WHEN** the user presses Esc during string editing
- **THEN** the edit SHALL be cancelled and the original value restored

### Requirement: Number property editing
The system SHALL allow inline editing of number properties with numeric validation.

#### Scenario: Activate number editing
- **WHEN** the cursor is on a number property and the user presses Enter
- **THEN** a textinput SHALL appear pre-filled with the current value

#### Scenario: Confirm valid number
- **WHEN** the user enters a valid number and presses Enter
- **THEN** the property value SHALL be updated and saved

#### Scenario: Reject invalid number
- **WHEN** the user enters a non-numeric value and presses Enter
- **THEN** the edit SHALL be rejected and a toast notification SHALL display the validation error

### Requirement: Date property editing
The system SHALL allow inline editing of date properties with format validation.

#### Scenario: Activate date editing
- **WHEN** the cursor is on a date property and the user presses Enter
- **THEN** a textinput SHALL appear pre-filled with the current value in YYYY-MM-DD format

#### Scenario: Confirm valid date
- **WHEN** the user enters a valid YYYY-MM-DD date and presses Enter
- **THEN** the property value SHALL be updated and saved

#### Scenario: Reject invalid date
- **WHEN** the user enters an invalid date format and presses Enter
- **THEN** the edit SHALL be rejected and a toast notification SHALL display the validation error

### Requirement: Datetime property editing
The system SHALL allow inline editing of datetime properties with ISO 8601 validation.

#### Scenario: Activate datetime editing
- **WHEN** the cursor is on a datetime property and the user presses Enter
- **THEN** a textinput SHALL appear pre-filled with the current value in ISO 8601 format

#### Scenario: Confirm valid datetime
- **WHEN** the user enters a valid ISO 8601 datetime and presses Enter
- **THEN** the property value SHALL be updated and saved

#### Scenario: Reject invalid datetime
- **WHEN** the user enters an invalid datetime format and presses Enter
- **THEN** the edit SHALL be rejected and a toast notification SHALL display the validation error

### Requirement: URL property editing
The system SHALL allow inline editing of URL properties with scheme validation.

#### Scenario: Activate URL editing
- **WHEN** the cursor is on a url property and the user presses Enter
- **THEN** a textinput SHALL appear pre-filled with the current value

#### Scenario: Confirm valid URL
- **WHEN** the user enters a URL with http:// or https:// scheme and presses Enter
- **THEN** the property value SHALL be updated and saved

#### Scenario: Reject invalid URL
- **WHEN** the user enters a URL without a valid scheme and presses Enter
- **THEN** the edit SHALL be rejected and a toast notification SHALL display the validation error

### Requirement: Checkbox property toggle
The system SHALL allow toggling checkbox properties directly without a textinput.

#### Scenario: Toggle checkbox with Enter
- **WHEN** the cursor is on a checkbox property and the user presses Enter
- **THEN** the value SHALL toggle between ☐ (false) and ☑ (true) and save immediately

#### Scenario: Toggle checkbox with Space
- **WHEN** the cursor is on a checkbox property and the user presses Space
- **THEN** the value SHALL toggle between ☐ (false) and ☑ (true) and save immediately

#### Scenario: Checkbox display format
- **WHEN** a checkbox property is displayed
- **THEN** it SHALL show ☐ for false/nil and ☑ for true

### Requirement: Select property editing
The system SHALL allow editing select properties via an option list picker.

#### Scenario: Activate select picker
- **WHEN** the cursor is on a select property and the user presses Enter
- **THEN** an option list SHALL appear showing all available options from the schema

#### Scenario: Navigate select options
- **WHEN** the select option list is shown
- **THEN** the user SHALL navigate options with j/k or arrow keys

#### Scenario: Confirm select option
- **WHEN** the user presses Enter on a select option
- **THEN** the property value SHALL be updated to the selected option and saved

#### Scenario: Cancel select picker
- **WHEN** the user presses Esc in the select option list
- **THEN** the picker SHALL close and the original value SHALL be preserved

#### Scenario: Current value highlighted
- **WHEN** the select option list opens
- **THEN** the currently selected option SHALL be highlighted

### Requirement: Multi-select property editing
The system SHALL allow editing multi_select properties via an option multi-picker.

#### Scenario: Activate multi-select picker
- **WHEN** the cursor is on a multi_select property and the user presses Enter
- **THEN** an option list SHALL appear showing all available options with checkmarks for selected items

#### Scenario: Toggle multi-select option
- **WHEN** the user presses Space on an option in the multi-select picker
- **THEN** the option SHALL be toggled (selected ↔ unselected)

#### Scenario: Confirm multi-select
- **WHEN** the user presses Enter in the multi-select picker
- **THEN** the property value SHALL be updated to all selected options and saved

#### Scenario: Cancel multi-select picker
- **WHEN** the user presses Esc in the multi-select picker
- **THEN** the picker SHALL close and the original value SHALL be preserved

### Requirement: Description property editing
The system SHALL allow inline editing of the `description` system property using a textinput component.

#### Scenario: Edit description
- **WHEN** the cursor is on the description property and the user presses Enter
- **THEN** a textinput SHALL appear pre-filled with the current description

#### Scenario: Confirm description edit
- **WHEN** the user edits the description and presses Enter
- **THEN** the description SHALL be updated and saved

### Requirement: Property validation via core
The TUI SHALL validate property values using `core.ValidatePropertyValue()` before accepting edits.

#### Scenario: Validation error shown as toast
- **WHEN** a property edit fails validation
- **THEN** a toast notification SHALL display the validation error message

#### Scenario: Valid input accepted
- **WHEN** a property edit passes validation
- **THEN** the value SHALL be written to `obj.Properties` and saved

### Requirement: Auto-save on property edit
The system SHALL automatically save the object after each property edit is confirmed.

#### Scenario: Save after edit confirm
- **WHEN** the user confirms a property edit (Enter)
- **THEN** the object SHALL be saved to disk immediately

#### Scenario: No save on cancel
- **WHEN** the user cancels a property edit (Esc)
- **THEN** no save SHALL occur

### Requirement: Locked objects disable property editing
The TUI property editor SHALL NOT activate for locked objects. When a locked object is displayed, the properties panel SHALL show all properties in read-only mode without the cursor indicator.

#### Scenario: Property editor not available for locked objects
- **WHEN** the user views a locked object and presses Tab to focus properties
- **THEN** the properties panel SHALL remain in read-only display mode
- **AND** no property cursor shall be shown

#### Scenario: Toast notification on edit attempt of locked object
- **WHEN** the user attempts to activate property editing on a locked object
- **THEN** a toast notification SHALL display "Object is locked. Unlock to edit."
