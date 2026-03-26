---
title: Relations
description: How Objects connect to each other through typed links.
sidebar:
  order: 4
---

A Relation is a structured link between two Objects. It's how you express that a book has an author, a project has members, or an idea connects to a meeting.

## What is a Relation?

A Relation is a property of type `relation` defined in a Type schema. Unlike a plain hyperlink or wiki-link, a Relation is:

- **Named** — it has a property name like `author` or `books`, not just a URL
- **Typed** — it knows the target Type (e.g. `author` points to a `person`)
- **Optionally bidirectional** — linking a book to an author can automatically link the author back to the book

## Unidirectional vs. bidirectional

A **unidirectional** Relation is a one-way link. Setting `author` on a book points to a person, but the person Object doesn't know about it.

A **bidirectional** Relation keeps both sides in sync. When you link a book's `author` to a person, the person's `books` property is automatically updated — and vice versa.

```yaml
# book.yaml
- name: author
  type: relation
  target: person
  bidirectional: true
  inverse: books

# person.yaml
- name: books
  type: relation
  target: book
  multiple: true
  bidirectional: true
  inverse: author
```

## Single vs. multiple

- **Single-value** — A book has one `author`. Re-linking overwrites the previous value.
- **Multiple-value** — A person can have many `books`. Linking appends to the list.

This is controlled by the `multiple` field in the schema.

## Relation fields

| Field | Description |
|-------|-------------|
| `target` | Target object's type name |
| `multiple` | Whether the property holds multiple values (array). Single-value relations are overwritten on re-link; multi-value relations append. |
| `bidirectional` | Auto-sync the inverse side when linking |
| `inverse` | Property name on the target type's schema |

## Using Relations

### Create a link

```bash
tmd relation link book/golang-in-action author person/alan-donovan
```

When `bidirectional: true`, this automatically updates both the book's `author` and the person's `books` property.

### Remove a link

```bash
tmd relation unlink book/golang-in-action author person/alan-donovan --both
```

Use `--both` to remove the inverse side as well. This only takes effect when the relation property has `bidirectional: true` and an `inverse` field defined in the schema.

## Display indicators

When viewing an Object's properties, Relations use directional indicators to show the link direction:

| Indicator | Meaning | Example |
|-----------|---------|---------|
| `→` | Forward relation (this object links to another) | `→ person/alan-donovan` |
| `←` | Reverse relation (linked via inverse property) | `← book/clean-code` |
| `⟵` | Backlink (from wiki-links, not schema relations) | `⟵ note/my-note` |

## Relations vs. Links

TypeMD also supports [Links](/concepts/links) (`[[type/slug]]`) for informal inline references. Relations are structured and schema-defined; Links are freeform and live in the Markdown body. See the [Links](/concepts/links) page for a detailed comparison.
