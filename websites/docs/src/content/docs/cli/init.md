---
title: tmd init
description: Initialize a new vault.
sidebar:
  order: 2
---

Initializes a new vault. Creates `.typemd/` directory structure, SQLite database, and a `.typemd/.gitignore` that excludes `index.db`.

After creating the vault structure, an interactive starter type selector is displayed. Starter types are common type schemas (idea, note, book) that help you get started quickly. All are selected by default — use the keyboard to customize your selection:

- **↑↓** or **j/k** — move cursor
- **Space** — toggle selection
- **a** — select all
- **n** — deselect all
- **Enter** — confirm selection
- **Esc** — skip (select none)

Selected starter types are written as regular `.typemd/types/*.yaml` files, fully owned and editable.

```bash
tmd init
```

### Flags

| Flag | Description |
|------|-------------|
| `--no-starters` | Skip the starter type selector (empty vault) |

```bash
# Non-interactive: skip starter types
tmd init --no-starters
```

### Vault configuration

When starter types are selected, `tmd init` also creates `.typemd/config.yaml` with a default type for quick object creation:

- If **idea** is selected → `cli.default_type: idea`
- If **note** is selected (without idea) → `cli.default_type: note`
- If only **book** is selected → no config file is created

This enables `tmd object create "Some Thought"` without specifying a type. See [tmd object create](/cli/create) for details.

Running `tmd init` on an already-initialized vault will return an error.
