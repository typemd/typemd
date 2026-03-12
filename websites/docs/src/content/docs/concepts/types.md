---
title: Types
description: How Types define the structure of your Objects.
sidebar:
  order: 3
---

A Type is a blueprint that defines what kind of Object you're working with and what properties it can have.

## What is a Type?

Think of a Type as a category with structure. A `book` Type knows that books have a title, a reading status, and a rating. A `person` Type knows that people have a name and a role.

Types are defined as YAML schema files in `.typemd/types/`:

```yaml
# .typemd/types/book.yaml
name: book
emoji: 📚
properties:
  - name: title
    type: string
  - name: status
    type: select
    options:
      - value: to-read
      - value: reading
      - value: done
    default: to-read
  - name: rating
    type: number
```

The optional `emoji` field provides a visual icon for the type in CLI and TUI output.

## Why Types matter

Types give your knowledge base **consistency** and **queryability**:

- **Consistency** — Every book Object has the same set of properties, so you don't end up with `status` on one book and `reading_state` on another.
- **Queryability** — Because properties are typed, you can query precisely: "show me all books where status is reading and rating > 4".
- **Validation** — TypeMD validates property values against the schema (e.g. select values must be in the allowed options), catching mistakes early.

## Built-in Types

TypeMD ships with three built-in Types to get you started:

| Type | Properties |
|------|------------|
| 📚 `book` | title (string), status (select: to-read/reading/done), rating (number) |
| 👤 `person` | role (string) |
| 📝 `note` | title (string), tags (string) |

You can modify these or create your own Types to fit your knowledge domain.

## Types vs. tags

Many tools use freeform tags to categorize notes. Tags are flexible but fragile — typos create phantom categories, and there's no structure to query against.

Types are explicit: a `book` is always a `book`, with defined properties you can filter, sort, and connect. You trade a bit of upfront definition for much more powerful organization.

## Property types

Each property in a Type schema has a data type:

| Type | Description | Example |
|------|-------------|---------|
| `string` | Text | `"Go in Action"` |
| `number` | Integer or float | `42`, `3.14` |
| `date` | Date in YYYY-MM-DD format | `"2026-01-15"` |
| `datetime` | ISO 8601 datetime | `"2026-01-15T10:30:00"` |
| `url` | URL with http/https scheme | `"https://example.com"` |
| `checkbox` | Boolean value | `true`, `false` |
| `select` | One of a fixed set of options | `"reading"` |
| `multi_select` | Multiple values from options | `["go", "programming"]` |
| `relation` | A link to another Object | `"person/alan-donovan"` |

For relation properties, see the [Relations](/concepts/relations) page for details on `target`, `bidirectional`, `inverse`, and `multiple` fields.

## Validation

TypeMD uses lenient validation:

- Only validates properties defined in the schema
- Extra properties (not in schema) are allowed
- Missing properties do not cause errors
- `select`/`multi_select` values must be in the defined `options` list
- `number` must be numeric
- `date` must be in YYYY-MM-DD format
- `url` must start with http:// or https://
- `relation` targets are checked for correct type
- Property names `name`, `description`, `created_at`, and `updated_at` are reserved for [system properties](/concepts/data-model#system-properties) and cannot be used in type schemas
