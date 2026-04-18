## Why

The 18 global TUI keybindings are hardcoded in `tui/keys.go`. Users have no way to change them, which causes real problems: `ctrl+s` (stats) conflicts with XOFF terminal flow control on many emulators, freezing the TUI. Beyond this specific conflict, users with different keyboard layouts, muscle memory from other tools, or accessibility needs cannot adapt the TUI to their workflow. Allowing keybinding customization through `.typemd/config.yaml` unblocks these users without forcing any change on the default experience.

## What Changes

- Add a `tui.keybindings` map in `.typemd/config.yaml` that overrides any of the 18 global keybindings by action name (e.g., `stats: "ctrl+d"`, `search: "ctrl+f"`).
- Merge user overrides into the default `keyMap` at TUI startup. Unset actions keep their defaults; empty string also means "use default".
- Register every keybinding as a `ConfigKeys` entry (`tui.keybindings.<action>`) so `tmd config get/set tui.keybindings.<action>` and the TUI Settings page work for them out of the box.
- Surface validation failures (unknown action, invalid key string, duplicate key across two actions) as startup toast warnings. Startup does not crash — the affected action falls back to its default binding.
- Update the help overlay and help bar so they show whichever key the user actually has configured, not the compile-time default.

## Capabilities

### New Capabilities
- `tui-keybindings`: User-configurable global TUI keybindings via `.typemd/config.yaml`, including merge rules, validation, duplicate detection, and help rendering.

### Modified Capabilities
<!-- None. tui-config-settings already covers the config page scaffolding; we only add keys to the registry. -->

## Impact

- **Code**: `tui/keys.go` (merge logic), `core/config.go` (new `TUIConfig.Keybindings` field + registry entries), `tui/help.go` and any component that reads `key.Binding.Help()` for rendering.
- **Config schema**: `.typemd/config.yaml` gains a new optional `tui.keybindings` map. Missing field keeps defaults; fully backwards-compatible.
- **Docs**: TUI keybinding documentation in `websites/docs` needs a note about customization and the list of action names.
- **Dependencies**: None. Uses existing `bubbles/v2/key` binding types and the existing config registry.
- **Out of scope**: Sub-model internal keys (129 hardcoded `msg.String()` checks), custom chord sequences, per-type overrides.
