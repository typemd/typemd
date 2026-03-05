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

This creates the `.typemd/` directory structure and SQLite database in the current directory.

## 2. Open the TUI

```bash
tmd
```

This launches the three-panel interface for browsing your vault.

## 3. Create Your First Object

Create a file at `objects/book/golang-in-action.md`:

```markdown
---
title: Go in Action
status: reading
rating: 4.5
---

# Notes

A great book about Go...
```

The TUI will automatically detect the new file and display it. If you cloned an existing vault, the database is populated automatically on first open — no manual `tmd reindex` needed.

## 4. Query and Search

```bash
# Filter by type and property
tmd query "type=book status=reading"

# Full-text search
tmd search "concurrency"
```
