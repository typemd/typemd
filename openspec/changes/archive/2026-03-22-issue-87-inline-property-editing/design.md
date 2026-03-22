## Context

The TUI properties panel (`tui/detail.go:renderProperties`) currently renders `DisplayProperty` entries as read-only text via `propsViewport`. The `editMode` flag exists but only activates body textarea editing — pressing `e` with `focusProps` sets `editMode = true` but has no further behavior. There is no per-property cursor or input widget.

The type system (`core/type_schema.go`) defines 9 property types: string, number, date, datetime, url, checkbox, select, multi_select, and relation. Validators exist in `core/type_schema_validate.go` for date, datetime, url, select, and multi_select. Relation editing is out of scope (handled by #88).

## Goals / Non-Goals

**Goals:**
- Per-property cursor navigation in the properties panel (j/k, skip read-only fields)
- Type-appropriate editing widget for each non-relation property type
- Input validation before accepting edits
- Auto-save on confirm (consistent with body edit behavior)
- Visual feedback: active field highlight, edit border color

**Non-Goals:**
- Relation property editing (separate issue #88)
- Undo/redo history
- Multi-field batch editing
- Editing system properties `name` (edited via title/rename), `created_at`, `updated_at` (immutable)
- Editing `tags` (relation to tag type, handled by #88)

## Decisions

### 1. Property editor as a sub-model (propEditor)

**Decision:** Create a `propEditor` sub-model similar to `typeEditor` and `templateEditor`, rather than inlining state into the main model.

**Rationale:** The main model already has 30+ fields. A sub-model encapsulates cursor position, active input widget, and edit state cleanly. The existing pattern (`typeEditor`, `templateEditor`, `viewEditor`) is well-established.

**Alternative considered:** Adding fields directly to `model` — rejected due to complexity and inconsistency with existing patterns.

### 2. Edit activation: Enter on focused property

**Decision:** When `focusProps` is active, j/k navigates properties. Enter activates inline editing for the current property. This replaces the current `e` key behavior for `focusProps`.

**Rationale:** j/k navigation is already the convention in the sidebar. Enter to activate is consistent with the type editor pattern. The current `e` + `focusProps` behavior (just setting `editMode = true` with no effect) is a placeholder that can be replaced.

**Alternative considered:** Using `e` to enter props edit mode, then j/k + Enter — rejected as an unnecessary extra step.

### 3. Widget mapping by property type

| Property Type | Widget | Behavior |
|---|---|---|
| string | textinput | Pre-filled with current value, Enter confirms, Esc cancels |
| number | textinput | Same as string but validates numeric input on confirm |
| date | textinput | Validates YYYY-MM-DD format on confirm |
| datetime | textinput | Validates ISO 8601 format on confirm |
| url | textinput | Validates http(s):// scheme on confirm |
| checkbox | direct toggle | Enter or Space toggles ☐ ↔ ☑, no textinput needed |
| select | option picker | Shows option list overlay, j/k navigate, Enter selects |
| multi_select | option multi-picker | Shows option list overlay, j/k navigate, Space toggles, Enter confirms |

**Rationale:** Textinput covers most types with type-specific validation. Checkbox and select types need distinct widgets because free-text input is inappropriate for constrained values.

### 4. Editable property filtering

**Decision:** The `propEditor` filters `DisplayProperty` to only include editable fields. A property is editable if:
- It is NOT `name` (edited via rename)
- It is NOT `created_at` or `updated_at` (immutable system properties)
- It is NOT `tags` (relation to tag type, handled by #88)
- It is NOT a relation, reverse relation, or backlink (`IsRelation`, `IsReverse`, `IsBacklink`)
- It is NOT pinned (`Pin > 0`) — pinned properties display in the body panel

Non-editable properties are still displayed but skipped during cursor navigation.

### 5. Validation approach: reuse existing core validators

**Decision:** Expose lightweight validation helpers from `core/` that the TUI can call before accepting input. These wrap the existing `validateDate`, `validateURL`, etc. functions with simpler signatures returning a single error string.

**Rationale:** Validators already exist in `core/type_schema_validate.go`. Duplicating validation logic in `tui/` would violate DRY. A thin adapter layer (`core.ValidatePropertyValue(propType string, options []Option, input string) error`) provides the right abstraction.

### 6. Save flow: immediate save on confirm

**Decision:** When a property edit is confirmed (Enter), update `obj.Properties`, set `dirty = true`, and save immediately (same as body edit exit). No explicit save step.

**Rationale:** Consistent with the existing body edit auto-save behavior. Keeps the UX simple.

### 7. Checkbox display: ☐ / ☑

**Decision:** Use ☐ (U+2610) for unchecked and ☑ (U+2611) for checked in `FormatValue()`.

**Rationale:** Decided during discussion. Traditional checkbox style, renders well in terminal fonts.

### 8. Props panel cursor is always active when focused

**Decision:** When `focusProps` is active, the property cursor is always visible (no separate "navigation mode" vs "browse mode"). The cursor highlights the current property. Enter activates editing, Esc returns to sidebar.

**Rationale:** Simpler mental model. The body panel already uses the same pattern — focus means interactive.

## Risks / Trade-offs

- **[Terminal font compatibility]** ☐/☑ characters may not render in all terminal fonts → Mitigation: These are Unicode 1.1 characters with wide support. Fallback to `[ ]`/`[x]` could be added later if needed.
- **[Select option list overlay]** The option picker needs to overlay on the properties panel without affecting layout → Mitigation: Use the existing `widget.OverlayPopup` pattern or a simple inline list within the property row area.
- **[DisplayProperty mutation]** `DisplayProperty` is currently a value type for display only; editing needs to map back to `obj.Properties` keys → Mitigation: `DisplayProperty.Key` already contains the property name, which maps directly to frontmatter keys.
