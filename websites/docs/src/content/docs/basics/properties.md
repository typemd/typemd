---
title: Properties
description: How properties define the structure of your Objects.
sidebar:
  order: 1
---

Properties are named fields on an Object, defined by its [Type schema](/concepts/types). They are stored in the YAML frontmatter of the Markdown file, giving each Object structured, queryable data alongside its free-form body content.

## System Properties

Every Object supports five system properties managed by TypeMD. These provide the baseline metadata that every knowledge management tool needs — identity, description, temporal tracking, and categorization — without requiring users to define them in every type schema. Not all system properties are present in every Object's frontmatter: `name` is auto-populated on creation, `created_at` and `updated_at` are set when using the CLI, while `description` and `tags` only appear when explicitly set by the user.

| Property | Description | Mutability | Why |
|----------|-------------|------------|-----|
| `name` | Display name, auto-populated from the slug on creation | User-authored | Decouples display from filename — supports spaces, casing, and renaming without moving files |
| `description` | Optional single-line summary for list displays and search results | User-authored | Provides a consistent summary field across all types, used in list views, search results, and API responses |
| `created_at` | Creation timestamp in RFC 3339 format (set once, never modified) | Auto-managed | Enables sorting by creation date and understanding the timeline of a vault |
| `updated_at` | Last-modified timestamp in RFC 3339 format (updated on every save) | Auto-managed | Enables sorting by recency and tracking the evolution of Objects |
| `tags` | Array of tag references (relation to the built-in `tag` type, multiple) | User-authored | Cross-cutting categorization that works across all types |

**User-authored** properties (`name`, `description`, `tags`) can be overridden by [object templates](/basics/templates). **Auto-managed** properties (`created_at`, `updated_at`) cannot be overridden — they always reflect the actual creation and modification times.

These names are reserved and cannot be used in type schemas or [shared properties](#shared-properties). The only exception is `name`, which can appear in `properties` with a `template` field for [name templates](/basics/templates#name-templates).

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
| `pin` | integer | Positive integer for prominent display at the top of the TUI body panel. Lower values appear first. Pinned properties are excluded from the Properties panel. |
| `default` | any | Default value assigned when creating a new object |

**Why emoji?** In compact UI contexts (TUI properties panel, table columns), property names consume significant horizontal space. Emojis provide a space-efficient visual identifier that is instantly recognizable.

**Why pin?** Without pinning, all properties are displayed equally in the Properties panel. Users must scan the full list to find key metadata like status or rating. Pinning lets type authors mark specific properties for prominent display at the top of the body panel.

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

If the same property appears in multiple types (e.g., `due_date` in both `project` and `task`), you can define it once in `.typemd/properties.yaml` and reference it with the `use` keyword.

### Defining shared properties

```yaml
# .typemd/properties.yaml
properties:
  - name: due_date
    type: date
    emoji: "\U0001F4C5"
  - name: priority
    type: select
    options:
      - value: high
      - value: medium
      - value: low
```

Each entry supports all the same fields as a regular type schema property: `name`, `type`, `emoji`, `pin`, `options`, `target`, `default`, `multiple`, `bidirectional`, and `inverse`. When referenced via `use:` in a type schema, only `pin`, `emoji`, and `description` can be overridden — all other fields are inherited from the shared definition.

The file is optional. If `.typemd/properties.yaml` does not exist or contains no `properties` array, TypeMD treats it as an empty set of shared properties.

### Referencing in type schemas

Use the `use` keyword to reference a shared property by name:

```yaml
# .typemd/types/project/schema.yaml
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
# .typemd/types/task/schema.yaml
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
| Defined in | `.typemd/properties.yaml` | `.typemd/types/<type>/schema.yaml` |
| Reusable | Yes — referenced via `use` | No — scoped to one type |
| Customizable per type | `pin`, `emoji`, and `description` | Fully customizable |
| Use case | Properties shared across multiple types | Properties unique to one type |

**Rule of thumb**: if a property appears in two or more types with the same definition, make it a shared property.

## See also

- [Types](/concepts/types) — how types define object structure
- [Relations](/concepts/relations) — linking objects together
- [Validation](/basics/validation) — how property values are validated
