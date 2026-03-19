---
title: tmd (TUI)
description: Interactive terminal interface for browsing your vault.
sidebar:
  order: 1
---

Launches the TUI interactive interface for browsing objects.

```bash
tmd
tmd --vault /path/to/vault
tmd --readonly
tmd --reindex
```

## Layout

The TUI uses a multi-panel layout:

| Panel | Description |
|-------|-------------|
| **Object list** (left) | Groups objects by type. Each group header shows type emoji (if defined), type plural name, and object count (e.g. `▼ 📚 books (3)`). All defined types appear even if they have no objects. A `+ New Type` row appears at the bottom. |
| **Title** (top-right) | Shows the selected object's type emoji, type name, and display name (e.g. `📖 book · Clean Code`). When a type header is selected, shows the type name instead. Hidden when nothing is selected. |
| **Body** (middle-right) | Displays pinned properties (if any) at the top, followed by the object's markdown body content. When a type header is selected, replaced by the **type editor**. |
| **Properties** (right) | Shows schema properties, relations, and wiki-link backlinks. Hidden by default; toggle with `p`. Auto-hides on narrow terminals (< 56 columns). Hidden when the type editor or template editor is active. |

The title panel spans the full width of the right side (body + properties area).

## Type Editor

Moving the cursor to a type group header automatically opens the **type editor** in the right panel. The type editor lets you manage the type schema without editing YAML files:

- **Meta fields** — Edit plural name, emoji, and unique constraint. Name is read-only (type rename not supported).
- **Properties** — Split into Pinned (Header) and Properties sections. Add new properties via a multi-step wizard, edit emoji, delete, or reorder with move mode.
- **Pin toggle** — Press `p` on a property to move it between Pinned and Properties sections.
- **Property detail** — Press `Enter` on a property to open a popup for editing metadata (emoji).
- **Delete type** — Press `D` (shift+d) to delete the type (with confirmation). The built-in `tag` type cannot be deleted.
- **Templates** — A Templates section lists available templates for the type. Press `Enter` on a template to open the **template editor**. Press `Enter` on `+ Add Template` to create a new template by name.

Changes are saved immediately on each operation (no explicit save step).

| Key | Action (in type editor) |
|-----|-------------------------|
| `e` | Edit meta field (Plural/Emoji: text input; Unique: toggle) |
| `Enter` | Open property detail popup / Open template / Start add wizard |
| `a` | Add property (starts wizard) |
| `d` | Delete property (with confirmation) |
| `D` | Delete type (with confirmation) |
| `m` | Enter move mode (`↑`/`↓` to reorder, `Enter`/`Esc` to exit) |
| `p` | Toggle pin (move property between Pinned/Properties sections) |
| `Esc` | Return focus to sidebar |

## Pinned Properties

Type schemas can mark properties with a `pin` value (positive integer) to highlight them at the top of the body panel. Pinned properties are sorted by pin value (lower = higher priority) and displayed as key-value lines with a separator before the body content. Properties with an emoji show the emoji alongside the value.

Pinned properties are **excluded** from the Properties panel to avoid duplication.

```yaml
# .typemd/types/book.yaml
properties:
  - name: status
    type: select
    emoji: 📋
    pin: 1        # displayed first in body panel
  - name: rating
    type: number
    emoji: ⭐
    pin: 2        # displayed second
  - name: title
    type: string  # no pin — shown in Properties panel
```

## Modes and Keybindings

The TUI operates in different modes depending on context. The current mode is shown in `[brackets]` in the status bar at the bottom. Each mode has its own set of active keybindings.

### `[VIEW]` — Normal Navigation

The default mode. Browse objects and types in the sidebar.

| Key | Action |
|-----|--------|
| `↑`/`k`, `↓`/`j` | Navigate list / scroll content |
| `Enter` | Select object / Focus type editor |
| `Space` | Toggle group expand/collapse |
| `Tab` | Cycle focus between panels |
| `v` | Open view mode for the current type |
| `n` | Create new object & edit body (in current type group) |
| `N` | Quick create — batch mode (in current type group) |
| `e` | Enter edit mode (when body panel focused) |
| `/` | Enter search mode |
| `p` | Toggle properties panel |
| `w` | Toggle soft wrap |
| `[`/`]` | Shrink/grow focused panel |
| `?`/`h` | Open help popup |
| `q`/`Ctrl+C` | Quit |

### View Mode

Press `v` on a type group header or object to enter **view mode**. This replaces the three-panel layout with a full-width display of objects. View mode supports two layouts:

- **List layout** — Each row shows the object's emoji and name, optionally followed by inline property values separated by " · ". This is the default layout for new views.
- **Table layout** — A columnar display with a NAME column followed by property columns. Column headers include sort indicators (↑/↓) when active. The number of columns adjusts to the terminal width.

