## ADDED Requirements

### Requirement: Property supports optional emoji field

The Property struct SHALL support an optional `emoji` field that stores a string value. When a property definition in a type schema YAML includes an `emoji` field, it SHALL be parsed and stored. When the field is omitted, the emoji SHALL default to an empty string.

#### Scenario: Property with emoji defined
- **WHEN** a type schema property definition contains `emoji: 👤`
- **THEN** the loaded Property SHALL have its Emoji field set to "👤"

#### Scenario: Property without emoji defined
- **WHEN** a type schema property definition does not contain an `emoji` field
- **THEN** the loaded Property SHALL have its Emoji field set to an empty string

### Requirement: Property emojis unique within type scope

Within a single type schema, no two properties SHALL have the same non-empty emoji. Schema validation SHALL reject duplicate property emojis.

#### Scenario: Unique property emojis accepted
- **WHEN** a type schema has properties with emojis "👤" and "⭐"
- **THEN** schema validation SHALL accept it without error

#### Scenario: Duplicate property emojis rejected
- **WHEN** a type schema has two properties both with emoji "👤"
- **THEN** schema validation SHALL return an error indicating duplicate property emoji "👤"

#### Scenario: Empty emojis do not conflict
- **WHEN** a type schema has three properties where two have no emoji and one has emoji "📝"
- **THEN** schema validation SHALL accept it without error

#### Scenario: Same emoji allowed across different types
- **WHEN** type "book" has a property with emoji "📝" and type "person" also has a property with emoji "📝"
- **THEN** schema validation SHALL accept both types without error
