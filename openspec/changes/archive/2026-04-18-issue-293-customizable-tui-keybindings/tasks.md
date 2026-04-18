## 1. Config model & registry

- [x] 1.1 Add unit tests for `validateKeyString` covering valid keys (`a`, `ctrl+s`, `ctrl+shift+a`, `esc`, `enter`), invalid keys (`crtl+s` typo, `ctrl+`, empty segment, `ctrl++s`), and modifier normalization (`shift+ctrl+a` acceptable or rejected consistently)
- [x] 1.2 Implement `validateKeyString` in `tui/keys.go` (or new `tui/keybindings.go`) to make 1.1 pass
- [x] 1.3 Add unit tests asserting `ConfigKeysInfo()` includes one entry per rebindable action under `tui.keybindings.` with non-empty `Description` and the compile-time default in `Default`
- [x] 1.4 Introduce `defaultKeybindings` table (action name → default key + description) shared between `core/config.go` registration and `tui/keys.go` defaults
- [x] 1.5 Add `Keybindings map[string]string` field to `TUIConfig` in `core/config.go` with `yaml:"keybindings,omitempty"` and register `tui.keybindings.<action>` entries programmatically so 1.3 passes
- [x] 1.6 Add unit tests for `tmd config set tui.keybindings.stats ctrl+d` end-to-end: writes `.typemd/config.yaml`, `GetConfigValue` returns `ctrl+d`, re-open vault sees the value

## 2. Keybinding merge layer

- [x] 2.1 Add unit tests for `buildKeyMap` covering: no config → all defaults; override of one action; empty string → default; unknown action → `UnknownAction` issue; invalid key → `InvalidKey` issue and action keeps default; duplicate key across two actions → `DuplicateKey` issue with both action names
- [x] 2.2 Refactor `tui/keys.go` to expose `defaultKeyMap()` (same content as today's `var keys`) and `buildKeyMap(cfg *core.VaultConfig) (keyMap, []KeybindingIssue)` to make 2.1 pass
- [x] 2.3 Keep a package-level `var keys = defaultKeyMap()` alias so pre-existing `tui/app_test.go` tests that read `keys.X.Help().Key` still work. Production code reads through `model.keys` (the resolved keyMap) — see startup wiring in task 3.3.

## 3. Dispatch integration

- [x] 3.1 Add unit tests in `tui/keybindings_test.go` covering: overridden key triggers the expected action on the resolved keyMap; default key no longer triggers it; unset action still matches its default; `actionFor` returns the correct action enum for each default and overridden key. (BDD is deferred — `tui/features/` does not exist yet; the codebase establishes TUI tests via `tui/*_test.go`.)
- [x] 3.2 Refactor `tui/update.go` outer `switch msg.String()` for rebindable actions to route through `m.keys.actionFor(msg.String())` + `switch action`, leaving non-rebindable literals (`tab`, `esc`, arrows inside edit modes) as hardcoded string cases — make 3.1 tests pass
- [x] 3.3 Wire TUI startup to call `buildKeyMap(vault.Config())`, store the resolved keyMap on the model, and queue one `widget.ToastWarning` per `KeybindingIssue` on the first frame

## 4. Help & rendering

- [x] 4.1 Add unit tests asserting `helpEntries` reflects overridden keys when the model's keyMap is built from a config with overrides
- [x] 4.2 Refactor `tui/help.go` to read from a caller-supplied `keyMap` (or a model accessor) instead of the package-level `keys`, make 4.1 pass
- [x] 4.3 Verify help bar rendering in `tui/view.go` (or equivalent) uses the same resolved keyMap

## 5. Tests & docs

- [x] 5.1 Update `tui/app_test.go` to use the new `keyMap` construction path
- [x] 5.2 Add a unit test asserting `tui.keybindings.<action>` keys appear under the TUI config page (visible through `newConfigEditor`/category list)
- [x] 5.3 Update `websites/docs` TUI keybinding page with a "Customization" section listing action names, example config, and behavior on invalid input
- [x] 5.4 Run `make test` — go build + go test + go vet all clean
