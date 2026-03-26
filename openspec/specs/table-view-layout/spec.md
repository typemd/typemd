## MODIFIED Requirements

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
