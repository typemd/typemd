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

The `VaultConfig` struct SHALL organize settings under interface-layer keys: `cli`, `tui` (future). Each interface layer SHALL have its own sub-struct.

#### Scenario: CLI config with default_type

- **WHEN** config contains `cli:\n  default_type: note`
- **THEN** `config.CLI.DefaultType` SHALL be `"note"`

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

### Requirement: Init creates config.yaml when starter types are selected

When `tmd init` writes starter types and the user selects the `idea` or `note` type, it SHALL also create `.typemd/config.yaml` with `cli.default_type` set to `idea` (if selected) or `note` (if selected without idea). If neither is selected, no config file SHALL be created.

#### Scenario: Init with idea starter selected

- **WHEN** `tmd init` is run and user selects the `idea` starter type
- **THEN** `.typemd/config.yaml` SHALL be created with `cli:\n  default_type: idea`

#### Scenario: Init with note starter selected but not idea

- **WHEN** `tmd init` is run and user selects the `note` starter type but not `idea`
- **THEN** `.typemd/config.yaml` SHALL be created with `cli:\n  default_type: note`

#### Scenario: Init with no quick-suitable starter selected

- **WHEN** `tmd init` is run and user selects only the `book` starter type
- **THEN** no `.typemd/config.yaml` SHALL be created
