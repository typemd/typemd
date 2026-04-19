> 🌐 [English](README.md) | [繁體中文](README.zh-TW.md)

<p align="center">
  <img src="websites/docs/src/assets/icon.svg" width="120" alt="TypeMD icon">
</p>

<h1 align="center">TypeMD</h1>

<p align="center">
  A local-first CLI knowledge management tool inspired by <a href="https://anytype.io">Anytype</a> and <a href="https://capacities.io">Capacities</a>.
</p>

<p align="center">
  <a href="https://typemd.io">Website</a> · <a href="https://docs.typemd.io">Docs</a> · <a href="https://github.com/typemd/typemd">GitHub</a>
</p>

Your knowledge base is made of **Objects** — not files. Markdown is just the storage format.

## Philosophy

Most note-taking tools make you think like a computer: files, folders, hierarchies.

TypeMD lets you think in **Objects** — books, people, ideas, meetings — connected by **Relations**. The structure emerges from your knowledge, not from a folder tree.

## Features

- **Typed Objects** — define schemas for each type (Book, Person, Idea, etc.)
- **Structured Relations** — connect objects with named, optionally bidirectional links
- **Wiki-links & Backlinks** — link objects inline with `[[type/name-ulid]]` syntax, with automatic backlink tracking
- **Full-text search** — find anything across your vault
- **Structured queries** — filter objects by type, property, or relation
- **TUI** — Three-panel interface powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea), with inline property editing and auto-refresh on file changes
- **Web UI** — Browser-based interface via `tmd serve`, with three-panel layout, property editing, and object creation
- **AI Assist** — generate descriptions, suggest tags, and explore schemas with Claude AI in the TUI
- **MCP Server** — integrate with AI assistants via Model Context Protocol
- **Local-first** — everything lives on your machine as plain Markdown files

## Data Structure

```
vault/
├── .typemd/
│   ├── index.db            # SQLite index (auto-updated)
│   └── tui-state.yaml      # TUI session state (auto-saved)
├── types/                  # user-defined type schemas (directory format)
│   ├── book/
│   │   └── schema.yaml     # you create these
│   └── person/
│       └── schema.yaml
├── properties/
│   ├── <name>.yaml         # shared property definitions (optional, one per property)
│   └── ...
├── templates/              # object templates by type (optional)
│   └── book/
│       └── review.md       # template with default frontmatter + body
└── objects/
    ├── book/
    │   └── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
    └── person/
        └── alan-donovan-01jqr3k8yznw2a4dbx6t7c9fpq.md
```

Objects are stored as Markdown files with YAML frontmatter. Each directory under `objects/` is a **type namespace** — different types can share the same slug.

