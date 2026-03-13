## MODIFIED Requirements

### Requirement: Property type allowlist expanded

The schema validation SHALL accept the following property types: `string`, `number`, `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`, `relation`. The type `enum` SHALL be rejected with a message directing users to use `select` instead. The property name `name` SHALL be allowed in the `properties` array ONLY when it contains exclusively the `template` field. All other system property names SHALL be rejected as reserved.

#### Scenario: New property types accepted
- **WHEN** a type schema defines properties with types `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`
- **THEN** schema validation SHALL accept all of them

#### Scenario: Enum type rejected with guidance
- **WHEN** a type schema defines a property with `type: enum`
- **THEN** schema validation SHALL return an error message indicating to use `select` instead

#### Scenario: Property named "name" with template accepted
- **WHEN** a type schema defines a property with `name: name` and `template: "{{ date:YYYY-MM-DD }}"`
- **THEN** schema validation SHALL accept it without error

#### Scenario: Property named "name" without template rejected
- **WHEN** a type schema defines a property with `name: name` and `type: string`
- **THEN** schema validation SHALL return an error indicating only `template` is allowed on the `name` system property entry

#### Scenario: Property named "description" rejected
- **WHEN** a type schema defines a property with `name: description`
- **THEN** schema validation SHALL return an error indicating that "description" is a reserved system property
