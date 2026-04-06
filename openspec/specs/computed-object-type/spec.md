## ADDED Requirements

### Requirement: GetProperty resolves object_type from Object.Type
The `Object.GetProperty("object_type")` method SHALL return the object's type, derived from its file path convention `objects/<type>/<name>.md`. The value SHALL equal `Object.Type`.

#### Scenario: GetProperty returns object_type for a book object
- **WHEN** a book object exists
- **THEN** `GetProperty("object_type")` SHALL return `"book"` and `true`

#### Scenario: GetProperty returns object_type for a person object
- **WHEN** a person object exists
- **THEN** `GetProperty("object_type")` SHALL return `"person"` and `true`

### Requirement: GetProperty resolves stored properties
The `Object.GetProperty(name)` method SHALL return stored frontmatter properties when the property is not derived or computed.

#### Scenario: GetProperty returns a stored property value
- **WHEN** an object has property "title" set to "Go in Action"
- **THEN** `GetProperty("title")` SHALL return `"Go in Action"` and `true`

#### Scenario: GetProperty returns false for missing stored property
- **WHEN** an object does not have property "rating"
- **THEN** `GetProperty("rating")` SHALL return `nil` and `false`

### Requirement: object_type is read-only
Setting `object_type` via `SetProperty` SHALL be rejected with an error mentioning "non-stored system property".

#### Scenario: SetProperty rejects object_type
- **WHEN** a user calls `SetProperty("object_type", "page", schema)`
- **THEN** the operation SHALL return an error containing "non-stored system property"

### Requirement: object_type is not stored in frontmatter
When an object file contains `object_type` in its YAML frontmatter, the reconciler SHALL strip it on save.

#### Scenario: Frontmatter strips object_type on save
- **WHEN** a raw object file contains `object_type: book` in frontmatter
- **AND** the object is saved
- **THEN** the saved file SHALL NOT contain `object_type` in frontmatter
