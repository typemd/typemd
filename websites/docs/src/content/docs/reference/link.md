---
title: tmd relation link
description: Create a relation between two objects.
sidebar:
  order: 6
---

Creates a relation between two objects. If the schema defines `bidirectional: true`, the inverse property is automatically updated.

```bash
tmd relation link book/golang-in-action author person/alan-donovan
```

For single-value relations (without `multiple: true`), linking overwrites the previous value. For multi-value relations, the new target is appended to the list.

