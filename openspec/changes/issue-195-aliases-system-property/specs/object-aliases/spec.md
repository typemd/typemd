## ADDED Requirements

### Requirement: aliases is a stored system property

The system SHALL support `aliases` as a stored system property of type `list[text]`. It SHALL be optional and default to absent (treated as empty). When present, it SHALL be written to object frontmatter after `tags` and before schema properties.

#### Scenario: Object with aliases serializes them after tags

- **WHEN** an object has `tags: [go]` and `aliases: ["Go 語言", "Golang"]`
- **THEN** the frontmatter contains `aliases` immediately after `tags`

#### Scenario: Object without aliases omits the field

- **WHEN** an object has no aliases set
- **THEN** the frontmatter does not contain an `aliases` key

#### Scenario: Empty aliases array is treated as absent

- **WHEN** an object frontmatter contains `aliases: []`
- **THEN** the object's aliases are treated as empty (no alias lookups)

### Requirement: aliases is a user-authored system property

The system SHALL treat `aliases` as user-authored — it can be set and modified by users. It SHALL NOT be auto-managed by the system.

#### Scenario: User can set aliases on create

- **WHEN** an object is created with `aliases: ["Alt Name"]` in frontmatter
- **THEN** the object's aliases contain "Alt Name"

#### Scenario: User can update aliases on save

- **WHEN** an object's frontmatter is updated to add a new alias
- **THEN** the updated alias is reflected after the next sync

### Requirement: aliases is reserved in type schemas

The system SHALL reject any type schema or shared property file that defines a property named `aliases`. This name is reserved as a system property.

#### Scenario: Type schema with aliases property is rejected

- **WHEN** a type schema defines a property named `aliases`
- **THEN** schema validation reports an error indicating `aliases` is a reserved system property name

#### Scenario: Shared property named aliases is rejected

- **WHEN** a shared property file is named `aliases.yaml`
- **THEN** schema validation reports an error indicating `aliases` is a reserved system property name

### Requirement: aliases are indexed for name lookup

The system SHALL index each alias string as an additional lookup key for the object in the name index. Alias lookup uses the same slugified matching as name and slug lookups.

#### Scenario: Object is found by alias in name index

- **WHEN** an object has alias "Go 語言" and the name index is built
- **THEN** the slugified form of "Go 語言" maps to that object's ID in the name index

#### Scenario: Duplicate alias across objects is treated as ambiguous

- **WHEN** two objects both have alias "Golang"
- **THEN** the name index treats the alias as ambiguous (same behavior as duplicate names)

#### Scenario: Alias matching object name of another object uses name priority

- **WHEN** object A has alias "Clean Code" and object B has name "Clean Code"
- **THEN** wiki-link resolution for `[[clean-code]]` resolves to object B (exact name match takes priority)