The full Object ID is `type/<slug>-<ulid>`, e.g. `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`. A [ULID](https://github.com/ulid/spec) is automatically appended to every new object to guarantee uniqueness.

## Installation

```bash
brew install typemd/tap/typemd-cli
```

Or from source:

```bash
go install github.com/typemd/typemd/cmd/tmd@latest
```

Pre-built binaries are also available on [GitHub Releases](https://github.com/typemd/typemd/releases).

## Usage

```bash
# Initialize a new vault
tmd init

# Open TUI
tmd

# Create objects (names are auto-slugified, ULID appended)
tmd object create book "Clean Code"
tmd object create "Some Thought"          # uses default type from config

# Search and explore
tmd search "concurrency"
tmd object show book/clean-code           # prefix matching, no full ULID needed
tmd object list

# Connect objects (name references in frontmatter auto-expand on sync)
tmd relation link book/golang-in-action author person/alan-donovan

# Manage templates
tmd template list                         # list all templates
tmd template show book/review             # view template content
tmd template create book/review           # create and open in $EDITOR
tmd template delete book/review           # delete with confirmation

# Maintenance
tmd format                                # normalize frontmatter & schema formatting
tmd doctor                                # vault health check
tmd stats                                 # vault-wide statistics

# Web UI
tmd serve                                 # starts web UI at http://localhost:3000
tmd serve -p 8080                         # custom port

# MCP server for AI integration
tmd mcp
```

See `tmd --help` and [docs](https://docs.typemd.io) for the full command reference.

### TUI

```
┌─ Objects ─────────┐  ┌─ Body ─────────────┐  ┌─ Properties ──────┐
│ ▼ book (2)        │  │ # Notes            │  │ title: Go in      │
│   golang-in-action│  │ A great book about │  │   Action          │
│   clean-code      │  │ Go...              │  │ status: reading   │
│ ▶ person (1)      │  │                    │  │ author:           │
│ ▶ note (3)        │  │                    │  │   → person/alan   │
│                   │  │                    │  │                   │
│                   │  │                    │  │                   │
│                   │  │                    │  │                   │
└───────────────────┘  └────────────────────┘  └───────────────────┘
```

Press `?` in the TUI for the full keybinding reference.

## Type Schema

Define your types in `types/` (`tag` and `page` are built-in — all others are user-defined):

```yaml
# types/book/schema.yaml
name: book
plural: books
emoji: 📚
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
  - name: status
    type: select
    options:
      - value: to-read
      - value: reading
      - value: done
    default: to-read
  - name: rating
    type: number
```

Relations are defined as `type: relation` properties. Use `bidirectional` and `inverse` to auto-sync both sides. See [docs](https://docs.typemd.io) for full schema reference.

## MCP Server

Run `tmd mcp` to start a [Model Context Protocol](https://modelcontextprotocol.io) server over stdio. AI clients (e.g. Claude Code) can read and write your vault through these tools:

| Tool | Description |
|------|-------------|
| `search` | Full-text search objects, returns ID, type, and filename |
| `get_object` | Get full object detail by ID, including properties and body |
| `list_types` | List all available type schemas with metadata |
| `vault_overview` | One-call vault summary: per-type count, emoji, description, and recent objects |
| `list_objects` | List object summaries with optional `type` filter and pagination |
| `query_objects` | Structured query with `filters`, optional `sort`, and pagination |
| `list_backlinks` | Return wiki-link and typed-relation backlinks for an object |
| `vault_stats` | Per-property distribution stats for a single type (filled, fill rate) |
| `create_object` | Create a new object with type, name, optional template, properties, and body |
| `update_object` | Update an object's properties (merge) and/or body (replace) |
| `link_objects` | Create a relation between two objects |
| `unlink_objects` | Remove a relation between two objects |

## Architecture

TypeMD is a monorepo with a shared Go core and multiple interfaces:

```
typemd/
├── core/       # Core library — objects, types, relations, index
├── cmd/        # CLI commands (Cobra)
├── tui/        # Terminal UI (Bubble Tea)
├── mcp/        # MCP server for AI integration
├── web/        # Web UI — Go HTTP server + React frontend
├── site/       # Official website (Astro) → typemd.io
├── docs/       # Documentation (Starlight) → docs.typemd.io
└── app/        # Desktop app (planned)
```

All interfaces share the same `core` library.

## Tech Stack

- **Language**: Go
- **TUI**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Web UI**: React + [Tailwind CSS](https://tailwindcss.com) + [Vite](https://vite.dev) (embedded in Go binary)
- **MCP**: [mcp-go](https://github.com/mark3labs/mcp-go) — Model Context Protocol server
- **Index**: SQLite with FTS5 full-text search
- **Storage**: Markdown + YAML frontmatter

## Resources

- [CHANGELOG](CHANGELOG.md)
- [CONTRIBUTING](CONTRIBUTING.md)
- [Blog](https://blog.typemd.io)

## Inspiration

- [Anytype](https://anytype.io) — encrypted, local-first alternative to cloud apps
- [Capacities](https://capacities.io) — object-based knowledge studio
