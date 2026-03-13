## Purpose

Type schemas define the structure and metadata for object types in typemd. This specification defines how type schemas are loaded, validated, and applied.

## Requirements

### Requirement: Type schema supports optional emoji field

The TypeSchema struct SHALL support an optional `emoji` field that stores a string value. When a type schema YAML file includes an `emoji` field, it SHALL be parsed and stored. When the field is omitted, the emoji SHALL default to an empty string.

#### Scenario: Type schema with emoji defined

- **WHEN** a type schema YAML file contains `emoji: 📚`
- **THEN** the loaded TypeSchema SHALL have its Emoji field set to "📚"

#### Scenario: Type schema without emoji defined

- **WHEN** a type schema YAML file does not contain an `emoji` field
- **THEN** the loaded TypeSchema SHALL have its Emoji field set to an empty string

### Requirement: Custom type emoji overrides built-in default

When a custom type schema defines its own emoji, it SHALL override the built-in default emoji for that type. Since only `tag` remains as a built-in type, this override behavior only applies to the `tag` type.

#### Scenario: Custom tag type with different emoji

- **WHEN** a custom `tag.yaml` defines `emoji: 🔖`
- **THEN** the loaded tag type SHALL have emoji "🔖" instead of the built-in "🏷️"

### Requirement: CLI type show displays emoji

The `tmd type show` command SHALL display the emoji alongside the type name when an emoji is defined.

#### Scenario: Show type with emoji

- **WHEN** user runs `tmd type show book` and the book type has emoji "📚"
- **THEN** the output SHALL include the emoji in the type display

#### Scenario: Show type without emoji

- **WHEN** user runs `tmd type show` for a type with no emoji
- **THEN** the output SHALL display the type without any emoji prefix

### Requirement: CLI type list displays emoji

The `tmd type list` command SHALL display emoji alongside type names when available.

#### Scenario: List types with emoji

- **WHEN** user runs `tmd type list` and types have emoji defined
- **THEN** each type with an emoji SHALL show the emoji alongside its name

### Requirement: Property struct supports options field

The Property struct SHALL support an `options` field containing an array of Option objects. Each Option SHALL have a required `value` string field and an optional `label` string field. When `label` is omitted, it SHALL default to the `value`. The Property struct SHALL also support an optional `emoji` string field for compact visual identification.

#### Scenario: Property with options defined
- **WHEN** a type schema property has `options: [{value: reading, label: Reading}]`
- **THEN** the loaded Property SHALL have Options with one entry where Value is "reading" and Label is "Reading"

#### Scenario: Property with options without label
- **WHEN** a type schema property has `options: [{value: reading}]`
- **THEN** the loaded Property SHALL have Options with one entry where Value is "reading" and Label defaults to "reading"

#### Scenario: Property with emoji and options
- **WHEN** a type schema property has `emoji: 📊` and `options: [{value: active}]`
- **THEN** the loaded Property SHALL have Emoji "📊" and Options with one entry where Value is "active"

### Requirement: Property type allowlist expanded

The schema validation SHALL accept the following property types: `string`, `number`, `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`, `relation`. The type `enum` SHALL be rejected with a message directing users to use `select` instead. The property name `name` SHALL be allowed in the `properties` array ONLY when it contains exclusively the `template` field. All other system property names SHALL be rejected as reserved.

#### Scenario: New property types accepted
- **WHEN** a type schema defines properties with types `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`
- **THEN** schema validation SHALL accept all of them

#### Scenario: Enum type rejected with guidance
- **WHEN** a type schema defines a property with `type: enum`
- **THEN** schema validation SHALL return an error message indicating to use `select` instead

#### Scenario: Property named "name" with template accepted
- **WHEN** a type schema defines a property with `name: name` and `template: "{{ date:YYYY-MM-DD }}"`
- **THEN** schema validation SHALL accept it without error

#### Scenario: Property named "name" without template rejected
- **WHEN** a type schema defines a property with `name: name` and `type: string`
- **THEN** schema validation SHALL return an error indicating only `template` is allowed on the `name` system property entry

#### Scenario: Property named "description" rejected
- **WHEN** a type schema defines a property with `name: description`
- **THEN** schema validation SHALL return an error indicating that "description" is a reserved system property

