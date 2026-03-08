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
- **TUI** — Three-panel interface powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea), with auto-refresh on file changes
- **MCP Server** — integrate with AI assistants via Model Context Protocol
- **Local-first** — everything lives on your machine as plain Markdown files

## Data Structure

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
        └── alan-donovan-01jqr3k8yznw2a4dbx6t7c9fpq.md
```

Objects are stored as Markdown files with YAML frontmatter. Each directory under `objects/` is a **type namespace** — different types can share the same slug.

The full Object ID is `type/<slug>-<ulid>`, e.g. `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`. A [ULID](https://github.com/ulid/spec) is automatically appended to every new object to guarantee uniqueness.

## Usage

```bash
# Initialize a new vault
tmd init

# Open TUI (current directory)
tmd

# Open TUI with specific vault path
tmd --vault /path/to/vault

# Create a new object (ULID is appended automatically)
tmd object create book clean-code
# → Created book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz

# Show object detail (use the full ID from create output)
tmd object show book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz

# List all objects
tmd object list
tmd object list --json

# Query by type and property
tmd query "type=book status=reading"
tmd query "type=book" --json

# Full-text search
tmd search "concurrency"

# Link two objects (use full IDs)
tmd relation link book/golang-in-action-01jqr3k5mp... author person/alan-donovan-01jqr3k8yz...

# Unlink (with --both to remove inverse side too)
tmd relation unlink book/golang-in-action-01jqr3k5mp... author person/alan-donovan-01jqr3k8yz... --both

# Sync files to DB and rebuild search index (only needed after manual edits)
tmd reindex

# Validate schemas, objects, and relations
tmd type validate

# Show type schema details
tmd type show book

# List all available types
tmd type list

# Start MCP server for AI integration
tmd mcp
tmd mcp --vault /path/to/vault
```

### `tmd object show` Output

```
book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz

Properties
──────────
  title: Go in Action
  status: reading
  rating: 4.5
  author: → person/alan-donovan-01jqr3k8yznw2a4dbx6t7c9fpq

Body
────
  # Notes
  A great book about Go...
```

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

The properties panel is hidden by default and can be toggled with `p`. On narrow terminals (< 56 columns), it auto-hides.

### TUI Controls

| Key | Action |
|-----|--------|
| `↑`/`k`, `↓`/`j` | Navigate object list |
| `Enter`/`Space` | Select object / Toggle group |
| `Tab` | Cycle focus between panels |
| `e` | Enter edit mode (body or properties panel) |
| `/` | Search (FTS5 full-text search) |
| `Esc` | Exit edit mode (auto-saves if changed) / Clear search results |
| `p` | Toggle properties panel |
| `w` | Toggle soft wrap |
| `[`/`]` | Shrink/grow focused panel |
| `?`/`h` | Open help popup |
| `q`/`Ctrl+C` | Quit |

The status bar shows the current mode: `[VIEW]` for normal navigation and `[EDIT]` when editing is active.

When exiting edit mode, changes are automatically saved to the `.md` file and the SQLite index is updated. If the file was modified externally while editing, a `[CONFLICT]` prompt appears — press `y` to overwrite, `n` to reload from disk, or `Esc` to cancel.

The TUI automatically watches the `objects/` directory and refreshes when files are created, modified, or deleted.

## Type Schema

Define your own types in `.typemd/types/`:

```yaml
# .typemd/types/book.yaml
name: book
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
  - name: status
    type: enum
    values: [to-read, reading, done]
    default: to-read
  - name: rating
    type: number
```

Properties support an optional `default` field to specify a default value.

## Relations

Relations are defined as `type: relation` properties within type schemas. Use `bidirectional` and `inverse` to auto-sync both sides:

```yaml
# .typemd/types/person.yaml
name: person
properties:
  - name: name
    type: string
  - name: books
    type: relation
    target: book
    multiple: true
    bidirectional: true
    inverse: author
```

When `bidirectional: true`, linking a book to a person via `author` automatically updates both the book's `author` and the person's `books` property.

## MCP Server

Run `tmd mcp` to start a [Model Context Protocol](https://modelcontextprotocol.io) server over stdio. AI clients (e.g. Claude Code) can query your vault through these tools:

| Tool | Description |
|------|-------------|
| `search` | Full-text search objects, returns ID, type, and filename |
| `get_object` | Get full object detail by ID, including properties and body |

## Architecture

TypeMD is a monorepo with a shared Go core and multiple interfaces:

```
typemd/
├── core/       # Core library — objects, types, relations, index
├── cmd/        # CLI commands (Cobra)
├── tui/        # Terminal UI (Bubble Tea)
├── mcp/        # MCP server for AI integration
├── web/        # Web UI API (planned)
├── site/       # Official website (Astro) → typemd.io
├── docs/       # Documentation (Starlight) → docs.typemd.io
└── app/        # Desktop app (planned)
```

All interfaces share the same `core` library.

## Tech Stack

- **Language**: Go
- **TUI**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss)
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
