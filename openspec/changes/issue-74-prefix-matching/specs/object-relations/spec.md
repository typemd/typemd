## MODIFIED Requirements

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
