---
title: tmd relation link
description: Create a relation between two objects.
sidebar:
  order: 7
---

Creates a relation between two objects. If the schema defines `bidirectional: true`, the inverse property is automatically updated.

Object IDs support prefix matching — you can omit the ULID suffix if the prefix uniquely identifies an object. If a prefix matches multiple objects in an interactive terminal, a picker is shown to select the intended object.

```bash
tmd relation link book/golang-in-action author person/alan-donovan
```

For single-value relations (without `multiple: true`), linking overwrites the previous value. For multi-value relations, the new target is appended to the list.
