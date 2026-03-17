# type-color Specification

## Purpose
TBD - created by archiving change issue-228-type-schema-metadata. Update Purpose after archive.
## Requirements
### Requirement: Type schema supports optional color field

The TypeSchema struct SHALL support an optional `color` field that stores a string value. When a type schema YAML file includes a `color` field, it SHALL be parsed and stored. When the field is omitted, the color SHALL default to an empty string.

#### Scenario: Type schema with preset color defined

- **WHEN** a type schema YAML file contains `color: green`
- **THEN** the loaded TypeSchema SHALL have its Color field set to "green"

#### Scenario: Type schema with hex color defined

- **WHEN** a type schema YAML file contains `color: "#FF5733"`
- **THEN** the loaded TypeSchema SHALL have its Color field set to "#FF5733"

#### Scenario: Type schema without color defined

- **WHEN** a type schema YAML file does not contain a `color` field
- **THEN** the loaded TypeSchema SHALL have its Color field set to an empty string

### Requirement: Color validates against preset names

Schema validation SHALL accept the following preset color names: `red`, `blue`, `green`, `yellow`, `purple`, `orange`, `pink`, `cyan`, `gray`, `brown`. Preset names SHALL be case-sensitive (lowercase only).

#### Scenario: Valid preset color accepted

- **WHEN** a type schema has `color: blue`
- **THEN** schema validation SHALL accept without error

#### Scenario: All preset colors accepted

- **WHEN** type schemas are validated with each of `red`, `blue`, `green`, `yellow`, `purple`, `orange`, `pink`, `cyan`, `gray`, `brown`
- **THEN** schema validation SHALL accept all without error

#### Scenario: Unknown preset color rejected

- **WHEN** a type schema has `color: magenta`
- **THEN** schema validation SHALL return an error indicating invalid color value

#### Scenario: Uppercase preset color rejected

- **WHEN** a type schema has `color: Red`
- **THEN** schema validation SHALL return an error indicating invalid color value

### Requirement: Color validates hex format

Schema validation SHALL accept hex color values in `#RRGGBB` (6-digit) or `#RGB` (3-digit) format. Hex digits SHALL be case-insensitive.

#### Scenario: 6-digit hex accepted

- **WHEN** a type schema has `color: "#FF5733"`
- **THEN** schema validation SHALL accept without error

#### Scenario: 3-digit hex accepted

- **WHEN** a type schema has `color: "#F53"`
- **THEN** schema validation SHALL accept without error

#### Scenario: Lowercase hex accepted

- **WHEN** a type schema has `color: "#ff5733"`
- **THEN** schema validation SHALL accept without error

#### Scenario: Mixed case hex accepted

- **WHEN** a type schema has `color: "#Ff5733"`
- **THEN** schema validation SHALL accept without error

#### Scenario: Invalid hex length rejected

- **WHEN** a type schema has `color: "#FF57"`
- **THEN** schema validation SHALL return an error indicating invalid color value

#### Scenario: Hex without hash rejected

- **WHEN** a type schema has `color: FF5733`
- **THEN** schema validation SHALL return an error indicating invalid color value

#### Scenario: Invalid hex characters rejected

- **WHEN** a type schema has `color: "#GGGGGG"`
- **THEN** schema validation SHALL return an error indicating invalid color value

### Requirement: Color included in YAML serialization

When a TypeSchema with a non-empty color is serialized to YAML, the output SHALL include the `color` field. When color is empty, it SHALL be omitted from output.

#### Scenario: Non-empty color serialized

- **WHEN** a TypeSchema with Color "green" is serialized to YAML
- **THEN** the output SHALL contain `color: green`

#### Scenario: Empty color omitted

- **WHEN** a TypeSchema with empty Color is serialized to YAML
- **THEN** the output SHALL NOT contain a `color:` line

### Requirement: Custom type color overrides built-in default

When a custom type schema defines its own color, it SHALL override the built-in default (which has no color). This applies to any built-in type when overridden by a custom schema.

#### Scenario: Custom tag type with color

- **WHEN** a custom `tag.yaml` defines `color: purple`
- **THEN** the loaded tag type SHALL have color "purple"

### Requirement: Color presets are queryable

A `ValidColorPresets()` function SHALL return the list of valid preset color names. This enables UI components to offer a color picker with preset options.

#### Scenario: ValidColorPresets returns all presets

- **WHEN** `ValidColorPresets()` is called
- **THEN** it SHALL return a slice containing exactly: `red`, `blue`, `green`, `yellow`, `purple`, `orange`, `pink`, `cyan`, `gray`, `brown`

