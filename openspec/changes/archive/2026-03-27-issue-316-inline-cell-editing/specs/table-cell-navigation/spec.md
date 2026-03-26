## ADDED Requirements

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
