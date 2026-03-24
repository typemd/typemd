## MODIFIED Requirements

### Requirement: AI availability detection on Vault open

The Vault SHALL detect AI availability during `Open()`. If `ai.enabled` is `true` in config, the Vault SHALL resolve the active provider from `ai.default` and `ai.providers`. For `cli` type providers, it SHALL check for the `claude` binary via `exec.LookPath`. For `openai-compatible` type providers, it SHALL instantiate an `OpenAICompatible` provider (no upfront connectivity check). The Vault SHALL expose an `AIService()` method that returns the service or nil.

#### Scenario: AI enabled with cli provider and claude binary found

- **WHEN** config has `ai.enabled: true` and the resolved provider has `type: cli` and `claude` is in PATH
- **THEN** `vault.AIService()` SHALL return a non-nil `*ai.AIService` backed by a `ClaudeCLI` provider

#### Scenario: AI enabled with cli provider but claude binary not found

- **WHEN** config has `ai.enabled: true` and the resolved provider has `type: cli` but `claude` is not in PATH
- **THEN** `vault.AIService()` SHALL return nil

#### Scenario: AI enabled with openai-compatible provider

- **WHEN** config has `ai.enabled: true` and the resolved provider has `type: openai-compatible` with a valid `base_url`
- **THEN** `vault.AIService()` SHALL return a non-nil `*ai.AIService` backed by an `OpenAICompatible` provider

#### Scenario: AI enabled with legacy flat config (no providers section)

- **WHEN** config has `ai.enabled: true` and `ai.model: some-model` but no `ai.providers`
- **THEN** the Vault SHALL auto-migrate to a `cli` provider and behave as if `type: cli` was configured
- **AND** `vault.AIService()` SHALL return a non-nil `*ai.AIService` if `claude` is in PATH

#### Scenario: AI not enabled

- **WHEN** config has `ai.enabled: false` or the key is absent
- **THEN** `vault.AIService()` SHALL return nil
- **AND** no provider resolution SHALL be performed
