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
```

## Layout

The TUI uses a multi-panel layout:

| Panel | Description |
|-------|-------------|
| **Object list** (left) | Groups objects by type. Each group header shows type emoji (if defined), type plural name, and object count (e.g. `▼ 📚 books (3)`). All defined types appear even if they have no objects. A `+ New Type` row appears at the bottom. |
| **Title** (top-right) | Shows the selected object's type emoji, type name, and display name (e.g. `📖 book · Clean Code`). When a type header is selected, shows the type name instead. Hidden when nothing is selected. |
| **Body** (middle-right) | Displays pinned properties (if any) at the top, followed by the object's markdown body content with markdown rendering (syntax markers hidden, styled text displayed). When a type header is selected, replaced by the **type editor**. |
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
# types/book/schema.yaml
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
| `r` | Rename object (inline in title panel) |
| `/` | Enter search mode |
| `p` | Toggle properties panel |
| `w` | Toggle soft wrap |
| `-`/`=` | Shrink/grow focused panel |
| `.` | Toggle focus mode (full-width body) |
| `g` | AI generate (describe / tag suggestions) |
| `Ctrl+E` | AI schema explore |
| `,` | Open config settings page |
| `?`/`h` | Open help popup |
| `q`/`Ctrl+C` | Quit |

### AI Assist

Requires `ai.enabled: true` in `.typemd/config.yaml` and the `claude` CLI installed and authenticated.

- **`g`** — Opens an AI action picker popup with two options:
  - **Generate description** — AI analyzes the object's name, properties, and body to generate a description. Shown as ghost text preview; press `Tab` to accept or `Esc` to reject.
  - **Suggest tags** — AI suggests relevant tags based on object content and existing tags. Shows a selectable popup with checkboxes; `Space` to toggle, `Enter` to apply, `Esc` to cancel. Tags are classified as existing or new (★).
- **`Ctrl+E`** — Enters schema explore mode (from sidebar). Select a type, and AI analyzes its objects to suggest schema improvements (add/modify/remove properties). Press `Enter` to accept or `s` to skip each suggestion.

### Toast Notifications

The TUI shows transient toast notifications in the bottom-right corner for events like sync warnings (unresolved references), AI errors, and property validation errors. Toasts auto-dismiss after 3 seconds (configurable via `tui.toast.duration_ms`). Press `Esc` to dismiss manually — while a toast is visible, `Esc` is consumed by the toast and does not propagate to other panels. See [Configuration](/basics/configuration/) for all toast settings.

### View Mode

Press `v` on a type group header or object to enter **view mode**. This replaces the three-panel layout with a full-width display of objects. View mode supports two layouts:

- **List layout** — Each row shows the object's emoji and name, optionally followed by inline property values separated by " · ". This is the default layout for new views.
- **Table layout** — A columnar display with a NAME column followed by property columns. Column headers include sort indicators (↑/↓) when active. The number of columns adjusts to the terminal width. Supports inline cell editing with crosshair highlighting.

**Common keybindings (both layouts):**

| Key | Action |
|-----|--------|
| `↑`/`k`, `↓`/`j` | Navigate rows |
| `e` | Open view editor (right split panel) |
| `p` | Toggle preview panel (right side) |
| `Esc` | Return to sidebar |
| `q`/`Ctrl+C` | Quit |

**List layout:**

| Key | Action |
|-----|--------|
| `Enter`/`Space` | Open object detail / Toggle group expand/collapse |

**Table layout** adds cell navigation and inline editing:

| Key | Action |
|-----|--------|
| `←`/`h`, `→`/`l` | Navigate columns |
| `Enter`/`Space` | Edit focused cell (toggle checkbox directly) / Toggle group expand/collapse |
| `o` | Open object detail view |
| `Tab` | Move to next editable cell (skips read-only) |
| `Shift+Tab` | Move to previous editable cell |

A crosshair highlight shows the active cell: the cursor row gets a dim background, the cursor column header gets a tint, and the active cell gets a strong highlight. Editing uses type-appropriate widgets: textinput for string/number/datetime/url, date picker for date (segmented input or inline calendar, toggle with `c`), option picker for select, multi-picker for multi\_select, and direct toggle for checkbox. Read-only columns (relations, created\_at, updated\_at) are navigable but not editable. Edits auto-save on confirm (`Enter`); press `Esc` to cancel.

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

| Key | Action |
|-----|--------|
| `Tab` | Next section |
| `Shift+Tab` | Previous section |
| `Enter` | Edit field / Toggle layout / Add rule |
| `x`/`d` | Delete selected rule |
| `K` | Move selected rule up |
| `J` | Move selected rule down |
| `D` | Delete entire view (with confirmation) |
| `Esc` | Close editor and return to view list |

Changes are auto-saved on every edit.

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

### `[PROPS]` — Property Navigation and Editing

Active when the properties panel is focused. Press `Tab` from the body panel to enter this mode. A cursor (▸) highlights the current editable property. Read-only properties (`created_at`, `updated_at`, relations, backlinks, tags, local properties) are displayed dimmed and skipped during navigation. Local properties (frontmatter keys not defined in the type schema) appear after a `── Local Properties ──` separator.

| Key | Action |
|-----|--------|
| `↑`/`k`, `↓`/`j` | Navigate between editable properties |
| `Enter` | Activate editing for current property |
| `Space` | Toggle checkbox property (when on checkbox) |
| `Esc` | Return focus to sidebar |
| `Tab` | Switch focus to another panel |

#### Editing by property type

| Property Type | Widget | Behavior |
|---|---|---|
| string, number, datetime, url | Textinput | Pre-filled with current value. `Enter` confirms, `Esc` cancels. |
| date | Date picker | Two modes: **segmented input** (YYYY-MM-DD segments, `←`/`→` to switch segments, `↑`/`↓` to adjust, type digits directly, shows day of week) and **calendar** (mini month grid, `h`/`j`/`k`/`l` to navigate days, `H`/`L` to switch months, `t` to jump to today). Press `c` to toggle modes. `Enter` confirms, `Esc` cancels. |
| checkbox | Toggle | `Enter` or `Space` toggles between ☐ and ☑. Saves immediately. |
| select | Option picker | Shows available options. `↑`/`↓` navigate, `Enter` selects, `Esc` cancels. |
| multi_select | Multi-picker | Shows options with checkboxes. `Space` toggles, `Enter` confirms all, `Esc` cancels. |

**Input validation:** When confirming a text edit, the value is validated:
- **Number** — must be a valid integer or decimal
- **Date** — must be YYYY-MM-DD format
- **Datetime** — must be valid ISO 8601 (e.g. `2024-01-15T10:30:00`)
- **URL** — must start with `http://` or `https://`
- **Select** — must match one of the defined options

