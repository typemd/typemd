---
title: tmd (TUI)
description: Interactive terminal interface for browsing your vault.
sidebar:
  order: 1
---

Launches the TUI interactive interface — a three-panel layout for browsing objects (list, body, and properties).

```bash
tmd
tmd --vault /path/to/vault
tmd --readonly
```

The properties panel is hidden by default and can be toggled with `p`. On narrow terminals (< 56 columns), it auto-hides. The panel displays schema properties, relations, and wiki-link backlinks.

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
