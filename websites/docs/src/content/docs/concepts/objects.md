---
title: Objects
description: The basic unit of knowledge in TypeMD.
sidebar:
  order: 2
---

An Object is the fundamental building block of a TypeMD knowledge base. Everything you track — books, people, ideas, projects — is an Object.

## What is an Object?

An Object is a Markdown file with two parts:

- **Frontmatter** — YAML metadata at the top (system properties like `name` and timestamps, plus schema-defined properties)
- **Body** — free-form Markdown content (notes, thoughts, references)

```markdown
---
name: Go in Action
created_at: "2026-03-09T10:30:00+08:00"
updated_at: "2026-03-11T18:00:00+08:00"
status: reading
rating: 4.5
author: person/alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz
---

# Notes

A great book about Go concurrency patterns.

See also: [[note/go-concurrency-patterns]]
```

## Name is not filename

An Object's display name (`name` property) is independent of its filename. The filename is a stable storage identifier (slug + ULID); the name is a human-readable label that can contain spaces, mixed casing, and special characters.

This decoupling means you can rename an Object without moving files, and auto-generate names from templates (e.g., `日記 2026-03-14`) without worrying about filesystem constraints.

## Object ID

Every Object has an ID in the format `type/slug-ulid`, for example:

```
book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz
```

| Part | Description |
|------|-------------|
| `book` | The Type this Object belongs to |
| `golang-in-action` | The Slug — a human-readable name |
| `01jqr3k5mpbvn8e0f2g7h9txyz` | The ULID — a unique identifier |

The ULID is automatically appended when you use `tmd object create`. If you create files manually, the ULID is optional — TypeMD is backward compatible with plain slugs.

## Objects are not files

In most note-taking tools, you organize knowledge by putting files into folders. In TypeMD, you organize knowledge by giving Objects a **Type** and connecting them with **Relations**.

The file system is just a storage layer. A book Object and a person Object can reference each other regardless of where they sit in the directory tree. The structure comes from your data model, not your folder hierarchy.

## Where Objects live

Objects are stored under the `objects/` directory in your Vault, organized by Type:

```
objects/
├── book/
│   ├── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
│   └── designing-data-intensive-applications-01jqr4m7npcwo9f3g8h0jtuyab.md
├── person/
│   └── alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz.md
└── note/
    └── go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc.md
```

Each Type has its own subdirectory. This keeps the file system tidy, but the real organization happens through Types and Relations.
