---
title: tmd init
description: Initialize a new vault.
sidebar:
  order: 2
---

Initializes a new vault. Creates `.typemd/` directory structure, SQLite database, and a `.typemd/.gitignore` that excludes `index.db`.

```bash
tmd init
```

Running `tmd init` on an already-initialized vault will return an error.
