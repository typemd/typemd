## 1. AI Config Extension

- [x] 1.1 Write BDD scenarios for AI config in vault-config (ai.enabled, ai.model, ai.prompts.*, ai.explore.*)
- [x] 1.2 Implement BDD step definitions for AI config scenarios
- [x] 1.3 Add AIConfig, PromptsConfig, ExploreConfig structs to VaultConfig and extend configKeyRegistry
- [x] 1.4 Add unit tests for AI config edge cases (missing keys, invalid values, defaults)

## 2. AI Provider Interface and ClaudeCLI Implementation

- [x] 2.1 Write BDD scenarios for AI provider (completion with/without JSON schema, binary not found, non-zero exit, context cancellation)
- [x] 2.2 Implement BDD step definitions for AI provider scenarios
- [x] 2.3 Create core/ai/ package with Provider interface, request/response types (provider.go)
- [x] 2.4 Implement ClaudeCLI provider (claude_cli.go) — subprocess invocation with -p, --output-format json, --tools "", --json-schema, --model flags
- [x] 2.5 Add unit tests for ClaudeCLI edge cases (stderr parsing, timeout, argument construction)

## 3. AIService Domain Operations

- [x] 3.1 Write BDD scenarios for AIService (Describe, SuggestTags, ExploreSchema — prompt assembly, response parsing)
- [x] 3.2 Implement BDD step definitions for AIService scenarios
- [x] 3.3 Create AIService with Describe method — prompt context assembly (name, properties, body, schema descriptions), JSON schema for response, response parsing
- [x] 3.4 Create AIService SuggestTags method — prompt with existing tags list, TagSuggestion/SuggestedTag response types, JSON schema
- [x] 3.5 Create AIService ExploreSchema method — object sampling (configurable count), body truncation, SchemaExploration/SchemaSuggestion response types, JSON schema
- [x] 3.6 Create built-in prompt templates (prompts.go) for describe, tag, and explore operations
- [x] 3.7 Add unit tests for prompt assembly (correct context included) and response parsing edge cases

## 4. Vault AI Integration

- [x] 4.1 Write BDD scenarios for AI availability detection (enabled + binary found, enabled + binary missing, disabled)
- [x] 4.2 Implement BDD step definitions for AI availability scenarios
- [x] 4.3 Add AIService field to Vault, initialize during Open() with exec.LookPath check, expose AIService() accessor
- [x] 4.4 Add unit tests for Vault AI initialization edge cases

## 5. TUI: Auto-Describe

- [x] 5.1 Add ctrl+g keybinding to keyMap (AIGenerate)
- [x] 5.2 Implement ctrl+g handler in Update() — detect focused field (description vs tags vs other), dispatch to appropriate AI flow
- [x] 5.3 Create AI loading state and spinner rendering for description field
- [x] 5.4 Implement inline preview component (ghost text style) with Tab-to-accept, Esc-to-reject
- [x] 5.5 Wire Describe result into inline preview, handle accept (write to object, mark dirty) and reject
- [x] 5.6 Handle AI errors (display inline error message, dismissable with any key)

## 6. TUI: Auto-Tag

- [x] 6.1 Create tag suggestion popup component (selectable list with checkboxes, existing/new markers)
- [x] 6.2 Implement ctrl+g on tags field handler — invoke SuggestTags, show loading state
- [x] 6.3 Wire SuggestTags result into popup — filter out already-assigned tags, classify existing vs new
- [x] 6.4 Implement popup navigation (up/down, Space to toggle, Enter to confirm, Esc to cancel)
- [x] 6.5 Implement confirm action — link existing tags, create + link new tags, mark object dirty
- [x] 6.6 Handle AI errors in tag flow

## 7. TUI: Schema Explore

- [x] 7.1 Add ctrl+e keybinding to keyMap (SchemaExplore)
- [x] 7.2 Add panelSchemaExplore to rightPanelMode enum and schemaExplorer sub-model to model struct
- [x] 7.3 Create type selection prompt for schema explore entry (list available types)
- [x] 7.4 Create schemaExplorer sub-model with Update()/View() — loading state, suggestion list rendering
- [x] 7.5 Implement ctrl+e handler — type selection, invoke ExploreSchema, transition to panelSchemaExplore
- [x] 7.6 Implement suggestion navigation (up/down), accept (Enter — modify schema file), skip (s)
- [x] 7.7 Implement schema modification on accept (add/modify/remove property in TypeSchema, save file)
- [x] 7.8 Implement completion summary view and Esc to return to sidebar

## 8. Integration and Polish

- [x] 8.1 Add AI keybinding hints to help overlay (conditional on AIService availability)
- [x] 8.2 Verify all AI features are hidden when ai.enabled is false or claude binary missing
- [x] 8.3 Run full test suite (`go test ./...`) and fix any regressions
