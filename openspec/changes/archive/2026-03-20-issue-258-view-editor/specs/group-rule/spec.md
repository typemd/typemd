## ADDED Requirements

### Requirement: GroupRule struct for multi-level grouping
The system SHALL define a `GroupRule` struct with a `Property` field (string). `ViewConfig.GroupBy` SHALL be typed as `[]GroupRule` instead of `string`.

#### Scenario: ViewConfig with single group rule
- **WHEN** a ViewConfig has `GroupBy: []GroupRule{{Property: "genre"}}`
- **THEN** the system SHALL group objects by the "genre" property value

#### Scenario: ViewConfig with multiple group rules
- **WHEN** a ViewConfig has `GroupBy: []GroupRule{{Property: "genre"}, {Property: "status"}}`
- **THEN** the system SHALL group objects first by "genre", then sub-group by "status" within each genre group

#### Scenario: ViewConfig with empty group rules
- **WHEN** a ViewConfig has `GroupBy: nil` or `GroupBy: []GroupRule{}`
- **THEN** the system SHALL display objects as a flat list without grouping

### Requirement: YAML serialization for GroupRule
The system SHALL serialize `GroupBy` as a YAML array of objects with `property` keys. An empty or nil `GroupBy` SHALL be omitted from YAML output.

#### Scenario: Serialize single group rule
- **WHEN** a ViewConfig with `GroupBy: []GroupRule{{Property: "genre"}}` is marshaled to YAML
- **THEN** the YAML output SHALL contain `group_by:\n- property: genre`

#### Scenario: Serialize multiple group rules
- **WHEN** a ViewConfig with `GroupBy: []GroupRule{{Property: "genre"}, {Property: "status"}}` is marshaled to YAML
- **THEN** the YAML output SHALL contain both entries under `group_by`

#### Scenario: Omit empty group rules
- **WHEN** a ViewConfig with `GroupBy: nil` is marshaled to YAML
- **THEN** the YAML output SHALL NOT contain a `group_by` key

### Requirement: Backward-compatible YAML deserialization
The system SHALL accept both legacy string format (`group_by: "genre"`) and new array format (`group_by: [{property: genre}]`) when loading view YAML files.

#### Scenario: Load legacy string format
- **WHEN** a view YAML file contains `group_by: "genre"`
- **THEN** the system SHALL parse it as `[]GroupRule{{Property: "genre"}}`

#### Scenario: Load new array format
- **WHEN** a view YAML file contains `group_by:\n- property: genre`
- **THEN** the system SHALL parse it as `[]GroupRule{{Property: "genre"}}`

#### Scenario: Load empty group_by
- **WHEN** a view YAML file does not contain a `group_by` key
- **THEN** the system SHALL set `GroupBy` to `nil`

### Requirement: Multi-level group display uses compound labels
When multiple `GroupRule` entries are defined, the TUI SHALL display compound group header labels by joining property values with " · " separator.

#### Scenario: Compound group label
- **WHEN** a view has `GroupBy: []GroupRule{{Property: "genre"}, {Property: "status"}}` and an object has genre "sci-fi" and status "reading"
- **THEN** the group header label SHALL be "sci-fi · reading"

#### Scenario: Missing sub-group value
- **WHEN** a view has multi-level grouping and an object is missing the second group property value
- **THEN** the group header SHALL use "(none)" for the missing value (e.g., "sci-fi · (none)")

### Requirement: DefaultView returns empty GroupBy
The `DefaultView()` function SHALL return a ViewConfig with `GroupBy: nil` (no grouping).

#### Scenario: Default view has no grouping
- **WHEN** `DefaultView("book")` is called and no saved default.yaml exists
- **THEN** the returned ViewConfig SHALL have `GroupBy` as `nil`
