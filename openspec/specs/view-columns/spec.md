### Requirement: ViewConfig supports a columns field
ViewConfig SHALL support an optional `columns` field containing an ordered list of property names to display.

#### Scenario: Columns field in YAML
- **WHEN** a view YAML contains `columns: [status, rating]`
- **THEN** ViewConfig.Columns SHALL be `["status", "rating"]`

#### Scenario: Missing columns field
- **WHEN** a view YAML has no `columns` field
- **THEN** ViewConfig.Columns SHALL be nil/empty

### Requirement: Default columns depend on layout
When `columns` is empty, the default column behavior SHALL depend on the layout.

#### Scenario: List layout default columns
- **WHEN** `layout: list` and `columns` is empty
- **THEN** no columns SHALL be shown (name only)

#### Scenario: Table layout default columns
- **WHEN** `layout: table` and `columns` is empty
- **THEN** all schema properties SHALL be shown as columns (pinned first, then unpinned)

### Requirement: Explicit columns override defaults
When `columns` is non-empty, both layouts SHALL use exactly the specified properties in order.

#### Scenario: Explicit columns in table layout
- **WHEN** `layout: table` and `columns: [rating, status]`
- **THEN** the table SHALL show NAME, rating, status columns in that order (not all properties)

#### Scenario: Explicit columns in list layout
- **WHEN** `layout: list` and `columns: [status]`
- **THEN** each row SHALL show name followed by the status value inline

### Requirement: View editor supports columns editing
The view editor SHALL provide a Columns section for adding, removing, and reordering columns.

#### Scenario: Add column via editor
- **WHEN** the user adds property `rating` to the columns list in the view editor
- **THEN** the view SHALL save with `columns: [..., rating]` and the display SHALL update immediately

#### Scenario: Remove column via editor
- **WHEN** the user removes property `status` from the columns list in the view editor
- **THEN** the view SHALL save without `status` in columns and the display SHALL update immediately
