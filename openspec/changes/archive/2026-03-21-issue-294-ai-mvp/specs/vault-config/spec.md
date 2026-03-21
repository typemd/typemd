## MODIFIED Requirements

### Requirement: Config struct uses interface-layer namespacing

The `VaultConfig` struct SHALL organize settings under interface-layer keys: `cli`, `tui`, `ai`. Each layer SHALL have its own sub-struct. The `AIConfig` sub-struct SHALL include:
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
