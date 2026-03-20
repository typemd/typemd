## 1. Core: Config Key

- [x] 1.1 Write BDD scenario for `tui.stats_type_layout` config key (load fullscreen, popup, default)
- [x] 1.2 Implement step definitions for stats_type_layout BDD scenarios
- [x] 1.3 Add `StatsTypeLayout` field to `TUIConfig` struct and register `tui.stats_type_layout` in config key registry
- [x] 1.4 Add unit test for invalid stats_type_layout values and default behavior

## 2. TUI: Stats Mode Sub-Model Foundation

- [x] 2.1 Add `panelStats` constant to `rightPanelMode` in `app.go`
- [x] 2.2 Create `stats_mode.go` with `statsMode` struct, `newStatsMode()`, and stub `Update()` / `View()` methods
- [x] 2.3 Add `ctrl+s` keybinding in `keys.go` and wire it to `openStatsMsg` in `app.go` Update
- [x] 2.4 Wire `panelStats` rendering in `app_render.go` to call `statsMode.View()` fullscreen

## 3. TUI: Vault Overview Screen

- [x] 3.1 Implement `statsMode.loadVaultStats()` to call `vault.VaultStats()` and sort by count descending
- [x] 3.2 Implement Vault Overview `View()` rendering — header (total objects/types) + type list with emoji, name, count, relative last-updated
- [x] 3.3 Implement j/k and arrow key navigation for the type list
- [x] 3.4 Implement Enter to transition to Type Detail screen
- [x] 3.5 Implement Esc to exit stats mode and restore previous panel state

## 4. TUI: Type Detail Screen

- [x] 4.1 Implement `statsMode.loadTypeStats(typeName)` to call `vault.TypeStats(typeName)`
- [x] 4.2 Implement bar chart rendering helper — proportional horizontal bars using `█▏▎▍▌▋▊▉` block characters
- [x] 4.3 Implement Type Detail `View()` rendering — header + per-property stats (number, select, checkbox, date, relation, string fill rate)
- [x] 4.4 Implement j/k scrolling for long property lists on Type Detail
- [x] 4.5 Implement Esc to return from Type Detail to Vault Overview

## 5. TUI: Type Detail Layout (Fullscreen / Popup)

- [x] 5.1 Read `tui.stats_type_layout` config on stats mode entry and set initial layout
- [x] 5.2 Implement fullscreen Type Detail rendering (replaces vault overview)
- [x] 5.3 Implement popup Type Detail rendering using `widget.CenteredPopup`
- [x] 5.4 Implement `t` key to toggle between fullscreen and popup at runtime

## 6. TUI: Stats Refresh

- [x] 6.1 Implement `r` key to refresh stats — recompute VaultStats on overview, TypeStats on detail

## 7. TUI: State Persistence

- [x] 7.1 Add stats mode fields to `SessionState` struct (`StatsCursor`, `StatsScroll`, `StatsTypeName`)
- [x] 7.2 Implement `captureState()` logic for stats mode (save cursor, scroll, type name when in stats)
- [x] 7.3 Implement `restoreStatsMode()` to rebuild stats mode from saved state on startup
- [x] 7.4 Add unit test for stats state precedence over view state

## 8. TUI: Help Text and Integration

- [x] 8.1 Add `ctrl+s` to help overlay key list
- [x] 8.2 Add stats-mode-specific help keys (j/k, Enter, Esc, t, r) to help overlay when in stats mode
