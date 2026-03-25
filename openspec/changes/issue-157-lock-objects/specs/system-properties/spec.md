## MODIFIED Requirements

### Requirement: Name is always first in frontmatter
When writing an object's frontmatter, `name` SHALL always appear as the first key, followed by `description` (if present), then other system properties (`created_at`, `updated_at`, `tags`) if present, then `locked` (if true), then schema-defined properties.

#### Scenario: Frontmatter key ordering with locked
- **WHEN** an object with `name: "Clean Code"`, `description: "A handbook"`, `created_at: "2026-03-01T10:00:00+08:00"`, `updated_at: "2026-03-11T18:00:00+08:00"`, `locked: true`, `author: "Robert Martin"`, and `rating: 5` is saved
- **THEN** the frontmatter SHALL have keys in order: `name`, `description`, `created_at`, `updated_at`, `tags`, `locked`, `author`, `rating`

#### Scenario: Frontmatter ordering without locked
- **WHEN** an object with `name: "Clean Code"`, `description: "A handbook"`, `created_at: "2026-03-01T10:00:00+08:00"`, `updated_at: "2026-03-11T18:00:00+08:00"`, `author: "Robert Martin"`, and `rating: 5` is saved
- **THEN** the frontmatter SHALL have keys in order: `name`, `description`, `created_at`, `updated_at`, `author`, `rating`
- **AND** no `locked` key SHALL appear
