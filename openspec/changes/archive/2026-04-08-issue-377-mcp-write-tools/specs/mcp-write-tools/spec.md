## ADDED Requirements

### Requirement: Create object via MCP

The MCP server SHALL provide a `create_object` tool that creates a new object in the vault. The tool SHALL accept `type` (required), `name` (required), `template` (optional), `properties` (optional JSON object), and `body` (optional string). On success, the tool SHALL return the created object's ID, type, and filename.

#### Scenario: Create object with type and name

- **WHEN** client calls `create_object` with `type: "book"` and `name: "clean-code"`
- **THEN** the tool creates a new book object and returns its ID, type, and filename

#### Scenario: Create object with properties and body

- **WHEN** client calls `create_object` with `type: "book"`, `name: "clean-code"`, `properties: {"status": "reading"}`, and `body: "A handbook of agile software craftsmanship."`
- **THEN** the tool creates the object with the specified properties and body content

#### Scenario: Create object with template

- **WHEN** client calls `create_object` with `type: "book"`, `name: "clean-code"`, and `template: "review"`
- **THEN** the tool creates the object using the specified template

#### Scenario: Create object with invalid type

- **WHEN** client calls `create_object` with `type: "nonexistent"` and `name: "test"`
- **THEN** the tool returns an error indicating the type is invalid

### Requirement: Update object via MCP

The MCP server SHALL provide an `update_object` tool that updates an existing object's properties and/or body. The tool SHALL accept `id` (required), `properties` (optional JSON object), and `body` (optional string). Properties SHALL be merged with existing properties — only provided keys are updated, existing keys not in the input are preserved. If `body` is provided, it SHALL replace the existing body entirely. On success, the tool SHALL return the updated object's ID.

#### Scenario: Update object properties with merge semantics

- **WHEN** an object exists with properties `{"status": "reading", "rating": "5"}`
- **AND** client calls `update_object` with `id` and `properties: {"status": "completed"}`
- **THEN** the object's properties become `{"status": "completed", "rating": "5"}`

#### Scenario: Update object body

- **WHEN** client calls `update_object` with `id` and `body: "New content"`
- **THEN** the object's body is replaced with "New content"

#### Scenario: Update locked object

- **WHEN** an object is locked
- **AND** client calls `update_object` with that object's `id`
- **THEN** the tool returns an error indicating the object is locked

#### Scenario: Update nonexistent object

- **WHEN** client calls `update_object` with `id: "book/nonexistent"`
- **THEN** the tool returns an error indicating the object was not found

### Requirement: Link objects via MCP

The MCP server SHALL provide a `link_objects` tool that creates a relation between two objects. The tool SHALL accept `from_id` (required), `relation` (required), and `to_id` (required). On success, the tool SHALL return a confirmation message.

#### Scenario: Link two objects

- **WHEN** two objects exist and a relation property is defined
- **AND** client calls `link_objects` with `from_id`, `relation`, and `to_id`
- **THEN** the relation is created between the objects

#### Scenario: Link with invalid relation name

- **WHEN** client calls `link_objects` with a relation name that doesn't exist in the type schema
- **THEN** the tool returns an error

### Requirement: Unlink objects via MCP

The MCP server SHALL provide an `unlink_objects` tool that removes a relation between two objects. The tool SHALL accept `from_id` (required), `relation` (required), `to_id` (required), and `both` (optional boolean, default false). When `both` is true, the relation SHALL be removed in both directions. On success, the tool SHALL return a confirmation message.

#### Scenario: Unlink two objects

- **WHEN** two objects are linked by a relation
- **AND** client calls `unlink_objects` with `from_id`, `relation`, and `to_id`
- **THEN** the relation is removed

#### Scenario: Unlink in both directions

- **WHEN** two objects are linked by a relation
- **AND** client calls `unlink_objects` with `from_id`, `relation`, `to_id`, and `both: true`
- **THEN** the relation is removed in both directions
