## Requirements

### Requirement: Table layout displays objects in columnar format
The table layout SHALL display objects as rows with a NAME column followed by property columns, with a header row and separator. The table SHALL support per-cell styling for crosshair highlighting when a column cursor is active.

#### Scenario: Table with default columns
- **WHEN** a view has `layout: table` and no `columns` field
- **THEN** the table SHALL display NAME plus all schema properties as columns (pinned first, then unpinned)

#### Scenario: Table with configured columns
- **WHEN** a view has `layout: table` and `columns: [status, rating]`
- **THEN** the table SHALL display NAME, status, and rating columns in that order
- **AND** other properties SHALL NOT be shown as columns

#### Scenario: Per-cell rendering for crosshair
- **WHEN** a column cursor is active in table view
- **THEN** each cell SHALL be independently styled based on its position relative to the cursor row and cursor column

### Requirement: Horizontal cell navigation in table view
The table view SHALL support horizontal cell navigation using left/right arrow keys or h/l keys, in addition to existing vertical row navigation.

#### Scenario: Move cursor right
- **WHEN** the table view is focused and the user presses right arrow or `l`
- **THEN** the column cursor SHALL move to the next column

#### Scenario: Move cursor left
- **WHEN** the table view is focused and the user presses left arrow or `h`
- **THEN** the column cursor SHALL move to the previous column

#### Scenario: Clamp at right boundary
- **WHEN** the column cursor is on the last column and the user presses right
- **THEN** the column cursor SHALL remain on the last column

#### Scenario: Clamp at left boundary
- **WHEN** the column cursor is on the NAME column (index 0) and the user presses left
- **THEN** the column cursor SHALL remain on the NAME column

#### Scenario: Column cursor on group header
- **WHEN** the row cursor is on a group header row and the user presses left or right
- **THEN** the column cursor SHALL not move (left/right is a no-op on headers)

### Requirement: Column cursor index mapping
The column cursor SHALL use index 0 for the NAME column and index 1..N for property columns mapped from `viewColumns()`.

#### Scenario: Column 0 is NAME
- **WHEN** the column cursor is at index 0
- **THEN** the active cell SHALL be the NAME column

#### Scenario: Column 1+ maps to properties
- **WHEN** the column cursor is at index 1
- **THEN** the active cell SHALL correspond to the first property column from `viewColumns()`

### Requirement: Tab navigation between editable cells
The table view SHALL support Tab to move to the next editable cell and Shift+Tab to move to the previous editable cell.

#### Scenario: Tab moves to next editable cell
- **WHEN** the user presses Tab in the table view
- **THEN** the cursor SHALL move to the next editable cell, skipping read-only cells

#### Scenario: Tab wraps to next row
- **WHEN** the cursor is on the last editable cell of a row and the user presses Tab
- **THEN** the cursor SHALL move to the first editable cell of the next object row

#### Scenario: Shift+Tab moves to previous editable cell
- **WHEN** the user presses Shift+Tab in the table view
- **THEN** the cursor SHALL move to the previous editable cell, skipping read-only cells

### Requirement: Crosshair highlighting
The table view SHALL display a crosshair highlight centered on the active cell: the cursor row gets a dim background, the cursor column header gets a dim tint, and the active cell gets a strong highlight.

#### Scenario: Active cell highlight
- **WHEN** a cell is at the intersection of the cursor row and cursor column
- **THEN** it SHALL be rendered with a strong foreground and background highlight

#### Scenario: Cursor row highlight
- **WHEN** a row is the cursor row but a cell is not the cursor column
- **THEN** the cell SHALL be rendered with a dim background highlight

#### Scenario: Cursor column header highlight
- **WHEN** a column is the cursor column
- **THEN** the column header SHALL be rendered with a dim tint to indicate the active column

#### Scenario: Non-cursor cells
- **WHEN** a cell is not on the cursor row or cursor column
- **THEN** it SHALL be rendered with default styling

### Requirement: Key remapping for object detail
The `o` key SHALL open the object detail view, replacing the current Enter behavior.

#### Scenario: Open object detail with o
- **WHEN** the user presses `o` on an object row in the table view
- **THEN** the object detail view SHALL open for that object

#### Scenario: Enter no longer opens object detail
- **WHEN** the user presses Enter on an object row in the table view
- **THEN** the cell editor SHALL activate (not the object detail view)

### Requirement: Navigation disabled during editing
Cell navigation keys SHALL be disabled while a cell edit is active.

#### Scenario: Arrow keys during editing
- **WHEN** a cell edit is active and the user presses arrow keys
- **THEN** the arrow keys SHALL be handled by the edit widget (e.g., textinput cursor movement), not cell navigation

### Requirement: Navigation disabled when preview or editor is open
Cell navigation and editing SHALL be disabled when the preview panel or view editor is open.

#### Scenario: Cell editing disabled with preview open
- **WHEN** the preview panel is open and the user presses Enter on a cell
- **THEN** the cell editor SHALL NOT activate

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
