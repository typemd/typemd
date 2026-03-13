### Requirement: All objects have a required name property

Every object SHALL have a `name` property that serves as its primary display title. The `name` property is implicit — it does not need to be declared in type schemas. The `name` property SHALL be stored in the object's YAML frontmatter. When a type schema defines a name template, the `name` property SHALL be auto-populated from the evaluated template at creation time if no explicit name is provided.

#### Scenario: New object has name populated from slug

- **WHEN** a new object is created with slug "golang-in-action"
- **THEN** the object's frontmatter SHALL contain `name: golang-in-action`

#### Scenario: Object with explicit name

- **WHEN** an object's frontmatter contains `name: "Go 語言實戰"`
- **THEN** `GetName()` SHALL return "Go 語言實戰"

#### Scenario: New object has name populated from template

- **WHEN** a new object is created with type "journal" that has name template "日記 {{ date:YYYY-MM-DD }}"
- **AND** no explicit name is provided
- **AND** the current date is 2026-03-14
- **THEN** the object's frontmatter SHALL contain `name: 日記 2026-03-14`

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

### Requirement: Name is always first in frontmatter

When writing an object's frontmatter, `name` SHALL always appear as the first key, before any schema-defined properties.

#### Scenario: Frontmatter key ordering

- **WHEN** an object with `name: "Clean Code"`, `author: "Robert Martin"`, and `rating: 5` is saved
- **THEN** the frontmatter SHALL have `name` as the first key, followed by schema-ordered properties

### Requirement: Sync migrates existing objects without name

During vault sync, objects that lack a `name` property SHALL have one automatically added using the value from `DisplayName()` (ULID-stripped filename).

#### Scenario: Sync adds name to existing object

- **WHEN** vault sync encounters an object without a `name` property
- **AND** the object's filename is "clean-code-01jqr3k5mpbvn8e0f2g7h9txyz"
- **THEN** the object SHALL be updated with `name: clean-code` in its frontmatter

#### Scenario: Sync preserves existing name

- **WHEN** vault sync encounters an object that already has `name: "Clean Code"`
- **THEN** the `name` property SHALL remain unchanged

### Requirement: Name property is editable

Users SHALL be able to modify the `name` property value in the object's markdown file. The updated value SHALL be reflected in all display contexts after the next sync or reload.

#### Scenario: User edits name in frontmatter

- **WHEN** a user changes `name: clean-code` to `name: "Clean Code: A Handbook"` in the markdown file
- **AND** the vault is synced
- **THEN** `GetName()` SHALL return "Clean Code: A Handbook"
