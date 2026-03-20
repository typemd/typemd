## ADDED Requirements

### Requirement: View editor opens from view mode
The TUI SHALL open a view editor when the user presses `e` in view mode. The view editor SHALL appear as the right split panel, replacing the preview panel if active.

#### Scenario: Open view editor
- **WHEN** the user is in view mode and presses `e`
- **THEN** the right panel SHALL display the view editor with the current view's filter, sort, and group rules

#### Scenario: Editor replaces preview
- **WHEN** the user has the preview panel open and presses `e`
- **THEN** the preview panel SHALL close and the view editor SHALL open in its place

#### Scenario: Close view editor
- **WHEN** the view editor is open and the user presses `Esc`
- **THEN** the view editor SHALL close and the view mode SHALL return to full-width table (or preview if previously active)

### Requirement: View editor displays three rule sections
The view editor SHALL display three sections vertically: Filter, Sort, and Group By. Each section SHALL list its current rules and a "+ Add" action row.

#### Scenario: Editor with existing rules
- **WHEN** the view has filter `[{property: status, operator: is, value: reading}]`, sort `[{property: name, direction: asc}]`, and group_by `[{property: genre}]`
- **THEN** the editor SHALL display the filter rule "status is reading", the sort rule "name asc", and the group rule "genre"

#### Scenario: Editor with no rules
- **WHEN** the view has no filter, sort, or group rules
- **THEN** the editor SHALL display three empty sections each with only a "+ Add" action row

### Requirement: Navigation between sections and rules
The user SHALL navigate between sections with Tab and within a section with ↑/↓ keys.

#### Scenario: Tab cycles sections
- **WHEN** the user is in the Filter section and presses Tab
- **THEN** the focus SHALL move to the Sort section

#### Scenario: Tab wraps around
- **WHEN** the user is in the Group By section and presses Tab
- **THEN** the focus SHALL wrap to the Filter section

#### Scenario: Arrow keys within section
- **WHEN** the user presses ↓ within a section that has 2 rules and cursor is on rule 1
- **THEN** the cursor SHALL move to rule 2

### Requirement: Add filter rule
The user SHALL add a filter rule via the "+ Add Filter" action. The flow SHALL be: select property → select operator → enter value → confirm.

#### Scenario: Add filter rule flow
- **WHEN** the user presses Enter on "+ Add Filter"
- **THEN** a property picker SHALL appear listing all properties of the type schema
- **AND** after selecting a property, an operator picker SHALL appear with type-appropriate operators
- **AND** after selecting an operator, a value text input SHALL appear (unless operator is `is_empty` or `is_not_empty`)
- **AND** after confirming the value (Enter), the filter rule SHALL be added and the view SHALL re-query

#### Scenario: Add filter with no-value operator
- **WHEN** the user selects operator `is_empty` or `is_not_empty`
- **THEN** the value input step SHALL be skipped and the rule SHALL be added immediately

### Requirement: Add sort rule
The user SHALL add a sort rule via the "+ Add Sort" action. The flow SHALL be: select property → toggle direction → confirm.

#### Scenario: Add sort rule flow
- **WHEN** the user presses Enter on "+ Add Sort"
- **THEN** a property picker SHALL appear
- **AND** after selecting a property, a direction toggle SHALL appear defaulting to "asc"
- **AND** the user MAY toggle between "asc" and "desc" with ↑/↓ or Space
- **AND** after confirming (Enter), the sort rule SHALL be added and the view SHALL re-query

### Requirement: Add group rule
The user SHALL add a group rule via the "+ Add Group" action. The flow SHALL be: select property → confirm.

#### Scenario: Add group rule flow
- **WHEN** the user presses Enter on "+ Add Group"
- **THEN** a property picker SHALL appear
- **AND** after selecting a property (Enter), the group rule SHALL be added and the view SHALL re-query

### Requirement: Edit existing rule
The user SHALL edit an existing rule by pressing Enter while the cursor is on that rule. The edit flow SHALL reuse the same pickers as the add flow, pre-populated with the current values.

#### Scenario: Edit filter rule
- **WHEN** the cursor is on a filter rule and the user presses Enter
- **THEN** the property picker SHALL appear pre-populated with the current property
- **AND** after selecting property, the operator picker SHALL appear pre-populated with the current operator
- **AND** after selecting operator, the value input SHALL appear pre-populated with the current value
- **AND** after confirming, the rule SHALL be updated in place and the view SHALL re-query

#### Scenario: Edit sort rule
- **WHEN** the cursor is on a sort rule and the user presses Enter
- **THEN** the property picker SHALL appear pre-populated with the current property
- **AND** the direction toggle SHALL appear pre-populated with the current direction
- **AND** after confirming, the rule SHALL be updated in place and the view SHALL re-query

#### Scenario: Edit group rule
- **WHEN** the cursor is on a group rule and the user presses Enter
- **THEN** the property picker SHALL appear pre-populated with the current property
- **AND** after confirming, the rule SHALL be updated in place and the view SHALL re-query

### Requirement: Delete rule
The user SHALL delete a rule by pressing `x` or `d` while the cursor is on that rule. Deletion SHALL be immediate (no confirmation dialog).

#### Scenario: Delete filter rule
- **WHEN** the cursor is on a filter rule and the user presses `x`
- **THEN** the filter rule SHALL be removed from the view and the view SHALL re-query

#### Scenario: Delete sort rule
- **WHEN** the cursor is on a sort rule and the user presses `d`
- **THEN** the sort rule SHALL be removed and the view SHALL re-query

