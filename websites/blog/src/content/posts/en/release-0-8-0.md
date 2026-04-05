---
title: "TypeMD 0.8.0 — Your Vault, at a Glance"
description: "v0.8.0 moves type schemas and properties to where you can see them, adds object archival, graph export, and a template CLI — making your knowledge base more transparent and manageable."
date: 2026-04-05
tags: [release, devlog]
---

0.7.0 let you edit everything right in the TUI. But some things you might not even know where to find — type schemas and shared properties were buried inside the `.typemd/` directory, invisible unless you dug through config files. This release brings them out into the open, making your vault's structure crystal clear.

## Type Schemas Move to Vault Root

This is a **breaking change**. The `types/` and `properties/` directories move from `.typemd/` to the vault root:

```
# Before
my-vault/
├── .typemd/
│   ├── types/book/schema.yaml
│   └── properties/favorite.yaml
└── objects/

# After
my-vault/
├── types/book/schema.yaml
├── properties/favorite.yaml
├── objects/
└── .typemd/          ← internal state only
```

The principle is simple: **files you edit live at root, internal state managed by TypeMD lives in `.typemd/`**. Migration happens automatically when you open the vault — no manual file moves needed. (#362)

## Shared Properties Split into Individual Files

Previously, all shared properties were crammed into a single `properties.yaml`. Now each property gets its own file:

```yaml
# properties/due_date.yaml
type: date
emoji: "📅"
```

```yaml
# properties/priority.yaml
type: select
options:
  - value: high
  - value: medium
  - value: low
```

One property, one file — easier to find, easier to manage, better for version control. The old format migrates automatically. (#363)

## Object Archival: Soft Delete Without Regret

Objects you no longer need but can't bring yourself to permanently delete? Now you can archive them:

```bash
tmd object archive book/old-notes
```

Archived objects disappear from default queries and the TUI, but the files stay. Want to see everything? Add `--archived`. Want to restore? `tmd object unarchive`. (#34)

## Relation Graph Export

Want to visualize the relationships in your vault?

```bash
tmd graph | dot -Tpng -o graph.png
```

`tmd graph` exports all relations and wiki-links as DOT format, ready for Graphviz rendering. Use `--type` to filter by specific types. (#25)

## Template CLI

No more opening files manually to manage templates:

```bash
tmd template list                    # list all templates
tmd template show book/quick-note    # show template content
tmd template create book new-review  # create a new template
tmd template delete book/quick-note  # delete a template
```

(#248)

## `tmd log`: Git History for Objects

Want to know what changed in an object?

```bash
tmd log book/golang-in-action
```

See the git commit history for that object's file directly, without having to figure out the file path yourself. (#24)

## Markdown Rendering

The TUI body panel now renders Markdown — headings, bold, italic, code blocks, blockquotes, and lists all get proper styling and colors. (#103)

## What Else

- **`tmd type validate --watch`** — continuously monitors schemas and objects, re-validates whenever a file changes (#200)
- **SQLite fallback** — when the index file is corrupted or locked, automatically falls back to filesystem queries so your work isn't interrupted (#180)
- **TUI config page** — press `,` to open a settings page, no more manual config.yaml editing (#355)

## Migration Notes

The `types/` and `properties/` directories move to vault root automatically — TypeMD handles the migration when you open the vault. But if your `.gitignore` specifically excludes these paths, remember to update it.

## What's Next

With the structure in place, 0.9.0 will expand the type system — computed properties, aliases, and richer relation semantics.

Upgrade:

```bash
go install github.com/typemd/typemd/cmd/tmd@v0.8.0
```

Or download pre-built binaries from [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.8.0).
