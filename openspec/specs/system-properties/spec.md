### Requirement: All objects have a required name property

Every object SHALL have a `name` property that serves as its primary display title. The `name` property is implicit — it does not need to be declared in type schemas. The `name` property SHALL be stored in the object's YAML frontmatter.

#### Scenario: New object has name populated from slug

- **WHEN** a new object is created with slug "golang-in-action"
- **THEN** the object's frontmatter SHALL contain `name: golang-in-action`

#### Scenario: Object with explicit name

- **WHEN** an object's frontmatter contains `name: "Go 語言實戰"`
- **THEN** `GetName()` SHALL return "Go 語言實戰"

### Requirement: GetName method provides centralized name access

The Object struct SHALL have a `GetName()` method that returns the object's display name. This method SHALL read the `name` property from the object's properties map. If the `name` property is missing or empty, it SHALL fall back to `DisplayName()` (ULID-stripped filename).

#### Scenario: GetName returns name property

- **WHEN** an object has `name: "Clean Code"` in its properties
- **THEN** `GetName()` SHALL return "Clean Code"

#### Scenario: GetName falls back to DisplayName when name is missing

- **WHEN** an object has no `name` property in its properties
- **THEN** `GetName()` SHALL return the result of `DisplayName()`

#### Scenario: GetName falls back to DisplayName when name is empty

- **WHEN** an object has `name: ""` in its properties
- **THEN** `GetName()` SHALL return the result of `DisplayName()`

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

### Requirement: Name is always first in frontmatter

When writing an object's frontmatter, `name` SHALL always appear as the first key, followed by `description` (if present), then other system properties (`created_at`, `updated_at`) if present, then schema-defined properties.

#### Scenario: Frontmatter key ordering

- **WHEN** an object with `name: "Clean Code"`, `description: "A handbook of agile software craftsmanship"`, `created_at: "2026-03-01T10:00:00+08:00"`, `updated_at: "2026-03-11T18:00:00+08:00"`, `author: "Robert Martin"`, and `rating: 5` is saved
- **THEN** the frontmatter SHALL have keys in order: `name`, `description`, `created_at`, `updated_at`, `author`, `rating`

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

### Requirement: Name property is editable

Users SHALL be able to modify the `name` property value in the object's markdown file. The updated value SHALL be reflected in all display contexts after the next sync or reload.

#### Scenario: User edits name in frontmatter

- **WHEN** a user changes `name: clean-code` to `name: "Clean Code: A Handbook"` in the markdown file
- **AND** the vault is synced
- **THEN** `GetName()` SHALL return "Clean Code: A Handbook"

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
