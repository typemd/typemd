## ADDED Requirements

### Requirement: Type schema supports optional emoji field

The TypeSchema struct SHALL support an optional `emoji` field that stores a string value. When a type schema YAML file includes an `emoji` field, it SHALL be parsed and stored. When the field is omitted, the emoji SHALL default to an empty string.

#### Scenario: Type schema with emoji defined

- **WHEN** a type schema YAML file contains `emoji: 📚`
- **THEN** the loaded TypeSchema SHALL have its Emoji field set to "📚"

#### Scenario: Type schema without emoji defined

- **WHEN** a type schema YAML file does not contain an `emoji` field
- **THEN** the loaded TypeSchema SHALL have its Emoji field set to an empty string

### Requirement: Built-in default types include emoji

Built-in default types SHALL include predefined emoji values for visual identification.

#### Scenario: Book default type has emoji

- **WHEN** the built-in "book" type is loaded
- **THEN** its emoji SHALL be "📚"

#### Scenario: Person default type has emoji

- **WHEN** the built-in "person" type is loaded
- **THEN** its emoji SHALL be "👤"

#### Scenario: Note default type has emoji

- **WHEN** the built-in "note" type is loaded
- **THEN** its emoji SHALL be "📝"

### Requirement: Custom type emoji overrides built-in default

When a custom type schema defines its own emoji, it SHALL override the built-in default emoji for that type.

#### Scenario: Custom book type with different emoji

- **WHEN** a custom `book.yaml` defines `emoji: 📖`
- **THEN** the loaded book type SHALL have emoji "📖" instead of the built-in "📚"

### Requirement: CLI type show displays emoji

The `tmd type show` command SHALL display the emoji alongside the type name when an emoji is defined.

#### Scenario: Show type with emoji

- **WHEN** user runs `tmd type show book` and the book type has emoji "📚"
- **THEN** the output SHALL include the emoji in the type display

#### Scenario: Show type without emoji

- **WHEN** user runs `tmd type show` for a type with no emoji
- **THEN** the output SHALL display the type without any emoji prefix

### Requirement: CLI type list displays emoji

The `tmd type list` command SHALL display emoji alongside type names when available.

#### Scenario: List types with emoji

- **WHEN** user runs `tmd type list` and types have emoji defined
- **THEN** each type with an emoji SHALL show the emoji alongside its name

### Requirement: Property struct supports options field

The Property struct SHALL support an `options` field containing an array of Option objects. Each Option SHALL have a required `value` string field and an optional `label` string field. When `label` is omitted, it SHALL default to the `value`. The Property struct SHALL also support an optional `emoji` string field for compact visual identification.

#### Scenario: Property with options defined
- **WHEN** a type schema property has `options: [{value: reading, label: Reading}]`
- **THEN** the loaded Property SHALL have Options with one entry where Value is "reading" and Label is "Reading"

#### Scenario: Property with options without label
- **WHEN** a type schema property has `options: [{value: reading}]`
- **THEN** the loaded Property SHALL have Options with one entry where Value is "reading" and Label defaults to "reading"

#### Scenario: Property with emoji and options
- **WHEN** a type schema property has `emoji: 📊` and `options: [{value: active}]`
- **THEN** the loaded Property SHALL have Emoji "📊" and Options with one entry where Value is "active"

### Requirement: Property type allowlist expanded

The schema validation SHALL accept the following property types: `string`, `number`, `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`, `relation`. The type `enum` SHALL be rejected with a message directing users to use `select` instead.

#### Scenario: New property types accepted
- **WHEN** a type schema defines properties with types `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`
- **THEN** schema validation SHALL accept all of them

#### Scenario: Enum type rejected with guidance
- **WHEN** a type schema defines a property with `type: enum`
- **THEN** schema validation SHALL return an error message indicating to use `select` instead
