---
title: Objects & Types
description: Understanding Objects and Types in TypeMD.
sidebar:
  order: 1
---

## Objects

An Object is the basic unit of TypeMD. Each object is stored as a Markdown file with YAML frontmatter (properties) and body content.

### Object ID

Object IDs follow the format `type/<slug>-<ulid>`, e.g. `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`. Each directory under `objects/` is a **type namespace** — different types can share the same slug. When created via the CLI, a 26-character lowercase ULID is automatically appended to the slug for uniqueness. Objects created manually (without the CLI) do not require a ULID — this is backward compatible.

### Data Structure

```
vault/
├── .typemd/
│   ├── types/              # type schema definitions (YAML)
│   │   ├── book.yaml
│   │   └── person.yaml
│   └── index.db            # SQLite index (auto-updated)
└── objects/
    ├── book/
    │   └── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
    └── person/
        └── alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz.md
```

## Types

Every object belongs to a type. Types define property names, data types, and validation rules via schema files stored in `.typemd/types/`.

### Built-in Types

| Type | Properties |
|------|-----------|
| `book` | title (string), status (enum: to-read/reading/done), rating (number) |
| `person` | name (string), role (string) |
| `note` | title (string), tags (string) |

### Property Types

| Type | Description | Example |
|------|-------------|---------|
| `string` | Text | `"Go in Action"` |
| `number` | Integer or float | `42`, `3.14` |
| `enum` | Enumerated value, requires `values` | `"reading"` |
| `relation` | Link to another object | `"person/alan"` |
