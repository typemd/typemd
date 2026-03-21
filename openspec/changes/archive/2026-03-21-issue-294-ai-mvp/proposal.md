## Why

typemd stores knowledge as structured Markdown files, but users still manually write descriptions, assign tags, and evolve type schemas. Adding AI assistance for these repetitive tasks reduces friction and helps users discover structure in their knowledge base. Claude CLI (`claude -p`) is already available on most users' machines, making this a zero-configuration integration — no API keys to manage, no SDK dependencies to add.

## What Changes

- Add an AI provider abstraction layer (`core/ai/`) with a Claude CLI implementation that invokes `claude -p` as a subprocess
- Add AI configuration namespace (`ai.*`) to vault config for enabling/disabling AI features, model selection, and custom prompts
- Add **auto-describe** feature: press `ctrl+g` on the `description` field in TUI object detail to generate a description from object content, shown as inline preview (Tab to accept, Esc to reject)
- Add **auto-tag** feature: press `ctrl+g` on the `tags` field to get AI-suggested tags (marked as existing or new), shown as a selectable popup list
- Add **schema explore** feature: press `ctrl+e` from sidebar to enter a full-width AI analysis panel that examines objects of a type and suggests schema modifications (add/modify/remove properties)
- AI prompts use type and property `description` fields as semantic context, giving users control over AI behavior through schema documentation

## Capabilities

### New Capabilities
- `ai-provider`: AI provider interface and Claude CLI implementation — subprocess invocation, structured JSON output via `--json-schema`, configuration, and error handling
- `ai-auto-describe`: TUI inline AI description generation for objects — context assembly, preview-before-apply interaction, accept/reject flow
- `ai-auto-tag`: TUI AI tag suggestion with existing/new tag classification — tag list popup, multi-select, apply flow
- `ai-schema-explore`: TUI AI-powered schema analysis and suggestion — object sampling, structured suggestions, per-suggestion accept/skip, schema modification

### Modified Capabilities
- `vault-config`: Add `ai` namespace with `ai.enabled`, `ai.model`, `ai.prompts.describe`, `ai.prompts.tag`, `ai.prompts.explore`, `ai.explore.sample_count`, `ai.explore.body_truncate` keys

## Impact

- **New package**: `core/ai/` (Provider interface + ClaudeCLI implementation)
- **core/**: Extend `VaultConfig` with `AIConfig` struct, extend `configKeyRegistry` with `ai.*` keys
- **tui/**: New keybindings (`ctrl+g`, `ctrl+e`), new AI-related UI components (inline preview, tag popup, schema explore panel), new right panel mode for schema explore
- **Dependencies**: No new Go module dependencies (uses `os/exec` to invoke `claude` CLI)
- **External dependency**: Requires `claude` CLI installed and authenticated on user's machine
- **No breaking changes**: AI features are opt-in (`ai.enabled: true`) and hidden when disabled or `claude` binary not found
