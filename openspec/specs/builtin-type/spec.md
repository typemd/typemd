### Requirement: Tag is a built-in type

The `defaultTypes` map SHALL contain a `tag` type entry. The built-in `tag` type SHALL have emoji "🏷️", plural "tags", unique `true`, and properties: color (string, emoji "🎨") and icon (string, emoji "✨"). The `tag` type backs the `tags` system property.

#### Scenario: Tag type loads without custom schema

- **WHEN** no `.typemd/types/tag.yaml` exists
- **AND** `LoadType("tag")` is called
- **THEN** it SHALL return the built-in tag type schema with emoji "🏷️", plural "tags", and unique true

#### Scenario: Tag type has color and icon properties

- **WHEN** a tag object is created
- **THEN** the object SHALL have property "color" and property "icon"

### Requirement: Page is a built-in type

The `defaultTypes` map SHALL contain a `page` type entry. The built-in `page` type SHALL have emoji "📄", plural "pages", and unique `false`. It SHALL have no custom properties (only system properties apply).

#### Scenario: Page type loads without custom schema

- **WHEN** no `.typemd/types/page.yaml` exists
- **AND** `LoadType("page")` is called
- **THEN** it SHALL return the built-in page type schema with emoji "📄", plural "pages", and unique false

#### Scenario: Page type has no custom properties

- **WHEN** the built-in page type is loaded
- **THEN** its Properties list SHALL be empty

### Requirement: Built-in types cannot be deleted

The system SHALL prevent deletion of any built-in type. Attempting to delete a built-in type SHALL return an error containing "cannot delete built-in type".

#### Scenario: Deleting tag type returns error

- **WHEN** user attempts to delete the `tag` type
- **THEN** the system SHALL return an error containing "cannot delete built-in type"

#### Scenario: Deleting page type returns error

- **WHEN** user attempts to delete the `page` type
- **THEN** the system SHALL return an error containing "cannot delete built-in type"

### Requirement: Built-in types appear in type listing

All built-in types SHALL appear in the type listing alongside custom types.

#### Scenario: ListTypes includes all built-in types

- **WHEN** no custom types exist beyond built-in defaults
- **THEN** `ListTypes()` SHALL return a list containing both "tag" and "page"

#### Scenario: Custom type overrides built-in in listing

- **WHEN** a custom `.typemd/types/page.yaml` exists
- **THEN** `ListTypes()` SHALL include "page" exactly once (no duplicates)

### Requirement: Custom schema overrides built-in

When a custom `.typemd/types/<name>.yaml` exists for a built-in type, it SHALL take precedence over the built-in definition. The custom schema completely replaces the built-in defaults.

#### Scenario: Custom tag schema with different emoji

- **WHEN** `.typemd/types/tag.yaml` exists with `emoji: 🔖`
- **AND** `LoadType("tag")` is called
- **THEN** it SHALL return the custom schema with emoji "🔖" instead of the built-in "🏷️"

#### Scenario: Custom page schema with properties

- **WHEN** `.typemd/types/page.yaml` exists with custom properties
- **AND** `LoadType("page")` is called
- **THEN** it SHALL return the custom schema, not the built-in default
