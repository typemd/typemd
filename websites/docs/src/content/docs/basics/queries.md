---
title: Queries
description: Structured filtering of Objects by property values.
sidebar:
  order: 5
---

Queries let you filter Objects by their property values using structured filter rules. Unlike [search](/basics/search), which matches free-text across all content, queries target specific properties with typed operators.

## How it works

Queries use structured `FilterRule` conditions to match Objects. Each rule specifies a property, an operator, and a value. When multiple rules are provided, they are combined with AND logic — only Objects matching all conditions are returned.

The `type` property is a special filter that matches the Object's type name. All other properties match against frontmatter property values.

## Listing objects

Use the CLI to list all objects or search:

```bash
# List all objects
tmd object list
tmd object list --json

# Full-text search
tmd search "concurrency"
```

## Filter operators

Views support type-aware filter operators. Each property type has a defined set of valid operators:

| Property type | Operators |
|---------------|-----------|
| `string`, `url` | `is`, `is_not`, `contains`, `does_not_contain`, `starts_with`, `ends_with`, `is_empty`, `is_not_empty` |
| `number` | `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `is_empty`, `is_not_empty` |
| `date`, `datetime` | `eq`, `before`, `after`, `on_or_before`, `on_or_after`, `is_empty`, `is_not_empty` |
| `select` | `is`, `is_not`, `is_empty`, `is_not_empty` |
| `multi_select`, `relation` | `contains`, `does_not_contain`, `is_empty`, `is_not_empty` |
| `checkbox` | `is`, `is_not` |

Filter rules are defined in view configurations:

```yaml
filter:
  - property: status
    operator: is
    value: reading
  - property: rating
    operator: gt
    value: "3"
```

Multiple filters are combined with AND logic.

## Fallback mode

If the SQLite index is unavailable (missing or corrupt), queries automatically fall back to filesystem scanning with in-memory filter matching. The fallback path supports all the same operators listed above, ensuring consistent results. Performance is degraded (O(n) file reads instead of indexed queries), but operations complete successfully. A warning is logged when fallback is triggered.

## Sorting

Queries support sorting by property values. When used via the TUI [view mode](/tui/tui#view-mode--table-display), sort rules are defined in the view configuration:

```yaml
sort:
  - property: rating
    direction: desc
```

Sort direction is either `asc` (ascending) or `desc` (descending). Multiple sort rules are applied in order (first rule is the primary sort).

## See also

- [Search](/basics/search) — full-text search across all content
- [Views](/advanced/file-structure#views) — view configuration format
- [Data Model](/developers/data-model#querying-architecture) — querying details
