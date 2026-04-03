## ADDED Requirements

### Requirement: Archived system property

The system SHALL support an `archived` boolean system property on objects. When `archived` is `true`, the object is considered soft-deleted — hidden from default queries but not removed from the filesystem. When `archived` is `false` or absent, the property SHALL be omitted from frontmatter (omitempty behavior).

#### Scenario: Object with archived true

- **WHEN** an object has `archived: true` in its frontmatter
- **THEN** `IsArchived()` SHALL return `true`

#### Scenario: Object without archived property

- **WHEN** an object has no `archived` property in its frontmatter
- **THEN** `IsArchived()` SHALL return `false`

#### Scenario: Archived false is omitted from frontmatter

- **WHEN** an object is unarchived (archived set to false)
- **THEN** the `archived` key SHALL be removed from frontmatter (not written as `archived: false`)

### Requirement: Archive and unarchive operations

The system SHALL provide `SetArchived(id, archived)` on `ObjectService` to set or clear the archived flag on an object. This method SHALL emit an `ObjectUpserted` domain event.

#### Scenario: Archive an object

- **WHEN** `SetArchived(id, true)` is called on a non-archived object
- **THEN** the object's `archived` property SHALL be set to `true` and the file SHALL be updated

#### Scenario: Unarchive an object

- **WHEN** `SetArchived(id, false)` is called on an archived object
- **THEN** the `archived` property SHALL be removed from the object and the file SHALL be updated

#### Scenario: Archive an already-archived object

- **WHEN** `SetArchived(id, true)` is called on an already-archived object
- **THEN** the operation SHALL be a no-op (no error, no file write)

#### Scenario: Unarchive a non-archived object

- **WHEN** `SetArchived(id, false)` is called on a non-archived object
- **THEN** the operation SHALL be a no-op (no error, no file write)

#### Scenario: Archive a locked object

- **WHEN** `SetArchived(id, true)` is called on a locked object
- **THEN** the operation SHALL succeed (archive bypasses the lock guard)

### Requirement: Default query exclusion

The `QueryService.Query()` method SHALL exclude objects with `archived: true` by default. An `IncludeArchived` option SHALL allow callers to opt in to including archived objects.

#### Scenario: Default query excludes archived objects

- **WHEN** a query is executed without specifying include-archived
- **THEN** objects with `archived: true` SHALL NOT appear in the results

#### Scenario: Query with include-archived returns all objects

- **WHEN** a query is executed with include-archived enabled
- **THEN** all objects (including archived ones) SHALL appear in the results

#### Scenario: Direct GetObject returns archived objects

- **WHEN** `GetObject(id)` is called with an archived object's ID
- **THEN** the object SHALL be returned (GetObject is not affected by archive filtering)

### Requirement: CLI archive command

The system SHALL provide `tmd object archive <id>` to archive an object and `tmd object unarchive <id>` to unarchive an object. Both commands SHALL use `resolveIDInteractive` for prefix matching.

#### Scenario: Archive an object via CLI

- **WHEN** `tmd object archive book/my-book-<ulid>` is run
- **THEN** the object SHALL be archived and the command SHALL print "Archived book/my-book-<ulid>"

#### Scenario: Unarchive an object via CLI

- **WHEN** `tmd object unarchive book/my-book-<ulid>` is run
- **THEN** the object SHALL be unarchived and the command SHALL print "Unarchived book/my-book-<ulid>"

#### Scenario: Archive already-archived object via CLI

- **WHEN** `tmd object archive <id>` is run on an already-archived object
- **THEN** the command SHALL print "Object <id> is already archived" and exit successfully

#### Scenario: Unarchive non-archived object via CLI

- **WHEN** `tmd object unarchive <id>` is run on a non-archived object
- **THEN** the command SHALL print "Object <id> is not archived" and exit successfully

### Requirement: Include-archived flag on list and query

The `tmd object list` and `tmd object query` commands SHALL support a `--include-archived` flag that includes archived objects in results.

#### Scenario: List with include-archived flag

- **WHEN** `tmd object list --include-archived` is run
- **THEN** archived objects SHALL be included in the output

#### Scenario: List without include-archived flag

- **WHEN** `tmd object list` is run
- **THEN** archived objects SHALL be excluded from the output

#### Scenario: Query with include-archived flag

- **WHEN** `tmd object query --include-archived -f type=book` is run
- **THEN** archived books SHALL be included in the results

### Requirement: Frontmatter property ordering

The `archived` property SHALL appear after `locked` in the canonical frontmatter order (last among system properties).

#### Scenario: Archived property position in frontmatter

- **WHEN** an object has both `locked: true` and `archived: true`
- **THEN** in the frontmatter output, `locked` SHALL appear before `archived`
