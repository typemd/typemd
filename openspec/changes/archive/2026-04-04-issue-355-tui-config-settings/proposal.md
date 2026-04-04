## Why

Config values can only be viewed and edited via CLI commands (`tmd config get`, `tmd config set`, `tmd config list`). Users must leave the TUI to change settings — breaking their flow. A dedicated config settings page lets users browse and edit `.typemd/config.yaml` directly in the terminal.

## What Changes

- Add `Description` field to `configKeyRegistry` entries in `core/config.go`, providing human-readable descriptions for all 23 config keys
- Add a new `configEditor` TUI sub-model (`tui/config_editor.go`) with two-column layout: category navigation (left) and settings list (right)
- Add `panelConfig` right panel mode, activated by `,` keybinding
- Use `widget.CenteredPopup` for editing individual settings, showing key name, description, default value, and input field
- Each edit saves immediately via `SetConfigValue()` on Enter

## Capabilities

### New Capabilities
- `tui-config-settings`: TUI config settings page with category browsing, two-column layout, and popup editing of vault config values

### Modified Capabilities

## Impact

- **`core/config.go`**: Add `Description` field to `configKeyEntry` struct and populate for all 23 registry entries. New exported function to list keys with metadata.
- **`tui/`**: New `config_editor.go` file, new `panelConfig` mode in `app.go`, new `,` keybinding in `keys.go`, new render path in `app_render.go`, new help entry in `help.go`.
- **No breaking changes**: This is purely additive. Existing config CLI commands and the config file format are unchanged.
