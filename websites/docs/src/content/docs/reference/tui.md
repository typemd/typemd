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
```

The properties panel is hidden by default and can be toggled with `p`. On narrow terminals (< 56 columns), it auto-hides.

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up (navigate list / scroll detail) |
| `↓` / `j` | Move down (navigate list / scroll detail) |
| `Enter` / `Space` | Select object / Toggle group |
| `Tab` | Cycle focus between panels |
| `/` | Enter search mode |
| `Esc` | Exit search / Clear results |
| `p` | Toggle properties panel |
| `w` | Toggle soft wrap |
| `[` / `]` | Shrink / grow focused panel |
| `q` / `Ctrl+C` | Quit |

## Auto-refresh

The TUI watches the `objects/` directory via fsnotify. When files are created, modified, or deleted, it automatically syncs the database and refreshes the view (200ms debounce), preserving the current selection when possible.
