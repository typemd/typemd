---
title: tmd object show
description: Display an object's full information.
sidebar:
  order: 3
---

Displays an object's full information: properties (including relations and backlinks) and body.

Supports prefix matching — you can omit the ULID suffix if the prefix uniquely identifies an object. If multiple objects match in an interactive terminal, a picker is shown to select the intended object. In non-interactive mode (piped input), the command returns an error listing the matching IDs.

```bash
tmd object show book/golang-in-action
```

Example output:

```
book/golang-in-action

Properties
──────────
  title: Go in Action
  status: reading
  rating: 4.5
  author: → person/alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz
  backlinks: ⟵ note/reading-list-01jqr4a2bcdef0123456789xyz

Local Properties
────────────────
  mood: reflective

Body
────
  # Notes
  A great book about Go...
```

Properties are displayed in schema-defined order. Relation properties use `→` for forward and `←` for inverse links. Wiki-link backlinks use `⟵` to indicate incoming references from other objects.

If the object contains frontmatter properties not defined in the type schema ("local properties"), they are displayed in a separate **Local Properties** section. This section only appears when local properties exist.
