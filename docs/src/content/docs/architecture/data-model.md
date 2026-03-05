---
title: Data Model
description: How TypeMD stores and indexes data.
sidebar:
  order: 1
---

## Storage

Objects are stored as Markdown files with YAML frontmatter under `objects/<type>/`. The full Object ID is `type/filename`, e.g. `book/golang-in-action`.

```
vault/
├── .typemd/
│   ├── types/              # type schema definitions (YAML)
│   │   ├── book.yaml
│   │   └── person.yaml
│   └── index.db            # SQLite index (auto-updated)
└── objects/
    ├── book/
    │   └── golang-in-action.md
    └── person/
        └── alan-donovan.md
```

## Indexing

TypeMD uses SQLite with FTS5 for indexing. The index is stored at `.typemd/index.db` and contains:

- Object metadata (type, filename, properties)
- Full-text search index over filenames, properties, and body content

The index is automatically updated when using the TUI or CLI commands. Use `tmd reindex` to rebuild after manual file edits.

## Architecture

TypeMD is a monorepo with a shared Go core and multiple interfaces:

```
typemd/
├── core/       # Core library — objects, types, relations, index
├── cmd/        # CLI commands (Cobra)
├── tui/        # Terminal UI (Bubble Tea)
├── mcp/        # MCP server for AI integration
├── site/       # Official website (Astro)
└── docs/       # Documentation (Starlight)
```

All interfaces share the same `core` library.

## Tech Stack

- **Language**: Go
- **TUI**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **MCP**: [mcp-go](https://github.com/mark3labs/mcp-go) — Model Context Protocol server
- **Index**: SQLite with FTS5 full-text search
- **Storage**: Markdown + YAML frontmatter