### Requirement: Property supports optional pin field

The Property struct SHALL support an optional `pin` field that stores a positive integer value. When a property definition in a type schema YAML includes a `pin` field, it SHALL be parsed and stored. When the field is omitted, the pin SHALL default to zero (not pinned). Pinned properties are displayed prominently at the top of the TUI body panel rather than in the Properties panel.

#### Scenario: Property with pin defined
- **WHEN** a type schema property definition contains `pin: 1`
- **THEN** the loaded Property SHALL have its Pin field set to 1

#### Scenario: Property without pin defined
- **WHEN** a type schema property definition does not contain a `pin` field
- **THEN** the loaded Property SHALL have its Pin field set to 0

### Requirement: Pin values must be positive integers

When a property has a pin value set, it SHALL be a positive integer (greater than zero). Schema validation SHALL reject negative pin values.

#### Scenario: Positive pin value accepted
- **WHEN** a type schema property has `pin: 3`
- **THEN** schema validation SHALL accept it without error

#### Scenario: Negative pin value rejected
- **WHEN** a type schema property has `pin: -1`
- **THEN** schema validation SHALL return an error indicating invalid pin value

### Requirement: Pin values unique within type scope

Within a single type schema, no two properties SHALL have the same non-zero pin value. Schema validation SHALL reject duplicate pin values.

#### Scenario: Unique pin values accepted
- **WHEN** a type schema has properties with pin values 1 and 2
- **THEN** schema validation SHALL accept it without error

#### Scenario: Duplicate pin values rejected
- **WHEN** a type schema has two properties both with `pin: 1`
- **THEN** schema validation SHALL return an error indicating duplicate pin value 1

#### Scenario: Unpinned properties do not conflict
- **WHEN** a type schema has three properties where two have no pin and one has `pin: 1`
- **THEN** schema validation SHALL accept it without error

### Requirement: Type schema supports optional unique field

The TypeSchema struct SHALL support an optional `unique` field that stores a boolean value. When a type schema YAML file includes `unique: true`, it SHALL be parsed and stored. When the field is omitted, unique SHALL default to `false`.

#### Scenario: Type schema with unique true

- **WHEN** a type schema YAML file contains `unique: true`
- **THEN** the loaded TypeSchema SHALL have its Unique field set to true

#### Scenario: Type schema with unique false

- **WHEN** a type schema YAML file contains `unique: false`
- **THEN** the loaded TypeSchema SHALL have its Unique field set to false

#### Scenario: Type schema without unique field

- **WHEN** a type schema YAML file does not contain a `unique` field
- **THEN** the loaded TypeSchema SHALL have its Unique field set to false

### Requirement: Custom type schema can override built-in unique setting

When a custom type schema overrides a built-in type (e.g., `tag.yaml`), the `unique` field from the custom schema SHALL take effect. If the custom schema omits `unique`, it SHALL default to `false`, overriding the built-in default.

#### Scenario: Custom tag schema with unique true

- **WHEN** `.typemd/types/tag.yaml` exists with `unique: true`
- **THEN** the loaded tag type SHALL have Unique true

#### Scenario: Custom tag schema without unique field

- **WHEN** `.typemd/types/tag.yaml` exists without a `unique` field
- **THEN** the loaded tag type SHALL have Unique false (overriding built-in default)

### Requirement: Only tag is a built-in type

The `defaultTypes` map SHALL contain only the `tag` type. All other types MUST be defined via `.typemd/types/*.yaml` files. The built-in `tag` type SHALL have `Unique: true` in its default schema.

#### Scenario: Loading an undefined type returns error

- **WHEN** no `.typemd/types/book.yaml` exists and no built-in `book` type is defined
- **AND** `LoadType("book")` is called
- **THEN** it SHALL return an error containing "unknown type: book"

#### Scenario: Tag type loads without custom schema

- **WHEN** no `.typemd/types/tag.yaml` exists
- **AND** `LoadType("tag")` is called
- **THEN** it SHALL return the built-in tag type schema with emoji "🏷️" and Unique true

#### Scenario: User-defined type loads from YAML

- **WHEN** `.typemd/types/book.yaml` exists with a valid schema
- **AND** `LoadType("book")` is called
- **THEN** it SHALL return the schema from the YAML file

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
