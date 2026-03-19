## MODIFIED Requirements

### Requirement: QueryService aggregation methods
The system SHALL expose vault statistics through `QueryService` methods, making aggregation logic reusable across CLI, MCP, and Web consumers. Internally, these methods SHALL use `[]FilterRule` (via `TypeFilter()`) instead of legacy filter strings when querying the index.

#### Scenario: VaultStats method
- **WHEN** a consumer calls `QueryService.VaultStats()`
- **THEN** it SHALL return a `VaultStats` struct with per-type summaries and total count

#### Scenario: TypeStats method
- **WHEN** a consumer calls `QueryService.TypeStats("book")`
- **THEN** it SHALL return a `TypeStats` struct with property-level aggregations for all aggregable properties of the `book` type

#### Scenario: Internal query uses structured filter
- **WHEN** `VaultStats()` or `TypeStats()` queries the index for objects of a specific type
- **THEN** it SHALL pass `TypeFilter(typeName)` (i.e. `[]FilterRule{{Property: "type", Operator: "is", Value: typeName}}`) instead of a filter string
