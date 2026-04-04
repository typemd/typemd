## 1. Core: Config key descriptions and metadata API

- [x] 1.1 Write BDD scenarios for config key descriptions (all keys have descriptions, ConfigKeysInfo returns metadata)
- [x] 1.2 Implement step definitions for config key description scenarios
- [x] 1.3 Add `Description` field to `configKeyEntry`, populate all 23 keys, add `ConfigKeyInfo` struct and `ConfigKeysInfo()` function
- [x] 1.4 Add unit tests for ConfigKeysInfo edge cases (unset values show empty, set values show current)

## 2. TUI: Keybinding and panel mode setup

- [x] 2.1 Add `Settings` keybinding (`,`) to `keyMap` in `tui/keys.go`
- [x] 2.2 Add `panelConfig` to `rightPanelMode` enum in `tui/app.go`
- [x] 2.3 Add `,` help entry to `helpEntries()` in `tui/help.go`

## 3. TUI: Config editor sub-model

- [x] 3.1 Create `tui/config_editor.go` with `configEditor` struct — categories, settings list, two-column layout, cursor navigation
- [x] 3.2 Implement `configEditor.Update()` — j/k navigation, Tab column switching, Enter to open edit popup, Esc to exit
- [x] 3.3 Implement `configEditor.View()` — two-column render with category list (left) and settings list (right)
- [x] 3.4 Implement edit popup — CenteredPopup with key name, description, default, text input; boolean cycling; save on Enter
- [x] 3.5 Implement `configEditor.HelpBar()` — contextual navigation hints

## 4. TUI: Integration with app lifecycle

- [x] 4.1 Wire `,` keypress in `app.Update()` to create `configEditor` and set `panelConfig`
- [x] 4.2 Add `panelConfig` render path in `app_render.go` `View()`
- [x] 4.3 Handle Esc from configEditor to restore previous panel mode

## 5. Testing and verification

- [x] 5.1 Add unit tests for configEditor navigation (category switching, settings cursor, column focus)
- [x] 5.2 Add unit tests for configEditor edit flow (open popup, save value, cancel, clear to default)
- [x] 5.3 Verify all existing tests still pass
