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

Use the CLI to create an object:

```bash
tmd create book golang-in-action
# Created: book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz
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

The TUI will automatically detect the new file and display it. If you cloned an existing vault, the database is populated automatically on first open — no manual `tmd reindex` needed.

## 4. Query and Search

```bash
# Filter by type and property
tmd query "type=book status=reading"

# Full-text search
tmd search "concurrency"
```
