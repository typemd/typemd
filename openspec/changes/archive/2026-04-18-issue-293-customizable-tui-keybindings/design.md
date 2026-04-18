## Context

The 18 global TUI keybindings are declared in `tui/keys.go` as a package-level `var keys = keyMap{…}` initialized at import time. Bindings are created with `key.NewBinding(key.WithKeys("ctrl+s"), key.WithHelp("ctrl+s", "stats"))`.

Two important facts shape the design:

1. **Dispatch does not use `key.Matches(msg, keys.X)`** — it uses `switch msg.String()` with literal strings like `"ctrl+s"`, `"ctrl+e"`, `","` in `tui/update.go`. The `keys` variable is read today by exactly two call sites: `tui/help.go` (for rendering help text) and `tui/app_test.go` (tests). So the merge layer has to affect both the `msg.String()` dispatch paths AND the help labels.

2. **`core/config.go` already has a registry** with getters/setters per dotted key. We can register each action as `tui.keybindings.<action>` to integrate with `tmd config set`, the TUI Settings page, and `ConfigKeysInfo()` for free.

Related work: `tui-config-settings` capability already covers the config page; we only add new keys. Issue description (#293) pre-commits to toast warnings for validation errors, fallback to defaults, and empty string ≡ unset.

## Goals / Non-Goals

**Goals:**
- Users can override any of the 18 global keybindings through `.typemd/config.yaml` under `tui.keybindings.<action>`.
- Every keybinding is also reachable through `tmd config set tui.keybindings.<action> <key>` and the TUI Settings page.
- Invalid configurations (unknown action, bad key string, duplicate key) produce a visible warning and fall back to defaults — startup never crashes.
- Help overlay and help bar show whichever key is actually active.

**Non-Goals:**
- Sub-model internal keys (property editor, cell editor, etc. — the 129 hardcoded `msg.String()` checks outside the top-level switch). These stay hardcoded.
- Multi-key chords / sequences.
- Per-type or per-panel keybindings.
- Removing the default value for an action (every action must always have a working binding).

## Decisions

### 1. Storage shape: `map[string]string` keyed by action name

`TUIConfig.Keybindings map[string]string` — action names are `snake_case` (e.g. `stats`, `schema_explore`, `toggle_props`). This matches the style already used for config keys and makes `tui.keybindings.<action>` read naturally.

**Alternative considered:** a nested struct with one field per action. Rejected because adding an action then requires two code edits (struct + registry); a map plus a single action→default table keeps changes in one place.

### 2. `keys` becomes a function of config, not a package-level `var`

Replace the global `var keys = keyMap{…}` with:
- `defaultKeyMap()` — returns the canonical `keyMap`. This is the single source of truth for defaults.
- `buildKeyMap(cfg *core.VaultConfig) (keyMap, []KeybindingIssue)` — merges user overrides over defaults, returns issues for the caller to toast.

The model holds the resolved `keyMap` as a field (`m.keys`). `tui/help.go` and any other consumer reads from the model, not from a package-level `keys`. `tui/update.go`'s `switch msg.String()` stays literal — the key strings it compares against are derived from `m.keys` via a small helper like `m.keys.matches(action, msg.String())` when needed for actions that can be re-bound. The switch cases that are NOT user-bindable (e.g. `"tab"` for focus switching, `"esc"` for cancel) stay as hardcoded string literals.

**Alternative considered:** keep the package-level `keys` and mutate it at startup. Rejected: mutable package-level state makes tests flaky and hides the dependency.

### 3. Dispatch migration: action table, not per-case rewrite

Introduce an `action` enum (e.g. `actionStats`, `actionSchemaExplore`, `actionNewObject`, …) and a helper:

```go
func (km keyMap) actionFor(s string) (action, bool)
```

`update.go`'s outer `switch msg.String()` is refactored to:

```go
if act, ok := m.keys.actionFor(msg.String()); ok {
    switch act {
    case actionStats: …
    case actionSchemaExplore: …
    }
}
// non-rebindable literals still compared directly below
```

This keeps each case's behavior intact while routing through the resolved key map. Non-rebindable literals (`tab`, `esc`, `shift+tab`, character input inside edit modes, etc.) are unchanged.

**Alternative considered:** keep the literal `switch msg.String()` and post-translate overridden keys back to defaults. Rejected — it's fragile when the user binds an action to a key that's a literal for something else, and it hides what's rebindable.

### 4. Validation and fallback

At `buildKeyMap` time:
- **Unknown action name** → collect `KeybindingIssue{Kind: UnknownAction, Action: "…"}`. Action is not in the default table → ignore it (nothing to override).
- **Empty string value** → treat as "use default". No issue reported.
- **Invalid key string** (fails `validateKeyString`) → collect `KeybindingIssue{Kind: InvalidKey, Action, Value}`. Action keeps its default.
- **Duplicate key** (two actions resolve to the same key after merge) → collect `KeybindingIssue{Kind: DuplicateKey, Key, Actions: [a, b]}`. Both actions keep their resolved keys (user is warned, can fix).

`validateKeyString` accepts the same grammar that `key.NewBinding(key.WithKeys(...))` accepts: single chars, named keys (`esc`, `tab`, `enter`, `up`, `down`, `left`, `right`, `backspace`, `delete`, `home`, `end`, `pgup`, `pgdown`, `space`), modifier prefixes (`ctrl+`, `alt+`, `shift+`) in any order, joined with `+`. Reject unknown names, empty segments, and unsupported modifier combos.

At startup the TUI renders one toast per issue via the existing `widget.ToastModel` with severity `ToastWarning`.

### 5. Registry integration

Instead of adding 18 hand-written `configKeyRegistry` entries, we add them programmatically at package-init time via a helper:

```go
func registerKeybindingConfigKeys() {
    for action, def := range defaultKeybindings {
        configKeyRegistry["tui.keybindings."+action] = configKeyEntry{
            Get: func(cfg *VaultConfig) string { return cfg.TUI.Keybindings[action] },
            Set: func(cfg *VaultConfig, value string) {
                if cfg.TUI.Keybindings == nil {
                    cfg.TUI.Keybindings = map[string]string{}
                }
                cfg.TUI.Keybindings[action] = value
            },
            Default: def.key,
            Description: "TUI keybinding for " + def.desc,
        }
    }
}
```

called from `init()` in `core/config.go`. `defaultKeybindings` is the canonical table shared between the default `keyMap` and the registry so the two cannot drift.

### 6. Help rendering

`tui/help.go` already reads `keys.X.Help()`. After the migration it reads `m.keys.X.Help()` through a small accessor. The `Help().Key` string comes from whatever the resolved binding carries, so a user who sets `stats` to `ctrl+d` sees `ctrl+d` in both the popup and the help bar.

## Risks / Trade-offs

- **Risk:** Refactoring `update.go` dispatch to go through an action table may change behavior subtly (e.g. order of case evaluation). **Mitigation:** keep case bodies untouched, wrap only the `switch msg.String()` envelope; unit-test a handful of actions to confirm dispatch still reaches each handler.
- **Risk:** Closure-captured `action` in `registerKeybindingConfigKeys` is a classic Go pitfall. **Mitigation:** loop-local variable capture (`action := action`) or use `maps.Keys`.
- **Risk:** A user binds two actions to the same key, gets warned, but the outer switch still dispatches the alphabetically-first case. **Trade-off accepted:** both actions keep the key, last-registered wins for dispatch; the warning is surfaced so the user can fix it. We don't silently drop bindings.
- **Risk:** A user binds an action to a key that's also a non-rebindable literal (e.g., binds `stats` to `tab`). **Trade-off accepted:** the literal case wins because it's checked outside `actionFor`; the user sees no dispatch for their rebinding. We document that `tab`, `esc`, `shift+tab`, and arrow keys inside edit modes are reserved.

## Open Questions

None — the issue body pre-committed the main decisions (toast-warnings, defaults-on-failure, empty = unset). The only new commitment here is the `action` enum refactor in `update.go`, which is self-contained.
