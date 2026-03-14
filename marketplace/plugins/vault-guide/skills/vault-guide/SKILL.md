---
name: vault-guide
description: typemd reference guide. Use when working in a typemd vault — provides knowledge of CLI commands, type schemas, object format, wiki-links, and vault structure.
---

# typemd Reference Guide

You are working in a **typemd vault** — a local-first knowledge management system where objects (books, people, ideas) are stored as Markdown files with YAML frontmatter, connected by relations, and indexed in SQLite.

## Vault Structure

```
project-root/
├── .typemd/
│   ├── index.db             # SQLite index (auto-managed, gitignored)
│   ├── tui-state.yaml       # TUI session state (persisted on quit)
│   ├── types/               # Type schema definitions
│   │   ├── book.yaml
│   │   └── person.yaml
│   └── properties.yaml      # Shared property definitions (optional)
├── templates/               # Object templates (optional)
│   └── book/
│       └── review.md        # Template with frontmatter overrides + body
├── objects/                  # Object files organized by type
│   ├── book/
│   │   └── clean-code-01kk39c30x27ck7ahyc7ct4nyn.md
│   └── person/
│       └── robert-martin-01kk39c30y47xb1dvbs8ywqv50.md
└── ...
```

## CLI Commands

### Vault

| Command | Description |
|---------|-------------|
| `tmd init` | Initialize a new vault in the current directory |
| `tmd` (no args) | Launch the TUI |
| `tmd --vault <path>` | Specify vault directory (default: current directory) |
| `tmd --reindex` | Force reindex before running |
| `tmd --readonly` | Open vault in read-only mode (TUI only) |

### Objects

| Command | Description |
|---------|-------------|
| `tmd object create <type> [name] [-t tmpl]` | Create a new object (ULID auto-appended; name optional with name template; `-t` selects object template) |
| `tmd object show <id>` | Display object properties, relations, and body |
| `tmd object list [--json]` | List all objects |

Object IDs support **prefix matching** — `tmd object show book/clean` works if unambiguous.

### Types

| Command | Description |
|---------|-------------|
| `tmd type list` | List all type schemas (with emoji) |
| `tmd type show <type>` | Display schema definition with all properties |
| `tmd type validate` | Validate schemas, objects, relations, and wiki-links |

### Relations

| Command | Description |
|---------|-------------|
| `tmd relation link <from> <relation> <to>` | Create a relation between objects |
| `tmd relation unlink <from> <relation> <to> [--both]` | Remove a relation (`--both` removes inverse too) |

### Search & Query

| Command | Description |
|---------|-------------|
| `tmd search <keyword> [--json]` | Full-text search across filename, properties, and body |
| `tmd query "<key=value> ..." [--json]` | Filter objects by property values (AND logic) |

Examples:
```bash
tmd search "concurrency"
tmd query "type=book status=reading"
tmd query "type=book" --json
```

### Migration

| Command | Description |
|---------|-------------|
| `tmd migrate [--dry-run]` | Migrate all type schemas (e.g., `enum` → `select`) |
| `tmd migrate <type> [--dry-run]` | Migrate objects to match current schema |
| `tmd migrate <type> --rename old:new` | Rename a property across all objects of a type |

### MCP

| Command | Description |
|---------|-------------|
| `tmd mcp` | Start MCP server via stdio (for AI client integration) |

## Object File Format

Objects are Markdown files with YAML frontmatter at `objects/<type>/<slug>-<ulid>.md`.

```markdown
---
name: Clean Code
description: A handbook of agile software craftsmanship
created_at: 2024-01-15T10:30:00Z
updated_at: 2024-03-13T14:22:15Z
tags:
  - tag/programming-01kk39c30z58yc2ewct9zxrw61
  - tag/refactoring-01kk39c31a69zd3fxdu0aysx72
title: "Clean Code: A Handbook of Agile Software Craftsmanship"
author: person/robert-martin-01kk39c30y47xb1dvbs8ywqv50
status: reading
rating: 4
favorite: false
---

A practical guide to writing clean, maintainable code.

See also [[person/robert-martin-01kk39c30y47xb1dvbs8ywqv50|Robert C. Martin]]'s other works.
```

### System Properties

Always appear first in frontmatter, in this order:

| Property | Description | Auto-set |
|----------|-------------|----------|
| `name` | Display name (from slug: `clean-code` → `Clean Code`, or auto-generated from name template if defined) | Yes |
| `description` | Optional user-authored summary | No |
| `created_at` | ISO 8601 datetime, set on creation | Yes, immutable |
| `updated_at` | ISO 8601 datetime, updated on save | Yes |
| `tags` | Relation to built-in tag type, multiple | No |

These are **reserved** — type schemas cannot define properties named `name`, `description`, `created_at`, `updated_at`, or `tags`.

### Object ID Format

- **Full ID**: `book/clean-code-01kk39c30x27ck7ahyc7ct4nyn` (with ULID)
- **Display ID**: `book/clean-code` (ULID stripped)
- **ULID**: 26-character lowercase alphanumeric, sortable by timestamp

## Type Schema Format

Type schemas live in `.typemd/types/<type>.yaml`:

