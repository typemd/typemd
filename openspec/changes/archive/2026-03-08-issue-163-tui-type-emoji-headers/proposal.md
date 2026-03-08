## Why

Type schemas now support an `emoji` field (#145), but the TUI object list panel doesn't use it. Adding emoji prefixes to group headers (e.g., `▼ 📚 book (4)`) provides instant visual identification of type groups.

## What Changes

- Modify TUI object list group header rendering to include the type emoji when defined
- Pass type schema emoji data to the list rendering layer
- Types without emoji display unchanged (no placeholder or extra spacing)

## Capabilities

### New Capabilities

_None — this change enhances existing TUI rendering._

### Modified Capabilities

_None — no spec-level behavioral requirements are changing._

## Impact

- `tui/list.go` — group header format string
- `tui/app.go` — `typeGroup` struct may need an emoji field, or emoji lookup at render time
- No API changes, no dependency changes, no breaking changes
