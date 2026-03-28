## ADDED Requirements

### Requirement: Local property identification
The system SHALL mark properties in `DisplayProperty` as local (`IsLocal = true`) when the property exists in the object's frontmatter but is not defined in the type schema and is not a system property.

#### Scenario: Schema-defined property is not local
- **WHEN** an object has a property that is defined in its type schema
- **THEN** the `DisplayProperty` for that property SHALL have `IsLocal = false`

#### Scenario: System property is not local
- **WHEN** an object has a system property (name, description, created_at, updated_at, tags, locked)
- **THEN** the `DisplayProperty` for that property SHALL have `IsLocal = false`

#### Scenario: Extra property is local
- **WHEN** an object has a property in its frontmatter that is not in the type schema and not a system property
- **THEN** the `DisplayProperty` for that property SHALL have `IsLocal = true`

#### Scenario: Object with no local properties
- **WHEN** all properties in the object's frontmatter are defined in the schema or are system properties
- **THEN** no `DisplayProperty` SHALL have `IsLocal = true`

#### Scenario: Object with no schema
- **WHEN** an object's type has no schema (schema is nil)
- **THEN** all non-system properties SHALL have `IsLocal = true`

### Requirement: TUI local property separator
The TUI property panel SHALL render a visual separator (`── Local Properties ──`) between schema-defined properties and local properties.

#### Scenario: Separator before first local property
- **WHEN** an object has both schema-defined and local properties
- **THEN** a separator line `── Local Properties ──` SHALL appear between the last schema property and the first local property

#### Scenario: No separator when no local properties
- **WHEN** an object has no local properties
- **THEN** no separator SHALL be rendered

#### Scenario: Separator when only local properties
- **WHEN** an object has only local properties (no schema-defined properties beyond system properties)
- **THEN** the separator SHALL appear before the first local property

### Requirement: TUI local properties are read-only
Local properties in the TUI property panel SHALL be displayed in dim style and the cursor SHALL skip them during navigation.

#### Scenario: Local property not editable
- **WHEN** a property has `IsLocal = true`
- **THEN** the property editor SHALL treat it as non-editable (cursor skips it)

#### Scenario: Local property dim rendering
- **WHEN** a local property is displayed in the property panel
- **THEN** it SHALL be rendered with dim styling (same as reverse relations and backlinks)

### Requirement: CLI local property section
The `tmd object show` command SHALL display local properties in a separate section with a "Local Properties" header.

#### Scenario: CLI shows local properties separately
- **WHEN** an object has both schema-defined and local properties
- **THEN** `tmd object show` SHALL display schema properties under "Properties" and local properties under "Local Properties" with a separator

#### Scenario: CLI no local section when none exist
- **WHEN** an object has no local properties
- **THEN** `tmd object show` SHALL only show the "Properties" section

#### Scenario: CLI local property format
- **WHEN** local properties are displayed in `tmd object show`
- **THEN** each local property SHALL be formatted as `key: value` (same as schema properties)
