---
title: Tags
description: How tags provide cross-cutting categorization for your Objects.
sidebar:
  order: 2
---

Tags let you categorize Objects across Types. A `book` and a `note` can share the same tag, making tags ideal for cross-cutting themes like "golang", "important", or "to-review".

## Tags as first-class Objects

Unlike simple string labels, tags in TypeMD are full Objects backed by the built-in `tag` type. Each tag has its own Markdown file under `objects/tag/` with two optional properties:

| Property | Type | Description |
|----------|------|-------------|
| `color` | string | Display color for the tag |
| `icon` | string | Display icon for the tag |

Because tags are Objects, they can have descriptions, body content, and wiki-links — just like any other Object in your vault.

## Tag uniqueness

The built-in `tag` type has `unique: true` in its schema, which means tag names must be unique. You cannot have two tag Objects with the same `name` value. This uniqueness constraint is enforced at creation time and validated by `tmd type validate`.

Any user-defined type can also enable name uniqueness by adding `unique: true` to its schema. See [Data Model](/developers/data-model#unique-constraint-enforcement) for details.

## Using tags

Every Object supports a `tags` [system property](/basics/properties#system-properties) that holds references to tag Objects. The `tags` property only appears in frontmatter when explicitly set. In frontmatter, tag references can use either the full Object ID or the tag name:

```yaml
tags:
  - tag/go
  - tag/programming-01jqr3k5mpbvn8e0f2g7h9txyz
```

Both forms are valid. Name-based references (like `tag/go`) are resolved to their full IDs during sync.

## Tag reference resolution

During sync, TypeMD resolves name-based tag references to their full Object IDs. If a referenced tag does not exist, TypeMD auto-creates it with an empty body. This means you can freely add tags to your frontmatter without having to create the tag Object first — TypeMD handles it for you.

## Tags vs. Types

Tags and Types serve different purposes:

- **Types** define structure — a `book` always has title, status, and rating.
- **Tags** categorize across Types — a `book` and a `note` can share the same tag.

Use Types for structural consistency; use tags for cross-cutting categorization.

## See also

- [Properties](/basics/properties#system-properties) — the `tags` system property
- [Types](/concepts/types#built-in-types) — the built-in `tag` type
- [Validation](/basics/validation) — name uniqueness validation for tags
