## Context

The TUI currently has five right panel modes (`panelEmpty`, `panelObject`, `panelTypeEditor`, `panelTemplate`, `panelView`). The `panelView` mode introduced the pattern for fullscreen sub-models: a dedicated struct with its own `Update()` / `View()` methods that replaces the three-panel layout.

The core layer already provides `QueryService.VaultStats()` and `QueryService.TypeStats(typeName)` methods that power the CLI `tmd stats` command. The TUI stats mode will consume these same APIs, requiring no core changes beyond a new config key.

## Goals / Non-Goals

**Goals:**
- Add a `panelStats` mode and `statsMode` sub-model following the `viewMode` architectural pattern
- Provide an interactive vault overview with navigable type list
- Provide per-type property statistics with visual bar charts
- Support configurable Type Detail layout (fullscreen vs popup)
- Persist stats mode state across TUI sessions

**Non-Goals:**
- Chart library integration (termdash, termui) — use pure text block characters
- Relationship/link graph visualization — deferred to Web UI
- Real-time streaming stats updates — stats are computed on entry/refresh
- Modifying core stats computation logic

## Decisions

### 1. New `panelStats` mode vs extending `panelView`

**Decision:** Add a new `panelStats` right panel mode with a dedicated `statsMode` sub-model.

**Rationale:** `viewMode` is about browsing filtered/sorted object lists. `statsMode` is about displaying aggregate statistics. They serve fundamentally different purposes with different navigation stacks (vault overview → type detail vs object list → object detail). Mixing them into one sub-model would create confusing state management.

**Alternative considered:** Extending `panelView` with a stats sub-mode. Rejected because it would bloat viewMode with unrelated state and require complex mode disambiguation.

### 2. Pure text bar charts

**Decision:** Use Unicode block characters (`█▏▎▍▌▋▊▉`) for bar charts, rendered with lipgloss styling.

**Rationale:** The CLI already uses `strings.Repeat("█", count)` for select distributions. External chart libraries (termdash, termui) are standalone TUI frameworks incompatible with Bubble Tea's architecture. Pure text charts maintain visual consistency and add zero dependencies.

### 3. Type Detail layout: fullscreen vs popup

**Decision:** Support both layouts via `tui.stats_type_layout` config key (default: `fullscreen`). The `t` key toggles between them at runtime (runtime toggle is ephemeral, not persisted to config).

**Rationale:** Fullscreen provides maximum space for property statistics. Popup allows seeing the vault overview behind the detail. Making it configurable satisfies both preferences without forcing one.

**Implementation:** Fullscreen replaces the vault overview entirely. Popup uses the existing `widget.CenteredPopup` component (already used in the codebase) overlaid on the vault overview.

### 4. Stats data loading

**Decision:** Load stats data on entry into stats mode and on explicit refresh (`r` key). No automatic re-computation while in stats mode.

**Rationale:** Stats computation scans all objects, which could be expensive on large vaults. Computing once on entry and allowing manual refresh is predictable and avoids performance surprises.

### 5. File organization

**Decision:** Single file `tui/stats_mode.go` for the sub-model, following the pattern of `tui/view_mode.go`.

**Rationale:** The statsMode is simpler than viewMode (no editor sub-model, no preview split). If the file grows beyond ~500 lines, it can be split into `stats_mode_update.go` and `stats_mode_render.go` later.

## Risks / Trade-offs

- **[Performance on large vaults]** → `VaultStats()` queries all types; `TypeStats()` loads all objects of a type into memory. Mitigation: Stats are computed once on entry, not continuously. For extremely large vaults, we accept a brief loading delay.
- **[Terminal width for bar charts]** → Narrow terminals may truncate bar charts. Mitigation: Scale bar width relative to available terminal width; use proportional bars (longest option = max width) rather than absolute counts.
- **[`ctrl+s` conflict]** → Some terminal emulators intercept `ctrl+s` for XOFF flow control. Mitigation: Document this in help text; users can disable flow control with `stty -ixon`. This is a known limitation shared by many TUI apps.
