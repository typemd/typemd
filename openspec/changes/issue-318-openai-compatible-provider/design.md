## Context

typemd's AI features (Describe, SuggestTags, ExploreSchema) are built on a `Provider` interface (`core/ai/provider.go`), with a single implementation: `ClaudeCLI` that invokes the `claude` CLI binary as a subprocess. The `Vault.initAI()` method hardcodes this provider — it checks for the `claude` binary and creates a `ClaudeCLI` instance.

The AI architecture is already provider-agnostic at the service layer (`AIService` uses `Provider` interface), but the config and initialization are tightly coupled to Claude CLI. Extending this to support OpenAI-compatible HTTP APIs requires changes at two layers: config structure and provider instantiation.

Most local LLM servers (Ollama, LM Studio, vLLM, LocalAI) expose an OpenAI-compatible `/v1/chat/completions` endpoint, making this a single implementation that covers many backends.

## Goals / Non-Goals

**Goals:**

- Support multiple named AI providers in config with easy switching via `ai.default`
- Implement OpenAI-compatible HTTP provider using `/v1/chat/completions` with structured JSON output
- Maintain full backward compatibility with existing `ai.enabled` / `ai.model` flat config
- Use only stdlib `net/http` — no external HTTP client dependencies
- Clear error messages for common failure modes (server not running, model not found, structured output not supported)

**Non-Goals:**

- Non-OpenAI-compatible APIs (custom protocols like Ollama's native `/api/generate`)
- Streaming responses (non-streaming is simpler and sufficient for current use cases)
- API key management UI or secure credential storage
- TUI provider selection (future enhancement)
- Provider health checks or auto-fallback between providers

## Decisions

### 1. Multi-provider config with `ai.providers` map + `ai.default` selector

**Choice:** Named provider map with a `default` selector, rather than a single flat provider config or ordered list.

**Rationale:** Users may want multiple providers configured (e.g., Ollama for local, Claude for cloud) and switch between them without editing provider details. A map keyed by user-chosen names is the most ergonomic.

**Alternative considered:** Ordered list with index-based selection — less readable and harder to reference in config.

```yaml
ai:
  enabled: true
  default: ollama
  providers:
    claude:
      type: cli
      model: claude-sonnet-4-20250514
    ollama:
      type: openai-compatible
      base_url: http://localhost:11434
      model: qwen3-coder:30b
```

### 2. Explicit `type` field for provider routing

**Choice:** Each provider config has a `type: cli | openai-compatible` field that determines which `Provider` implementation to instantiate.

**Rationale:** Makes routing explicit and extensible. Future provider types (e.g., `bedrock`, `vertex`) can be added without ambiguity.

**Alternative considered:** Auto-detect from fields present (e.g., presence of `base_url` means HTTP) — fragile and error-prone.

### 3. `ProviderConfig` as a flat struct with type-specific fields

**Choice:** Single `ProviderConfig` struct with all fields (`Type`, `BaseURL`, `Model`, `APIKey`). Fields unused by a given type are simply ignored.

**Rationale:** Simplest approach. A discriminated union (interface-per-type) would require custom YAML unmarshaling and adds complexity for only two types. If we add more provider types with very different configs, we can refactor then.

### 4. Structured output via `response_format: { type: "json_schema" }`

**Choice:** Use the OpenAI `response_format` field with `type: "json_schema"` to request structured JSON output.

**Rationale:** This is the standard approach for structured output in the OpenAI API. Ollama, vLLM, and LM Studio all support this. The existing `CompletionRequest.JSONSchema` field maps directly to this.

**Fallback:** If a model doesn't support structured output, the API will return an error. We surface this as a clear error message rather than attempting text-based JSON extraction.

### 5. Backward-compatible config migration in `loadVaultConfig`

**Choice:** When loading config, if `ai.model` is set but `ai.providers` is empty, auto-create a `claude` provider entry from the flat fields. This migration happens in memory only — the file is not rewritten.

**Rationale:** Existing users should not need to update their config. The migration is deterministic and lossless.

### 6. HTTP client with configurable timeout

**Choice:** Use `net/http` with a 60-second default timeout on the `OpenAICompatible` struct. No external HTTP client library.

**Rationale:** Local LLM inference can be slow (especially first load). 60s covers most cases. Users running very large models can adjust if we expose timeout config later.

## Risks / Trade-offs

- **[Risk] Structured output quality varies by model** → Mitigation: We rely on the API's `response_format` enforcement. If the model produces invalid JSON despite the schema, the existing JSON parse error handling in `AIService` will surface a clear error.

- **[Risk] Connection refused when Ollama is not running** → Mitigation: The `OpenAICompatible.Complete()` method wraps connection errors with a hint: "connection refused — is the server running at <base_url>?"

- **[Risk] Config migration edge cases** → Mitigation: Migration only triggers when `ai.providers` is empty AND `ai.model` is set. If both old and new fields are present, new structure wins. Unit tests cover all combinations.

- **[Trade-off] No streaming** → Simpler implementation, but users won't see incremental output. Acceptable since AI operations already show a loading spinner and results are short (a description, a few tags, a few suggestions).

- **[Trade-off] Flat ProviderConfig struct** → Unused fields per type (e.g., `base_url` unused for `cli`). Acceptable for two types; refactor if we reach 4+ types.
