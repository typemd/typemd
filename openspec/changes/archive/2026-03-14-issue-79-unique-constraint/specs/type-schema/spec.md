## MODIFIED Requirements

### Requirement: Only tag is a built-in type

The `defaultTypes` map SHALL contain only the `tag` type. All other types MUST be defined via `.typemd/types/*.yaml` files. The built-in `tag` type SHALL have `Unique: true` in its default schema.

#### Scenario: Loading an undefined type returns error

- **WHEN** no `.typemd/types/book.yaml` exists and no built-in `book` type is defined
- **AND** `LoadType("book")` is called
- **THEN** it SHALL return an error containing "unknown type: book"

#### Scenario: Tag type loads without custom schema

- **WHEN** no `.typemd/types/tag.yaml` exists
- **AND** `LoadType("tag")` is called
- **THEN** it SHALL return the built-in tag type schema with emoji "🏷️" and Unique true

#### Scenario: User-defined type loads from YAML

- **WHEN** `.typemd/types/book.yaml` exists with a valid schema
- **AND** `LoadType("book")` is called
- **THEN** it SHALL return the schema from the YAML file

## ADDED Requirements

### Requirement: Type schema supports optional unique field

The TypeSchema struct SHALL support an optional `unique` field that stores a boolean value. When a type schema YAML file includes `unique: true`, it SHALL be parsed and stored. When the field is omitted, unique SHALL default to `false`.

#### Scenario: Type schema with unique true

- **WHEN** a type schema YAML file contains `unique: true`
- **THEN** the loaded TypeSchema SHALL have its Unique field set to true

#### Scenario: Type schema with unique false

- **WHEN** a type schema YAML file contains `unique: false`
- **THEN** the loaded TypeSchema SHALL have its Unique field set to false

#### Scenario: Type schema without unique field

- **WHEN** a type schema YAML file does not contain a `unique` field
- **THEN** the loaded TypeSchema SHALL have its Unique field set to false

### Requirement: Custom type schema can override built-in unique setting

When a custom type schema overrides a built-in type (e.g., `tag.yaml`), the `unique` field from the custom schema SHALL take effect. If the custom schema omits `unique`, it SHALL default to `false`, overriding the built-in default.

#### Scenario: Custom tag schema with unique true

- **WHEN** `.typemd/types/tag.yaml` exists with `unique: true`
- **THEN** the loaded tag type SHALL have Unique true

#### Scenario: Custom tag schema without unique field

- **WHEN** `.typemd/types/tag.yaml` exists without a `unique` field
- **THEN** the loaded tag type SHALL have Unique false (overriding built-in default)
