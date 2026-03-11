### Requirement: System property registry defines all system-managed properties

The core package SHALL maintain a registry of system properties as a package-level slice. Each entry SHALL declare a property name and type. The registry SHALL define properties in display order: `name`, `created_at`, `updated_at`.

#### Scenario: Registry contains all system properties

- **WHEN** the system property registry is queried
- **THEN** it SHALL contain entries for `name` (text), `created_at` (datetime), and `updated_at` (datetime) in that order

### Requirement: IsSystemProperty identifies reserved property names

The `IsSystemProperty(name)` function SHALL return `true` for any property name defined in the system property registry, and `false` for all other names.

#### Scenario: Recognized system property

- **WHEN** `IsSystemProperty("name")` is called
- **THEN** it SHALL return `true`

#### Scenario: Recognized system property created_at

- **WHEN** `IsSystemProperty("created_at")` is called
- **THEN** it SHALL return `true`

#### Scenario: Recognized system property updated_at

- **WHEN** `IsSystemProperty("updated_at")` is called
- **THEN** it SHALL return `true`

#### Scenario: Non-system property

- **WHEN** `IsSystemProperty("title")` is called
- **THEN** it SHALL return `false`

### Requirement: SystemPropertyNames returns ordered list

The `SystemPropertyNames()` function SHALL return a slice of all system property names in their defined registry order.

#### Scenario: Property names in order

- **WHEN** `SystemPropertyNames()` is called
- **THEN** it SHALL return `["name", "created_at", "updated_at"]`

### Requirement: Type schema validation rejects all system property names

`ValidateSchema` SHALL reject any property whose name matches a system property in the registry. The error message SHALL indicate the property is a reserved system property.

#### Scenario: Schema defines created_at property

- **WHEN** a type schema defines a property named `created_at`
- **THEN** validation SHALL return an error containing "reserved system property"

#### Scenario: Schema defines updated_at property

- **WHEN** a type schema defines a property named `updated_at`
- **THEN** validation SHALL return an error containing "reserved system property"

### Requirement: Shared property validation rejects all system property names

`ValidateSharedProperties` SHALL reject any shared property whose name matches a system property in the registry. The error message SHALL indicate the property is a reserved system property.

#### Scenario: Shared property named created_at

- **WHEN** a shared properties file defines a property named `created_at`
- **THEN** validation SHALL return an error containing "reserved system property"

#### Scenario: Shared property named updated_at

- **WHEN** a shared properties file defines a property named `updated_at`
- **THEN** validation SHALL return an error containing "reserved system property"

### Requirement: SyncIndex preserves all system properties

During property filtering in `SyncIndex`, all system properties present in the object's frontmatter SHALL be preserved in the filtered property set, regardless of type schema definitions.

#### Scenario: Sync preserves created_at and updated_at

- **WHEN** an object has `created_at` and `updated_at` in its frontmatter
- **AND** the type schema does not define these properties
- **THEN** after sync, the indexed properties SHALL include `created_at` and `updated_at`

### Requirement: Frontmatter orders system properties first

`OrderedPropKeys` SHALL place system properties before schema-defined properties, in registry order (name, created_at, updated_at). Schema-defined properties follow in schema order. Extra properties are appended alphabetically.

#### Scenario: Full property ordering

- **WHEN** an object has properties `name`, `created_at`, `updated_at`, `title`, and `rating`
- **AND** the schema defines `title` then `rating`
- **THEN** `OrderedPropKeys` SHALL return `["name", "created_at", "updated_at", "title", "rating"]`

#### Scenario: System properties absent

- **WHEN** an object has properties `name` and `title` (no timestamps)
- **AND** the schema defines `title`
- **THEN** `OrderedPropKeys` SHALL return `["name", "title"]`
