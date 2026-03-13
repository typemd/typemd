## ADDED Requirements

### Requirement: Creation rejects duplicate name for unique types

When a type schema has `unique: true`, creating an object SHALL be rejected if another object of the same type already has an identical `name` property value. The error message SHALL include the conflicting name and type.

#### Scenario: First object with a name succeeds

- **WHEN** type "person" has `unique: true`
- **AND** no existing person object has name "john-doe"
- **THEN** creating a person with name "john-doe" SHALL succeed

#### Scenario: Duplicate name rejected

- **WHEN** type "person" has `unique: true`
- **AND** a person object with name "john-doe" already exists
- **THEN** creating another person with name "john-doe" SHALL fail with an error containing "already exists"

#### Scenario: Same name allowed for different types

- **WHEN** type "person" has `unique: true`
- **AND** type "character" has `unique: true`
- **AND** a person object with name "john-doe" exists
- **THEN** creating a character with name "john-doe" SHALL succeed

#### Scenario: Duplicate name allowed when unique is false

- **WHEN** type "book" does not have `unique: true`
- **AND** a book object with name "clean-code" already exists
- **THEN** creating another book with name "clean-code" SHALL succeed

### Requirement: Tag uniqueness uses the generalized mechanism

The built-in `tag` type SHALL have `unique: true` in its default schema. Tag name uniqueness SHALL be enforced through the same mechanism as any other unique type — no special-case logic.

#### Scenario: Tag rejects duplicate name

- **WHEN** a tag with name "golang" already exists
- **THEN** creating another tag with name "golang" SHALL fail with an error containing "already exists"

#### Scenario: Tag allows unique names

- **WHEN** a tag with name "golang" already exists
- **THEN** creating a tag with name "rust" SHALL succeed

### Requirement: Validation detects name uniqueness violations for all unique types

The `tmd type validate` command SHALL check all types with `unique: true` for duplicate `name` values among their objects. Each violation SHALL report the duplicate name, the type, and the IDs of the conflicting objects.

#### Scenario: Validation passes with no duplicates

- **WHEN** type "person" has `unique: true`
- **AND** all person objects have distinct name values
- **THEN** validation SHALL report no uniqueness errors

#### Scenario: Validation reports duplicates

- **WHEN** type "person" has `unique: true`
- **AND** two person objects have name "john-doe"
- **THEN** validation SHALL report a uniqueness violation for "john-doe" in type "person" with both object IDs

#### Scenario: Validation skips non-unique types

- **WHEN** type "book" does not have `unique: true`
- **AND** two book objects have name "clean-code"
- **THEN** validation SHALL NOT report a uniqueness violation for books

### Requirement: Uniqueness comparison is exact string match

The uniqueness check SHALL compare `name` property values using exact string matching. No normalization (case folding, whitespace trimming) SHALL be applied.

#### Scenario: Different case is not a duplicate

- **WHEN** type "tag" has `unique: true`
- **AND** a tag with name "Golang" exists
- **THEN** creating a tag with name "golang" SHALL succeed

#### Scenario: Exact match is a duplicate

- **WHEN** type "tag" has `unique: true`
- **AND** a tag with name "golang" exists
- **THEN** creating another tag with name "golang" SHALL fail
