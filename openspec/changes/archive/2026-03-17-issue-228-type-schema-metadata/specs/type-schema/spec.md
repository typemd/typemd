## MODIFIED Requirements

### Requirement: Use entries only allow pin and emoji overrides

A `use` property entry SHALL only contain the fields `use`, `pin`, `emoji`, and `description`. Any other fields (type, options, default, target, etc.) SHALL be rejected by validation.

#### Scenario: Use with pin override accepted

- **WHEN** a type schema contains `- use: due_date` with `pin: 1`
- **THEN** schema validation SHALL accept without error

#### Scenario: Use with emoji override accepted

- **WHEN** a type schema contains `- use: due_date` with `emoji: 🗓️`
- **THEN** schema validation SHALL accept without error

#### Scenario: Use with pin and emoji overrides accepted

- **WHEN** a type schema contains `- use: due_date` with `pin: 1` and `emoji: 🗓️`
- **THEN** schema validation SHALL accept without error

#### Scenario: Use with description override accepted

- **WHEN** a type schema contains `- use: due_date` with `description: "Project deadline"`
- **THEN** schema validation SHALL accept without error

#### Scenario: Use with all allowed overrides accepted

- **WHEN** a type schema contains `- use: due_date` with `pin: 1`, `emoji: 🗓️`, and `description: "Project deadline"`
- **THEN** schema validation SHALL accept without error

#### Scenario: Use with type field rejected

- **WHEN** a type schema contains `- use: due_date` with `type: string`
- **THEN** schema validation SHALL return an error indicating only `pin`, `emoji`, and `description` overrides are allowed on `use` entries

#### Scenario: Use with options field rejected

- **WHEN** a type schema contains `- use: priority` with `options: [{value: urgent}]`
- **THEN** schema validation SHALL return an error indicating only `pin`, `emoji`, and `description` overrides are allowed on `use` entries
