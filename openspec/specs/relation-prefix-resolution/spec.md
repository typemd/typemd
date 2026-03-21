### Requirement: Relation values without ULID suffix are resolved by name

During Projector sync, when a relation property value does not end with a ULID suffix, the system SHALL treat it as a `type/name` reference and attempt to resolve it to a full object ID using a shared name-to-ID resolution method.

#### Scenario: Name resolves to a unique match

- **WHEN** a book's `author` property contains `person/john-doe` and exactly one person object has name `john-doe`
- **THEN** the reference SHALL be resolved to the full ID (e.g., `person/john-doe-01abc...`)

#### Scenario: Name with no matches

- **WHEN** a book's `author` property contains `person/nobody` and no person objects have name `nobody`
- **THEN** the reference SHALL NOT be resolved
- **AND** the unresolved reference SHALL be reported in `SyncResult`

#### Scenario: Name with multiple matches (ambiguous)

- **WHEN** a book's `author` property contains `person/john` and two person objects have name `john`
- **THEN** the reference SHALL NOT be resolved
- **AND** the ambiguous reference SHALL be reported in `SyncResult` with the list of candidates

#### Scenario: Value with ULID suffix is treated as a full ID

- **WHEN** a book's `author` property contains `person/john-doe-01jqr3k5mpbvn8e0f2g7h9txyz`
- **THEN** the value SHALL be treated as a full ID without name resolution

#### Scenario: Relation target type constrains resolution

- **WHEN** a book's `author` relation has `target: person` and the value `person/john-doe` is provided
- **THEN** name resolution SHALL only search within the `person` type

### Requirement: Resolved prefixes are auto-expanded in the object file

When a prefix is successfully resolved to a unique full ID, the Projector SHALL write the expanded ID back to the object's frontmatter file, replacing the prefix with the full ID.

#### Scenario: Prefix is expanded in file after sync

- **WHEN** a book file contains `author: person/john-doe` and sync resolves it to `person/john-doe-01abc...`
- **THEN** the book file SHALL be updated to contain `author: person/john-doe-01abc...`

#### Scenario: Multiple relation properties expanded in one file

- **WHEN** a book file contains `author: person/john-doe` and `editor: person/jane-smith` and both resolve uniquely
- **THEN** both properties SHALL be expanded in a single file write

#### Scenario: Unresolvable prefix is left unchanged

- **WHEN** a book file contains `author: person/nobody` and no match is found
- **THEN** the file SHALL NOT be modified

### Requirement: Multi-value relation prefixes are resolved individually

For relation properties with `multiple: true`, each value in the array SHALL be resolved independently.

#### Scenario: Array with mixed full IDs and prefixes

- **WHEN** a person's `books` property contains `[book/clean-code-01abc..., book/refactoring]`
- **THEN** the full ID SHALL be kept as-is and the prefix SHALL be resolved independently
- **AND** the file SHALL be updated with only the resolved prefix expanded

#### Scenario: Array with one ambiguous prefix

- **WHEN** a person's `books` property contains `[book/clean-code, book/refactoring]` and `book/clean-code` matches multiple objects
- **THEN** `book/clean-code` SHALL be left as-is (ambiguous)
- **AND** `book/refactoring` SHALL be resolved if it matches uniquely

### Requirement: Name-to-ID resolution is a shared method usable by relations and wiki-links

The system SHALL provide a shared name resolution method that builds a per-type name index from walked objects and resolves `type/name` references to full object IDs. This method SHALL be usable by both relation prefix resolution and future wiki-link shorthand resolution.

#### Scenario: Name index is built from walked objects

- **WHEN** the Projector walks objects and finds person `john-doe-01abc...` with name `john-doe`
- **THEN** the name index SHALL contain an entry mapping `person` + `john-doe` to `person/john-doe-01abc...`

#### Scenario: Name index handles slugified names

- **WHEN** an object was created with name `John Doe` and has slug `john-doe`
- **THEN** the name index SHALL map both the slug `john-doe` and the original name `John Doe` to the object's full ID

#### Scenario: Duplicate names within same type produce ambiguous entries

- **WHEN** two person objects both have name `john-doe`
- **THEN** the name index SHALL mark `person` + `john-doe` as ambiguous

### Requirement: SyncResult reports prefix resolution outcomes

`SyncResult` SHALL include fields reporting the number of prefixes expanded and the list of unresolved references.

#### Scenario: SyncResult after successful expansions

- **WHEN** sync resolves 3 prefixes and finds 1 ambiguous reference
- **THEN** `SyncResult.Expanded` SHALL be 3
- **AND** `SyncResult.Unresolved` SHALL contain 1 entry with the prefix and reason

#### Scenario: SyncResult with no prefix references

- **WHEN** all relation values already have ULID suffixes
- **THEN** `SyncResult.Expanded` SHALL be 0
- **AND** `SyncResult.Unresolved` SHALL be empty
