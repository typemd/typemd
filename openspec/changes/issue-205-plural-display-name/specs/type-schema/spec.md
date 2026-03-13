## MODIFIED Requirements

### Requirement: Only tag is a built-in type

The `defaultTypes` map SHALL contain only the `tag` type. All other types MUST be defined via `.typemd/types/*.yaml` files. The built-in `tag` type SHALL include `Plural: "tags"`.

#### Scenario: Loading an undefined type returns error

- **WHEN** no `.typemd/types/book.yaml` exists and no built-in `book` type is defined
- **AND** `LoadType("book")` is called
- **THEN** it SHALL return an error containing "unknown type: book"

#### Scenario: Tag type loads without custom schema

- **WHEN** no `.typemd/types/tag.yaml` exists
- **AND** `LoadType("tag")` is called
- **THEN** it SHALL return the built-in tag type schema with emoji "🏷️" and plural "tags"

#### Scenario: User-defined type loads from YAML

- **WHEN** `.typemd/types/book.yaml` exists with a valid schema
- **AND** `LoadType("book")` is called
- **THEN** it SHALL return the schema from the YAML file
