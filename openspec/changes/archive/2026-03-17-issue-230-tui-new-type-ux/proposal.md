## Why

The TUI's "New Type" creation flow is minimal — it only prompts for a type name via an inline text input at the bottom of the sidebar, then creates an empty schema and opens the type editor. Users must manually navigate the type editor to set emoji, plural, and other fields. In contrast, the "New Object" flow (PR #242) already uses a title panel with multi-field input and live preview, providing a much better UX. The "New Type" flow should match this level of polish.

## What Changes

- Move "New Type" input from sidebar inline text to the **title panel**, matching the "New Object" creation pattern
- Add **three input fields** in the title panel: emoji (optional), name (required), plural (optional)
- Add **Tab navigation** between fields
- Add **live preview** of the type schema in the right panel as the user types
- Add **inline validation** with error messages (duplicate names, empty name)
- After creation, automatically open the type editor with emoji/plural already populated

## Capabilities

### New Capabilities
- `tui-type-creation`: Covers the TUI "New Type" creation wizard — title panel layout, multi-field input, live preview, validation, and post-creation behavior

### Modified Capabilities
<!-- No existing spec-level behavior changes needed. The type-schema core creation logic remains the same. -->

## Impact

- **tui/app.go** — Replace `newTypeMode`/`newTypeName` fields with a `createTypeState` struct; update `startNewType()`
- **tui/update.go** — Replace `updateNewType()` with new handler using title panel fields
- **tui/detail.go** — Add title panel rendering for type creation mode
- **tui/list.go** — Minor: sidebar rendering during type creation
- **No core/ changes** — `SaveType()` and `TypeSchema` remain unchanged
- **No breaking changes** — Pure UX improvement
