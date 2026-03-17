## MODIFIED Requirements

### Requirement: Init always creates config.yaml with page as default type

`tmd init` SHALL always create `.typemd/config.yaml` with `cli.default_type` set to `page`. The built-in `page` type serves as the default for quick object creation, regardless of which starter types are selected.

#### Scenario: Init creates config with page default

- **WHEN** `tmd init` is run
- **THEN** `.typemd/config.yaml` SHALL be created with `cli:\n  default_type: page`

#### Scenario: Init with --no-starters still creates config

- **WHEN** `tmd init --no-starters` is run
- **THEN** `.typemd/config.yaml` SHALL be created with `cli:\n  default_type: page`
