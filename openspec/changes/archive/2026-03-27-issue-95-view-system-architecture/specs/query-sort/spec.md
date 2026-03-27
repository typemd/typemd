## ADDED Requirements

### Requirement: Query accepts variadic sort rules

`QueryService.Query()` SHALL accept a filter string followed by variadic `SortRule` parameters: `Query(filter string, sort ...SortRule)`. Existing callers that pass only a filter string SHALL continue to work without modification.

#### Scenario: Query with filter only

- **WHEN** `Query("type=book")` is called with no sort rules
- **THEN** the result SHALL be equivalent to the previous behavior

#### Scenario: Query with empty filter

- **WHEN** `Query("")` is called
- **THEN** the result SHALL return all objects

### Requirement: Query supports sort by property

When `QueryOptions.Sort` contains one or more `SortRule`, the query result SHALL be sorted by the specified properties in the given order. The first SortRule is the primary sort, subsequent rules break ties.

#### Scenario: Sort by single property ascending

- **WHEN** `Query(QueryOptions{Filter: "type=book", Sort: [{Property: "name", Direction: "asc"}]})` is called
- **THEN** the returned objects SHALL be sorted by the "name" property in ascending alphabetical order

#### Scenario: Sort by single property descending

- **WHEN** `Query(QueryOptions{Filter: "type=book", Sort: [{Property: "rating", Direction: "desc"}]})` is called
- **THEN** the returned objects SHALL be sorted by the "rating" property in descending order

#### Scenario: Sort by multiple properties

- **WHEN** `Query(QueryOptions{Sort: [{Property: "type", Direction: "asc"}, {Property: "name", Direction: "asc"}]})` is called
- **THEN** the returned objects SHALL be sorted by type ascending first, then by name ascending within each type

### Requirement: Sort handles system properties

Sort SHALL support system properties: `name`, `created_at`, `updated_at`. System properties SHALL be extracted from the properties JSON the same way as custom properties.

#### Scenario: Sort by created_at

- **WHEN** `Query(QueryOptions{Filter: "type=book", Sort: [{Property: "created_at", Direction: "desc"}]})` is called
- **THEN** the returned objects SHALL be sorted by creation date, most recent first

#### Scenario: Sort by name

- **WHEN** `Query(QueryOptions{Filter: "type=book", Sort: [{Property: "name", Direction: "asc"}]})` is called
- **THEN** the returned objects SHALL be sorted alphabetically by name

### Requirement: Sort with missing property values

When an object does not have the sort property, it SHALL be treated as having a null/empty value. Objects with missing sort values SHALL appear after objects with values when sorting ascending, and before when sorting descending.

#### Scenario: Missing property in ascending sort

- **WHEN** sorting by "rating" ascending and some objects have no "rating" property
- **THEN** objects without "rating" SHALL appear after objects with "rating" values

#### Scenario: Missing property in descending sort

- **WHEN** sorting by "rating" descending and some objects have no "rating" property
- **THEN** objects without "rating" SHALL appear after objects with "rating" values

### Requirement: SQLiteObjectIndex generates ORDER BY clause

The `SQLiteObjectIndex` SHALL translate `SortRule` into SQL `ORDER BY` clauses using `json_extract(properties, '$.<property>')`. Multiple sort rules SHALL produce multiple ORDER BY expressions in the specified order.

#### Scenario: Single sort to SQL

- **WHEN** a query has sort [{property: "rating", direction: "desc"}]
- **THEN** the SQL SHALL include `ORDER BY json_extract(properties, '$.rating') DESC`

#### Scenario: Multiple sorts to SQL

- **WHEN** a query has sort [{property: "type", direction: "asc"}, {property: "name", direction: "asc"}]
- **THEN** the SQL SHALL include `ORDER BY json_extract(properties, '$.type') ASC, json_extract(properties, '$.name') ASC`

#### Scenario: No sort produces no ORDER BY

- **WHEN** a query has no sort rules
- **THEN** the SQL SHALL NOT include an ORDER BY clause
