## ADDED Requirements

### Requirement: Config key registry with struct-aware validation

The system SHALL maintain a registry of known config keys that maps dot-notation strings to VaultConfig struct fields. Unknown keys SHALL be rejected with an error listing known keys.

#### Scenario: Known key is accepted

- **WHEN** `SetConfigValue("cli.default_type", "idea")` is called
- **THEN** the operation SHALL succeed
- **AND** `GetConfigValue("cli.default_type")` SHALL return `("idea", nil)`

#### Scenario: Unknown key is rejected

- **WHEN** `SetConfigValue("foo.bar", "baz")` is called
- **THEN** the operation SHALL return an error
- **AND** the error message SHALL contain `unknown config key`
- **AND** the error message SHALL list known keys

#### Scenario: Get unknown key is rejected

- **WHEN** `GetConfigValue("nonexistent")` is called
- **THEN** the result SHALL be `("", error)`
- **AND** the error message SHALL contain `unknown config key`

### Requirement: SetConfigValue creates config file if missing

When `SetConfigValue` is called and no `.typemd/config.yaml` exists, the system SHALL create the file with the specified key-value pair.

#### Scenario: Set on vault without config file

- **GIVEN** a vault with no `.typemd/config.yaml`
- **WHEN** `SetConfigValue("cli.default_type", "note")` is called
- **THEN** `.typemd/config.yaml` SHALL be created
- **AND** the file SHALL contain `cli:\n  default_type: note`

#### Scenario: Set on vault with existing config

- **GIVEN** a vault with `.typemd/config.yaml` containing `cli:\n  default_type: idea`
- **WHEN** `SetConfigValue("cli.default_type", "note")` is called
- **THEN** the config SHALL be updated to `cli:\n  default_type: note`

### Requirement: GetConfigValue returns empty for unset keys

When a known key is not set in the config, `GetConfigValue` SHALL return an empty string with `nil` error (key is known but unset).

#### Scenario: Get on unset known key

- **GIVEN** a vault with empty config
- **WHEN** `GetConfigValue("cli.default_type")` is called
- **THEN** the result SHALL be `("", nil)`

### Requirement: ConfigKeys returns all known keys

`ConfigKeys()` SHALL return a sorted list of all known config key strings.

#### Scenario: List known keys

- **WHEN** `ConfigKeys()` is called
- **THEN** the result SHALL contain `["cli.default_type"]`
- **AND** the list SHALL be sorted alphabetically

### Requirement: `tmd config set` CLI command

`tmd config set <key> <value>` SHALL set the specified config value and write to `.typemd/config.yaml`.

#### Scenario: Set a valid key

- **WHEN** `tmd config set cli.default_type idea` is run
- **THEN** `.typemd/config.yaml` SHALL contain `cli:\n  default_type: idea`
- **AND** exit code SHALL be 0

#### Scenario: Set with unknown key

- **WHEN** `tmd config set unknown.key value` is run
- **THEN** an error SHALL be printed to stderr
- **AND** the error SHALL list known keys
- **AND** exit code SHALL be non-zero

### Requirement: `tmd config get` CLI command

`tmd config get <key>` SHALL print the value of the specified config key to stdout.

#### Scenario: Get a set key

- **GIVEN** config has `cli.default_type: idea`
- **WHEN** `tmd config get cli.default_type` is run
- **THEN** stdout SHALL contain `idea`
- **AND** exit code SHALL be 0

#### Scenario: Get an unset key

- **GIVEN** config has no `cli.default_type`
- **WHEN** `tmd config get cli.default_type` is run
- **THEN** stdout SHALL be empty
- **AND** exit code SHALL be 0

#### Scenario: Get with unknown key

- **WHEN** `tmd config get unknown.key` is run
- **THEN** an error SHALL be printed to stderr
- **AND** exit code SHALL be non-zero

### Requirement: `tmd config list` CLI command

`tmd config list` SHALL print all set (non-empty) config values in `key: value` format.

#### Scenario: List with set values

- **GIVEN** config has `cli.default_type: idea`
- **WHEN** `tmd config list` is run
- **THEN** stdout SHALL contain `cli.default_type: idea`
- **AND** exit code SHALL be 0

#### Scenario: List with no config file

- **GIVEN** no `.typemd/config.yaml` exists
- **WHEN** `tmd config list` is run
- **THEN** stdout SHALL be empty
- **AND** exit code SHALL be 0

#### Scenario: List with empty config

- **GIVEN** `.typemd/config.yaml` exists but is empty
- **WHEN** `tmd config list` is run
- **THEN** stdout SHALL be empty
- **AND** exit code SHALL be 0
