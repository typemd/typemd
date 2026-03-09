## ADDED Requirements

### Requirement: Supported property types

The system SHALL support the following property types: `string`, `number`, `date`, `datetime`, `url`, `checkbox`, `select`, `multi_select`, `relation`.

#### Scenario: Valid property type accepted
- **WHEN** a type schema defines a property with `type: date`
- **THEN** schema validation SHALL accept it without error

#### Scenario: Invalid property type rejected
- **WHEN** a type schema defines a property with `type: formula`
- **THEN** schema validation SHALL return an error indicating the type is not supported

#### Scenario: Legacy enum type rejected
- **WHEN** a type schema defines a property with `type: enum`
- **THEN** schema validation SHALL return an error indicating `enum` has been replaced by `select`

### Requirement: String type validation

A property with `type: string` SHALL accept any string value.

#### Scenario: String value accepted
- **WHEN** an object has a string property with value `"hello"`
- **THEN** object validation SHALL accept it

#### Scenario: Non-string value rejected
- **WHEN** an object has a string property with a non-string value (e.g., numeric `42`)
- **THEN** object validation SHALL return a type mismatch error

### Requirement: Number type validation

A property with `type: number` SHALL accept integer and floating-point values.

#### Scenario: Integer value accepted
- **WHEN** an object has a number property with value `42`
- **THEN** object validation SHALL accept it

#### Scenario: Float value accepted
- **WHEN** an object has a number property with value `4.5`
- **THEN** object validation SHALL accept it

#### Scenario: String value rejected for number
- **WHEN** an object has a number property with value `"not a number"`
- **THEN** object validation SHALL return a type mismatch error

### Requirement: Date type validation

A property with `type: date` SHALL accept values in `YYYY-MM-DD` format only.

#### Scenario: Valid date accepted
- **WHEN** an object has a date property with value `"2026-03-09"`
- **THEN** object validation SHALL accept it

#### Scenario: Datetime value rejected for date type
- **WHEN** an object has a date property with value `"2026-03-09T10:30:00"`
- **THEN** object validation SHALL return an error indicating the value must be a date without time

#### Scenario: Invalid date format rejected
- **WHEN** an object has a date property with value `"03/09/2026"`
- **THEN** object validation SHALL return a format error

#### Scenario: YAML auto-parsed time.Time accepted for date
- **WHEN** YAML parsing produces a `time.Time` value for a date property
- **THEN** object validation SHALL accept it and format it as `YYYY-MM-DD`

### Requirement: Datetime type validation

A property with `type: datetime` SHALL accept values in ISO 8601 format with time component.

#### Scenario: Datetime with seconds accepted
- **WHEN** an object has a datetime property with value `"2026-03-09T10:30:00"`
- **THEN** object validation SHALL accept it

#### Scenario: Datetime with timezone accepted
- **WHEN** an object has a datetime property with value `"2026-03-09T10:30:00+08:00"`
- **THEN** object validation SHALL accept it

#### Scenario: Date-only value rejected for datetime type
- **WHEN** an object has a datetime property with value `"2026-03-09"`
- **THEN** object validation SHALL return an error indicating the value must include a time component

#### Scenario: YAML auto-parsed time.Time accepted for datetime
- **WHEN** YAML parsing produces a `time.Time` value for a datetime property
- **THEN** object validation SHALL accept it and format it as ISO 8601

### Requirement: URL type validation

A property with `type: url` SHALL accept values starting with `http://` or `https://`.

#### Scenario: HTTPS URL accepted
- **WHEN** an object has a url property with value `"https://example.com"`
- **THEN** object validation SHALL accept it

#### Scenario: HTTP URL accepted
- **WHEN** an object has a url property with value `"http://example.com/page"`
- **THEN** object validation SHALL accept it

#### Scenario: Non-URL string rejected
- **WHEN** an object has a url property with value `"not-a-url"`
- **THEN** object validation SHALL return a validation error

#### Scenario: FTP URL rejected
- **WHEN** an object has a url property with value `"ftp://files.example.com"`
- **THEN** object validation SHALL return a validation error indicating only http/https are supported

### Requirement: Checkbox type validation

A property with `type: checkbox` SHALL accept boolean values only.

#### Scenario: Boolean true accepted
- **WHEN** an object has a checkbox property with value `true`
- **THEN** object validation SHALL accept it

#### Scenario: Boolean false accepted
- **WHEN** an object has a checkbox property with value `false`
- **THEN** object validation SHALL accept it

#### Scenario: String "true" rejected for checkbox
- **WHEN** an object has a checkbox property with string value `"true"`
- **THEN** object validation SHALL return a type mismatch error

### Requirement: Select type with options

A property with `type: select` SHALL define options as an array of objects with `value` (required) and `label` (optional) fields. The property value MUST match one of the defined `options[].value` entries.

#### Scenario: Valid select value accepted
- **WHEN** a select property has options `[{value: reading}, {value: done}]` and the object value is `"reading"`
- **THEN** object validation SHALL accept it

#### Scenario: Invalid select value rejected
- **WHEN** a select property has options `[{value: reading}, {value: done}]` and the object value is `"abandoned"`
- **THEN** object validation SHALL return an error indicating the value is not in the allowed options

#### Scenario: Select schema without options rejected
- **WHEN** a type schema defines a select property without an `options` field
- **THEN** schema validation SHALL return an error indicating options are required for select type

#### Scenario: Option with label
- **WHEN** a select property has options `[{value: in-progress, label: In Progress}]`
- **THEN** schema validation SHALL accept it and the label SHALL be available for display

### Requirement: Multi-select type with options

A property with `type: multi_select` SHALL define options as an array of objects with `value` and optional `label` fields. The property value MUST be a list, and each element MUST match one of the defined `options[].value` entries.

#### Scenario: Valid multi_select values accepted
- **WHEN** a multi_select property has options `[{value: fiction}, {value: sci-fi}, {value: classic}]` and the object value is `["fiction", "sci-fi"]`
- **THEN** object validation SHALL accept it

#### Scenario: Invalid multi_select value rejected
- **WHEN** a multi_select property has options `[{value: fiction}, {value: sci-fi}]` and the object value is `["fiction", "romance"]`
- **THEN** object validation SHALL return an error indicating `"romance"` is not in the allowed options

#### Scenario: Single string value coerced to list
- **WHEN** a multi_select property receives a single string value `"fiction"` instead of a list
- **THEN** object validation SHALL accept it, treating it as `["fiction"]`

#### Scenario: Multi_select schema without options rejected
- **WHEN** a type schema defines a multi_select property without an `options` field
- **THEN** schema validation SHALL return an error indicating options are required

### Requirement: Enum to select migration

The `tmd migrate` command SHALL support migrating type schemas from `type: enum` with `values` to `type: select` with `options`.

#### Scenario: Migrate enum to select
- **WHEN** a type schema has `type: enum` with `values: [to-read, reading, done]`
- **AND** user runs `tmd migrate`
- **THEN** the schema SHALL be updated to `type: select` with `options: [{value: to-read}, {value: reading}, {value: done}]`

#### Scenario: Migrate dry-run shows changes
- **WHEN** user runs `tmd migrate --dry-run` with enum schemas present
- **THEN** the output SHALL show the planned enum → select conversions without modifying files
