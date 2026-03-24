## Why

AI features (Describe, SuggestTags, ExploreSchema) currently only work via the `claude` CLI binary. Users who prefer local LLMs (Ollama, LM Studio, vLLM, LocalAI) or self-hosted models cannot use any AI functionality. Most local LLM servers implement the OpenAI-compatible `/v1/chat/completions` API, making it a practical standard to support.

## What Changes

- Introduce a multi-provider config structure (`ai.providers` map + `ai.default` selector) so users can define multiple AI backends and switch between them
- Add `OpenAICompatible` provider implementing the `Provider` interface via HTTP calls to `/v1/chat/completions` with `response_format: { type: "json_schema" }` for structured output
- Each provider config has an explicit `type` field (`cli` or `openai-compatible`) to determine which implementation to use
- Migrate old flat config (`ai.model`) to the new structure transparently — existing configs continue to work
- Update `Vault.initAI()` to resolve `ai.default` → provider config → instantiate the correct `Provider` implementation

## Capabilities

### New Capabilities

- `openai-compatible-provider`: OpenAI-compatible HTTP provider implementing the `Provider` interface, supporting `/v1/chat/completions` with structured JSON output, optional API key auth, connection error handling, and model validation
- `multi-provider-config`: Multi-provider AI configuration structure (`ai.providers` + `ai.default`), provider type routing (`cli` vs `openai-compatible`), and backward-compatible migration from flat config format

### Modified Capabilities

- `ai-provider`: Update AI availability detection in `Vault.initAI()` to resolve named providers from the new config structure instead of hardcoding `ClaudeCLI`
- `vault-config`: Extend `AIConfig` struct with `Default` and `Providers` fields; add config key registry entries for new provider settings

## Impact

- **`core/ai/`**: New `openai_compatible.go` file with `OpenAICompatible` struct + HTTP client logic
- **`core/config.go`**: Extended `AIConfig` with `Default string`, `Providers map[string]ProviderConfig`; new `ProviderConfig` struct; migration logic for old format; new config key registry entries
- **`core/vault.go`**: Updated `initAI()` to resolve provider by name and type
- **Dependencies**: `net/http` (stdlib only, no new external dependencies)
- **No TUI changes**: TUI already uses `vault.AIService()` — transparent
- **No CLI changes**: CLI commands already delegate to vault AI service
