---
title: tmd show
description: Display an object's full information.
sidebar:
  order: 3
---

Displays an object's full information: properties (including relations) and body.

```bash
tmd show book/golang-in-action-01jqr3k5mp...
```

Example output:

```
book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz

Properties
──────────
  title: Go in Action
  status: reading
  rating: 4.5
  author: → person/alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz

Body
────
  # Notes
  A great book about Go...
```

Properties are displayed in schema-defined order. Relation properties use `→` for forward and `←` for inverse links.
