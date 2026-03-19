### Requirement: Vault-wide statistics summary
The system SHALL provide a vault-wide statistics summary showing each type's emoji, display name (plural if available), object count, and most recent update time.

#### Scenario: Vault with multiple types
- **WHEN** user runs `tmd stats` on a vault containing books (3), ideas (2), and tags (1)
- **THEN** the output SHALL list each type with emoji, plural name, count, and last updated date, followed by a total count

#### Scenario: Empty vault
- **WHEN** user runs `tmd stats` on a vault with no objects
- **THEN** the output SHALL display a total of 0 with no type rows

#### Scenario: Built-in types without schema files
- **WHEN** user runs `tmd stats` and the vault contains `tag` objects but no `.typemd/types/tag.yaml`
- **THEN** the output SHALL use the built-in type definition (emoji 🏷️, plural "tags")

### Requirement: Single-type property statistics
The system SHALL provide per-property aggregate statistics for a specified type, showing the property name, type, fill rate, and type-appropriate aggregations.

#### Scenario: Number property aggregation
- **WHEN** user runs `tmd stats --type book` and the `book` type has a `rating` number property with values [3, 4, 5]
- **THEN** the output SHALL show sum (12), average (4.0), min (3), max (5), and filled count (3)

#### Scenario: Select property distribution
- **WHEN** user runs `tmd stats --type book` and the `book` type has a `status` select property
- **THEN** the output SHALL show the count for each option value (e.g., reading: 2, finished: 3)

#### Scenario: Multi-select property distribution
- **WHEN** user runs `tmd stats --type book` and the `book` type has a `genres` multi_select property
- **THEN** the output SHALL show the count for each selected option across all objects

#### Scenario: Checkbox property ratio
- **WHEN** user runs `tmd stats --type book` and the `book` type has a `read` checkbox property
- **THEN** the output SHALL show the count of true and false values

#### Scenario: Date property range
- **WHEN** user runs `tmd stats --type book` and the `book` type has a `published` date property with values
- **THEN** the output SHALL show the earliest and latest dates

#### Scenario: Relation property count
- **WHEN** user runs `tmd stats --type book` and the `book` type has an `author` relation property
- **THEN** the output SHALL show the total number of relation links

#### Scenario: Property with no values
- **WHEN** user runs `tmd stats --type book` and no objects have a value for the `rating` property
- **THEN** the output SHALL show filled 0 out of total and skip type-specific stats

#### Scenario: Non-existent type
- **WHEN** user runs `tmd stats --type nonexistent`
- **THEN** the system SHALL return an error indicating the type does not exist

### Requirement: JSON output format
The system SHALL support a `--json` flag that outputs statistics as structured JSON instead of human-readable text.

#### Scenario: Vault-wide JSON output
- **WHEN** user runs `tmd stats --json`
- **THEN** the output SHALL be valid JSON containing an array of type summaries and a total count

#### Scenario: Single-type JSON output
- **WHEN** user runs `tmd stats --type book --json`
- **THEN** the output SHALL be valid JSON containing the type name, object count, and per-property statistics with type-specific aggregation data

### Requirement: QueryService aggregation methods
The system SHALL expose vault statistics through `QueryService` methods, making aggregation logic reusable across CLI, MCP, and Web consumers.

#### Scenario: VaultStats method
- **WHEN** a consumer calls `QueryService.VaultStats()`
- **THEN** it SHALL return a `VaultStats` struct with per-type summaries and total count

#### Scenario: TypeStats method
- **WHEN** a consumer calls `QueryService.TypeStats("book")`
- **THEN** it SHALL return a `TypeStats` struct with property-level aggregations for all aggregable properties of the `book` type
