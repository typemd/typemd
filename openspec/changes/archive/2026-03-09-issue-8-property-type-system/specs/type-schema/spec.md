## MODIFIED Requirements

### Requirement: Type schema supports optional emoji field

The TypeSchema struct SHALL support an optional `emoji` field that stores a string value. When a type schema YAML file includes an `emoji` field, it SHALL be parsed and stored. When the field is omitted, the emoji SHALL default to an empty string.

#### Scenario: Type schema with emoji defined

- **WHEN** a type schema YAML file contains `emoji: 📚`
- **THEN** the loaded TypeSchema SHALL have its Emoji field set to "📚"

#### Scenario: Type schema without emoji defined

- **WHEN** a type schema YAML file does not contain an `emoji` field
- **THEN** the loaded TypeSchema SHALL have its Emoji field set to an empty string

## ADDED Requirements

### Requirement: Property struct supports options field

The Property struct SHALL support an `options` field containing an array of Option objects. Each Option SHALL have a required `value` string field and an optional `label` string field. When `label` is omitted, it SHALL default to the `value`.

#### Scenario: Property with options defined
- **WHEN** a type schema property has `options: [{value: reading, label: Reading}]`
- **THEN** the loaded Property SHALL have Options with one entry where Value is "reading" and Label is "Reading"

#### Scenario: Property with options without label
- **WHEN** a type schema property has `options: [{value: reading}]`
- **THEN** the loaded Property SHALL have Options with one entry where Value is "reading" and Label defaults to "reading"

### Requirement: Property type allowlist expanded

The schema validation SHALL accept the following property types: `string`, `number`, `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`, `relation`. The type `enum` SHALL be rejected with a message directing users to use `select` instead.

#### Scenario: New property types accepted
- **WHEN** a type schema defines properties with types `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`
- **THEN** schema validation SHALL accept all of them

#### Scenario: Enum type rejected with guidance
- **WHEN** a type schema defines a property with `type: enum`
- **THEN** schema validation SHALL return an error message indicating to use `select` instead

## REMOVED Requirements

### Requirement: Property Values field for enum type
**Reason**: Replaced by `options` field with richer structure (value + label).
**Migration**: Use `tmd migrate` to convert `values: [a, b]` to `options: [{value: a}, {value: b}]`.
