## ADDED Requirements

### Requirement: All objects may have an optional description property

Every object MAY have a `description` property that serves as a brief, single-line summary. The `description` property is a stored system property — it is written to the object's YAML frontmatter. It does not need to be declared in type schemas.

#### Scenario: Object with description

- **WHEN** an object's frontmatter contains `description: "A practical guide to Go programming"`
- **THEN** the object's properties SHALL include `description` with value "A practical guide to Go programming"

#### Scenario: Object without description

- **WHEN** an object's frontmatter does not contain a `description` property
- **THEN** the object SHALL load successfully without error
- **AND** the object's properties SHALL not contain a `description` key

### Requirement: Description is not auto-populated

Unlike `name` (auto-populated from slug) and timestamps (auto-set on creation), `description` SHALL NOT be automatically set when creating a new object. It is entirely user-authored.

#### Scenario: New object has no description

- **WHEN** a new object is created via `NewObject`
- **THEN** the object's frontmatter SHALL NOT contain a `description` property

### Requirement: Description is editable

Users SHALL be able to add or modify the `description` property in the object's markdown file. The updated value SHALL be reflected after the next sync or reload.

#### Scenario: User adds description to frontmatter

- **WHEN** a user adds `description: "A practical guide to Go"` to an object's frontmatter
- **AND** the vault is synced
- **THEN** the indexed object's properties SHALL include `description` with value "A practical guide to Go"

### Requirement: Sync does not add description to existing objects

During vault sync, objects that lack a `description` property SHALL NOT have one automatically added. This is consistent with `created_at`/`updated_at` sync behavior.

#### Scenario: Sync preserves absence of description

- **WHEN** vault sync encounters an object without a `description` property
- **THEN** the object file SHALL not be modified to add a `description` property

## MODIFIED Requirements

### Requirement: Name is always first in frontmatter

When writing an object's frontmatter, `name` SHALL always appear as the first key, followed by `description` (if present), then other system properties (`created_at`, `updated_at`) if present, then schema-defined properties.

#### Scenario: Frontmatter key ordering

- **WHEN** an object with `name: "Clean Code"`, `description: "A handbook of agile software craftsmanship"`, `created_at: "2026-03-01T10:00:00+08:00"`, `updated_at: "2026-03-11T18:00:00+08:00"`, `author: "Robert Martin"`, and `rating: 5` is saved
- **THEN** the frontmatter SHALL have keys in order: `name`, `description`, `created_at`, `updated_at`, `author`, `rating`
