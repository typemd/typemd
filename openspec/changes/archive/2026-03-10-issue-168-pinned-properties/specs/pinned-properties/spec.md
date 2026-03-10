## ADDED Requirements

### Requirement: Property supports optional pin field

The Property struct SHALL support an optional `pin` field that stores a positive integer value. When a property definition in a type schema YAML includes a `pin` field, it SHALL be parsed and stored. When the field is omitted, the pin SHALL default to zero (not pinned).

#### Scenario: Property with pin defined
- **WHEN** a type schema property definition contains `pin: 1`
- **THEN** the loaded Property SHALL have its Pin field set to 1

#### Scenario: Property without pin defined
- **WHEN** a type schema property definition does not contain a `pin` field
- **THEN** the loaded Property SHALL have its Pin field set to 0

### Requirement: Pin values must be positive integers

When a property has a pin value set, it SHALL be a positive integer (greater than zero). Schema validation SHALL reject zero or negative pin values.

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

### Requirement: Pinned properties displayed at top of body panel

In the TUI detail view, properties with a non-zero pin value SHALL be rendered at the top of the body panel, above the markdown body content. Pinned properties SHALL be sorted by pin value ascending (lower number first).

#### Scenario: Pinned property rendered with emoji
- **GIVEN** a type schema with property `status` having `emoji: 📋` and `pin: 1`
- **AND** an object with `status: reading`
- **WHEN** the TUI body panel is rendered
- **THEN** the body panel SHALL display `📋 status: reading` at the top

#### Scenario: Pinned property rendered without emoji
- **GIVEN** a type schema with property `rating` having no emoji and `pin: 2`
- **AND** an object with `rating: 5`
- **WHEN** the TUI body panel is rendered
- **THEN** the body panel SHALL display `rating: 5` at the top (after any pin: 1 properties)

#### Scenario: Separator between pinned properties and body
- **GIVEN** an object with at least one pinned property and non-empty body content
- **WHEN** the TUI body panel is rendered
- **THEN** a horizontal separator SHALL appear between the pinned properties and the body content

#### Scenario: No separator when no body content
- **GIVEN** an object with pinned properties but empty body content
- **WHEN** the TUI body panel is rendered
- **THEN** the pinned properties SHALL be displayed without a trailing separator

### Requirement: Pinned properties excluded from Properties panel

Properties with a non-zero pin value SHALL NOT appear in the Properties panel. Only unpinned properties (pin = 0) SHALL be displayed in the Properties panel.

#### Scenario: Pinned property absent from Properties panel
- **GIVEN** a type schema with property `status` having `pin: 1`
- **WHEN** the Properties panel is rendered
- **THEN** `status` SHALL NOT appear in the Properties panel

#### Scenario: Unpinned property remains in Properties panel
- **GIVEN** a type schema with property `title` having no pin value
- **WHEN** the Properties panel is rendered
- **THEN** `title` SHALL appear in the Properties panel as usual
