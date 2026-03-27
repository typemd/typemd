## MODIFIED Requirements

### Requirement: Type schema supports optional emoji field

The TypeSchema struct SHALL support an optional `emoji` field that stores a string value. When a type schema YAML file includes an `emoji` field, it SHALL be parsed and stored. When the field is omitted, the emoji SHALL default to an empty string. The type schema MAY be located at either `.typemd/types/<name>.yaml` (legacy) or `.typemd/types/<name>/schema.yaml` (directory format).

#### Scenario: Type schema with emoji defined

- **WHEN** a type schema YAML file contains `emoji: 📚`
- **THEN** the loaded TypeSchema SHALL have its Emoji field set to "📚"

#### Scenario: Type schema without emoji defined

- **WHEN** a type schema YAML file does not contain an `emoji` field
- **THEN** the loaded TypeSchema SHALL have its Emoji field set to an empty string

#### Scenario: Type schema loaded from directory format

- **WHEN** `.typemd/types/book/schema.yaml` contains `emoji: 📚`
- **THEN** the loaded TypeSchema SHALL have its Emoji field set to "📚"

### Requirement: Custom type emoji overrides built-in default

When a custom type schema defines its own emoji, it SHALL override the built-in default emoji for that type. This applies to any built-in type (`tag`, `page`) when a custom schema is defined. The custom schema MAY be in either single-file or directory format.

#### Scenario: Custom tag type with different emoji

- **WHEN** a custom `tag.yaml` or `tag/schema.yaml` defines `emoji: 🔖`
- **THEN** the loaded tag type SHALL have emoji "🔖" instead of the built-in "🏷️"
