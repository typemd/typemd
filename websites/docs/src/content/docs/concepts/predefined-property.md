---
title: Predefined Property Types
description: The built-in property types available for type schemas.
sidebar:
  order: 3.5
---

Every property in a type schema has a `type` that determines what values it accepts and how they are validated.

## Available Types

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

## Basic types

### string

Accepts any text value.

```yaml
- name: title
  type: string
```

### number

Accepts integers and floating-point numbers.

```yaml
- name: rating
  type: number
```

### date

Accepts dates in `YYYY-MM-DD` format. YAML auto-parsed dates (e.g., `2026-03-09` without quotes) are also accepted.

```yaml
- name: published
  type: date
```

### datetime

Accepts ISO 8601 date-time values with a time component. Must include hours and minutes at minimum.

```yaml
- name: created_at
  type: datetime
```

### url

Accepts URLs starting with `http://` or `https://`. Other schemes (ftp, ssh, etc.) are not supported.

```yaml
- name: website
  type: url
```

### checkbox

Accepts boolean `true` or `false` only. String values like `"true"` are rejected.

```yaml
- name: favorite
  type: checkbox
```

## Choice types

### select

Single choice from a predefined list. Requires an `options` array where each entry has a `value` (required) and optional `label` for display.

```yaml
- name: status
  type: select
  options:
    - value: to-read
      label: To Read
    - value: reading
      label: Reading
    - value: done
      label: Done
```

When `label` is omitted, it defaults to the `value`.

### multi_select

Multiple choices from a predefined list. Same `options` format as `select`. The property value is a list of strings.

```yaml
- name: tags
  type: multi_select
  options:
    - value: fiction
      label: Fiction
    - value: non-fiction
      label: Non-Fiction
    - value: classic
      label: Classic
```

A single string value (e.g., `tags: fiction`) is automatically treated as a one-element list (`["fiction"]`).

## Relation type

Links an Object to another Object. See the [Relations](/concepts/relations) page for details on `target`, `bidirectional`, `inverse`, and `multiple` fields.

```yaml
- name: author
  type: relation
  target: person
  bidirectional: true
  inverse: books
```

## Property attributes

Beyond `name` and `type`, properties support optional attributes:

| Attribute | Type | Description |
|-----------|------|-------------|
| `emoji` | string | Visual icon for compact display (must be unique within the type) |
| `pin` | integer | Positive integer for prominent display at the top of the TUI body panel. Lower values appear first. Pinned properties are excluded from the Properties panel. |
| `default` | any | Default value assigned when creating a new object |

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

## Validation

TypeMD uses lenient validation — only properties defined in the schema are validated:

- Extra properties (not in schema) are allowed
- Missing properties do not cause errors
- `select` / `multi_select` values must be in the `options` list
- `number` must be numeric
- `date` must be `YYYY-MM-DD` format
- `datetime` must include a time component
- `url` must start with `http://` or `https://`
- `checkbox` must be boolean
- `relation` targets are checked for correct type
