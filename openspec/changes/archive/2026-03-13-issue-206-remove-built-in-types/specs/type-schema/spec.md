## REMOVED Requirements

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

**Reason**: Built-in types `book`, `person`, `note` are removed. Only `tag` remains as a built-in type. Users define all other types via `.typemd/types/*.yaml`.
**Migration**: Create `.typemd/types/book.yaml`, `.typemd/types/person.yaml`, `.typemd/types/note.yaml` with the desired schema definitions.

## MODIFIED Requirements

### Requirement: Custom type emoji overrides built-in default

~~When a custom type schema defines its own emoji, it SHALL override the built-in default emoji for that type.~~

~~#### Scenario: Custom book type with different emoji~~

~~- **WHEN** a custom `book.yaml` defines `emoji: 📖`~~
~~- **THEN** the loaded book type SHALL have emoji "📖" instead of the built-in "📚"~~

When a custom type schema defines its own emoji, it SHALL override the built-in default emoji for that type. Since only `tag` remains as a built-in type, this override behavior only applies to the `tag` type.

#### Scenario: Custom tag type with different emoji

- **WHEN** a custom `tag.yaml` defines `emoji: 🔖`
- **THEN** the loaded tag type SHALL have emoji "🔖" instead of the built-in "🏷️"

## ADDED Requirements

### Requirement: Only tag is a built-in type

The `defaultTypes` map SHALL contain only the `tag` type. All other types MUST be defined via `.typemd/types/*.yaml` files.

#### Scenario: Loading an undefined type returns error

- **WHEN** no `.typemd/types/book.yaml` exists and no built-in `book` type is defined
- **AND** `LoadType("book")` is called
- **THEN** it SHALL return an error containing "unknown type: book"

#### Scenario: Tag type loads without custom schema

- **WHEN** no `.typemd/types/tag.yaml` exists
- **AND** `LoadType("tag")` is called
- **THEN** it SHALL return the built-in tag type schema with emoji "🏷️"

#### Scenario: User-defined type loads from YAML

- **WHEN** `.typemd/types/book.yaml` exists with a valid schema
- **AND** `LoadType("book")` is called
- **THEN** it SHALL return the schema from the YAML file
