---
title: tmd (TUI)
description: Interactive terminal interface for browsing your vault.
sidebar:
  order: 1
---

Launches the TUI interactive interface — a lazygit-style two-panel layout for browsing objects.

```bash
tmd
tmd --vault /path/to/vault
```

## Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Move up (navigate list / scroll detail) |
| `↓` / `j` | Move down (navigate list / scroll detail) |
| `Enter` / `Space` | Select object / Toggle group |
| `Tab` | Switch focus between panels |
| `/` | Enter search mode |
| `Esc` | Exit search / Clear results |
| `q` / `Ctrl+C` | Quit |

## Auto-refresh

The TUI watches the `objects/` directory via fsnotify. When files are created, modified, or deleted, it automatically syncs the database and refreshes the view (200ms debounce), preserving the current selection when possible.
