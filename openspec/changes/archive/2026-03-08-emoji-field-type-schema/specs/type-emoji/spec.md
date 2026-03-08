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
