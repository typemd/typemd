## ADDED Requirements

### Requirement: Source type loads with default schema
The system SHALL provide a built-in `source` type with emoji 📥, plural "sources", and unique false. The source type SHALL be available without a YAML file.

#### Scenario: Source type loads with default emoji and plural
- **WHEN** I load type "source"
- **THEN** the loaded schema SHALL have emoji "📥"
- **AND** the loaded schema SHALL have plural "sources"
- **AND** the loaded schema SHALL have unique false

### Requirement: Source type has provenance properties
The built-in `source` type SHALL include three properties for tracking provenance: `url` (string, 🔗), `author` (string, ✍️), and `ingested_at` (date, 📅).

#### Scenario: Source type has url, author, and ingested_at properties
- **WHEN** I load type "source"
- **THEN** the loaded schema SHALL have a property "url" of type "string" with emoji "🔗"
- **AND** the loaded schema SHALL have a property "author" of type "string" with emoji "✍️"
- **AND** the loaded schema SHALL have a property "ingested_at" of type "date" with emoji "📅"

### Requirement: Source type cannot be deleted
The system SHALL reject attempts to delete the built-in `source` type with an error message containing "cannot delete built-in type".

#### Scenario: Delete built-in source type is rejected
- **WHEN** I delete type "source"
- **THEN** an error SHALL occur
- **AND** the error message SHALL contain "cannot delete built-in type"

### Requirement: Source type can be overridden
Users SHALL be able to override the built-in `source` type by providing a custom `types/source/schema.yaml` file. The custom schema takes precedence over the built-in default.

#### Scenario: Custom source schema overrides built-in default
- **WHEN** a custom "source" type schema with emoji "📦" exists
- **AND** I load type "source"
- **THEN** the loaded schema SHALL have emoji "📦"
