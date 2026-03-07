---
title: Relations
description: Connecting objects with typed, bidirectional links.
sidebar:
  order: 2
---

Relations are defined as `relation`-type properties in type schemas. They let you connect objects with named links.

## Defining Relations

```yaml
# .typemd/types/book.yaml
name: book
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
```

```yaml
# .typemd/types/person.yaml
name: person
properties:
  - name: name
    type: string
  - name: books
    type: relation
    target: book
    multiple: true
    bidirectional: true
    inverse: author
```

## Relation Fields

| Field | Description |
|-------|-------------|
| `target` | Target object's type name |
| `multiple` | Whether the property holds multiple values (array) |
| `bidirectional` | Auto-sync the inverse side when linking |
| `inverse` | Property name on the target type's schema |

## Using Relations

### Create a Link

```bash
tmd link book/golang-in-action author person/alan-donovan
```

When `bidirectional: true`, this automatically updates both the book's `author` and the person's `books` property.

### Remove a Link

```bash
tmd unlink book/golang-in-action author person/alan-donovan --both
```

Use `--both` to remove the inverse side as well. This only takes effect when the relation property has `bidirectional: true` and an `inverse` field defined in the schema.
