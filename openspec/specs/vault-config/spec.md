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

The `VaultConfig` struct SHALL organize settings under interface-layer keys: `cli`, `tui`, `ai`. Each layer SHALL have its own sub-struct. The `TUIConfig` sub-struct SHALL include:
- `debounce_ms` (int) — file watcher debounce interval
- `stats_type_layout` (string) — stats detail layout mode
- `toast` (ToastConfig) — toast notification settings with sub-keys:
  - `position` (string, default `"bottom-right"`) — toast display position (`"bottom-right"` or `"help-bar"`)
  - `duration_ms` (int, default `3000`) — auto-dismiss duration in milliseconds
  - `dismiss_key` (string, default `"esc"`) — key to manually dismiss toast
  - `show_warnings` (*bool, default `true`) — whether to show warning-level toasts
  - `show_success` (*bool, default `false`) — whether to show info/success-level toasts

The `AIConfig` sub-struct SHALL include:
- `enabled` (bool, default `false`) — opt-in toggle for AI features
- `model` (string, default `"claude-sonnet-4-6-20250627"`) — Claude model to use
- `prompts` (PromptsConfig) — customizable system prompts with sub-keys `describe`, `tag`, `explore`
- `explore` (ExploreConfig) — schema explore parameters with sub-keys `sample_count` (int, default 10), `body_truncate` (int, default 500)

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

#### Scenario: TUI toast config with custom duration

- **WHEN** config contains `tui:\n  toast:\n    duration_ms: 5000`
- **THEN** `config.TUI.Toast.DurationMs` SHALL be `5000`

#### Scenario: TUI toast config with position

- **WHEN** config contains `tui:\n  toast:\n    position: help-bar`
- **THEN** `config.TUI.Toast.Position` SHALL be `"help-bar"`

#### Scenario: TUI toast config with dismiss_key

- **WHEN** config contains `tui:\n  toast:\n    dismiss_key: q`
- **THEN** `config.TUI.Toast.DismissKey` SHALL be `"q"`

#### Scenario: TUI toast config with show_warnings disabled

- **WHEN** config contains `tui:\n  toast:\n    show_warnings: false`
- **THEN** `config.TUI.Toast.ShowWarnings` SHALL point to `false`

#### Scenario: TUI toast config with show_success enabled

- **WHEN** config contains `tui:\n  toast:\n    show_success: true`
- **THEN** `config.TUI.Toast.ShowSuccess` SHALL point to `true`

#### Scenario: TUI toast config defaults when absent

- **WHEN** config does not contain a `tui.toast` section
- **THEN** `config.TUI.Toast.Position` SHALL be `""` (empty; runtime default is `"bottom-right"`)
- **AND** `config.TUI.Toast.DurationMs` SHALL be `0` (runtime default is `3000`)
- **AND** `config.TUI.Toast.DismissKey` SHALL be `""` (runtime default is `"esc"`)
- **AND** `config.TUI.Toast.ShowWarnings` SHALL be `nil` (runtime default is `true`)
- **AND** `config.TUI.Toast.ShowSuccess` SHALL be `nil` (runtime default is `false`)

#### Scenario: AI config with enabled flag

- **WHEN** config contains `ai:\n  enabled: true`
- **THEN** `config.AI.Enabled` SHALL be `true`

#### Scenario: AI config with model override

- **WHEN** config contains `ai:\n  enabled: true\n  model: claude-haiku-4-5-20251001`
- **THEN** `config.AI.Model` SHALL be `"claude-haiku-4-5-20251001"`

#### Scenario: AI config with custom prompts

- **WHEN** config contains `ai:\n  prompts:\n    describe: "Custom describe prompt"`
- **THEN** `config.AI.Prompts.Describe` SHALL be `"Custom describe prompt"`

#### Scenario: AI config with explore settings

- **WHEN** config contains `ai:\n  explore:\n    sample_count: 20\n    body_truncate: 1000`
- **THEN** `config.AI.Explore.SampleCount` SHALL be `20`
- **AND** `config.AI.Explore.BodyTruncate` SHALL be `1000`

#### Scenario: AI config defaults when absent

- **WHEN** config does not contain an `ai` section
- **THEN** `config.AI.Enabled` SHALL be `false`
- **AND** `config.AI.Model` SHALL be `""` (empty; provider uses built-in default)

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
