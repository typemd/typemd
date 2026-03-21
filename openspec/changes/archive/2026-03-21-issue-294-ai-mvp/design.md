## Context

typemd is a local-first CLI knowledge management tool where objects (Markdown files with YAML frontmatter) are organized by type schemas. Users currently manage descriptions, tags, and schema evolution manually. Issue #294 introduces AI assistance for these tasks.

The codebase follows Clean Architecture with CQRS. Configuration uses interface-layer namespacing (`cli.*`, `tui.*`) in `.typemd/config.yaml`. The TUI uses Bubble Tea with a three-panel layout and a right panel mode system for different views.

Key constraint: users already have Claude CLI (`claude`) installed and authenticated. Rather than adding an SDK dependency and requiring API key management, we invoke `claude -p` as a subprocess — zero additional configuration.

## Goals / Non-Goals

**Goals:**
- Provide AI-assisted description generation, tag suggestion, and schema exploration in the TUI
- Use Claude CLI as the AI backend with no additional dependencies or API key setup
- Follow existing architecture patterns (Clean Architecture, CQRS, sub-model pattern for TUI panels)
- Make AI features opt-in and gracefully hidden when unavailable
- Use structured JSON output (`--json-schema`) for reliable response parsing

**Non-Goals:**
- Direct Anthropic API integration (future — can add as another Provider implementation)
- Streaming responses in Phase 1 (will use one-shot JSON; streaming is a Phase 2 enhancement)
- CLI commands for AI features (`tmd ai ...`) — TUI-only for MVP
- AI-assisted body writing or content generation beyond descriptions
- Multi-provider support (Ollama, OpenAI, etc.) — architecture supports it but MVP is Claude CLI only

## Decisions

### Decision 1: Claude CLI as AI Provider

**Choice:** Invoke `claude -p --output-format json --json-schema <schema> --tools ""` as a subprocess.

**Alternatives considered:**
- **Anthropic SDK (`anthropic-sdk-go`)**: Requires new dependency, API key management, token counting. More control but more setup.
- **MCP client**: typemd already has an MCP server, but adding an MCP client to call Claude adds unnecessary indirection.

**Rationale:** Claude CLI handles auth, model selection, rate limiting, and retries. Users who have Claude Code installed get AI features for free. `--json-schema` ensures structured output without prompt engineering for format. `--tools ""` disables tool use for pure text generation.

### Decision 2: core/ai/ Sub-Package

**Choice:** Create `core/ai/` as a new sub-package with:
- `provider.go` — `Provider` interface and request/response types
- `claude_cli.go` — `ClaudeCLI` struct implementing `Provider`
- `service.go` — `AIService` with use-case methods (`Describe`, `SuggestTags`, `ExploreSchema`)

**Alternatives considered:**
- **Flat in core/**: Follows existing convention but core/ already has 30+ files. AI is a distinct concern.

**Rationale:** AI is a self-contained domain with clear boundaries. The sub-package keeps `core/` manageable while the `Provider` interface allows future implementations. Dependency direction: `core/ai/` imports from `core/` (for `Object`, `TypeSchema`, etc.), never the reverse. `Vault` holds an optional `*ai.AIService` initialized during `Open()`.

### Decision 3: Config Structure

**Choice:** Add `AIConfig` struct under `ai` namespace in `VaultConfig`:

```yaml
ai:
  enabled: true
  model: claude-sonnet-4-6-20250627
  prompts:
    describe: "custom system prompt for descriptions"
    tag: "custom system prompt for tags"
    explore: "custom system prompt for schema exploration"
  explore:
    sample_count: 10
    body_truncate: 500
```

All fields optional. `ai.enabled` defaults to `false`. Model defaults to `claude-sonnet-4-6-20250627`. Prompts default to built-in prompts defined in `core/ai/prompts.go`.

Config keys registered in `configKeyRegistry`: `ai.enabled`, `ai.model`, `ai.prompts.describe`, `ai.prompts.tag`, `ai.prompts.explore`, `ai.explore.sample_count`, `ai.explore.body_truncate`.

### Decision 4: TUI Integration Pattern

**Choice:** Context-sensitive `ctrl+g` keybinding + dedicated `ctrl+e` for schema explore.

- `ctrl+g` on `description` field → inline preview (ghost text style, Tab to accept, Esc to reject)
- `ctrl+g` on `tags` field → popup list with checkboxes (existing tags marked, new tags starred)
- `ctrl+g` on other fields → no response (silent, no error)
- `ctrl+e` from sidebar → new `panelSchemaExplore` right panel mode

**UI states during AI call:**
- Loading spinner in the target area while waiting for response
- Error message displayed inline if `claude` binary not found or call fails

### Decision 5: Prompt Context Assembly

**Choice:** Each AI use case assembles its own prompt context:

**Auto-describe:**
- Object name, all properties (key-value), body content
- Type schema with property descriptions (semantic context)
- System prompt: "Generate a concise description for this object based on its content and type context."

**Auto-tag:**
- Object name, properties, body content
- Complete list of existing tags (with descriptions if available)
- Type schema context
- System prompt: "Suggest relevant tags for this object. Mark each as existing or new."

**Schema explore:**
- Sample N objects of the selected type (configurable, default 10; if fewer exist, take all)
- Body truncated to M chars (configurable, default 500)
- Current type schema with all property definitions and descriptions
- System prompt: "Analyze these objects and suggest schema improvements: properties to add, modify, or remove."

### Decision 6: AI Availability Detection

**Choice:** On Vault.Open(), if `ai.enabled` is true:
1. Check if `claude` binary exists in PATH (`exec.LookPath("claude")`)
2. If found, initialize `AIService` with `ClaudeCLI` provider
3. If not found, log a warning and leave `AIService` nil

TUI checks `vault.AIService() != nil` before showing AI keybinding hints. AI keybindings are silently no-op when AIService is nil.

## Risks / Trade-offs

- **[Subprocess latency]** → Claude CLI cold start adds 1-2s overhead per call. Mitigation: show loading spinner, consider caching for repeated calls in same session. Phase 2 streaming will improve perceived latency.
- **[Claude CLI version dependency]** → We depend on `claude` CLI's `--json-schema` and `--output-format json` flags existing. Mitigation: version-check on startup, graceful degradation if flags unsupported.
- **[Token costs opaque]** → Users may not realize AI calls cost money via their Claude account. Mitigation: show a one-time notice when AI is first used in a session.
- **[Large vault sampling]** → Schema explore truncates data, may miss patterns in large vaults. Mitigation: configurable `sample_count` and `body_truncate` settings.
- **[Error handling]** → CLI errors (auth expired, rate limited, network down) need clean user-facing messages. Mitigation: parse stderr for known error patterns, show actionable messages.
