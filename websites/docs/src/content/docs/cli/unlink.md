---
title: tmd relation unlink
description: Remove a relation between two objects.
sidebar:
  order: 7
---

Removes a relation from the source object's frontmatter.

By default, only the forward side is removed. Use `--both` to also remove the inverse relation on the target object. This only applies to bidirectional relations (those with `bidirectional: true` and an `inverse` field in the schema).

Object IDs support prefix matching — you can omit the ULID suffix if the prefix uniquely identifies an object. If a prefix matches multiple objects in an interactive terminal, a picker is shown to select the intended object.

```bash
# Remove only the forward relation (book's "author" property)
tmd relation unlink book/golang-in-action author person/alan-donovan

# Remove both sides: the book's "author" and the person's inverse "books" property
tmd relation unlink book/golang-in-action author person/alan-donovan --both
```

## Flags

| Flag | Description |
|------|-------------|
| `--both` | Also remove the inverse relation on the target object (requires a bidirectional relation with an `inverse` field) |