```yaml
name: book
plural: books            # optional: display name in collection contexts
emoji: 📚
unique: true             # optional: enforce name uniqueness across objects of this type
properties:
  - name: name           # optional: name entry with template for auto-generated names
    template: "{{.title}} by {{.author}}"
  - name: title
    type: string
    emoji: 📖
  - name: status
    type: select
    emoji: 📋
    pin: 1
    options:
      - value: to-read
      - value: reading
      - value: done
  - name: rating
    type: number
    emoji: ⭐
    pin: 2
  - name: published
    type: date
    emoji: 📅
  - use: favorite
  - name: author
    type: relation
    emoji: ✍️
    target: person
    bidirectional: true
    inverse: books
```

### Property Types

| Type | Value | Notes |
|------|-------|-------|
| `string` | Any text | — |
| `number` | Integer or float | — |
| `date` | `YYYY-MM-DD` | — |
| `datetime` | ISO 8601 / RFC 3339 | — |
| `url` | `http://` or `https://` | — |
| `checkbox` | `true` / `false` | — |
| `select` | One of `options[].value` | Single choice |
| `multi_select` | Array of `options[].value` | Multiple choices |
| `relation` | Object ID(s) | See relation options below |

### Schema-Level Options

| Option | Description |
|--------|-------------|
| `plural` | Plural display name for collection contexts (e.g., `books`) |
| `unique` | When `true`, enforces name uniqueness across objects of this type |

### Property Options

| Option | Description |
|--------|-------------|
| `emoji` | Icon displayed in TUI and CLI |
| `pin: <int>` | Display order in TUI body panel (1 = top) |
| `default` | Default value for new objects |
| `use: <name>` | Reference a shared property from `properties.yaml` |
| `template` | Only valid on the `name` property — Go template for auto-generated names (e.g., `"{{.title}} by {{.author}}"`) |

### Relation Options

| Option | Description |
|--------|-------------|
| `target: <type>` | Target type name (required) |
| `multiple: true` | Allow linking to multiple objects |
| `bidirectional: true` | Create inverse relation on target |
| `inverse: <name>` | Property name for inverse relation (required if bidirectional) |

## Shared Properties

Define reusable properties in `.typemd/properties.yaml`:

```yaml
properties:
  - name: source
    type: url
    emoji: 🔗
  - name: favorite
    type: checkbox
    emoji: ❤️
```

Reference in type schemas with `use:`. Only `pin` and `emoji` can be overridden:

```yaml
properties:
  - use: source
  - use: favorite
    pin: 1
```

## Object Templates

Object templates live at `templates/<type>/<name>.md`. They are Markdown files with frontmatter property overrides and body content that are applied during `tmd object create`.

- **Single template** for a type → auto-applied on creation
- **Multiple templates** → user is prompted to select, or use `-t <name>` flag
- Templates can set default property values and provide starter body content

## Built-in Tag Type

Tags are managed as objects of the built-in `tag` type (`unique: true`, `plural: "tags"`, emoji: 🏷️). When you reference a tag in an object's `tags` property, typemd automatically creates the tag object if it doesn't already exist. Tag objects live at `objects/tag/<name>-<ulid>.md` like any other object.

The default tag schema has two properties:

| Property | Type | Emoji | Description |
|----------|------|-------|-------------|
| `color` | `string` | 🎨 | Tag color |
| `icon` | `string` | ✨ | Tag icon |

You can override the built-in tag schema by creating `.typemd/types/tag.yaml`.

## Wiki-Links

Link to other objects in the markdown body:

| Syntax | Renders as |
|--------|-----------|
| `[[book/clean-code-01kk...]]` | `book/clean-code` |
| `[[book/clean-code-01kk...\|Clean Code]]` | `Clean Code` |

Backlinks are tracked automatically — objects know what links to them.

## TUI

The terminal UI has a three-panel layout (sidebar, body, properties) with a **right panel mode** system:

1. **Sidebar (Left)**: Object list grouped by type (collapsible), with `+ New Type` at the bottom
2. **Body (Center)**: Pinned properties + markdown body (when viewing objects), or type schema editor (when viewing types)
3. **Properties (Right)**: All properties, relations, backlinks (toggleable)
4. **Status Bar**: Mode, unsaved indicator, messages

The right panel follows the sidebar cursor:
- Selecting an **object** → shows object detail (body + properties)
- Selecting a **type header** → shows the type schema editor with full CRUD support (add/edit/reorder/delete properties)

### Key Keybindings

| Key | Action |
|-----|--------|
| `j`/`k` or `↑`/`↓` | Navigate |
| `enter`/`space` | Select / expand group |
| `tab` | Switch panel focus |
| `e` | Enter edit mode |
| `esc` | Exit edit mode (saves) |
| `n` | Create new object (on type header or object) |
| `/` | Search |
| `p` | Toggle properties panel |
| `w` | Toggle soft-wrap |
| `[`/`]` | Resize panel |
| `?` or `h` | Help |
| `q` | Quit |

## Tips for Working with typemd

- **Always read type schemas first** before creating or modifying objects — check `.typemd/types/*.yaml`
- **Use `tmd type validate`** after making changes to catch schema/object/relation errors
- **Object IDs include ULIDs** — don't guess them, use `tmd object list` to find exact IDs
- **Relations are properties** — set them in frontmatter, not as wiki-links (wiki-links are for body text references)
- **System properties are managed** — you can set `name` and `description`, but `created_at` is immutable and `updated_at` auto-updates on save
- **Prefix matching works** in CLI — `tmd object show book/clean` resolves if unambiguous
