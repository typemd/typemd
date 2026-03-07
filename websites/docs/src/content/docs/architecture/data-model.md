---
title: Data Model
description: How TypeMD stores and indexes data.
sidebar:
  order: 1
---

## Storage

Objects are stored as Markdown files with YAML frontmatter under `objects/<type>/`. The full Object ID is `type/<slug>-<ulid>`, e.g. `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`.

When created via the CLI (`tmd create`), a 26-character lowercase ULID is automatically appended to the slug for uniqueness. Objects created manually (without the CLI) do not require a ULID — this is backward compatible.

```
vault/
├── .typemd/
│   ├── types/              # type schema definitions (YAML)
│   │   ├── book.yaml
│   │   └── person.yaml
│   └── index.db            # SQLite index (auto-updated)
└── objects/
    ├── book/
    │   └── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
    └── person/
        └── alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz.md
```

## Indexing

TypeMD uses SQLite with FTS5 for indexing. The index is stored at `.typemd/index.db` and contains:

- Object metadata (type, filename, properties)
- Wiki-link records extracted from object body content (for backlink tracking)
- Full-text search index over filenames, properties, and body content

The index is automatically synced when opening a vault with an empty or missing database (e.g. after a fresh clone). It is also kept up-to-date when using the TUI or CLI commands. Use `tmd reindex` to rebuild after manual file edits outside of TypeMD.

## Architecture

TypeMD is a monorepo with a shared Go core and multiple interfaces:

```
typemd/
├── core/       # Core library — objects, types, relations, index
├── cmd/        # CLI commands (Cobra)
├── tui/        # Terminal UI (Bubble Tea)
├── mcp/        # MCP server for AI integration
├── web/        # Web UI API (planned)
├── app/        # Desktop app (planned)
├── site/       # Official website (Astro) → typemd.io
└── docs/       # Documentation (Starlight) → docs.typemd.io
```

All interfaces share the same `core` library.

## Tech Stack

- **Language**: Go
- **TUI**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **MCP**: [mcp-go](https://github.com/mark3labs/mcp-go) — Model Context Protocol server
- **Index**: SQLite with FTS5 full-text search
- **Storage**: Markdown + YAML frontmatter
