## ADDED Requirements

### Requirement: New objects have created_at timestamp

When a new object is created via `NewObject`, it SHALL have a `created_at` property set to the current time in RFC 3339 format with local timezone offset. This property SHALL never be modified after creation.

#### Scenario: New object gets created_at

- **WHEN** a new object is created at 2026-03-11T18:30:00+08:00
- **THEN** the object's frontmatter SHALL contain `created_at: "2026-03-11T18:30:00+08:00"`

#### Scenario: created_at is not modified on save

- **WHEN** an object with `created_at: "2026-03-01T10:00:00+08:00"` is saved via `SaveObject`
- **THEN** the `created_at` property SHALL remain `"2026-03-01T10:00:00+08:00"`

### Requirement: New objects have updated_at timestamp

When a new object is created via `NewObject`, it SHALL have an `updated_at` property set to the current time in RFC 3339 format with local timezone offset.

#### Scenario: New object gets updated_at

- **WHEN** a new object is created at 2026-03-11T18:30:00+08:00
- **THEN** the object's frontmatter SHALL contain `updated_at: "2026-03-11T18:30:00+08:00"`

### Requirement: SaveObject updates updated_at timestamp

When an object is saved via `SaveObject` or `SetProperty`, the `updated_at` property SHALL be updated to the current time in RFC 3339 format with local timezone offset.

#### Scenario: SaveObject updates updated_at

- **WHEN** an object is saved via `SaveObject` at 2026-03-11T20:00:00+08:00
- **THEN** the object's `updated_at` property SHALL be `"2026-03-11T20:00:00+08:00"`

#### Scenario: SetProperty updates updated_at

- **WHEN** `SetProperty` is called on an object at 2026-03-11T20:00:00+08:00
- **THEN** the object's `updated_at` property SHALL be `"2026-03-11T20:00:00+08:00"`

### Requirement: Existing objects without timestamps work gracefully

Objects that do not have `created_at` or `updated_at` properties SHALL continue to function normally. No migration is performed during SyncIndex for these properties.

#### Scenario: Object without timestamps loads successfully

- **WHEN** an object file has frontmatter with only `name: "Clean Code"` and `title: "Clean Code"`
- **THEN** `GetObject` SHALL return the object without error
- **AND** the object's properties SHALL not contain `created_at` or `updated_at`

#### Scenario: SyncIndex does not add timestamps to existing objects

- **WHEN** SyncIndex processes an object without `created_at` or `updated_at`
- **THEN** the object file SHALL not be modified to add these properties

## MODIFIED Requirements

### Requirement: Name is always first in frontmatter

When writing an object's frontmatter, `name` SHALL always appear as the first key, followed by other system properties (`created_at`, `updated_at`) if present, then schema-defined properties.

#### Scenario: Frontmatter key ordering

- **WHEN** an object with `name: "Clean Code"`, `created_at: "2026-03-01T10:00:00+08:00"`, `updated_at: "2026-03-11T18:00:00+08:00"`, `author: "Robert Martin"`, and `rating: 5` is saved
- **THEN** the frontmatter SHALL have keys in order: `name`, `created_at`, `updated_at`, `author`, `rating`

### Requirement: Sync migrates existing objects without name

During vault sync, objects that lack a `name` property SHALL have one automatically added using the value from `DisplayName()` (ULID-stripped filename). No migration is performed for `created_at` or `updated_at`.

#### Scenario: Sync adds name to existing object

- **WHEN** vault sync encounters an object without a `name` property
- **AND** the object's filename is "clean-code-01jqr3k5mpbvn8e0f2g7h9txyz"
- **THEN** the object SHALL be updated with `name: clean-code` in its frontmatter

#### Scenario: Sync preserves existing name

- **WHEN** vault sync encounters an object that already has `name: "Clean Code"`
- **THEN** the `name` property SHALL remain unchanged

#### Scenario: Sync does not add timestamps

- **WHEN** vault sync encounters an object without `created_at` or `updated_at`
- **THEN** the object file SHALL not be modified to add these properties
