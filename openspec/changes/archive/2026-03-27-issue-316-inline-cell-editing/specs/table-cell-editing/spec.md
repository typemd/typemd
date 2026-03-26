## ADDED Requirements

### Requirement: String cell editing
The table view SHALL allow inline editing of string property cells using a textinput component.

#### Scenario: Activate string cell editing
- **WHEN** the cursor is on a string property cell and the user presses Enter
- **THEN** a textinput SHALL appear in the cell pre-filled with the current value

#### Scenario: Confirm string cell edit
- **WHEN** the user edits the textinput and presses Enter
- **THEN** the property value SHALL be updated and saved immediately

#### Scenario: Cancel string cell edit
- **WHEN** the user presses Esc during string cell editing
- **THEN** the edit SHALL be cancelled and the original value restored

### Requirement: Number cell editing
The table view SHALL allow inline editing of number property cells with numeric validation.

#### Scenario: Reject invalid number in cell
- **WHEN** the user enters a non-numeric value in a number cell and presses Enter
- **THEN** the edit SHALL be rejected and a toast notification SHALL display the validation error

### Requirement: Date cell editing
The table view SHALL allow inline editing of date property cells with YYYY-MM-DD format validation.

#### Scenario: Reject invalid date in cell
- **WHEN** the user enters an invalid date format in a date cell and presses Enter
- **THEN** the edit SHALL be rejected and a toast notification SHALL display the validation error

### Requirement: Datetime cell editing
The table view SHALL allow inline editing of datetime property cells with ISO 8601 validation.

#### Scenario: Reject invalid datetime in cell
- **WHEN** the user enters an invalid datetime format in a datetime cell and presses Enter
- **THEN** the edit SHALL be rejected and a toast notification SHALL display the validation error

### Requirement: URL cell editing
The table view SHALL allow inline editing of URL property cells with scheme validation.

#### Scenario: Reject invalid URL in cell
- **WHEN** the user enters a URL without a valid scheme in a URL cell and presses Enter
- **THEN** the edit SHALL be rejected and a toast notification SHALL display the validation error

### Requirement: Checkbox cell toggle
The table view SHALL allow toggling checkbox property cells directly without entering edit mode.

#### Scenario: Toggle checkbox with Enter
- **WHEN** the cursor is on a checkbox property cell and the user presses Enter
- **THEN** the value SHALL toggle between false and true and save immediately

#### Scenario: Toggle checkbox with Space
- **WHEN** the cursor is on a checkbox property cell and the user presses Space
- **THEN** the value SHALL toggle between false and true and save immediately

### Requirement: Select cell editing
The table view SHALL allow editing select property cells via an option list picker.

#### Scenario: Activate select picker in cell
- **WHEN** the cursor is on a select property cell and the user presses Enter
- **THEN** an option list SHALL appear showing all available options from the schema

#### Scenario: Confirm select option
- **WHEN** the user presses Enter on a select option in the picker
- **THEN** the property value SHALL be updated to the selected option and saved

#### Scenario: Cancel select picker
- **WHEN** the user presses Esc in the select option picker
- **THEN** the picker SHALL close and the original value SHALL be preserved

### Requirement: Multi-select cell editing
The table view SHALL allow editing multi_select property cells via an option multi-picker.

#### Scenario: Activate multi-select picker in cell
- **WHEN** the cursor is on a multi_select property cell and the user presses Enter
- **THEN** an option list SHALL appear with checkmarks for currently selected items

#### Scenario: Toggle multi-select option
- **WHEN** the user presses Space on an option in the multi-select picker
- **THEN** the option SHALL be toggled (selected/unselected)

#### Scenario: Confirm multi-select
- **WHEN** the user presses Enter in the multi-select picker
- **THEN** the property value SHALL be updated to all selected options and saved

### Requirement: NAME column editing
The table view SHALL allow inline editing of object names in the NAME column.

#### Scenario: Edit object name in table
- **WHEN** the cursor is on the NAME column and the user presses Enter
- **THEN** a textinput SHALL appear pre-filled with the current object name

#### Scenario: Confirm name edit
- **WHEN** the user edits the name and presses Enter
- **THEN** the object name SHALL be updated and saved

### Requirement: Read-only cell detection
Cells for read-only properties SHALL be visually navigable but not editable.

#### Scenario: Enter on read-only cell
- **WHEN** the cursor is on a read-only cell (created_at, updated_at, relation) and the user presses Enter
- **THEN** no edit SHALL activate

#### Scenario: Tab skips read-only cells
- **WHEN** the user presses Tab
- **THEN** the cursor SHALL skip read-only cells and land on the next editable cell

### Requirement: Cell validation via core
The table view SHALL validate cell edits using `core.ValidatePropertyValue()` before accepting.

#### Scenario: Validation error shown as toast
- **WHEN** a cell edit fails validation
- **THEN** a toast notification SHALL display the validation error message and the edit SHALL remain active

### Requirement: Auto-save on cell edit
The table view SHALL automatically save the object after each cell edit is confirmed.

#### Scenario: Save after cell edit confirm
- **WHEN** the user confirms a cell edit (Enter)
- **THEN** the object SHALL be saved to disk immediately

#### Scenario: No save on cancel
- **WHEN** the user cancels a cell edit (Esc)
- **THEN** no save SHALL occur

### Requirement: Cancel edit on external file change
The table view SHALL cancel active cell editing when the underlying object file changes externally.

#### Scenario: File change during cell edit
- **WHEN** a cell edit is active and the file watcher detects a change to the object file
- **THEN** the edit SHALL be cancelled and a toast warning SHALL be shown
