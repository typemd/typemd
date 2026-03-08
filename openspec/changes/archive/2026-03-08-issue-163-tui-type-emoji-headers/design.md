## Context

The TUI object list panel groups objects by type and displays headers like `▼ book (4)`. Type schemas now have an `emoji` field (added in #145/#162). The `typeGroup` struct in `tui/app.go` currently only stores `Name`, `Objects`, and `Expanded` — no emoji data.

The rendering happens in `renderList()` in `tui/list.go`, which formats headers as `fmt.Sprintf(" %s %s (%d)", arrow, g.Name, len(g.Objects))`.

## Goals / Non-Goals

**Goals:**
- Display type emoji in group headers when available (e.g., `▼ 📚 book (4)`)
- No visual change for types without emoji

**Non-Goals:**
- Emoji fallback for types without emoji defined
- Making emoji configurable per-session or per-view
- Emoji in object rows (only group headers)

## Decisions

### Store emoji in `typeGroup` struct

Add an `Emoji` field to `typeGroup`. Populate it during `buildGroups()` by calling `vault.LoadType()` for each type. This avoids repeated schema lookups at render time.

**Alternative considered:** Look up emoji in `renderList()` — rejected because `renderList` is a pure function that doesn't have access to the vault, and passing vault through would complicate the rendering layer.

### Format: emoji between arrow and name

When emoji is present: `▼ 📚 book (4)`. When absent: `▼ book (4)` (unchanged). This keeps the arrow as the leftmost element for visual consistency.

## Risks / Trade-offs

- **Terminal emoji support**: Some terminals render emoji with inconsistent widths. This is a pre-existing concern (emoji already appears in other contexts) and not specific to this change.
