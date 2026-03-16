## Why

The TUI's "new object" creation flow is minimal — pressing `n` only prompts for a name in a basic inline text input, with no template selection, no name template support, and no post-creation guidance. This creates friction for users who have set up templates or name templates, and makes bulk object creation tedious.

## What Changes

- Add two distinct creation modes triggered by different keybindings: `n` for "Create & Edit" (single object, auto-enters body edit mode) and `N` for "Quick Create" (batch creation, stays in input mode for rapid entry)
- Add template selection step when a type has multiple templates (list picker UI); auto-apply when exactly one template exists; skip when none
- Auto-skip name input when type has a name template defined (in "Create & Edit" mode)
- Show inline validation feedback for `unique: true` constraint violations
- Show visual confirmation flash on successful creation (especially in batch mode)
- Update help bar to reflect new keybindings and available actions

## Capabilities

### New Capabilities
- `tui-object-creation`: Two creation modes (Create & Edit via `n`, Quick Create via `N`), template selection UI, name template auto-apply, unique validation feedback, and creation success confirmation

### Modified Capabilities
- `tui-object-list`: Add `N` keybinding for Quick Create mode alongside existing `n`

## Impact

- `tui/` — Main area of change: new creation modes, template picker UI component, input handling, visual feedback
- `tui/update.go` — Refactor `updateNewObject()` to support two modes and multi-step flow (template selection → name input → create)
- `tui/app.go` — New model fields for creation state (mode, selected template, step tracking)
- `tui/list.go` — Help bar updates for new keybindings
- No core/ changes needed — `ObjectService.Create()`, `ListTemplates()`, `LoadTemplate()`, and name template evaluation already support all required functionality
