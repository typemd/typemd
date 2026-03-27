## Context

The TUI property editor (`propEditor`) supports inline editing for all scalar property types (string, number, date, datetime, url, checkbox, select, multi_select) and the `description` system property. Relation properties — including forward relations, the `tags` system property, reverse relations, and backlinks — are currently displayed read-only. The `isPropertyEditable()` function explicitly skips them.

The existing picker patterns (select/multi_select) use a pre-loaded options list from the schema. Relations require dynamic loading — candidate objects must be fetched from the index at picker activation time, filtered by the relation's target type.

## Goals / Non-Goals

**Goals:**
- Enable editing of forward relation properties (single and multi-value) in the properties panel
- Enable editing of the `tags` system property via the same picker mechanism
- Provide a fuzzy-search text filter for finding target objects quickly
- Support clearing relation values (set to empty / deselect all)
- Reuse existing `ObjectService.Link/Unlink` for all relation mutations
- Follow established propEditor patterns for consistent UX

**Non-Goals:**
- Inline object creation from the picker (future enhancement)
- Editing reverse relations or backlinks (inherently read-only)
- Relation editing in table view cell editor (separate follow-up)
- Full-text body search for candidates (name substring is sufficient)

## Decisions

### 1. Picker as inline overlay, not popup

**Decision:** Render the relation picker inline within the properties panel area, similar to how select/multi_select pickers work — not as a centered popup overlay.

**Rationale:** The existing select picker renders a list of options inline below the property label. This keeps the context visible (other properties, body panel) and follows the established pattern. A popup overlay would be visually inconsistent with other property edit modes.

**Alternative considered:** `widget.OverlayPopup` for a floating picker. Rejected because it would introduce a different interaction model from other property types and require additional z-layer management.

### 2. Textinput filter for fuzzy search

**Decision:** Add a `textinput.Model` to the picker that filters candidate objects by substring match on the display name. The filter is case-insensitive.

**Rationale:** The view editor's property picker uses the same pattern (textinput + filtered list). Object counts per type are typically manageable (tens to hundreds), so client-side substring filtering is sufficient without needing FTS5. The pattern is already proven in this codebase.

### 3. Load candidates via QueryService.Query with type filter

**Decision:** At picker activation, query all objects of the target type using `QueryService.Query([]FilterRule{{Property: "type", Operator: "is", Value: targetType}})`. For relations without a target type constraint, query all objects.

**Rationale:** This leverages the existing SQLite index and returns `[]*Object` with all properties needed for display. No new core API methods are needed.

### 4. Single-select vs multi-select based on Multiple flag

**Decision:** If `relProp.Multiple == true`, use a multi-picker (Space to toggle, Enter to confirm, like multi_select). If `Multiple == false`, use single-select (Enter selects and confirms immediately). The `tags` system property is always `Multiple: true`.

**Rationale:** Matches the existing select/multi_select interaction patterns exactly. Users already know these keybindings.

### 5. Relation mutation via Link/Unlink, not direct property write

**Decision:** Apply relation changes through `Vault.LinkObjects()` / `Vault.UnlinkObjects()` rather than directly writing to `obj.Properties`. This is done via tea.Cmd returning messages, similar to how `applyPropertyValue` works for scalar edits.

**Rationale:** `LinkObjects` handles bidirectional inverse relations, type validation, duplicate detection, and index updates. Direct property writes would bypass these safeguards. The TUI already uses async commands for save operations.

### 6. Two new propEditor modes

**Decision:** Add `propModeRelationPick` (single-select relation) and `propModeRelationMultiPick` (multi-select relation) to the existing `propEditMode` enum.

**Rationale:** Mirrors the existing `propModeSelectPick` / `propModeMultiPick` separation, keeping the mode state machine clear and each mode's update/render logic focused.

### 7. Display name with type prefix for untyped relations

**Decision:** Picker items show `name` (stripped of ULID) for typed relations. For relations without a target type constraint, show `type/name` to disambiguate across types.

**Rationale:** When the target type is known, repeating it in every item is redundant. When target type is unconstrained, the type prefix provides necessary context.

### 8. "(none)" option for clearing single-value relations

**Decision:** Single-select relation picker includes a "(none)" option at the top of the list. Selecting it removes the relation. For multi-select, deselecting all items and confirming clears the relation.

**Rationale:** Users need a way to clear relations without switching to a different editing mode. The "(none)" option is a common UI pattern for nullable fields.

## Risks / Trade-offs

- **[Large candidate lists]** → Types with thousands of objects may produce a long picker list. Mitigation: substring filter reduces visible items; scrolling viewport limits rendering. Acceptable for v0.7.0; virtual scrolling can be added later if needed.
- **[Async link/unlink latency]** → `LinkObjects` writes to both file and index, which takes I/O time. Mitigation: the TUI already handles async saves; a brief delay between confirm and visual update is acceptable and consistent with scalar edits.
- **[Bidirectional relation side effects]** → Linking object A to B may create an inverse link on B. The file watcher will detect B's change and trigger a refresh. Mitigation: the existing file watcher debounce (200ms) already handles this.
