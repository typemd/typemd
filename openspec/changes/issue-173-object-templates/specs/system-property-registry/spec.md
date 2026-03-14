## MODIFIED Requirements

### Requirement: System property registry defines all system-managed properties

The core package SHALL maintain a registry of system properties as a package-level slice. Each entry SHALL declare a property name, type, and immutability. Entries with `type: relation` SHALL additionally declare `Target` and `Multiple` fields. The registry SHALL define properties in display order: `name`, `description`, `created_at`, `updated_at`, `tags`. Properties with `Immutable: true` are auto-managed and SHALL NOT be overridden by templates or user input during creation. Properties with `Immutable: false` are user-authored and MAY be overridden.

#### Scenario: Registry contains all system properties with immutability

- **WHEN** the system property registry is queried
- **THEN** it SHALL contain entries for `name` (text, mutable), `description` (text, mutable), `created_at` (datetime, immutable), `updated_at` (datetime, immutable), and `tags` (relation, target: tag, multiple: true, mutable) in that order

## ADDED Requirements

### Requirement: IsImmutableSystemProperty identifies auto-managed properties

The `IsImmutableSystemProperty(name)` function SHALL return `true` for any system property with `Immutable: true` in the registry (`created_at`, `updated_at`), and `false` for all other names (including mutable system properties and non-system properties).

#### Scenario: Immutable system property created_at

- **WHEN** `IsImmutableSystemProperty("created_at")` is called
- **THEN** it SHALL return `true`

#### Scenario: Immutable system property updated_at

- **WHEN** `IsImmutableSystemProperty("updated_at")` is called
- **THEN** it SHALL return `true`

#### Scenario: Mutable system property name

- **WHEN** `IsImmutableSystemProperty("name")` is called
- **THEN** it SHALL return `false`

#### Scenario: Mutable system property description

- **WHEN** `IsImmutableSystemProperty("description")` is called
- **THEN** it SHALL return `false`

#### Scenario: Mutable system property tags

- **WHEN** `IsImmutableSystemProperty("tags")` is called
- **THEN** it SHALL return `false`

#### Scenario: Non-system property

- **WHEN** `IsImmutableSystemProperty("title")` is called
- **THEN** it SHALL return `false`
