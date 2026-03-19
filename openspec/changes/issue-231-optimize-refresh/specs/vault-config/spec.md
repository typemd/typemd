## MODIFIED Requirements

### Requirement: Config struct uses interface-layer namespacing

The `VaultConfig` struct SHALL organize settings under interface-layer keys: `cli`, `tui`. Each interface layer SHALL have its own sub-struct.

#### Scenario: CLI config with default_type

- **WHEN** config contains `cli:\n  default_type: note`
- **THEN** `config.CLI.DefaultType` SHALL be `"note"`

#### Scenario: TUI config with debounce_ms

- **WHEN** config contains `tui:\n  debounce_ms: 500`
- **THEN** `config.TUI.DebounceMs` SHALL be `500`

#### Scenario: Unknown top-level keys are ignored

- **WHEN** config contains `unknown_key: value` alongside valid keys
- **THEN** the Vault SHALL load successfully, ignoring the unknown key
