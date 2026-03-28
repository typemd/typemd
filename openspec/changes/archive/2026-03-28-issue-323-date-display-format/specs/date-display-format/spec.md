## ADDED Requirements

### Requirement: Vault config defines date display formats
The system SHALL support `date_format` and `datetime_format` keys in `.typemd/config.yaml` that control how date and datetime values are displayed.

Format tokens: `YYYY` (4-digit year), `MM` (2-digit month), `DD` (2-digit day), `HH` (24-hour), `mm` (minute), `ss` (second).

Default values: `date_format: "YYYY-MM-DD"`, `datetime_format: "YYYY-MM-DD HH:mm:ss"`.

#### Scenario: Default date format
- **WHEN** no `date_format` is configured
- **THEN** date values SHALL display as `YYYY-MM-DD` (e.g., `2026-03-28`)

#### Scenario: Default datetime format
- **WHEN** no `datetime_format` is configured
- **THEN** datetime values SHALL display as `YYYY-MM-DD HH:mm:ss` (e.g., `2026-03-28 14:30:00`)

#### Scenario: Custom date format
- **WHEN** `date_format` is set to `"DD/MM/YYYY"`
- **THEN** date values SHALL display as `28/03/2026`

#### Scenario: Custom datetime format
- **WHEN** `datetime_format` is set to `"YYYY/MM/DD HH:mm"`
- **THEN** datetime values SHALL display as `2026/03/28 14:30`

### Requirement: Config keys are accessible via tmd config
The system SHALL register `date_format` and `datetime_format` in the config key registry so they can be read and set via `tmd config get` and `tmd config set`.

#### Scenario: Get date_format via CLI
- **WHEN** user runs `tmd config get date_format`
- **THEN** the current date_format value is returned (or empty if default)

#### Scenario: Set datetime_format via CLI
- **WHEN** user runs `tmd config set datetime_format "DD.MM.YYYY HH:mm:ss"`
- **THEN** the datetime_format is persisted in `.typemd/config.yaml`

### Requirement: FormatValue applies configured formats to date properties
`DisplayProperty.FormatValue()` SHALL use the configured date format when rendering properties of type `date`.

#### Scenario: Date property with custom format
- **WHEN** a DisplayProperty has type `date` and value `2026-03-28`
- **AND** DateFormat is `"MM/DD/YYYY"`
- **THEN** FormatValue() SHALL return `"03/28/2026"`

#### Scenario: Date property with empty format uses default
- **WHEN** a DisplayProperty has type `date` and value `2026-03-28`
- **AND** DateFormat is empty
- **THEN** FormatValue() SHALL return `"2026-03-28"`

### Requirement: FormatValue applies configured formats to datetime properties
`DisplayProperty.FormatValue()` SHALL use the configured datetime format when rendering properties of type `datetime`.

#### Scenario: Datetime property with custom format
- **WHEN** a DisplayProperty has type `datetime` and value `2026-03-28T14:30:00`
- **AND** DatetimeFormat is `"DD/MM/YYYY HH:mm:ss"`
- **THEN** FormatValue() SHALL return `"28/03/2026 14:30:00"`

#### Scenario: Datetime property with empty format uses default
- **WHEN** a DisplayProperty has type `datetime` and value `2026-03-28T14:30:00`
- **AND** DatetimeFormat is empty
- **THEN** FormatValue() SHALL return `"2026-03-28 14:30:00"`

### Requirement: System properties use datetime format
System properties `created_at` and `updated_at` SHALL be formatted using the configured `datetime_format`.

#### Scenario: created_at with custom datetime format
- **WHEN** `datetime_format` is `"YYYY/MM/DD HH:mm"`
- **AND** an object's `created_at` is `2026-03-28T14:30:00+08:00`
- **THEN** the displayed value SHALL be `"2026/03/28 14:30"`

### Requirement: Format tokens reuse existing dateFormatReplacer
The system SHALL reuse the `dateFormatReplacer` from `name_template.go` for converting user-friendly format tokens to Go time format strings. Unrecognized tokens pass through as literal text.

#### Scenario: Unrecognized tokens pass through
- **WHEN** `date_format` is `"YYYY年MM月DD日"`
- **AND** a date value is `2026-03-28`
- **THEN** FormatValue() SHALL return `"2026年03月28日"`
