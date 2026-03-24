## 1. Config: Multi-Provider Structure

- [x] 1.1 Write BDD scenarios for multi-provider config parsing (providers map, default field, ProviderConfig struct fields)
- [x] 1.2 Add ProviderConfig struct and extend AIConfig with Default and Providers fields in config.go
- [x] 1.3 Add config key registry entries for ai.default and ai.providers (make BDD scenarios pass)
- [x] 1.4 Add unit tests for ProviderConfig YAML round-trip edge cases

## 2. Config: Backward-Compatible Migration

- [x] 2.1 Write BDD scenarios for config migration (flat ai.model → providers map, precedence rules)
- [x] 2.2 Implement migrateAIConfig in config.go (auto-create cli provider from flat fields)
- [x] 2.3 Add unit tests for migration edge cases (both old and new present, empty model, no ai section)

## 3. Core: OpenAICompatible Provider

- [x] 3.1 Write unit tests for OpenAICompatible.Complete (structured output, no schema, model override, API key header, error cases)
- [x] 3.2 Implement OpenAICompatible struct in core/ai/openai_compatible.go (HTTP POST to /v1/chat/completions, response_format, auth header)
- [x] 3.3 Add unit tests for connection error wrapping (connection refused message, non-200 status, invalid JSON, context cancellation)

## 4. Core: Provider Resolution in Vault

- [x] 4.1 Write BDD scenarios for provider resolution (cli type, openai-compatible type, legacy flat config, validation errors)
- [x] 4.2 Update initAI in vault.go to resolve provider from ai.default + ai.providers, instantiate correct Provider by type
- [x] 4.3 Add unit tests for initAI edge cases (missing base_url, unknown type, default points to undefined provider)

## 5. Integration: End-to-End Verification

- [x] 5.1 Write BDD scenario for AI availability with openai-compatible provider configured
- [x] 5.2 Verify existing AI BDD scenarios still pass with legacy config format
- [x] 5.3 Update CLAUDE.md Key Files table with new openai_compatible.go entry