| Key | Action |
|-----|--------|
| `↑`/`k`, `↓`/`j` | Navigate rows |
| `Enter` | Open object detail (three-panel view) |
| `e` | Open view editor (right split panel) |
| `p` | Toggle preview panel (right side) |
| `Esc` | Return to sidebar |
| `q`/`Ctrl+C` | Quit |

When **preview** is active (`p`), the right side shows the selected object's properties and body content. The preview follows the cursor as you navigate. Preview and view editor are mutually exclusive.

If the type has multiple saved [views](/advanced/file-structure#views), a selection popup appears when pressing `v`. If only the default view exists, it opens directly.

The type editor also shows a **Views** section where you can browse saved views and create new ones with `+ Add View`.

#### View Editor

Press `e` in view mode to open the **view editor** as a right split panel. The editor has five sections:

| Section | Description |
|---------|-------------|
| **Layout** | Toggle between `list` and `table` with `Enter`. |
| **Columns** | Add, remove, and reorder properties to display. In list layout, selected properties appear as inline values. In table layout, they become column headers. |
| **Filter** | Add filter rules (property + operator + value) to narrow displayed objects. |
| **Sort** | Define sort rules (property + direction) controlling row order. |
| **Group By** | Group objects by one or more properties. |

Changes are auto-saved on every edit. Press `Esc` to close the editor and return to the view list.

### `[TYPE]` — Type Editor

Active when the type editor panel has focus. Press `Tab` or `Enter` on a type header to enter this mode.

| Key | Action |
|-----|--------|
| `↑`/`k`, `↓`/`j` | Navigate meta fields and properties |
| `Enter` | Open property detail popup / Start add wizard on `+ Add Property` |
| `e` | Edit meta field (Plural/Emoji: text input, Unique: toggle) |
| `a` | Add property (starts wizard) |
| `d` | Delete property (with confirmation) |
| `D` | Delete type (with confirmation) |
| `m` | Enter move mode |
| `p` | Toggle pin on property |
| `Tab`/`Esc` | Return focus to sidebar |
| `q`/`Ctrl+C` | Quit |

### `[EDIT]` — Body Editing

Active when editing the object body text. Press `e` in VIEW mode while the body panel is focused.

| Key | Action |
|-----|--------|
| Standard text keys | Edit body content |
| `Esc` | Exit edit mode (auto-saves if changed) |

All navigation keys (`j`/`k`, `Tab`) are intercepted and do not switch panels. The panel border changes to orange.

If the file was modified externally while editing, a **`[CONFLICT]`** warning appears. Press `y` to overwrite, `n` to reload from disk, or `Esc` to dismiss.

### `[MOVE]` — Property Reorder

Active within the type editor when reordering properties. Press `m` on a property in TYPE mode.

| Key | Action |
|-----|--------|
| `↑`/`k`, `↓`/`j` | Move property up/down (swaps with neighbor) |
| `Enter`/`Esc` | Confirm and save new order |

Moving a property across the Pinned/Properties boundary automatically toggles its pin value.

### `[ADD PROPERTY]` — Add Property Wizard

A multi-step wizard for adding a new property to a type. Press `a` in TYPE mode or `Enter` on `+ Add Property`.

| Step | Input | Keys |
|------|-------|------|
| 1. Name | Text input for property name | `Enter`: next, `Esc`: cancel |
| 2. Type | Select from list (string, number, date, ...) | `↑`/`↓`: select, `Enter`: next, `Esc`: back |
| 2b. Options | Comma-separated values (for select/multi_select) | `Enter`: create, `Esc`: back |
| 3. Relation | Target type, multiple, bidirectional, inverse name | `Tab`: next field, `Enter`/`Space`: toggle/confirm, `Esc`: back |

### `[PROPERTY]` — Property Detail Popup

A popup overlay for editing property metadata. Press `Enter` on a property in TYPE mode.

| Key | Action |
|-----|--------|
| `Enter`/`e` | Edit the selected field |
| `Esc` | Close popup |

When editing a field: `Enter` saves, `Esc` cancels (reverts).

### `[TEMPLATE]` — Template Editor

Active when viewing or editing a template. Press `Enter` on a template in the type editor to enter this mode.

| Key | Action |
|-----|--------|
| `e` | Edit template body (enters textarea edit mode) |
| `d` | Delete template (with confirmation) |
| `Tab` | Switch focus between body and properties panels |
| `↑`/`k`, `↓`/`j` | Scroll body (body focused) / Navigate properties (props focused) |
| `Enter` | Edit property value (when props focused) |
| `Esc` | Return to type editor |
| `q`/`Ctrl+C` | Quit |

When editing the body: `Esc` saves, `Ctrl+C` cancels. When editing a property: `Enter` confirms, `Esc` cancels. Clearing a property value removes it from the template frontmatter.

### `[NEW TEMPLATE]` — New Template Input

Shown when creating a template via `+ Add Template` in the type editor.

| Key | Action |
|-----|--------|
| Text keys | Template name |
| `Enter` | Create empty template |
| `Esc` | Cancel |

### `[DELETE]` / `[DELETE TYPE]` — Deletion Confirmation

Shown when deleting a property (`d`), type (`D`), or template (`d` in template editor).

| Key | Action |
|-----|--------|
| `y` | Confirm deletion |
| `n`/`Esc` | Cancel |

### `[NEW TYPE]` — New Type Creation

Shown when creating a new type via `+ New Type`. The **title panel** transforms into an inline creation form with three fields: emoji (optional), name (required), and plural (optional). The right panel shows a **live preview** of the type schema being created.

| Key | Action |
|-----|--------|
| Text keys | Type in the focused field |
| `Tab` | Cycle focus: name → plural → emoji |
| `Enter` | Create type and open type editor |
| `Esc` | Cancel |

### `[NEW OBJECT]` — Create & Edit Mode

Triggered by pressing `n`. The **title panel** transforms into an inline creation form with a name input and (when templates exist) a template cycling selector. The body and properties panels show a **live preview** of the selected template's content.

The form layout: `📚 book · [name█] 📝 review`

| Key | Action |
|-----|--------|
| Text keys | Object name |
| `Tab` | Switch focus between name and template fields |
| `←`/`→` | Cycle templates (when template field focused) |
| `Enter` | Create object & enter body edit mode |
| `Esc` | Cancel |

- If the type has **multiple templates**, `Tab` switches to the template selector where `←`/`→` cycles through templates plus a `(none)` option. A single template is auto-selected and shown as a static label.
- If the type defines a **name template** (e.g. `{{ date:YYYY-MM-DD }}`), the name is pre-filled with the evaluated value. You can edit it or press Enter to accept.
- Switching templates updates the body and properties panels in real time as a preview.
- Duplicate name errors (for types with `unique: true`) are shown inline in the title panel and clear when you modify the text.

### `[QUICK CREATE]` — Batch Creation Mode

Triggered by pressing `N`. Same title panel form as Create & Edit, but stays in input mode for rapid creation of multiple objects.

The template (if selected) persists across the batch — all objects use the same template. Name templates pre-fill but are always editable.

| Key | Action |
|-----|--------|
| Text keys | Object name |
| `Tab` | Switch focus between name and template fields |
| `←`/`→` | Cycle templates (when template field focused) |
| `Enter` | Create object, clear input, ready for next |
| `Esc` | Exit batch mode (selects last created object) |

A success flash (e.g. `✓ Created: my-book`) appears briefly in the title panel after each creation.

### `[READONLY]` — Read-Only Mode

Active when launched with `--readonly`. The `e`, `n`, and `N` keys are disabled, no write operations are performed, and the help popup hides edit-related keybindings.

## Session State

The TUI automatically saves your session state to `.typemd/tui-state.yaml` when you quit. On next launch, it restores:

- **Selected object or type** — cursor returns to the same object (by ID) or type header (by name)
- **Expanded groups** — type groups stay expanded/collapsed as you left them
- **Panel dimensions** — left panel and properties panel widths
- **Properties visibility** — whether the properties panel was shown
- **Scroll offset** — vertical scroll position in the object list
- **View mode** — if you were in a view when you quit, the TUI restores that view on next launch, including cursor position, scroll offset, and expanded groups within the view

If the previously selected object was deleted, the TUI falls back to the first object in the same type group, then to the first object overall. Focus always starts on the sidebar for a consistent experience, except when restoring view mode (which sets focus to the body panel).

If the saved view's type or view has been deleted, the TUI falls back to the default view for that type, or to the normal sidebar mode if no views exist.

Search state is not persisted — each session starts with a fresh view.

If the state file is missing or corrupt, the TUI silently falls back to default startup behavior.

## Reindex

The `--reindex` flag forces a full sync of the `objects/` directory to the database, cleans up orphaned relations, and rebuilds the full-text search index. It is a global flag that can be combined with any command.

> **Note:** When opening a vault, TypeMD automatically syncs the index if it is empty or missing. You only need `--reindex` when files have been edited while the vault was not open.

```bash
# Reindex and launch TUI
tmd --reindex

# Reindex and start MCP server
tmd mcp --reindex

# Reindex and list objects
tmd object list --reindex
```

### Orphaned relation cleanup

When an object is deleted from disk, any relations pointing to or from that object become orphaned. During reindex, these dangling references are automatically detected and removed from the index. A warning is displayed listing the affected relations:

```
Warning: Found 2 orphaned relation(s):
  book/golang-in-action -> person/deleted-author (relation: "author")
  person/deleted-author -> book/golang-in-action (relation: "books")
Orphaned relations have been removed from the index.
```

> **Note:** This only cleans up the SQLite index. The frontmatter in `.md` files is not modified — use `tmd relation unlink --both` to properly remove relations before deleting objects.

## Auto-refresh

The TUI watches the `objects/` directory via fsnotify. When files are created, modified, or deleted, it automatically syncs the database and refreshes the view (200ms debounce), preserving the current selection when possible.
