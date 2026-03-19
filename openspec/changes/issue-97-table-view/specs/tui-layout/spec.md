## ADDED Requirements

### Requirement: ViewConfig supports table layout constant
The core layer SHALL define a `ViewLayoutTable` constant with value `"table"` alongside the existing `ViewLayoutList`.

#### Scenario: Table layout in YAML
- **WHEN** a view YAML contains `layout: table`
- **THEN** ViewConfig.Layout SHALL be `ViewLayoutTable`

### Requirement: View editor supports layout selection
The view editor SHALL provide a Layout section allowing users to switch between `list` and `table` layouts.

#### Scenario: Switch from list to table
- **WHEN** the user selects `table` in the Layout section of the view editor
- **THEN** the view SHALL save with `layout: table` and the display SHALL switch to columnar table format immediately

#### Scenario: Switch from table to list
- **WHEN** the user selects `list` in the Layout section of the view editor
- **THEN** the view SHALL save with `layout: list` and the display SHALL switch to inline name format immediately
