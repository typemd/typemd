---
title: Types
description: How Types define the structure of your Objects.
sidebar:
  order: 3
---

A Type is a blueprint that defines what kind of Object you're working with and what properties it can have.

## What is a Type?

Think of a Type as a category with structure. A `book` Type knows that books have a title, a reading status, and a rating. A `person` Type knows that people have a name and a role.

Types are defined as YAML schema files in `types/`:

```yaml
# types/book/schema.yaml
name: book
plural: books
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

The optional `plural` field provides a grammatically correct display name for the type in collection contexts (e.g., "▼ 📚 books (3)" in TUI group headers). When omitted, the type name is used as fallback.

The optional `emoji` field provides a visual icon for the type in CLI and TUI output.

The optional `color` field assigns a color for visual theming in TUI and Web UI. Accepts preset names (`red`, `blue`, `green`, `yellow`, `purple`, `orange`, `pink`, `cyan`, `gray`, `brown`) or custom hex codes (`#RGB` or `#RRGGBB`).

The optional `description` field provides free-text documentation of the type's purpose (e.g., `description: "Books I've read or want to read"`). This is distinct from the `description` system property on objects.

The optional `version` field is a semver-style `"major.minor"` string for tracking schema evolution (e.g., `version: "1.0"`). When omitted, it defaults to `"0.0"` (unversioned). Increment the major number for breaking changes, and the minor number for backward-compatible changes, providing a foundation for future migration tooling.

The optional `unique` field, when set to `true`, enforces that no two objects of the same type can share the same name. See [Unique constraint](#unique-constraint) below for details.

## Why Types matter

Types give your knowledge base **consistency** and **queryability**:

- **Consistency** — Every book Object has the same set of properties, so you don't end up with `status` on one book and `reading_state` on another.
- **Queryability** — Because properties are typed, you can query precisely: "show me all books where status is reading and rating > 4".
- **Validation** — TypeMD validates property values against the schema (e.g. select values must be in the allowed options), catching mistakes early.

## Built-in Types

TypeMD has two built-in Types:

| Type | Properties | Purpose |
|------|------------|---------|
| 🏷️ `tag` | color (string), icon (string) | Backs the `tags` system property; has `unique: true` to enforce name uniqueness |
| 📄 `page` | _(none)_ | General-purpose content container for free-form writing |

Built-in types exist in every vault automatically and cannot be deleted. You can override them by creating a custom `types/<name>/schema.yaml` file with additional properties. TypeMD deliberately avoids shipping opinionated types like `book` or `note` — you design the types that fit your domain.

For details on tags, see [Tags](/basics/tags).

For details on object templates, see [Templates](/basics/templates).

## Unique constraint

Setting `unique: true` on a type schema enforces name uniqueness within that type. When enabled, TypeMD rejects the creation of a new object if another object of the same type already has the same name.

```yaml
# types/person/schema.yaml
name: person
emoji: 👤
unique: true
properties:
  - name: role
    type: string
```

Key behaviors:

- **Scoped to type** — The constraint applies within a single type. Two objects of *different* types can share the same name (e.g., a `person` named "john-doe" and a `character` named "john-doe" can coexist).
- **Case-sensitive** — Names are compared exactly. "Go" and "go" are treated as distinct names, so both are allowed even on a unique type.
- **Built-in `tag` type** — The `tag` type has `unique: true` by default, ensuring tag names are always unique.
- **Validation** — Running `tmd type validate` checks for duplicate names on unique types and reports any violations.

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

If a property appears in multiple types, you can define it once and reuse it. See [Shared Properties](/basics/properties#shared-properties) for details.

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
- Property names `name`, `description`, `created_at`, `updated_at`, `tags`, `locked`, `archived`, `object_type`, `links`, `backlinks`, `created_by`, and `updated_by` are reserved for [system properties](/basics/properties#system-properties) and cannot be used in type schemas. The only exception is `name`, which can appear in `properties` with a `template` field (for [name templates](/basics/templates#name-templates)).

## Views

Each type can have one or more **views** — saved configurations that control how objects are displayed (sorted, filtered, grouped). Views are stored alongside the type schema in `types/<name>/views/`.

Every type has an implicit default view (list layout, sorted by name). You can create custom views with different sort orders, filters, and grouping. In the TUI, press `v` on a type group to enter [view mode](/tui/tui#view-mode--table-display).

For details on view configuration format, see [File Structure: Views](/advanced/file-structure#views). For filter operators, see [Queries: Filter operators](/basics/queries#filter-operators).
