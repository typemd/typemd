## ADDED Requirements

### Requirement: Load shared properties from per-property files

The system SHALL load shared property definitions by scanning the `properties/` directory for `*.yaml` files. Each file SHALL be parsed as a single property definition. The property name SHALL be derived from the filename (without `.yaml` extension).

#### Scenario: Load single per-property file

- **WHEN** `properties/due_date.yaml` exists with content `type: date` and `emoji: "📅"`
- **THEN** `LoadSharedProperties()` returns one property with name `due_date`, type `date`, and emoji `📅`

#### Scenario: Load multiple per-property files

- **WHEN** `properties/due_date.yaml` and `properties/priority.yaml` both exist
- **THEN** `LoadSharedProperties()` returns two properties with names `due_date` and `priority`

#### Scenario: Empty properties directory

- **WHEN** `properties/` directory exists but contains no `.yaml` files
- **THEN** `LoadSharedProperties()` returns nil (no properties)

#### Scenario: Properties directory does not exist

- **WHEN** `properties/` directory does not exist
- **THEN** `LoadSharedProperties()` returns nil (no error)

#### Scenario: Non-YAML files are ignored

- **WHEN** `properties/` contains `.DS_Store` and `README.md` alongside `rating.yaml`
- **THEN** `LoadSharedProperties()` returns only the `rating` property

#### Scenario: Name field in file is ignored

- **WHEN** `properties/due_date.yaml` contains a `name: something_else` field
- **THEN** the property name SHALL be `due_date` (derived from filename, not file content)

### Requirement: Per-property file format

Each per-property YAML file SHALL contain the property definition fields directly (no wrapper key). The `name` field is not required — it is derived from the filename.

#### Scenario: Minimal property file

- **WHEN** `properties/rating.yaml` contains `type: number`
- **THEN** the property loads successfully with name `rating` and type `number`

#### Scenario: Full property file with all fields

- **WHEN** `properties/priority.yaml` contains type, emoji, description, options, default, and pin fields
- **THEN** all fields are loaded correctly alongside the filename-derived name

### Requirement: Validate per-property files

Validation rules for shared properties SHALL apply to per-property files. Property names (derived from filenames) SHALL NOT conflict with reserved system property names. Duplicate filenames are impossible (filesystem enforces uniqueness).

#### Scenario: Reserved system property name as filename

- **WHEN** `properties/created_at.yaml` exists
- **THEN** validation SHALL report an error that `created_at` is a reserved system property

#### Scenario: Invalid property type in per-property file

- **WHEN** `properties/rating.yaml` contains `type: invalid_type`
- **THEN** validation SHALL report an error for invalid property type

#### Scenario: Select property without options

- **WHEN** `properties/priority.yaml` contains `type: select` without `options`
- **THEN** validation SHALL report an error that select type requires options

### Requirement: Use references resolve from per-property files

Type schema `use` references SHALL resolve against per-property files exactly as they did with the single-file format. The resolution logic is unchanged — only the loading source changes.

#### Scenario: Use reference resolves from per-property file

- **WHEN** `properties/due_date.yaml` defines a date property
- **AND** a type schema has a property with `use: due_date`
- **THEN** the property resolves to the shared definition with overrides applied

### Requirement: Migrate legacy properties.yaml to per-property files

On vault open, the system SHALL detect `properties/properties.yaml` and automatically migrate it to per-property files. After migration, the legacy file SHALL be removed.

#### Scenario: Successful migration of legacy file

- **WHEN** `properties/properties.yaml` exists with two property definitions (`due_date` and `priority`)
- **THEN** `properties/due_date.yaml` and `properties/priority.yaml` are created
- **AND** `properties/properties.yaml` is removed

#### Scenario: Legacy file with empty properties list

- **WHEN** `properties/properties.yaml` exists with `properties: []` or `properties:` (empty)
- **THEN** `properties/properties.yaml` is removed
- **AND** no per-property files are created

#### Scenario: Legacy file coexists with per-property files

- **WHEN** `properties/properties.yaml` exists
- **AND** `properties/due_date.yaml` also exists
- **THEN** vault open SHALL return an error indicating conflicting property formats

#### Scenario: Migration preserves all property fields

- **WHEN** `properties/properties.yaml` contains a property with type, emoji, description, options, default, pin, target, multiple, bidirectional, and inverse fields
- **THEN** the migrated per-property file SHALL contain all those fields (except `name`, which becomes the filename)

### Requirement: File watcher monitors properties directory

The `tmd validate --watch` command SHALL watch the `properties/` directory recursively for changes, rather than watching a single `properties/properties.yaml` file.

#### Scenario: Watch detects new per-property file

- **WHEN** `tmd validate --watch` is running
- **AND** a new `properties/rating.yaml` file is created
- **THEN** validation is re-triggered

### Requirement: SharedPropertiesDir replaces SharedPropertiesPath

`Vault.SharedPropertiesPath()` SHALL be replaced with `Vault.SharedPropertiesDir()` returning the `properties/` directory path. `Vault.Init()` SHALL ensure the `properties/` directory exists.

#### Scenario: SharedPropertiesDir returns directory path

- **WHEN** a vault is opened at `/tmp/test-vault`
- **THEN** `SharedPropertiesDir()` returns `/tmp/test-vault/properties`
