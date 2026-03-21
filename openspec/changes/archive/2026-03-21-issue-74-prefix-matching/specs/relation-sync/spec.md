## ADDED Requirements

### Requirement: Projector syncs schema-defined relation properties to the relations table

During sync, the Projector SHALL read each object's frontmatter, identify relation properties defined in the object's type schema, and insert corresponding records into the SQLite `relations` table.

#### Scenario: Single-value relation synced from frontmatter

- **WHEN** a book file has `author: person/john-doe-01abc...` and the book schema defines `author` as `type: relation, target: person`
- **THEN** after sync, the `relations` table SHALL contain a record with `name: author`, `from_id: <book-id>`, `to_id: person/john-doe-01abc...`

#### Scenario: Multi-value relation synced from frontmatter

- **WHEN** a person file has `books: [book/a-01abc..., book/b-01xyz...]` and the person schema defines `books` as `type: relation, target: book, multiple: true`
- **THEN** after sync, the `relations` table SHALL contain two records for `books` linking the person to each book

#### Scenario: Relation value referencing non-existent object is skipped

- **WHEN** a book file has `author: person/gone-01abc...` but no such object exists on disk
- **THEN** the relation SHALL NOT be inserted into the `relations` table

#### Scenario: Non-relation properties are ignored

- **WHEN** a book file has a `genre` property with `type: select`
- **THEN** the `genre` value SHALL NOT be inserted into the `relations` table

### Requirement: Full sync clears and rebuilds non-tag relations

During a full `Projector.Sync()`, the system SHALL delete all non-tag relation records from the `relations` table before rebuilding from frontmatter. Tag relations SHALL continue to be managed by the existing `syncTagRelations` mechanism.

#### Scenario: Full sync rebuilds relations from scratch

- **WHEN** a full sync runs and a book previously linked to person A now has `author: person/b-01xyz...`
- **THEN** the old relation (book → person A) SHALL be removed and the new relation (book → person B) SHALL be inserted

#### Scenario: Tag relations are not affected by general relation sync

- **WHEN** a full sync runs
- **THEN** tag relations SHALL be managed by `syncTagRelations` as before, not by the general relation sync

### Requirement: Incremental sync rebuilds relations for changed objects

During `Projector.SyncFiles()`, the system SHALL delete relation records for changed objects and rebuild them from the updated frontmatter. Relations for unchanged objects SHALL remain intact.

#### Scenario: Incremental sync after editing one object

- **WHEN** a book's `author` property is changed from person A to person B and incremental sync runs
- **THEN** the old relation (book → person A) SHALL be removed
- **AND** the new relation (book → person B) SHALL be inserted
- **AND** relations for other objects SHALL be unaffected
