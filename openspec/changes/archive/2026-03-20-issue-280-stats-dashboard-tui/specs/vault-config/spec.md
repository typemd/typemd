## MODIFIED Requirements

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
