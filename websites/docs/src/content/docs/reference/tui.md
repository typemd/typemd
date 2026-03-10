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
| **Object list** (left) | Groups objects by type. Each group header shows type emoji (if defined), type name, and object count (e.g. `▼ 📚 book (3)`). |
| **Title** (top-right) | Shows the selected object's type emoji, type name, and display name (e.g. `📖 book · Clean Code`). Hidden when no object is selected. |
| **Body** (middle-right) | Displays pinned properties (if any) at the top, followed by the object's markdown body content. |
| **Properties** (right) | Shows schema properties, relations, and wiki-link backlinks. Hidden by default; toggle with `p`. Auto-hides on narrow terminals (< 56 columns). |

The title panel spans the full width of the right side (body + properties area).

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

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up (navigate list / scroll detail) |
| `↓` / `j` | Move down (navigate list / scroll detail) |
| `Enter` / `Space` | Select object / Toggle group |
| `Tab` | Cycle focus between panels |
| `e` | Enter edit mode (body or properties panel) |
| `/` | Enter search mode |
| `Esc` | Exit edit mode / Exit search / Clear results |
| `p` | Toggle properties panel |
| `w` | Toggle soft wrap |
| `[` / `]` | Shrink / grow focused panel |
| `?` / `h` | Open help popup |
| `q` / `Ctrl+C` | Quit |

## Edit Mode

Pressing `e` when the **body panel** is focused switches to in-place edit mode: the read-only viewport is replaced with a textarea pre-filled with the current object's body content. You can type, delete, and navigate freely using standard text-editing keys. The status bar shows `[EDIT]` and the panel border changes to orange.

All global navigation keys (`j`/`k`, `Tab`) are intercepted while editing — they do not switch panels. Press `Esc` to exit edit mode: changes are **auto-saved** to disk and the view returns to `[VIEW]`.

If the file was modified externally between when it was loaded and when you saved, a **conflict warning** appears in the status bar. Press `y` to overwrite, `n` to reload from disk, or `Esc` to cancel.

Pressing `e` on the properties panel also enters edit mode (visual indicator only; property editing is not yet supported).

## Read-Only Mode

Launch with `--readonly` to prevent any edits. In this mode:

- The `e` key is disabled — edit mode cannot be entered
- No write operations are performed to object files or the SQLite index
- The status bar shows `[READONLY]` instead of `[VIEW]`
- The help popup hides edit-related keybindings

This is useful for safely browsing a vault in shared terminals, presentations, or automated contexts.

## Auto-refresh

The TUI watches the `objects/` directory via fsnotify. When files are created, modified, or deleted, it automatically syncs the database and refreshes the view (200ms debounce), preserving the current selection when possible.
