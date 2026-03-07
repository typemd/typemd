---
title: Introduction
description: What is TypeMD and why it exists.
sidebar:
  order: 1
---

TypeMD is a local-first CLI knowledge management tool. Your knowledge base is made of **Objects** — not files. Markdown is just the storage format.

## Philosophy

Most note-taking tools make you think like a computer: files, folders, hierarchies.

TypeMD lets you think in **Objects** — books, people, ideas, meetings — connected by **Relations**. The structure emerges from your knowledge, not from a folder tree.

## Core Concepts

### Object

An Object is the basic unit of TypeMD. Each object is a Markdown file with YAML frontmatter (properties) and body content.

Object IDs follow the format `type/<slug>-<ulid>`, e.g. `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`. When created via the CLI, a 26-character ULID is automatically appended to the slug for uniqueness.

```markdown
---
title: Go in Action
status: reading
rating: 4.5
---

# Notes

A great book about Go...
```

### Type

Every object belongs to a type. Types define property names, data types, and validation rules via schema files.

TypeMD includes three built-in types:

| Type | Properties |
|------|-----------|
| `book` | title (string), status (enum: to-read/reading/done), rating (number) |
| `person` | name (string), role (string) |
| `note` | title (string), tags (string) |

Custom type schemas go in `.typemd/types/`.

### Relation

Relations are defined as `relation`-type properties in type schemas. They support:

- **Unidirectional / Bidirectional** — bidirectional relations auto-sync both sides
- **Single / Multiple values** — multiple values stored as YAML arrays
