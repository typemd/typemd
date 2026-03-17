## ADDED Requirements

### Requirement: Type schema supports optional description field

The TypeSchema struct SHALL support an optional `description` field that stores a string value. When a type schema YAML file includes a `description` field, it SHALL be parsed and stored. When the field is omitted, the description SHALL default to an empty string. This field describes the purpose of the type and is distinct from the `description` system property on objects.

#### Scenario: Type schema with description defined

- **WHEN** a type schema YAML file contains `description: "Slide decks and presentation materials"`
- **THEN** the loaded TypeSchema SHALL have its Description field set to "Slide decks and presentation materials"

#### Scenario: Type schema without description defined

- **WHEN** a type schema YAML file does not contain a `description` field
- **THEN** the loaded TypeSchema SHALL have its Description field set to an empty string

### Requirement: Type description included in YAML serialization

When a TypeSchema with a non-empty description is serialized to YAML, the output SHALL include the `description` field. When description is empty, it SHALL be omitted from output.

#### Scenario: Non-empty description serialized

- **WHEN** a TypeSchema with Description "Slide decks and presentation materials" is serialized to YAML
- **THEN** the output SHALL contain `description: Slide decks and presentation materials`

#### Scenario: Empty description omitted

- **WHEN** a TypeSchema with empty Description is serialized to YAML
- **THEN** the output SHALL NOT contain a `description:` line

### Requirement: Property supports optional description field

The Property struct SHALL support an optional `description` field that stores a string value. When a property definition includes a `description` field, it SHALL be parsed and stored. When the field is omitted, the description SHALL default to an empty string.

#### Scenario: Property with description defined

- **WHEN** a type schema property definition contains `description: "The person who gave this presentation"`
- **THEN** the loaded Property SHALL have its Description field set to "The person who gave this presentation"

#### Scenario: Property without description defined

- **WHEN** a type schema property definition does not contain a `description` field
- **THEN** the loaded Property SHALL have its Description field set to an empty string

### Requirement: Property description included in YAML serialization

When a Property with a non-empty description is serialized to YAML, the output SHALL include the `description` field. When description is empty, it SHALL be omitted from output.

#### Scenario: Property with description serialized

- **WHEN** a Property with Description "The person who gave this presentation" is serialized to YAML
- **THEN** the output SHALL contain `description: The person who gave this presentation`

#### Scenario: Property without description omitted

- **WHEN** a Property with empty Description is serialized to YAML
- **THEN** the output SHALL NOT contain a `description:` line

### Requirement: Shared property description can be overridden in use entries

A `use` property entry SHALL allow a `description` field as an override, alongside existing `pin` and `emoji` overrides. The overridden description SHALL replace the shared property's description for that type context.

#### Scenario: Use with description override accepted

- **WHEN** a type schema contains `- use: due_date` with `description: "Project deadline"`
- **THEN** schema validation SHALL accept without error

#### Scenario: Use entry resolved with description override

- **WHEN** a shared property `due_date` has `description: "A date something is due"` and a type schema contains `- use: due_date` with `description: "Project deadline"`
- **THEN** `LoadType()` SHALL return a Property with description "Project deadline"

#### Scenario: Use entry resolved without description override

- **WHEN** a shared property `due_date` has `description: "A date something is due"` and a type schema contains `- use: due_date` without a description override
- **THEN** `LoadType()` SHALL return a Property with description "A date something is due"
