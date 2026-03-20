## ADDED Requirements

### Requirement: Vault loads config from `.typemd/config.yaml`

The Vault SHALL load configuration from `.typemd/config.yaml` during `Open()`. The config file is optional — if the file does not exist or is empty, the Vault SHALL use zero-value defaults with no error.

#### Scenario: Config file exists with valid content

- **WHEN** `.typemd/config.yaml` contains `cli:\n  default_type: idea`
- **THEN** the Vault SHALL load the config with `CLI.DefaultType` set to `"idea"`

#### Scenario: Config file does not exist

- **WHEN** `.typemd/config.yaml` does not exist
- **THEN** the Vault SHALL use an empty config with all fields at zero values
- **AND** no error SHALL be returned

#### Scenario: Config file is empty

- **WHEN** `.typemd/config.yaml` exists but is empty
- **THEN** the Vault SHALL use an empty config with all fields at zero values
- **AND** no error SHALL be returned

#### Scenario: Config file has invalid YAML

- **WHEN** `.typemd/config.yaml` contains invalid YAML syntax
- **THEN** the Vault SHALL return an error during `Open()`

### Requirement: Config struct uses interface-layer namespacing

The `VaultConfig` struct SHALL organize settings under interface-layer keys: `cli`, `tui`. Each interface layer SHALL have its own sub-struct. The `TUIConfig` sub-struct SHALL include `debounce_ms` (int) and `stats_type_layout` (string, valid values: `"fullscreen"` or `"popup"`, default `"fullscreen"`).

#### Scenario: CLI config with default_type

- **WHEN** config contains `cli:\n  default_type: note`
- **THEN** `config.CLI.DefaultType` SHALL be `"note"`

#### Scenario: TUI config with debounce_ms

- **WHEN** config contains `tui:\n  debounce_ms: 500`
- **THEN** `config.TUI.DebounceMs` SHALL be `500`

#### Scenario: TUI config with stats_type_layout

- **WHEN** config contains `tui:\n  stats_type_layout: popup`
- **THEN** `config.TUI.StatsTypeLayout` SHALL be `"popup"`

#### Scenario: TUI config with default stats_type_layout

- **WHEN** config does not contain `tui.stats_type_layout`
- **THEN** the stats type detail layout SHALL default to `"fullscreen"`

#### Scenario: Unknown top-level keys are ignored

- **WHEN** config contains `unknown_key: value` alongside valid keys
- **THEN** the Vault SHALL load successfully, ignoring the unknown key

### Requirement: Vault exposes DefaultType accessor

The Vault SHALL provide a `DefaultType()` method that returns the configured `cli.default_type` value. If no default type is configured, it SHALL return an empty string.

#### Scenario: Default type is configured

- **WHEN** config has `cli.default_type: idea`
- **THEN** `vault.DefaultType()` SHALL return `"idea"`

#### Scenario: Default type is not configured

- **WHEN** config does not have `cli.default_type`
- **THEN** `vault.DefaultType()` SHALL return `""`

### Requirement: Init always creates config.yaml with page as default type

`tmd init` SHALL always create `.typemd/config.yaml` with `cli.default_type` set to `page`. The built-in `page` type serves as the default for quick object creation, regardless of which starter types are selected.

#### Scenario: Init creates config with page default

- **WHEN** `tmd init` is run
- **THEN** `.typemd/config.yaml` SHALL be created with `cli:\n  default_type: page`

#### Scenario: Init with --no-starters still creates config

- **WHEN** `tmd init --no-starters` is run
- **THEN** `.typemd/config.yaml` SHALL be created with `cli:\n  default_type: page`
