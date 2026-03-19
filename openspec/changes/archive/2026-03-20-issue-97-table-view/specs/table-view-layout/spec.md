## ADDED Requirements

### Requirement: Table layout displays objects in columnar format
The table layout SHALL display objects as rows with a NAME column followed by property columns, with a header row and separator.

#### Scenario: Table with default columns
- **WHEN** a view has `layout: table` and no `columns` field
- **THEN** the table SHALL display NAME plus all schema properties as columns (pinned first, then unpinned)

#### Scenario: Table with configured columns
- **WHEN** a view has `layout: table` and `columns: [status, rating]`
- **THEN** the table SHALL display NAME, status, and rating columns in that order
- **AND** other properties SHALL NOT be shown as columns

### Requirement: Table layout shows sort indicators in column headers
The table layout SHALL display `↑` or `↓` next to column names that are actively sorted.

#### Scenario: Ascending sort indicator
- **WHEN** a view has `layout: table` and `sort: [{property: status, direction: asc}]`
- **THEN** the STATUS column header SHALL display `STATUS ↑`

#### Scenario: Descending sort indicator
- **WHEN** a view has `layout: table` and `sort: [{property: rating, direction: desc}]`
- **THEN** the RATING column header SHALL display `RATING ↓`

#### Scenario: No sort indicator for unsorted columns
- **WHEN** a view has `layout: table` and sort on `status` only
- **THEN** other column headers (e.g., RATING) SHALL NOT display sort indicators

### Requirement: List layout displays objects as inline name with optional values
The list layout SHALL display each object as a single line: emoji (if defined) + name + inline column values separated by ` · `.

#### Scenario: List with no columns configured
- **WHEN** a view has `layout: list` and no `columns` field
- **THEN** each row SHALL display only the type emoji (if any) and object name

#### Scenario: List with columns configured
- **WHEN** a view has `layout: list` and `columns: [status]`
- **AND** an object has `status: reading`
- **THEN** the row SHALL display `📚 Clean Code · reading`

#### Scenario: List omits empty column values
- **WHEN** a view has `layout: list` and `columns: [status, rating]`
- **AND** an object has `status: reading` but no `rating`
- **THEN** the row SHALL display `📚 Clean Code · reading` (no trailing separator for empty rating)

### Requirement: Both layouts support grouping
Both list and table layouts SHALL display group headers when `group_by` is configured.

#### Scenario: Group headers in list layout
- **WHEN** a view has `layout: list` and `group_by: [{property: genre}]`
- **THEN** group headers SHALL appear as `── fiction ──` between groups (same format as table)

### Requirement: Both layouts support filter and sort
Both list and table layouts SHALL respect `filter` and `sort` rules from ViewConfig.

#### Scenario: Filtered list view
- **WHEN** a view has `layout: list` with a filter `status = reading`
- **THEN** only objects matching the filter SHALL appear in the list
