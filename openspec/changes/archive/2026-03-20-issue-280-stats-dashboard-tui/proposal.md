## Why

The TUI has no way to view vault or type statistics — users must exit to CLI (`tmd stats` / `tmd stats --type <name>`) to see aggregate data. Adding an interactive stats dashboard mode lets users explore vault-level summaries and per-type property analysis without leaving the TUI, using the same core `VaultStats` / `TypeStats` APIs that already power the CLI.

## What Changes

- Add a new `panelStats` right panel mode and `statsMode` sub-model to the TUI, following the same architectural pattern as `viewMode`
- `ctrl+s` global shortcut enters stats mode from any panel
- **Vault Overview** screen (fullscreen): total objects/types/tags summary, type list sorted by count with emoji + name + count + relative last-updated, navigable with j/k, Enter drills into type detail
- **Type Detail** screen (fullscreen or popup): per-property statistics — select frequency distribution with bar charts, number min/max/avg, checkbox ratios, date ranges, relation counts, fill rates
- `tui.stats_type_layout` config key (`fullscreen` | `popup`, default `fullscreen`) controls Type Detail layout; `t` key toggles at runtime
- Stats mode state persisted in `tui-state.yaml` and restored on next launch
- Pure text bar charts using block characters (`█▏▎▍▌▋▊▉`) — no external chart library dependency

## Capabilities

### New Capabilities
- `tui-stats-mode`: Interactive stats dashboard mode in the TUI — vault overview and type detail screens, navigation, keyboard shortcuts, bar chart rendering
- `tui-stats-persistence`: Stats mode state persistence in tui-state.yaml — selected type, cursor position, scroll offset, layout preference

### Modified Capabilities
- `vault-config`: Add `tui.stats_type_layout` config key (fullscreen | popup)

## Impact

- **TUI code** (`tui/`): New `stats_mode.go` sub-model + integration into `app.go` panel system, `app_render.go` fullscreen rendering, `state.go` persistence, `keys.go` keybinding
- **Core code** (`core/`): Add `stats_type_layout` to `TUIConfig` and config key registry in `vault_config.go`
- **No new dependencies**: Pure text rendering with lipgloss, no external chart libraries
- **No breaking changes**: Additive feature only