If validation fails, a toast notification shows the error and the edit remains active. Press `Esc` to cancel.

Edits are auto-saved on confirm. The panel border changes to orange during active editing, and the help bar shows `[EDIT]` for textinput, `[DATE]` for date segmented input, `[CAL]` for date calendar, or `[PICK]` for picker selection.

### Locked Objects

Objects with `locked: true` in their frontmatter are protected from accidental editing. Locked objects display a 🔒 badge right-aligned in the title panel.

When an object is locked:

- **`Tab`** (focus properties), **`e`** (edit body), and **`r`** (rename) are blocked — a toast notification shows "Object is locked".
- The object body and properties remain fully readable.
- Press **`L`** (uppercase) to toggle the lock state of the selected object.
- Alternatively, use the CLI: `tmd object lock <id>` / `tmd object unlock <id>`.

### Archived Objects

Objects with `archived: true` are hidden from the sidebar and all default queries. They act as soft-deleted — the file remains on disk but the object does not appear in the TUI. To restore an archived object, use the CLI: `tmd object unarchive <id>`.

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

### `[SETTINGS]` — Configuration Settings

Press `,` in VIEW mode to open a full-width config settings page. This lets you browse and edit `.typemd/config.yaml` values without leaving the TUI.

The page has a two-column layout:

