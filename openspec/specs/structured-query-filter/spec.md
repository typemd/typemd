# structured-query-filter Specification

## Purpose
Defines the structured filter parameter (`[]FilterRule`) that replaces legacy string-based filtering in the query pipeline. All query methods accept typed filter rules instead of raw filter strings.

## Requirements

### Requirement: Structured filter parameter in query pipeline
The `ObjectIndex.Query()`, `QueryService.Query()`, and `Vault.QueryObjects()` methods SHALL accept a `[]FilterRule` parameter instead of a `filter string` parameter. Each `FilterRule` specifies a property name, operator, and value. An empty slice (`[]FilterRule{}` or `nil`) SHALL return all objects (no filtering).

#### Scenario: Query with single type filter
- **WHEN** a caller invokes `Query([]FilterRule{{Property: "type", Operator: "is", Value: "book"}})`
- **THEN** the system SHALL return only objects of type "book"

#### Scenario: Query with multiple filters
- **WHEN** a caller invokes `Query` with `[]FilterRule{{Property: "type", Operator: "is", Value: "book"}, {Property: "name", Operator: "is", Value: "Go in Action"}}`
- **THEN** the system SHALL return only objects matching ALL filter conditions (AND logic)

#### Scenario: Query with empty filter
- **WHEN** a caller invokes `Query(nil)` or `Query([]FilterRule{})`
- **THEN** the system SHALL return all objects without filtering

#### Scenario: Query with sort rules
- **WHEN** a caller invokes `Query([]FilterRule{{Property: "type", Operator: "is", Value: "book"}}, SortRule{Field: "name", Desc: false})`
- **THEN** the system SHALL return filtered objects sorted by the specified rules

### Requirement: Type property maps to SQL column
When the `FilterRule.Property` is `"type"`, the `SQLiteObjectIndex` SHALL map it to the `type` column in the `objects` table instead of using `json_extract(properties, '$.type')`.

#### Scenario: Type filter uses column
- **WHEN** a query includes `FilterRule{Property: "type", Operator: "is", Value: "book"}`
- **THEN** the generated SQL SHALL use `type = ?` (column reference), not `json_extract(properties, '$.type') = ?`

#### Scenario: Non-type property uses json_extract
- **WHEN** a query includes `FilterRule{Property: "status", Operator: "is", Value: "reading"}`
- **THEN** the generated SQL SHALL use `json_extract(properties, '$.status') = ?`

### Requirement: TypeFilter convenience function
The system SHALL provide a `TypeFilter(typeName string) []FilterRule` convenience function that returns a single-element `[]FilterRule` with `Property: "type"`, `Operator: "is"`, and `Value` set to the given type name.

#### Scenario: TypeFilter returns correct filter
- **WHEN** a caller invokes `TypeFilter("book")`
- **THEN** the result SHALL be `[]FilterRule{{Property: "type", Operator: "is", Value: "book"}}`

### Requirement: Legacy filter code removed
The `parseFilter()` function, `filterCondition` type, and `tmd query` command SHALL be removed from the codebase.

#### Scenario: parseFilter removed
- **WHEN** searching the codebase for `parseFilter` or `filterCondition`
- **THEN** no production code references SHALL exist (test code referencing the removal is acceptable)

#### Scenario: tmd query command removed
- **WHEN** a user runs `tmd query "type=book"`
- **THEN** the system SHALL report an unknown command error
