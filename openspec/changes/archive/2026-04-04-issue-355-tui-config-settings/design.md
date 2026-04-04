## Context

The TUI currently supports multiple full-width panel modes (stats, schema explore, views) but has no config editing capability. Users must exit to the CLI (`tmd config set`) to change settings. The `configKeyRegistry` in `core/config.go` already provides type-safe get/set for all 23 config keys, but lacks human-readable descriptions needed for a browsable UI.

## Goals / Non-Goals

**Goals:**
- Provide a browsable, editable config settings page within the TUI
- Organize settings into logical categories matching the config.yaml structure
- Allow immediate save on edit (no separate "save" step)
- Follow established TUI patterns (panel mode, vim navigation, popup editing)

**Non-Goals:**
- AI providers map editing (complex nested `map[string]ProviderConfig` — defer to CLI)
- Config file creation wizard
- Undo/redo for config changes
- Live reload of TUI settings after edit (requires app restart for TUI-specific settings like debounce)

## Decisions

### 1. Two-column layout with category/settings split

Use a left column for category navigation (General, CLI, TUI, AI, Web) and a right column for settings within the selected category. This mirrors the familiar settings UI pattern and scales well as more config keys are added.

**Alternative considered:** Single flat list with all keys — rejected because 23+ keys without grouping is unwieldy.

### 2. Add `Description` field to `configKeyEntry`

Extend the existing `configKeyEntry` struct with a `Description string` field. This keeps metadata co-located with the get/set logic and avoids a separate description registry.

**Alternative considered:** Separate `configKeyDescriptions` map — rejected because it would diverge from the registry and require dual maintenance.

### 3. New `configEditor` sub-model following `statsMode` pattern

Create `tui/config_editor.go` with a `configEditor` struct that follows the same lifecycle as `statsMode`:
- Created on `,` keypress → assigned to `m.configEditor`
- `m.rightPanel = panelConfig`
- Has its own `Update(msg)` and `View()` methods
- Esc returns to previous panel mode

### 4. CenteredPopup for value editing

Reuse the existing `widget.CenteredPopup` / `renderOverlayPopup()` pattern for the edit popup. The popup shows: key name, description, current value, default value hint, and a text input. For boolean keys, cycle through true/false/unset instead of free-text input.

### 5. Exported `ConfigKeyInfo` and `ConfigKeysInfo()` function

Add an exported `ConfigKeyInfo` struct and `ConfigKeysInfo()` function to `core/config.go` that returns key metadata (name, description, default, current value) for TUI consumption. This keeps the TUI decoupled from internal registry details.

### 6. Category assignment via key prefix

Derive categories from the dot-notation key prefix:
- No prefix (`date_format`, `datetime_format`) → "General"
- `cli.*` → "CLI"
- `tui.*` → "TUI"
- `ai.*` → "AI"
- `web.*` → "Web"

This requires no additional metadata — category is implicit in the key name.

## Risks / Trade-offs

- **[TUI settings not live-reloaded]** → Some TUI config changes (e.g., `tui.debounce_ms`) only take effect after restarting the TUI. Mitigation: show a toast notification "Restart TUI to apply changes" for TUI-category settings.
- **[Boolean pointer tri-state complexity]** → `show_warnings` and `show_success` use `*bool` (nil/true/false). Mitigation: popup shows three options — "true", "false", "unset (default: X)" — cycling with Enter.
- **[Concurrent external edits]** → If config.yaml is edited externally while the config page is open, displayed values may be stale. Mitigation: re-read config from vault on each entry to the config page (not on every keypress).
