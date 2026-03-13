## ADDED Requirements

### Requirement: Type schema supports optional plural field

The TypeSchema struct SHALL support an optional `plural` field that stores a string value. When a type schema YAML file includes a `plural` field, it SHALL be parsed and stored. When the field is omitted, the plural SHALL default to an empty string.

#### Scenario: Type schema with plural defined

- **WHEN** a type schema YAML file contains `plural: books`
- **THEN** the loaded TypeSchema SHALL have its Plural field set to "books"

#### Scenario: Type schema without plural defined

- **WHEN** a type schema YAML file does not contain a `plural` field
- **THEN** the loaded TypeSchema SHALL have its Plural field set to an empty string

### Requirement: PluralName returns plural or falls back to name

The TypeSchema SHALL provide a `PluralName()` method that returns the `Plural` field when set, or falls back to `Name` when `Plural` is empty. No automatic pluralization (e.g., appending "s") SHALL be performed.

#### Scenario: PluralName with plural set

- **WHEN** a TypeSchema has Name "book" and Plural "books"
- **THEN** `PluralName()` SHALL return "books"

#### Scenario: PluralName without plural set

- **WHEN** a TypeSchema has Name "book" and Plural is empty
- **THEN** `PluralName()` SHALL return "book"

#### Scenario: PluralName with non-English name

- **WHEN** a TypeSchema has Name "筆記" and Plural is empty
- **THEN** `PluralName()` SHALL return "筆記" (no "s" appended)

### Requirement: Built-in tag type includes plural

The built-in `tag` type in `defaultTypes` SHALL include `Plural: "tags"`.

#### Scenario: Built-in tag type has plural

- **WHEN** `LoadType("tag")` is called without a custom tag.yaml
- **THEN** the returned TypeSchema SHALL have Plural "tags"

### Requirement: CLI type show displays plural

The `tmd type show` command SHALL display the plural name when a type schema defines a `plural` field.

#### Scenario: Show type with plural

- **WHEN** user runs `tmd type show book` and the book type has plural "books"
- **THEN** the output SHALL include the plural form in the display

#### Scenario: Show type without plural

- **WHEN** user runs `tmd type show` for a type with no plural defined
- **THEN** the output SHALL display the type name without a plural line
