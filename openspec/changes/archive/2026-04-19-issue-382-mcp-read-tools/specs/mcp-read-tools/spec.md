## ADDED Requirements

### Requirement: Vault overview via MCP

The MCP server SHALL provide a `vault_overview` tool that returns a one-call summary of the vault's structure. The tool SHALL accept no required parameters. The response SHALL include, for every registered type (built-in and custom), the type name, plural form, emoji, description, object count, and a list of recent objects (ID, name, `updated_at`) capped at a small default (5 per type). This response enables an LLM to orient itself to the vault in a single call.

#### Scenario: Overview of empty vault returns built-in types

- **WHEN** the vault has no objects
- **AND** the client calls `vault_overview`
- **THEN** the response lists every built-in type with count `0` and an empty `recent` list

#### Scenario: Overview reports counts and recent objects per type

- **WHEN** the vault contains objects of multiple types
- **AND** the client calls `vault_overview`
- **THEN** each type entry reports its count
- **AND** the `recent` list is sorted by `updated_at` descending and capped at 5 entries

### Requirement: List objects via MCP

The MCP server SHALL provide a `list_objects` tool that returns object summaries with optional filtering and pagination. The tool SHALL accept `type` (optional string), `limit` (optional integer, default `50`, maximum `500`), and `offset` (optional integer, default `0`). The response SHALL include, for every matching object, its ID, type, name, and `updated_at`. Archived objects SHALL be excluded unless a future option enables them. The response SHALL include a `total` count reflecting all matches before pagination.

#### Scenario: List objects of a single type with pagination

- **WHEN** the vault contains 20 `book` objects
- **AND** the client calls `list_objects` with `type: "book"`, `limit: 10`, `offset: 0`
- **THEN** the response contains 10 object summaries and `total: 20`

#### Scenario: List objects without type filter returns all types

- **WHEN** the client calls `list_objects` with no `type`
- **AND** the vault contains objects of multiple types
- **THEN** the response contains summaries spanning every type

#### Scenario: List objects with unknown type returns empty

- **WHEN** the client calls `list_objects` with `type: "nonexistent"`
- **THEN** the response contains an empty list and `total: 0`

#### Scenario: List objects clamps oversized limit

- **WHEN** the client calls `list_objects` with `limit: 10000`
- **THEN** the server clamps `limit` to `500` and returns at most 500 summaries

### Requirement: Query objects via MCP

The MCP server SHALL provide a `query_objects` tool that runs a structured query against the vault using the existing `FilterRule` system. The tool SHALL accept `filters` (required array of `{property, operator, value}` objects), `sort` (optional array of `{property, direction}` objects where `direction` is `"asc"` or `"desc"`), `limit` (optional, default `50`, maximum `500`), and `offset` (optional, default `0`). The response SHALL include matching objects as summaries and a `total` count. If any `FilterRule` has an invalid `operator` or missing `property`, the tool SHALL return an error without executing the query.

#### Scenario: Query by property equality

- **WHEN** the client calls `query_objects` with `filters: [{property: "status", operator: "is", value: "reading"}]`
- **THEN** the response contains every object whose `status` equals `"reading"`

#### Scenario: Query with sort and limit

- **WHEN** the client calls `query_objects` with `filters: [...]`, `sort: [{property: "updated_at", direction: "desc"}]`, and `limit: 5`
- **THEN** the response contains at most 5 summaries in descending `updated_at` order

#### Scenario: Query with invalid filter returns error

- **WHEN** the client calls `query_objects` with a filter missing the `property` field
- **THEN** the tool returns an error describing the invalid filter

### Requirement: List backlinks via MCP

The MCP server SHALL provide a `list_backlinks` tool that returns all references pointing to a given object. The tool SHALL accept `id` (required) and SHALL resolve abbreviated IDs using the same prefix-matching rules as `get_object`. The response SHALL contain two lists: `wiki_backlinks` (objects with wiki-links pointing to the target) and `relation_backlinks` (objects with typed relations pointing to the target, including the relation name). Each entry SHALL include the source object's ID and name.

#### Scenario: List backlinks for an object with incoming wiki-links

- **WHEN** three objects contain `[[book/clean-code]]` in their body
- **AND** the client calls `list_backlinks` with `id: "book/clean-code"`
- **THEN** `wiki_backlinks` contains the three source objects

#### Scenario: List backlinks for an object with incoming typed relations

- **WHEN** two `review` objects have a `book` relation pointing to `book/clean-code`
- **AND** the client calls `list_backlinks` with `id: "book/clean-code"`
- **THEN** `relation_backlinks` contains the two reviews along with the relation name `"book"`

#### Scenario: List backlinks with abbreviated ID

- **WHEN** the client calls `list_backlinks` with `id: "book/clean-code"` (no ULID suffix)
- **AND** the prefix resolves to a single object
- **THEN** the tool returns the backlinks for the resolved object

#### Scenario: List backlinks for an object with no references

- **WHEN** no object references the target
- **AND** the client calls `list_backlinks` with the target's ID
- **THEN** both `wiki_backlinks` and `relation_backlinks` are empty lists

### Requirement: Vault stats via MCP

The MCP server SHALL provide a `vault_stats` tool that returns per-type property distribution statistics. The tool SHALL accept `type` (required string). The response SHALL include the type name, total object count, and for every schema-defined property, its name, number of objects with a non-empty value, and fill rate (0–1). The tool SHALL return an error if the type does not exist.

#### Scenario: Stats for a type with partial property coverage

- **WHEN** the `book` type has 10 objects and 6 have a non-empty `rating`
- **AND** the client calls `vault_stats` with `type: "book"`
- **THEN** the `rating` entry reports `filled: 6` and `fill_rate: 0.6`

#### Scenario: Stats for an unknown type returns error

- **WHEN** the client calls `vault_stats` with `type: "nonexistent"`
- **THEN** the tool returns an error indicating the type was not found

#### Scenario: Stats for a type with no objects

- **WHEN** the `book` type has no objects
- **AND** the client calls `vault_stats` with `type: "book"`
- **THEN** the response reports `count: 0` and each property has `filled: 0` and `fill_rate: 0`