| Column | Description |
|--------|-------------|
| **Categories** (left) | General, CLI, TUI, AI, Web — grouped by config key prefix |
| **Settings** (right) | Keys in the selected category with current values (or defaults if unset) |

| Key | Action |
|-----|--------|
| `↑`/`k`, `↓`/`j` | Navigate categories or settings |
| `Tab` | Switch focus between category and settings columns |
| `Enter` | Edit the selected setting (opens popup) |
| `Esc` | Exit config page (return to sidebar) |
| `?`/`h` | Show help |
| `q`/`Ctrl+C` | Quit |

**Editing settings:**

- **String/number settings** — A text input popup pre-filled with the current value. Press `Enter` to save, `Esc` to cancel. Saving an empty string removes the key (restores default).
- **Boolean settings** (e.g. `tui.toast.show_warnings`) — Cycles through `true` → `false` → `unset` with `Enter`/`↑`/`↓`. Press `Esc` to save and close.

Changes are saved immediately to `.typemd/config.yaml`. AI provider map settings (`ai.providers.*`) are not available in this page — use `tmd config` CLI or edit YAML directly.

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

## Theme

The TUI supports configurable colors via the `tui.theme` section in `.typemd/config.yaml`. Missing fields are silently ignored, keeping the defaults.

```yaml
tui:
  theme:
    focus_border: "63"    # focused panel border
    wiki_link: "33"       # wiki-link display text
    heading: "3"          # markdown headings (# through ######)
    bold: ""              # bold text (**text**), default: bold weight only
    italic: ""            # italic text (*text* / _text_), default: italic only
    inline_code: "245"    # inline code (`code`)
    code_block: "245"     # fenced code blocks (```)
    link: "33"            # markdown links ([text](url))
    blockquote: "8"       # blockquote lines (> text)
    hrule: "8"            # horizontal rules (---, ***, ___)
```

Values are ANSI color codes (0–255) or hex colors (`#RGB`, `#RRGGBB`). An empty string means no foreground color — only the text attribute (bold/italic) is applied. You can also use `tmd config set tui.theme.heading 196` to change individual colors.

### Markdown rendering

The body panel renders markdown content in view mode — syntax markers are hidden and styled text is displayed:

- **Headings** (`#` through `######`) — `#` markers hidden, text shown in heading color with bold
- **Bold** (`**text**`) — `**` markers hidden, text shown with bold weight
- **Italic** (`*text*`, `_text_`) — markers hidden, text shown with italic style; intra-word underscores (e.g. `snake_case`) are not treated as italic
- **Inline code** (`` `code` ``) — backticks hidden, text shown in distinct foreground color
- **Fenced code blocks** (` ``` `) — fence lines hidden, content shown in code block color; inline markdown not processed inside
- **Links** (`[text](url)`) — URL and brackets hidden, display text shown in link color
- **Blockquotes** (`> text`) — `>` marker replaced with `│` prefix, text shown in blockquote color
- **Horizontal rules** (`---`, `***`, `___`) — replaced with a styled `────────────────────` line

In edit mode (`e`), the raw markdown source is shown without rendering.

Wiki-links (`[[target]]`) are styled separately using the `wiki_link` color and are not affected by markdown rendering.

## Index sync

The index is automatically synced every time the vault is opened. TypeMD walks all Object files, upserts them into the SQLite index, and cleans up orphaned relations — no manual intervention needed, even after editing files outside of TypeMD.

When an object is deleted from disk, any relations pointing to or from that object become orphaned. These dangling references are automatically detected and removed from the index on the next vault open.

> **Note:** Index sync only updates the SQLite index. The frontmatter in `.md` files is not modified — use `tmd relation unlink --both` to properly remove relations before deleting objects.

## Auto-refresh

The TUI watches the `objects/` directory via fsnotify. When files are created, modified, or deleted, it performs an incremental index sync for the changed files and refreshes the view (200ms debounce), preserving the current selection when possible. The debounce interval can be customized via `tui.debounce_ms` in `.typemd/config.yaml`.

The TUI also watches the `types/` directory for schema changes. When type schemas are modified externally, the schema cache is invalidated and a full refresh is triggered.
