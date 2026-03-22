---
title: File Structure
description: How TypeMD organizes files and directories in a Vault.
sidebar:
  order: 1
---

TypeMD stores everything as plain files on disk. This page describes the directory layout and file formats you'll encounter when working with a Vault directly.

## Vault directory layout

A Vault is a regular directory with a specific structure:

```
vault/
├── .typemd/
│   ├── types/              # user-defined type schemas (directory format)
│   │   ├── book/
│   │   │   ├── schema.yaml # type schema definition
│   │   │   └── views/      # saved views for this type (optional)
│   │   │       ├── default.yaml
│   │   │       └── by-rating.yaml
│   │   └── person/
│   │       └── schema.yaml
│   ├── config.yaml         # vault configuration (optional)
│   ├── properties.yaml     # shared property definitions (optional)
│   ├── instructions/       # skill overrides (optional)
│   │   └── explore.md      # override for explore skill
│   ├── index.db            # SQLite index (auto-updated)
│   └── tui-state.yaml      # TUI session state (auto-saved)
├── templates/              # object templates by type (optional)
│   └── book/
│       └── review.md       # default frontmatter + body for new objects
└── objects/
    ├── book/
    │   └── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
    └── person/
        └── alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz.md
```

- `.typemd/` — configuration and internal state (vault config, type schemas, shared properties, skill overrides, index, TUI state)
- `templates/` — optional object templates organized by type
- `objects/` — all Object files organized by type

## Object files

Objects are stored as Markdown files with YAML frontmatter under `objects/<type>/`. Each file represents a single Object.

### Object ID

Every Object has a unique ID in the format `type/<slug>-<ulid>`, for example:

```
book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz
```

- **type** — the Object's type (matches the subdirectory name under `objects/`)
- **slug** — a human-readable identifier derived from the name
- **ulid** — a 26-character lowercase [ULID](https://github.com/ulid/spec) appended for uniqueness

### ULID

When you create an Object via the CLI (`tmd object create`), a ULID is automatically appended to the slug. If you create files manually (without the CLI), the ULID is optional — TypeMD will work with or without it. This makes it easy to add existing Markdown files to a Vault.

### System properties

All Objects have five system properties managed by TypeMD. These always appear first in the frontmatter in the following fixed order:

| Property | Description | Mutable |
|----------|-------------|---------|
| `name` | Display name, preserves the original input on creation (auto-derived from slug for pre-slugified names) | User-authored |
| `description` | Optional single-line summary for list displays and search results | User-authored |
| `created_at` | Creation timestamp in RFC 3339 format (set once, never modified) | Auto-managed |
| `updated_at` | Last-modified timestamp in RFC 3339 format (updated on every save) | Auto-managed |
| `tags` | Array of tag references (relation to the built-in `tag` type, multiple) | User-authored |

System properties are followed by schema-defined properties. These names are reserved and cannot be used in type schemas or shared properties.

**User-authored** properties (`name`, `description`, `tags`) can be overridden by object templates. **Auto-managed** properties (`created_at`, `updated_at`) cannot be overridden — they always reflect the actual creation and modification times.

For details on the frontmatter format and how to edit it manually, see [Frontmatter](/advanced/frontmatter).

## Type schema files

Each type is stored as a directory under `.typemd/types/<name>/` containing a `schema.yaml` file:

```yaml
# .typemd/types/book/schema.yaml
name: book
plural: books
color: blue
description: "Books I've read or want to read"
unique: false
properties:
  - name: status
    type: select
    options:
      - value: reading
      - value: completed
      - value: want-to-read
  - name: author
    type: relation
    target: person
```

Two types are built-in: `tag` (backs the `tags` system property, plural "tags", `unique: true`) and `page` (general-purpose content container, plural "pages", emoji 📄). Built-in types cannot be deleted but can be overridden by custom type schemas.

:::note
Legacy single-file format (`.typemd/types/book.yaml`) is automatically migrated to directory format (`book/schema.yaml`) on first load.
:::

## Views

Each type can have saved views that define how objects are filtered, sorted, and displayed. Views are stored as YAML files under `.typemd/types/<name>/views/`:

```yaml
# .typemd/types/book/views/by-rating.yaml
name: by-rating
layout: list
columns: [status, rating]
sort:
  - property: rating
    direction: desc
filter:
  - property: status
    operator: is
    value: reading
group_by:
  - property: genre
```

Two layouts are available:

- **`list`** — displays object names with optional inline property values. Default columns: none (name only).
- **`table`** — displays a columnar table with a NAME column followed by property columns with headers. Default columns: all properties.

The optional `columns` field specifies which properties to display. When set, both layouts use exactly those properties.

Every type has an implicit default view (list layout, sorted by name ascending). When you customize the default view, it is saved as `views/default.yaml`. If no view files exist, the type directory does not need a `views/` subdirectory.

For the full type schema format, see [Types](/concepts/types).

## Shared properties file

The optional file `.typemd/properties.yaml` defines reusable property definitions that can be referenced from type schemas via `use`:

```yaml
# .typemd/properties.yaml
properties:
  - name: status
    type: select
    options:
      - value: active
      - value: archived
```

For details on shared properties, see [Shared Properties](/basics/properties#shared-properties).

## Object templates

Each type can have one or more templates at `templates/<type>/<name>.md`. Templates are regular Markdown files (frontmatter + body) that provide default content when creating Objects.

```bash
# Single template — auto-applied
tmd object create book "Clean Code"

# Explicit template selection
tmd object create book "Clean Code" -t review

# Multiple templates, no flag — interactive selection
tmd object create book "Clean Code"
```

Template frontmatter properties override schema defaults. The template body becomes the initial Object body. Auto-managed system properties (`created_at`, `updated_at`) in templates are ignored.
