---
title: Properties
description: How properties define the structure of your Objects.
sidebar:
  order: 1
---

Properties are named fields on an Object, defined by its [Type schema](/concepts/types). They are stored in the YAML frontmatter of the Markdown file, giving each Object structured, queryable data alongside its free-form body content.

## System Properties

Every Object supports seven system properties managed by TypeMD. These provide the baseline metadata that every knowledge management tool needs — identity, description, temporal tracking, categorization, protection, and archival — without requiring users to define them in every type schema. Not all system properties are present in every Object's frontmatter: `name` is auto-populated on creation, `created_at` and `updated_at` are set when using the CLI, while `description`, `tags`, `locked`, and `archived` only appear when explicitly set by the user.

| Property | Description | Mutability | Why |
|----------|-------------|------------|-----|
| `name` | Display name, auto-populated from the slug on creation. When missing or empty, falls back to the [Display ID](/concepts/glossary#display-id) derived from the filename. | User-authored | Decouples display from filename — supports spaces, casing, and renaming without moving files |
| `description` | Optional single-line summary for list displays and search results | User-authored | Provides a consistent summary field across all types, used in list views, search results, and API responses |
| `created_at` | Creation timestamp in RFC 3339 format (set once, never modified) | Auto-managed | Enables sorting by creation date and understanding the timeline of a vault |
| `updated_at` | Last-modified timestamp in RFC 3339 format (updated on every save) | Auto-managed | Enables sorting by recency and tracking the evolution of Objects |
| `tags` | Array of tag references (relation to the built-in `tag` type, multiple) | User-authored | Cross-cutting categorization that works across all types |
| `locked` | Boolean that prevents editing when `true` | User-authored | Protects finished or archival Objects from accidental modification |
| `archived` | Soft-delete flag — hides object from default queries | User-authored | Removes clutter from default views without permanently deleting content |

**User-authored** properties (`name`, `description`, `tags`, `locked`, `archived`) can be overridden by [object templates](/basics/templates). **Auto-managed** properties (`created_at`, `updated_at`) cannot be overridden — they always reflect the actual creation and modification times.

These names are reserved and cannot be used in type schemas or [shared properties](#shared-properties). The only exception is `name`, which can appear in `properties` with a `template` field for [name templates](/basics/templates#name-templates).

## Derived and Computed Properties

In addition to the seven stored system properties above, TypeMD provides non-stored properties that are resolved at runtime — they never appear in the YAML frontmatter.

| Property | Category | Source | Description |
|----------|----------|--------|-------------|
| `object_type` | Derived | File path `objects/<type>/` | The object's type name |
| `created_by` | Derived | Git history | Author of the initial commit |
| `links` | Computed | Markdown body | Outgoing wiki-link references |
| `backlinks` | Computed | Index | Incoming wiki-link references from other objects |
| `updated_by` | Computed | Git history | Author of the most recent commit |

**Derived** properties have stable values inferred from structure or metadata — once an object exists, these values rarely change. **Computed** properties have dynamic values that may change whenever the object's content or the vault's state changes.

All non-stored properties are read-only. Attempting to set them via `SetProperty` returns an error, and any values found in frontmatter are stripped on save. These names are also reserved and cannot be used in type schemas.

## Property Types

Each property in a type schema has a `type` that determines what values it accepts and how they are validated.

| Type | Description | Example value |
|------|-------------|---------------|
| `string` | Plain text | `"Go in Action"` |
| `number` | Integer or floating-point | `42`, `3.14` |
| `date` | Date in `YYYY-MM-DD` format | `"2026-03-09"` |
| `datetime` | Date and time in ISO 8601 | `"2026-03-09T10:30:00+08:00"` |
| `url` | Web URL (http/https only) | `"https://example.com"` |
| `checkbox` | Boolean value | `true`, `false` |
| `select` | Single choice from options | `"reading"` |
| `multi_select` | Multiple choices from options | `["fiction", "sci-fi"]` |
| `relation` | Link to another Object | `"person/alan-donovan"` |

### Basic types

- **string** — Accepts any text value.
- **number** — Accepts integers and floating-point numbers.
- **date** — Accepts dates in `YYYY-MM-DD` format. YAML auto-parsed dates (e.g., `2026-03-09` without quotes) are also accepted.
- **datetime** — Accepts ISO 8601 date-time values. Must include hours and minutes at minimum.
- **url** — Accepts URLs starting with `http://` or `https://`. Other schemes (ftp, ssh, etc.) are not supported.
- **checkbox** — Accepts boolean `true` or `false` only. String values like `"true"` are rejected.

### Choice types

**select** — Single choice from a predefined list. Requires an `options` array where each entry has a `value` (required) and optional `label` for display:

```yaml
- name: status
  type: select
  options:
    - value: to-read
      label: To Read
    - value: reading
    - value: done
```

When `label` is omitted, it defaults to the `value`.

**multi_select** — Multiple choices from a predefined list. Same `options` format as `select`. The property value is a list of strings. A single string value (e.g., `genres: fiction`) is automatically treated as a one-element list (`["fiction"]`).

### Relation type

Links an Object to another Object. See the [Relations](/concepts/relations) page for details on `target`, `bidirectional`, `inverse`, and `multiple` fields.

```yaml
- name: author
  type: relation
  target: person
  bidirectional: true
  inverse: books
```

## Property Attributes

Beyond `name` and `type`, properties support optional attributes:

| Attribute | Type | Description |
|-----------|------|-------------|
| `emoji` | string | Visual icon for compact display (must be unique within the type) |
| `description` | string | Free-text documentation of the property's purpose |
| `pin` | integer | Positive integer for prominent display at the top of the TUI Properties panel. Lower values appear first. Pinned properties are sorted before unpinned ones and separated by a horizontal divider. |
| `default` | any | Default value assigned when creating a new object |

**Why emoji?** In compact UI contexts (TUI properties panel, table columns), property names consume significant horizontal space. Emojis provide a space-efficient visual identifier that is instantly recognizable.

**Why pin?** Without pinning, all properties are displayed equally in the Properties panel. Users must scan the full list to find key metadata like status or rating. Pinning lets type authors mark specific properties for prominent display at the top of the Properties panel, separated from the rest by a divider.

### Pin example

```yaml
- name: status
  type: select
  emoji: 📋
  pin: 1
  options:
    - value: to-read
    - value: reading
    - value: done
```

Pin values must be positive integers and unique within a type schema. Properties without `pin` (or `pin: 0`) appear in the Properties panel as usual.

## Shared Properties

When the same property appears in multiple types, defining it independently in each schema leads to duplication and inconsistency — a `due_date` in `project` might use `date` while `task` uses `datetime`. Shared properties let you define once and reference everywhere, ensuring consistent definitions across types.

If the same property appears in multiple types (e.g., `due_date` in both `project` and `task`), you can define it once in the `properties/` directory and reference it with the `use` keyword.

### Defining shared properties

Each shared property is its own file under `properties/<name>.yaml`, where `<name>` is the property name. No wrapper key is needed — the file contains the property definition directly.

```yaml
# properties/due_date.yaml
type: date
emoji: "\U0001F4C5"
```

```yaml
# properties/priority.yaml
type: select
options:
  - value: high
  - value: medium
  - value: low
```

Each file supports all the same fields as a regular type schema property: `type`, `emoji`, `pin`, `options`, `target`, `default`, `multiple`, `bidirectional`, and `inverse`. When referenced via `use:` in a type schema, only `pin`, `emoji`, and `description` can be overridden — all other fields are inherited from the shared definition.

The `properties/` directory is optional. If it does not exist or contains no `.yaml` files, TypeMD treats it as an empty set of shared properties.

### Referencing in type schemas

Use the `use` keyword to reference a shared property by name:

```yaml
# types/project/schema.yaml
name: project
emoji: "\U0001F4CB"
properties:
  - name: title
    type: string
  - use: due_date
  - use: priority
  - name: budget
    type: number
```

When the type is loaded, each `use` entry is resolved to the full property definition. You can mix `use` entries with regular property definitions — the order is preserved.

### Overriding fields

When referencing a shared property via `use`, you can override these fields:

| Field | Description |
|-------|-------------|
| `pin` | Override the pin position for this type |
| `emoji` | Override the display emoji for this type |
| `description` | Override the description for type-specific documentation |

All other fields cannot be overridden — they come from the shared definition.

```yaml
# types/task/schema.yaml
name: task
properties:
  - use: due_date
    pin: 1
    emoji: "\U0001F5D3\uFE0F"
  - use: priority
```

### Shared vs. inline properties

| | Shared Property | Inline Property |
|---|----------------|-----------------|
| Defined in | `properties/<name>.yaml` | `types/<type>/schema.yaml` |
| Reusable | Yes — referenced via `use` | No — scoped to one type |
| Customizable per type | `pin`, `emoji`, and `description` | Fully customizable |
| Use case | Properties shared across multiple types | Properties unique to one type |

**Rule of thumb**: if a property appears in two or more types with the same definition, make it a shared property.

## See also

- [Types](/concepts/types) — how types define object structure
- [Relations](/concepts/relations) — linking objects together
- [Validation](/basics/validation) — how property values are validated
