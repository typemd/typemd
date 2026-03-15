---
title: Quick Start
description: Get up and running with TypeMD in under a minute.
sidebar:
  order: 3
---

## 1. Initialize a Vault

```bash
tmd init
```

This creates the `.typemd/` directory structure and SQLite database in the current directory. You'll be prompted to select starter types (idea, note, book) — press **Enter** to accept the defaults, or customize your selection with **Space** to toggle and **Esc** to skip.

## 2. Open the TUI

```bash
tmd
tmd --readonly  # read-only mode (no editing)
```

This launches the interactive TUI for browsing your vault — objects are grouped by type in a list on the left, with a title, body, and properties panel on the right.

## 3. Define a Type

Create a type schema to define the structure of your objects:

```yaml
# .typemd/types/book.yaml
name: book
emoji: 📚
properties:
  - name: title
    type: string
  - name: status
    type: select
    options:
      - value: to-read
      - value: reading
      - value: done
  - name: rating
    type: number
```

## 4. Create Your First Object

Use the CLI to create an object:

```bash
tmd object create book golang-in-action
# Created: book/golang-in-action
```

The CLI automatically appends a ULID to the slug for uniqueness. Alternatively, you can create a file manually at `objects/book/golang-in-action.md` — manual files do not require a ULID and are backward compatible.

```markdown
---
title: Go in Action
status: reading
rating: 4.5
---

# Notes

A great book about Go...
```

The TUI will automatically detect the new file and display it. If you cloned an existing vault, the database is populated automatically on first open — no manual `tmd --reindex` needed.

## 5. Query and Search

```bash
# Filter by type and property
tmd query "type=book status=reading"

# Full-text search
tmd search "concurrency"
```
