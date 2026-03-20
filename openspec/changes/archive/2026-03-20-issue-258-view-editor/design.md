## Context

The TUI view system currently supports displaying objects filtered/sorted/grouped by a ViewConfig, but editing these settings requires manually editing YAML files. The type editor can create views with default settings and open them in view mode (full-width display), but there's no inline editor.

The view mode already has a split panel mechanism for preview (`p` key): left side shows the table, right side shows object preview using `lipgloss.JoinHorizontal`. The view editor will reuse this same pattern.

Current `ViewConfig.GroupBy` is a `string` (single property), which limits grouping to one level.

## Goals / Non-Goals

**Goals:**
- Add `GroupRule` struct and change `GroupBy` from `string` to `[]GroupRule`
- Backward-compatible YAML loading (auto-migrate `group_by: "genre"` → `group_by: [{property: genre}]`)
- Add `viewEditor` sub-model for inline editing of filter, sort, and group rules
- Integrate view editor as right split panel in view mode (reuse preview split pattern)
- Live reload: editing rules re-queries objects and updates the left table

**Non-Goals:**
- Layout editing (only `list` layout exists; future layouts are out of scope)
- View renaming (can delete and recreate)
- Drag-and-drop reordering of rules (use move up/down keys instead)
- Advanced filter value input (e.g., date pickers, relation pickers); text input is sufficient for now

## Decisions

### 1. GroupRule struct design

**Decision:** Minimal `GroupRule` with only `Property` field.

```go
type GroupRule struct {
    Property string `yaml:"property"`
}
```

**Alternative considered:** Adding `Direction` (asc/desc) to GroupRule. Rejected because group ordering follows the natural order of property values, and custom sort within groups can be achieved via SortRule. Keep it simple now, extend later if needed.

### 2. YAML backward compatibility for GroupBy

**Decision:** Custom `UnmarshalYAML` on `ViewConfig` that detects string vs array format.

When loading YAML:
- `group_by: "genre"` → `[]GroupRule{{Property: "genre"}}` (legacy migration)
- `group_by: [{property: genre}]` → `[]GroupRule{{Property: "genre"}}` (new format)
- `group_by:` absent → `nil` (no grouping)

When saving, always write the new array format. This is a one-way migration: once saved, files use the new format.

**Alternative considered:** Separate migration step on vault open. Rejected because lazy migration on load is simpler and doesn't require scanning all view files upfront.

### 3. View editor as split panel (not overlay)

**Decision:** Reuse the existing preview split pattern in view mode. When the editor is open, it replaces the preview panel on the right side.

```
┌────────────────────────────┬────────────────────────┐
│ Table (60%)                │ View Editor (40%)      │
│ (live query results)       │ Filter / Sort / Group  │
└────────────────────────────┴────────────────────────┘
```

**Rationale:** The preview split already works with `lipgloss.JoinHorizontal`. An overlay popup would require building a new dim-background overlay system in the terminal, which is complex and doesn't provide live preview of changes.

**Interaction:** `e` opens editor (closes preview if open), `Esc` closes editor. Preview (`p`) and editor (`e`) are mutually exclusive — both use the right panel.

### 4. View editor sub-model structure

**Decision:** New `viewEditor` struct following the `templateEditor` pattern.

```go
type viewEditor struct {
    typeName string
    viewName string
    view     *core.ViewConfig
    schema   *core.TypeSchema
    vault    *core.Vault

    // Section & cursor
    section  veSection    // filter, sort, group
    cursor   int          // within current section

    // Editing state
    mode     veMode       // view, addRule, editRule, deleteRule

    // Input fields for rule editing
    propInput    textinput.Model  // property name picker
    opInput      textinput.Model  // operator picker (filter only)
    valueInput   textinput.Model  // value input (filter only)
    dirInput     string           // "asc"/"desc" toggle (sort only)

    // Layout
    width, height int
    scroll        int
}
```

**Sections:** Three sections (filter, sort, group) displayed vertically. Tab cycles between sections. Within each section, ↑↓ navigates rules.

### 5. Picker design

**Property picker:** Text input + scrollable filtered list (hybrid). User types to filter, ↑/↓ to navigate the list, Enter to select. Shows all schema properties including system properties (name, description, tags, created_at, updated_at).

**Operator picker:** Pure scrollable list (no text input). Shows only operators valid for the selected property type from `validOperators` registry. ↑/↓ to navigate, Enter to select.

**Rationale:** Properties can be numerous, so text filtering helps. Operators are few (max 8), so a simple list is sufficient.

### 6. Rule editing flow

**Decision:** Inline step-by-step editing within the editor panel.

Adding a filter rule:
1. Move cursor to "+ Add Filter", press Enter
2. Property picker appears (text + list hybrid)
3. Operator picker appears (scrollable list, type-aware)
4. Value input appears (text input, Enter to confirm; skipped for is_empty/is_not_empty)
5. Rule added, view re-queries

Adding a sort rule:
1. Move cursor to "+ Add Sort", press Enter
2. Property picker appears
3. Direction toggle (asc/desc), Enter to confirm
4. Rule added, view re-queries

Adding a group rule:
1. Move cursor to "+ Add Group", press Enter
2. Property picker appears, Enter to confirm
3. Rule added, view re-queries

**Edit:** Enter on existing rule opens the same flow, pre-populated with current values.

**Delete:** `x` or `d` on a rule removes it immediately with re-query.

**Move:** Shift+K / Shift+J to move a rule up/down within its section. Auto-saves and re-queries.

### 6. Save behavior

**Decision:** Auto-save on every change. Each rule add/edit/delete immediately persists to YAML and re-queries objects.

**Rationale:** Follows the template editor pattern (saves on each edit). Provides instant feedback in the left table. No separate "Save" button needed.

### 7. Multi-level grouping in buildGroups

**Decision:** Nested iteration producing flat group labels with separator.

For `group_by: [{property: genre}, {property: status}]`:
- Group labels become compound: "sci-fi · reading", "sci-fi · done", "drama · reading"
- Implementation: iterate first group property, then sub-group within each group by second property
- Group headers display the compound label

This keeps the existing flat `[]viewGroup` structure without introducing nested tree structures.

## Risks / Trade-offs

**[Multi-level grouping complexity]** → Keep GroupRule minimal (property only). Compound group labels with "·" separator are simple to implement and display. Deeply nested grouping (3+ levels) may produce long labels — mitigate by truncating label display.

**[YAML breaking change]** → Lazy migration on load + always write new format. Existing single-string `group_by` files auto-convert on first load+save. Users who share vault files between old/new versions will see format differences — acceptable since typemd is pre-1.0.

**[Editor panel width on small terminals]** → Minimum width check: if terminal width < 80, editor takes full width (table hidden). Otherwise 60/40 split.

**[Property picker UX]** → Text-based filtering of property names is simple but may be awkward for types with many properties. Acceptable for v1; can add scrollable list picker later.
