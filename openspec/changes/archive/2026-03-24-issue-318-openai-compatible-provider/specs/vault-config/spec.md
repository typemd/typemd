## MODIFIED Requirements

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
- `model` (string) — legacy model field (used for backward-compatible migration)
- `default` (string) — name of the active provider from the `providers` map
- `providers` (map[string]ProviderConfig) — named AI provider configurations
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
- **AND** `config.AI.Default` SHALL be `""`
- **AND** `config.AI.Providers` SHALL be nil

#### Scenario: AI config with providers and default

- **WHEN** config contains `ai:\n  enabled: true\n  default: ollama\n  providers:\n    ollama:\n      type: openai-compatible\n      base_url: http://localhost:11434\n      model: qwen3-coder:30b`
- **THEN** `config.AI.Default` SHALL be `"ollama"`
- **AND** `config.AI.Providers["ollama"].Type` SHALL be `"openai-compatible"`
- **AND** `config.AI.Providers["ollama"].BaseURL` SHALL be `"http://localhost:11434"`
- **AND** `config.AI.Providers["ollama"].Model` SHALL be `"qwen3-coder:30b"`

#### Scenario: Unknown top-level keys are ignored

- **WHEN** config contains `unknown_key: value` alongside valid keys
- **THEN** the Vault SHALL load successfully, ignoring the unknown key