#### Scenario: Delete group rule
- **WHEN** the cursor is on a group rule and the user presses `x`
- **THEN** the group rule SHALL be removed and the view SHALL re-query

### Requirement: Move rule order
The user SHALL reorder rules within a section by pressing Shift+K (move up) or Shift+J (move down) while the cursor is on a rule. The move SHALL auto-save and re-query.

#### Scenario: Move sort rule up
- **WHEN** the cursor is on the second sort rule and the user presses Shift+K
- **THEN** the rule SHALL swap with the rule above it and the view SHALL re-query with the new sort order

#### Scenario: Move group rule down
- **WHEN** the cursor is on the first of two group rules and the user presses Shift+J
- **THEN** the rule SHALL swap with the rule below it and the view SHALL re-query with the new grouping order

#### Scenario: Move at boundary
- **WHEN** the cursor is on the first rule in a section and the user presses Shift+K
- **THEN** nothing SHALL happen (rule stays in place)

### Requirement: Auto-save on every change
The view editor SHALL persist changes to the view YAML file immediately after each rule add/edit/delete operation.

#### Scenario: Auto-save after adding rule
- **WHEN** a filter rule is added via the editor
- **THEN** the view YAML file SHALL be updated on disk immediately

#### Scenario: Auto-save after deleting rule
- **WHEN** a sort rule is deleted via the editor
- **THEN** the view YAML file SHALL be updated on disk immediately

### Requirement: Live table refresh after rule changes
After each rule change in the view editor, the left table panel SHALL re-query objects using the updated ViewConfig and refresh the display.

#### Scenario: Filter change updates table
- **WHEN** the user adds a filter rule `status is reading`
- **THEN** the left table SHALL immediately show only objects matching the updated filter

#### Scenario: Sort change updates table
- **WHEN** the user changes sort from `name asc` to `updated_at desc`
- **THEN** the left table SHALL immediately re-sort objects

### Requirement: Delete view action
The view editor SHALL provide a delete view action accessible via `D` (shift+d). Deletion SHALL require confirmation (y/n).

#### Scenario: Delete view with confirmation
- **WHEN** the user presses `D` in the view editor
- **THEN** a confirmation prompt SHALL appear: "Delete view '<name>'? [y/n]"
- **AND** pressing `y` SHALL delete the view file and exit view mode
- **AND** pressing `n` or `Esc` SHALL cancel deletion

### Requirement: View editor help bar
The view editor SHALL display a context-sensitive help bar at the bottom of the screen.

#### Scenario: Help bar in browse mode
- **WHEN** the view editor is in browse mode (navigating rules)
- **THEN** the help bar SHALL display "↑↓: navigate  J/K: move  tab: section  enter: edit  x: delete  D: delete view  esc: close"

#### Scenario: Help bar in add mode
- **WHEN** the view editor is in the process of adding a rule (property/operator/value picker)
- **THEN** the help bar SHALL display the current step (e.g., "Select property  ↑↓: navigate  enter: select  esc: cancel")

### Requirement: Property picker with text filtering and list
The property picker SHALL combine a text input with a scrollable list. The text input SHALL filter the list in real-time. The list SHALL display all properties defined in the type schema, including system properties (name, description, tags, created_at, updated_at). The user SHALL navigate the filtered list with ↑/↓ and confirm with Enter.

#### Scenario: Property picker content
- **WHEN** the property picker is shown for a type with properties [status, genre, rating]
- **THEN** the picker SHALL display a text input at top and a list below showing: name, description, tags, created_at, updated_at, status, genre, rating

#### Scenario: Property picker filtering
- **WHEN** the user types "sta" in the property picker text input
- **THEN** the list SHALL filter to show only properties containing "sta" (e.g., "status", "created_at", "updated_at")

#### Scenario: Property picker selection
- **WHEN** the user navigates to a property with ↑/↓ and presses Enter
- **THEN** the property SHALL be selected and the picker SHALL advance to the next step

#### Scenario: Property picker cancel
- **WHEN** the user presses Esc in the property picker
- **THEN** the add/edit flow SHALL be cancelled and the editor SHALL return to browse mode

### Requirement: Operator picker as scrollable list
The operator picker SHALL display a scrollable list of operators valid for the selected property's type, using the `validOperators` registry from `filter_operator.go`. The user SHALL navigate with ↑/↓ and confirm with Enter. No text filtering is needed.

#### Scenario: String property operators
- **WHEN** the user selects a property of type "string" in the filter add flow
- **THEN** the operator picker SHALL show a scrollable list: is, is_not, contains, does_not_contain, starts_with, ends_with, is_empty, is_not_empty

#### Scenario: Checkbox property operators
- **WHEN** the user selects a property of type "checkbox"
- **THEN** the operator picker SHALL show a scrollable list: is, is_not

#### Scenario: Operator picker cancel
- **WHEN** the user presses Esc in the operator picker
- **THEN** the add/edit flow SHALL return to the property picker step

### Requirement: Split panel layout
The view editor SHALL use a 60/40 split with the table on the left and the editor on the right, matching the existing preview split ratio.

#### Scenario: Normal terminal width
- **WHEN** the terminal width is 120 characters or more
- **THEN** the table SHALL occupy approximately 60% and the editor 40% of the width

#### Scenario: Narrow terminal
- **WHEN** the terminal width is less than 80 characters
- **THEN** the editor SHALL take the full width (table hidden)
