---
title: tmd create
description: Create a new object from a type schema.
sidebar:
  order: 2.5
---

Creates a new object file (Markdown + YAML frontmatter) based on the type schema.

```bash
tmd create book clean-code
tmd create person robert-martin
```

The command generates a `.md` file under `objects/<type>/` with all schema-defined properties set to their default values (or `null` if no default is specified). The object is also added to the SQLite index.

If the type does not exist or an object with the same name already exists, an error is returned.
