---
title: Type Schemas
description: Defining custom type schemas with properties, enums, and relations.
sidebar:
  order: 3
---

Type schemas are YAML files in `.typemd/types/` that define the structure of your objects.

## Basic Format

```yaml
# .typemd/types/book.yaml
name: book
properties:
  - name: title
    type: string
  - name: status
    type: enum
    values: [to-read, reading, done]
    default: to-read
  - name: rating
    type: number
```

## Property Types

| Type | Description | Example |
|------|-------------|---------|
| `string` | Text | `"Go in Action"` |
| `number` | Integer or float | `42`, `3.14` |
| `enum` | Enumerated value, requires `values` | `"reading"` |
| `relation` | Link to another object | `"person/alan"` |

## Relation Properties

See the [Relations](/guides/relations) guide for details on defining relation properties with `target`, `bidirectional`, `inverse`, and `multiple` fields.

## Validation

TypeMD uses lenient validation:

- Only validates properties defined in the schema
- Extra properties (not in schema) are allowed
- Missing properties do not cause errors
- `enum` values must be in the `values` list
- `number` must be numeric
- `relation` targets are checked for correct type
