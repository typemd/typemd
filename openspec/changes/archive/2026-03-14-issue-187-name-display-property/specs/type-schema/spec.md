## MODIFIED Requirements

### Requirement: Property type allowlist expanded

The schema validation SHALL accept the following property types: `string`, `number`, `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`, `relation`. The type `enum` SHALL be rejected with a message directing users to use `select` instead. The property name `name` SHALL be rejected as a reserved system property name.

#### Scenario: New property types accepted
- **WHEN** a type schema defines properties with types `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`
- **THEN** schema validation SHALL accept all of them

#### Scenario: Enum type rejected with guidance
- **WHEN** a type schema defines a property with `type: enum`
- **THEN** schema validation SHALL return an error message indicating to use `select` instead

#### Scenario: Property named "name" rejected
- **WHEN** a type schema defines a property with `name: name`
- **THEN** schema validation SHALL return an error indicating that "name" is a reserved system property
