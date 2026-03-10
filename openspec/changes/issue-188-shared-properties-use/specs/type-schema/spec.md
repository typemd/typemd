## ADDED Requirements

### Requirement: Type schema supports use keyword for shared properties

The Property struct SHALL support an optional `use` field. When a property entry has `use: <name>`, it references a shared property from `.typemd/properties.yaml`. The `use` and `name` fields are mutually exclusive — a property entry SHALL have exactly one of them.

#### Scenario: Type schema with use entry

- **WHEN** a type schema contains `- use: due_date` in its properties array
- **THEN** the parsed Property SHALL have its Use field set to "due_date" and Name field empty

#### Scenario: Property entry with both use and name rejected

- **WHEN** a type schema contains a property entry with both `use: due_date` and `name: my_date`
- **THEN** schema validation SHALL return an error indicating `use` and `name` are mutually exclusive

### Requirement: Use entries only allow pin and emoji overrides

A `use` property entry SHALL only contain the fields `use`, `pin`, and `emoji`. Any other fields (type, options, default, target, etc.) SHALL be rejected by validation.

#### Scenario: Use with pin override accepted

- **WHEN** a type schema contains `- use: due_date` with `pin: 1`
- **THEN** schema validation SHALL accept without error

#### Scenario: Use with emoji override accepted

- **WHEN** a type schema contains `- use: due_date` with `emoji: 🗓️`
- **THEN** schema validation SHALL accept without error

#### Scenario: Use with pin and emoji overrides accepted

- **WHEN** a type schema contains `- use: due_date` with `pin: 1` and `emoji: 🗓️`
- **THEN** schema validation SHALL accept without error

#### Scenario: Use with type field rejected

- **WHEN** a type schema contains `- use: due_date` with `type: string`
- **THEN** schema validation SHALL return an error indicating only `pin` and `emoji` overrides are allowed on `use` entries

#### Scenario: Use with options field rejected

- **WHEN** a type schema contains `- use: priority` with `options: [{value: urgent}]`
- **THEN** schema validation SHALL return an error indicating only `pin` and `emoji` overrides are allowed on `use` entries

### Requirement: Use must reference existing shared property

A `use` entry SHALL reference a property name that exists in `.typemd/properties.yaml`. Referencing a non-existent shared property SHALL be rejected by validation.

#### Scenario: Use references existing shared property

- **WHEN** `.typemd/properties.yaml` defines `due_date` and a type schema contains `- use: due_date`
- **THEN** schema validation SHALL accept without error

#### Scenario: Use references non-existent shared property

- **WHEN** `.typemd/properties.yaml` does not define `due_date` and a type schema contains `- use: due_date`
- **THEN** schema validation SHALL return an error indicating shared property "due_date" not found

### Requirement: Local property name must not conflict with shared property name

A type schema SHALL NOT define a `name` property that has the same name as any shared property in `.typemd/properties.yaml`, regardless of whether the type uses that shared property.

#### Scenario: Local property conflicts with shared property name

- **WHEN** `.typemd/properties.yaml` defines `due_date` and a type schema defines `- name: due_date` with `type: string`
- **THEN** schema validation SHALL return an error indicating "due_date" conflicts with a shared property name

#### Scenario: Local property with unique name accepted

- **WHEN** `.typemd/properties.yaml` defines `due_date` and a type schema defines `- name: title` with `type: string`
- **THEN** schema validation SHALL accept without error

### Requirement: LoadType resolves use entries

`LoadType()` SHALL resolve all `use` entries in a type schema by replacing them with fully resolved Property objects. The resolved Property SHALL have all fields from the shared definition, with `pin` and `emoji` overridden if specified in the `use` entry. After resolution, the `Use` field SHALL be empty.

#### Scenario: Use entry resolved with no overrides

- **WHEN** a shared property `due_date` has `type: date` and `emoji: 📅`, and a type schema contains `- use: due_date`
- **THEN** `LoadType()` SHALL return a Property with name "due_date", type "date", emoji "📅", and Use field empty

#### Scenario: Use entry resolved with pin override

- **WHEN** a shared property `due_date` has `type: date` and `emoji: 📅`, and a type schema contains `- use: due_date` with `pin: 1`
- **THEN** `LoadType()` SHALL return a Property with name "due_date", type "date", emoji "📅", pin 1, and Use field empty

#### Scenario: Use entry resolved with emoji override

- **WHEN** a shared property `due_date` has `type: date` and `emoji: 📅`, and a type schema contains `- use: due_date` with `emoji: 🗓️`
- **THEN** `LoadType()` SHALL return a Property with name "due_date", type "date", emoji "🗓️", and Use field empty

#### Scenario: Mixed use and name properties resolved

- **WHEN** a type schema has `[{name: title, type: string}, {use: due_date}, {name: budget, type: number}]`
- **THEN** `LoadType()` SHALL return three fully resolved Properties in original order: title, due_date, budget

### Requirement: Resolved properties must have unique names

After resolving all `use` entries, the type schema SHALL NOT have duplicate property names. This includes duplicates between `use`-resolved properties and `name`-defined properties, and between multiple `use` entries.

#### Scenario: Duplicate use entries rejected

- **WHEN** a type schema contains both `- use: due_date` and `- use: due_date`
- **THEN** schema validation SHALL return an error indicating duplicate property name "due_date"
