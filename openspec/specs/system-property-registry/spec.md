### Requirement: System property registry defines all system-managed properties

The core package SHALL maintain a registry of system properties as a package-level slice. Each entry SHALL declare a property name and type. Entries with `type: relation` SHALL additionally declare `Target` and `Multiple` fields. The registry SHALL define properties in display order: `name`, `description`, `created_at`, `updated_at`, `tags`, `locked`.

#### Scenario: Registry contains all system properties

- **WHEN** the system property registry is queried
- **THEN** it SHALL contain entries for `name` (text), `description` (text), `created_at` (datetime), `updated_at` (datetime), `tags` (relation, target: tag, multiple: true), and `locked` (checkbox) in that order

### Requirement: IsSystemProperty identifies reserved property names

The `IsSystemProperty(name)` function SHALL return `true` for any property name defined in the system property registry, and `false` for all other names.

#### Scenario: Recognized system property

- **WHEN** `IsSystemProperty("name")` is called
- **THEN** it SHALL return `true`

#### Scenario: Recognized system property description

- **WHEN** `IsSystemProperty("description")` is called
- **THEN** it SHALL return `true`

#### Scenario: Recognized system property created_at

- **WHEN** `IsSystemProperty("created_at")` is called
- **THEN** it SHALL return `true`

#### Scenario: Recognized system property updated_at

- **WHEN** `IsSystemProperty("updated_at")` is called
- **THEN** it SHALL return `true`

#### Scenario: Recognized system property tags

- **WHEN** `IsSystemProperty("tags")` is called
- **THEN** it SHALL return `true`

#### Scenario: Recognized system property locked

- **WHEN** `IsSystemProperty("locked")` is called
- **THEN** it SHALL return `true`

#### Scenario: Non-system property

- **WHEN** `IsSystemProperty("title")` is called
- **THEN** it SHALL return `false`

### Requirement: SystemPropertyNames returns ordered list

The `SystemPropertyNames()` function SHALL return a slice of all system property names in their defined registry order.

#### Scenario: Property names in order

- **WHEN** `SystemPropertyNames()` is called
- **THEN** it SHALL return `["name", "description", "created_at", "updated_at", "tags", "locked"]`

### Requirement: Type schema validation rejects all system property names

`ValidateSchema` SHALL reject any property whose name matches a system property in the registry. The error message SHALL indicate the property is a reserved system property.

#### Scenario: Schema defines description property

- **WHEN** a type schema defines a property named `description`
- **THEN** validation SHALL return an error containing "reserved system property"

#### Scenario: Schema defines created_at property

- **WHEN** a type schema defines a property named `created_at`
- **THEN** validation SHALL return an error containing "reserved system property"

#### Scenario: Schema defines updated_at property

- **WHEN** a type schema defines a property named `updated_at`
- **THEN** validation SHALL return an error containing "reserved system property"

#### Scenario: Schema defines tags property

- **WHEN** a type schema defines a property named `tags`
- **THEN** validation SHALL return an error containing "reserved system property"

#### Scenario: Schema defines locked property

- **WHEN** a type schema defines a property named `locked`
- **THEN** validation SHALL return an error containing "reserved system property"

### Requirement: Shared property validation rejects all system property names

`ValidateSharedProperties` SHALL reject any shared property whose name matches a system property in the registry. The error message SHALL indicate the property is a reserved system property.

#### Scenario: Shared property named description

- **WHEN** a shared properties file defines a property named `description`
- **THEN** validation SHALL return an error containing "reserved system property"

#### Scenario: Shared property named created_at

- **WHEN** a shared properties file defines a property named `created_at`
- **THEN** validation SHALL return an error containing "reserved system property"

#### Scenario: Shared property named updated_at

- **WHEN** a shared properties file defines a property named `updated_at`
- **THEN** validation SHALL return an error containing "reserved system property"

#### Scenario: Shared property named tags

- **WHEN** a shared properties file defines a property named `tags`
- **THEN** validation SHALL return an error containing "reserved system property"

#### Scenario: Shared property named locked

- **WHEN** a shared properties file defines a property named `locked`
- **THEN** validation SHALL return an error containing "reserved system property"

### Requirement: SyncIndex preserves all system properties

During property filtering in `SyncIndex`, all system properties present in the object's frontmatter SHALL be preserved in the filtered property set, regardless of type schema definitions.

#### Scenario: Sync preserves description, created_at and updated_at

- **WHEN** an object has `description`, `created_at` and `updated_at` in its frontmatter
- **AND** the type schema does not define these properties
- **THEN** after sync, the indexed properties SHALL include `description`, `created_at` and `updated_at`

### Requirement: System properties are classified by mutability

Each system property SHALL be classified as either **user-authored** (mutable by user or template) or **auto-managed** (immutable, always reflects actual system values). The `IsImmutableSystemProperty(name)` function SHALL return `true` for auto-managed properties (`created_at`, `updated_at`) and `false` for user-authored properties (`name`, `description`, `tags`, `locked`) and non-system properties.

#### Scenario: Auto-managed properties are immutable

- **WHEN** `IsImmutableSystemProperty("created_at")` is called
- **THEN** it SHALL return `true`

#### Scenario: Auto-managed property updated_at is immutable

- **WHEN** `IsImmutableSystemProperty("updated_at")` is called
- **THEN** it SHALL return `true`

#### Scenario: User-authored properties are not immutable

- **WHEN** `IsImmutableSystemProperty("name")` is called
- **THEN** it SHALL return `false`

#### Scenario: User-authored property description is not immutable

- **WHEN** `IsImmutableSystemProperty("description")` is called
- **THEN** it SHALL return `false`

#### Scenario: User-authored property tags is not immutable

- **WHEN** `IsImmutableSystemProperty("tags")` is called
- **THEN** it SHALL return `false`

#### Scenario: User-authored property locked is not immutable

- **WHEN** `IsImmutableSystemProperty("locked")` is called
- **THEN** it SHALL return `false`

#### Scenario: Non-system properties are not immutable

- **WHEN** `IsImmutableSystemProperty("title")` is called
- **THEN** it SHALL return `false`

### Requirement: Frontmatter orders system properties first

`OrderedPropKeys` SHALL place system properties before schema-defined properties, in registry order (name, description, created_at, updated_at, tags, locked). Schema-defined properties follow in schema order. Extra properties are appended alphabetically.

#### Scenario: Full property ordering

- **WHEN** an object has properties `name`, `description`, `created_at`, `updated_at`, `tags`, `locked`, `title`, and `rating`
- **AND** the schema defines `title` then `rating`
- **THEN** `OrderedPropKeys` SHALL return `["name", "description", "created_at", "updated_at", "tags", "locked", "title", "rating"]`

#### Scenario: System properties absent

- **WHEN** an object has properties `name` and `title` (no description, timestamps, or tags)
- **AND** the schema defines `title`
- **THEN** `OrderedPropKeys` SHALL return `["name", "title"]`
