### Requirement: ProviderConfig struct defines per-provider settings

The `ProviderConfig` struct in `core` SHALL include:
- `Type` (string, required) — provider type: `"cli"` or `"openai-compatible"`
- `BaseURL` (string) — HTTP endpoint base URL (required for `openai-compatible`)
- `Model` (string) — model identifier
- `APIKey` (string, optional) — authentication token

#### Scenario: CLI provider config

- **WHEN** a `ProviderConfig` has `type: cli` and `model: claude-sonnet-4-20250514`
- **THEN** `Type` SHALL be `"cli"` and `Model` SHALL be `"claude-sonnet-4-20250514"`
- **AND** `BaseURL` and `APIKey` SHALL be ignored

#### Scenario: OpenAI-compatible provider config

- **WHEN** a `ProviderConfig` has `type: openai-compatible`, `base_url: http://localhost:11434`, and `model: qwen3-coder:30b`
- **THEN** `Type` SHALL be `"openai-compatible"`, `BaseURL` SHALL be `"http://localhost:11434"`, and `Model` SHALL be `"qwen3-coder:30b"`

### Requirement: AIConfig supports multi-provider structure

The `AIConfig` struct SHALL include:
- `Default` (string, yaml `default`) — name of the active provider
- `Providers` (map[string]ProviderConfig, yaml `providers`) — named provider configurations

These fields coexist with existing fields (`Enabled`, `Model`, `Prompts`, `Explore`).

#### Scenario: Config with multiple providers and default

- **WHEN** config YAML contains `ai.default: ollama` and `ai.providers` with entries `claude` and `ollama`
- **THEN** `config.AI.Default` SHALL be `"ollama"`
- **AND** `config.AI.Providers` SHALL have two entries keyed `"claude"` and `"ollama"`

#### Scenario: Config with no providers section

- **WHEN** config YAML contains `ai.enabled: true` but no `ai.providers` or `ai.default`
- **THEN** `config.AI.Default` SHALL be `""`
- **AND** `config.AI.Providers` SHALL be nil

### Requirement: Backward-compatible config migration

When loading config, if `ai.providers` is nil or empty AND `ai.model` is non-empty, the system SHALL auto-create a provider entry to maintain backward compatibility. The migration SHALL happen in memory during config loading — the file SHALL NOT be rewritten.

#### Scenario: Old flat config with ai.model

- **WHEN** config contains `ai.enabled: true` and `ai.model: claude-haiku-4-5-20251001` but no `ai.providers`
- **THEN** `config.AI.Providers` SHALL contain a single entry named `"claude"` with `Type: "cli"` and `Model: "claude-haiku-4-5-20251001"`
- **AND** `config.AI.Default` SHALL be `"claude"`

#### Scenario: Old flat config without ai.model

- **WHEN** config contains `ai.enabled: true` but no `ai.model` and no `ai.providers`
- **THEN** `config.AI.Providers` SHALL contain a single entry named `"claude"` with `Type: "cli"` and empty `Model`
- **AND** `config.AI.Default` SHALL be `"claude"`

#### Scenario: New config takes precedence over old fields

- **WHEN** config contains both `ai.model: old-model` and `ai.providers` with entries
- **THEN** the `ai.providers` structure SHALL be used as-is
- **AND** the `ai.model` field SHALL be ignored for provider resolution

### Requirement: Provider config validation on vault open

The Vault SHALL validate provider configuration during `Open()` when AI is enabled.

#### Scenario: ai.default points to undefined provider

- **WHEN** `ai.default` is `"missing"` and `ai.providers` does not contain a `"missing"` key
- **THEN** the Vault SHALL log a warning and AI SHALL be unavailable (`AIService()` returns nil)

#### Scenario: openai-compatible provider missing base_url

- **WHEN** a provider has `type: openai-compatible` but no `base_url`
- **THEN** the Vault SHALL log a warning and AI SHALL be unavailable (`AIService()` returns nil)

#### Scenario: Unknown provider type

- **WHEN** a provider has `type: unknown-type`
- **THEN** the Vault SHALL log a warning and AI SHALL be unavailable (`AIService()` returns nil)
