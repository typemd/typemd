## ADDED Requirements

### Requirement: Query falls back to filesystem when index is unavailable

The system SHALL fall back to filesystem scanning via `ObjectRepository.Walk()` when `ObjectIndex.Query()` returns an error, applying filter rules and sort rules in memory to produce the same results as the indexed path.

#### Scenario: Query with type filter succeeds despite index unavailability

- **WHEN** the SQLite index is unavailable
- **AND** the user queries objects with a type filter
- **THEN** the system SHALL return all objects matching that type from the filesystem
- **AND** the system SHALL log a warning indicating fallback mode

#### Scenario: Query with property filter succeeds in fallback mode

- **WHEN** the SQLite index is unavailable
- **AND** the user queries objects with a property filter (e.g., "status is active")
- **THEN** the system SHALL return objects whose properties match the filter criteria

#### Scenario: Query with sort rules succeeds in fallback mode

- **WHEN** the SQLite index is unavailable
- **AND** the user queries objects with sort rules
- **THEN** the system SHALL return results sorted according to the specified sort rules

#### Scenario: Query with no filters returns all objects in fallback mode

- **WHEN** the SQLite index is unavailable
- **AND** the user queries with no filter rules
- **THEN** the system SHALL return all objects from the filesystem

### Requirement: Search falls back to substring matching when index is unavailable

The system SHALL fall back to filesystem scanning with case-insensitive substring matching when `ObjectIndex.Search()` returns an error. Matching SHALL be performed against the object's name, description, and body text.

#### Scenario: Search finds objects by name in fallback mode

- **WHEN** the SQLite index is unavailable
- **AND** the user searches for a keyword that appears in an object's name
- **THEN** the system SHALL return that object in the results

#### Scenario: Search finds objects by body text in fallback mode

- **WHEN** the SQLite index is unavailable
- **AND** the user searches for a keyword that appears in an object's body
- **THEN** the system SHALL return that object in the results

#### Scenario: Search returns empty when no match in fallback mode

- **WHEN** the SQLite index is unavailable
- **AND** the user searches for a keyword that does not appear in any object
- **THEN** the system SHALL return an empty result set

### Requirement: In-memory filter supports all standard operators

The system SHALL support in-memory evaluation of all filter operators that `FilterRuleToSQL` supports: `is`, `is_not`, `contains`, `does_not_contain`, `starts_with`, `ends_with`, `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `before`, `after`, `on_or_before`, `on_or_after`, `is_empty`, `is_not_empty`.

#### Scenario: Filter with "is" operator matches exact property value

- **WHEN** a filter rule with operator "is" and value "active" is applied in memory
- **AND** an object has the matching property set to "active"
- **THEN** the object SHALL be included in results

#### Scenario: Filter with "contains" operator matches substring

- **WHEN** a filter rule with operator "contains" and value "go" is applied in memory
- **AND** an object has the matching property set to "golang"
- **THEN** the object SHALL be included in results

#### Scenario: Filter with "is_empty" operator matches nil or empty properties

- **WHEN** a filter rule with operator "is_empty" is applied in memory
- **AND** an object does not have the matching property set (nil or empty string)
- **THEN** the object SHALL be included in results

#### Scenario: Filter with numeric "gt" operator compares numerically

- **WHEN** a filter rule with operator "gt" and value "5" is applied in memory
- **AND** an object has the matching property set to 10
- **THEN** the object SHALL be included in results

### Requirement: VaultStats and TypeStats work in fallback mode

`VaultStats()` and `TypeStats()` SHALL produce correct results when the index is unavailable, since they delegate to `Query()` which has fallback support.

#### Scenario: VaultStats returns correct counts in fallback mode

- **WHEN** the SQLite index is unavailable
- **AND** the vault contains objects of multiple types
- **THEN** `VaultStats()` SHALL return the correct total count and per-type summaries

#### Scenario: TypeStats returns property statistics in fallback mode

- **WHEN** the SQLite index is unavailable
- **AND** the vault contains objects of a specific type with various property values
- **THEN** `TypeStats()` SHALL return correct property statistics for that type

### Requirement: Fallback emits a warning

The system SHALL emit a `slog.Warn` log message each time a fallback path is taken, including the operation name (query or search).

#### Scenario: Warning logged on query fallback

- **WHEN** the SQLite index is unavailable
- **AND** the user performs a query
- **THEN** the system SHALL log a warning with message containing "index unavailable" and "fallback"

#### Scenario: Warning logged on search fallback

- **WHEN** the SQLite index is unavailable
- **AND** the user performs a search
- **THEN** the system SHALL log a warning with message containing "index unavailable" and "fallback"
