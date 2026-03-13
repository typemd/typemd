### Requirement: Shared properties file format

The system SHALL support an optional `.typemd/properties.yaml` file that defines shared property definitions. The file SHALL use a `properties` array with the same Property format as type schemas. When the file does not exist, the system SHALL treat it as an empty set of shared properties.

#### Scenario: Load shared properties file

- **WHEN** `.typemd/properties.yaml` contains a `properties` array with entries for `due_date` (type: date) and `priority` (type: select)
- **THEN** `LoadSharedProperties()` SHALL return two Property objects with the correct names and types

#### Scenario: Shared properties file does not exist

- **WHEN** `.typemd/properties.yaml` does not exist in the vault
- **THEN** `LoadSharedProperties()` SHALL return an empty slice without error

#### Scenario: Shared properties file is empty

- **WHEN** `.typemd/properties.yaml` exists but contains no `properties` array
- **THEN** `LoadSharedProperties()` SHALL return an empty slice without error

### Requirement: Shared properties support all property fields

Shared property definitions SHALL support all fields available in type schema properties: `name`, `type`, `emoji`, `pin`, `options`, `target`, `default`, `multiple`, `bidirectional`, `inverse`.

#### Scenario: Shared property with options

- **WHEN** a shared property defines `type: select` with `options: [{value: high}, {value: low}]`
- **THEN** the loaded shared property SHALL include the options array

#### Scenario: Shared property with relation fields

- **WHEN** a shared property defines `type: relation` with `target: person` and `multiple: true`
- **THEN** the loaded shared property SHALL include target and multiple fields

### Requirement: Shared properties names must be unique

Within `.typemd/properties.yaml`, no two properties SHALL have the same name. Validation SHALL reject duplicate property names.

#### Scenario: Unique shared property names accepted

- **WHEN** `.typemd/properties.yaml` defines properties named `due_date` and `priority`
- **THEN** validation SHALL accept without error

#### Scenario: Duplicate shared property names rejected

- **WHEN** `.typemd/properties.yaml` defines two properties both named `due_date`
- **THEN** validation SHALL return an error indicating duplicate shared property name "due_date"

### Requirement: Shared properties validated like type properties

Shared property definitions SHALL be validated using the same rules as type schema properties: valid property types, options required for select/multi_select, target required for relation, etc.

#### Scenario: Invalid property type in shared properties

- **WHEN** `.typemd/properties.yaml` defines a property with `type: invalid`
- **THEN** validation SHALL return an error indicating unknown property type

#### Scenario: Select without options in shared properties

- **WHEN** `.typemd/properties.yaml` defines a property with `type: select` and no options
- **THEN** validation SHALL return an error indicating options are required for select type

### Requirement: Shared property name cannot be reserved

Shared properties SHALL NOT use the reserved name `name`. Validation SHALL reject shared properties named `name`.

#### Scenario: Shared property named "name" rejected

- **WHEN** `.typemd/properties.yaml` defines a property with `name: name`
- **THEN** validation SHALL return an error indicating "name" is a reserved system property
