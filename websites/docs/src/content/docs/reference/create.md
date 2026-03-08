---
title: tmd object create
description: Create a new object from a type schema.
sidebar:
  order: 2.5
---

Creates a new object file (Markdown + YAML frontmatter) based on the type schema.

```bash
tmd object create book clean-code
tmd object create person robert-martin
```

The command generates a `.md` file under `objects/<type>/` with all schema-defined properties set to their default values (or `null` if no default is specified). The object is also added to the SQLite index.

Each object is assigned a unique ULID suffix, so multiple objects of the same type can share the same name (e.g. two `book/clean-code-<ulid>` files with different ULIDs). If the type does not exist, an error is returned.

