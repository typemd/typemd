## Context

Date properties currently fall through to the generic textinput widget in both `propEditor` (properties panel) and `cellEdit` (table view). Users must manually type `YYYY-MM-DD`, with no visual feedback until validation on Enter. All other property types with complex inputs (select, multi_select, relation) already have dedicated editing widgets.

Existing edit widget pattern: state fields embedded in the parent struct (`propEditor` or `cellEdit`), with mode constants controlling which widget is active. The `propEditor` has `propEditMode` (navigate, textInput, selectPick, multiPick, relationPick, relationMultiPick) and `cellEdit` has `cellEditMode` (textInput, selectPick, multiPick).

## Goals / Non-Goals

**Goals:**
- Provide an intuitive date editing experience with two modes: segmented input (default) and inline calendar
- Share the date editing widget between properties panel and table view cell editing
- Show live day-of-week feedback during segmented input
- Support keyboard-driven calendar navigation with today marker and jump
- Follow existing widget patterns (embedded state, mode constants)

**Non-Goals:**
- `datetime` property editing (remains textinput)
- Date range or recurring date selection
- Locale-specific date formats (always YYYY-MM-DD)
- Mouse/click interaction

## Decisions

### 1. Shared `dateEdit` struct as a self-contained sub-model

Create a new `dateEdit` struct in `tui/date_edit.go` that encapsulates all date editing state and logic (both modes). Both `propEditor` and `cellEdit` embed a `*dateEdit` field and delegate to it.

**Rationale:** The relation picker already uses a similar pattern — state fields in `propEditor` but rendering/update logic in `relation_picker.go`. A dedicated struct is cleaner because the date picker has two modes with significant state (segments, calendar cursor, focused month).

**Alternative considered:** Adding date-specific state fields directly to `propEditor` and `cellEdit` (like the select picker does). Rejected because it would require duplicating segment/calendar state in both structs.

### 2. Two internal modes: segment and calendar

The `dateEdit` struct has its own internal mode toggle:
- `dateSegmentMode`: Three text segments (year, month, day) with the focused segment highlighted. Arrow keys increment/decrement, digits replace.
- `dateCalendarMode`: 7-column month grid. Arrow/vim keys navigate days, `H`/`L` switch months.

Toggle between modes with `c`. Both modes produce the same output: a `time.Time` value.

**Rationale:** Segmented input is faster for known dates; calendar is better for browsing. The `c` toggle is consistent with the issue spec.

### 3. New mode constants for prop and cell editors

Add `propModeDateSegment` and `propModeDateCalendar` to `propEditMode`, and `cellModeDateSegment` and `cellModeDateCalendar` to `cellEditMode`.

**Rationale:** Existing help bar rendering and border color logic switch on mode constants. New modes allow `[DATE]` and `[CAL]` help bar text.

**Alternative considered:** A single `propModeDateEdit` mode that checks `dateEdit.mode` internally. Rejected because help bar needs to distinguish the two modes externally.

### 4. Pre-fill with today's date when empty

If the current date value is empty/nil, initialize the `dateEdit` with `time.Now()` (local time, truncated to date). The user sees today's date as a starting point but must confirm with Enter.

**Rationale:** Empty segments would be confusing. Today is the most useful default.

### 5. Segment increment with carry

Incrementing month past 12 carries to next year (and vice versa). Incrementing day past the month's last day carries to next month. This keeps navigation smooth across boundaries.

**Rationale:** Users expect up/down arrows to cycle naturally without hitting walls.

### 6. Calendar rendering as inline overlay

The calendar renders in the same space as the properties panel or table cell, using the existing `widget` layer/compositor pattern for overlay. The calendar grid is approximately 22 chars wide × 8 lines tall (header + 6 week rows).

**Rationale:** Consistent with how the relation picker and select picker render inline. No popup/modal needed.

## Risks / Trade-offs

- **Calendar size vs. table cell**: In table view, cells are narrow. The calendar will overflow the cell width → render as an overlay positioned at the cell location, similar to how select pickers work. Mitigation: use the compositor overlay pattern.
- **Day-of-week locale**: Using English abbreviations (`Mon`, `Tue`, etc.) hardcoded. Acceptable since the entire TUI is English.
- **Leap year edge case**: When navigating from Feb 29 to a non-leap year, clamp to Feb 28. Standard `time` package handles this naturally.
