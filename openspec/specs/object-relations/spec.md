## ADDED Requirements

### Requirement: Relation properties are defined in type schemas

A type schema SHALL support properties with `type: relation`. A relation property MUST declare a `target` type. It MAY declare `multiple: true` for multi-value, `bidirectional: true` with an `inverse` name for two-way linking.

#### Scenario: Valid relation property in schema

- **WHEN** a type schema defines a property with `type: relation` and `target: person`
- **THEN** the schema is valid and the property is recognized as a relation

#### Scenario: Relation property without target is invalid

- **WHEN** a type schema defines a property with `type: relation` but no `target`
- **THEN** schema validation reports an error: relation type requires target

### Requirement: Objects can be linked via relation properties

The system SHALL allow linking two objects through a named relation property using `LinkObjects`. The relation SHALL be persisted in both the source object's frontmatter and the `relations` database table. Additionally, relation values in frontmatter MAY use prefix form (without ULID suffix), which SHALL be resolved and expanded during Projector sync.

#### Scenario: Link two objects

- **WHEN** a book object is linked to a person object via the `author` relation
- **THEN** the book's `author` frontmatter property references the person's full ID
- **AND** a record is inserted into the `relations` table with the correct `name`, `from_id`, and `to_id`

#### Scenario: Link to non-existent object fails

- **WHEN** a link is attempted to an object ID that does not exist
- **THEN** an error is returned

#### Scenario: Link with unknown relation name fails

- **WHEN** a link is attempted using a relation name not defined in the source type schema
- **THEN** an error is returned

#### Scenario: Hand-edited relation with prefix is resolved during sync

- **WHEN** a user manually writes `author: person/john-doe` in a book's frontmatter
- **AND** Projector sync runs
- **THEN** the prefix SHALL be resolved to the full ID and written back to the file
- **AND** a relation record SHALL be inserted into the `relations` table

### Requirement: Target type is validated at link time

The system SHALL verify that the target object's type matches the relation property's `target` field. A type mismatch SHALL result in an error.

#### Scenario: Type mismatch is rejected

- **WHEN** a book's `author` relation (target: person) is linked to another book
- **THEN** an error is returned indicating target type mismatch

### Requirement: Single-value relations overwrite on re-link

When a relation property has `multiple: false` (default), linking to a new target SHALL overwrite the previous value.

#### Scenario: Overwrite single-value relation

- **WHEN** a book is linked to person A via `author`, then linked to person B via `author`
- **THEN** the book's `author` property references person B

### Requirement: Multiple-value relations append and reject duplicates

When a relation property has `multiple: true`, linking SHALL append the target to the existing array. Linking the same target again SHALL be rejected as a duplicate.

#### Scenario: Append to multiple-value relation

- **WHEN** person is linked to book A via `books`, then linked to book B via `books`
- **THEN** the person's `books` property contains both book A and book B

#### Scenario: Duplicate link is rejected

- **WHEN** a multiple-value relation already contains a target and the same link is attempted again
- **THEN** an error is returned indicating duplicate relation value

### Requirement: Bidirectional relations automatically create the inverse

When a relation property has `bidirectional: true` and an `inverse` name, linking A→B SHALL automatically create the inverse link B→A on the target object's inverse property.

#### Scenario: Bidirectional link creates both sides

- **WHEN** a book is linked to a person via `author` (bidirectional, inverse: `books`)
- **THEN** the book's `author` property references the person
- **AND** the person's `books` property contains the book

#### Scenario: Inverse property must exist in target schema

- **WHEN** a bidirectional relation declares inverse `books` but the target type schema has no such property
- **THEN** an error is returned

### Requirement: Objects can be unlinked

The system SHALL allow removing a relation between two objects using `UnlinkObjects`. The relation SHALL be removed from the source object's frontmatter and the `relations` database table.

#### Scenario: Unlink single direction

- **WHEN** a book linked to a person via `author` is unlinked without the `both` flag
- **THEN** the book's `author` property is cleared
- **AND** the person's `books` property still contains the book (inverse not removed)

#### Scenario: Unlink both directions

- **WHEN** a book linked to a person via `author` is unlinked with the `both` flag
- **THEN** the book's `author` property is cleared
- **AND** the person's `books` property no longer contains the book
- **AND** both forward and inverse records are removed from the `relations` table

#### Scenario: Unlink one from multiple-value relation

- **WHEN** a person has `books` containing book A and book B, and book A is unlinked with `both`
- **THEN** the person's `books` property contains only book B

### Requirement: Relations are queryable for an object

The system SHALL provide `ListRelations` to return all relations where the object is either the source (`from_id`) or target (`to_id`).

#### Scenario: List relations for a linked object

- **WHEN** a book is bidirectionally linked to a person via `author`
- **THEN** listing relations for the book returns 2 entries (forward `author` + inverse `books`)

#### Scenario: List relations for an unlinked object

- **WHEN** an object has no relations
- **THEN** listing relations returns an empty list

### Requirement: Tags system property uses the relation mechanism

The `tags` system property SHALL use the same relation linking and unlinking mechanism as schema-defined relation properties. Linking a tag to an object via the `tags` property SHALL follow the same `LinkObjects` path. Unlinking SHALL follow the same `UnlinkObjects` path.

#### Scenario: Link tag to object via tags system property

- **WHEN** a tag object "go" exists and a book object "golang-book" exists
- **AND** "golang-book" is linked to "go" via the `tags` property
- **THEN** the book's `tags` property SHALL contain a reference to the tag

#### Scenario: Unlink tag from object via tags system property

- **WHEN** a book "golang-book" has tag "go" linked via `tags`
- **AND** "golang-book" is unlinked from "go" via `tags` without the `both` flag
- **THEN** the book's `tags` property SHALL be empty

### Requirement: Reverse relations are displayed in object detail

The system SHALL display reverse relations (where the object is the `to_id`) in the object's display properties, shown after schema-defined properties with a `←` indicator.

#### Scenario: Reverse relation appears in display

- **WHEN** book A is linked to person B via `author`
- **THEN** person B's display properties include a reverse relation entry for `author` showing book A with `IsReverse: true`
