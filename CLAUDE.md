# CLAUDE.md

## Project Overview

typemd is a local-first CLI knowledge management tool. Objects (books, people, ideas) are stored as Markdown files with YAML frontmatter, connected by Relations. SQLite provides indexing.

## Architecture

- **core/** — Core library: objects, types, relations, index
- **cmd/** — CLI commands (Cobra)
- **tui/** — Terminal UI (Bubble Tea)
- **mcp/** — MCP server
- **web/** — Web UI: React + shadcn/ui (future)
- **app/** — Desktop app via Wails + shared React frontend (future)
- **websites/** — Non-Go websites (site, docs, blog)
- **marketplace/** — Claude Marketplace plugins (future)

## Data Model

- Objects identified by `type/<slug>-<ulid>` (e.g. `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`)
- All objects have system properties managed by typemd: `name` (auto-populated from slug, or from name template if defined), `description` (optional, user-authored), `created_at` (set on creation), `updated_at` (updated on save), `tags` (relation to built-in `tag` type, multiple). These appear first in frontmatter in that order.
- Type schemas: `.typemd/types/*.yaml` (cannot define properties named `description`, `created_at`, `updated_at`, or `tags` — they're reserved system properties; `name` can appear in `properties` with only a `template` field for auto-generated names). Type schemas support optional `plural` (for display in collection contexts) and `unique` (to enforce name uniqueness) fields.
- Shared properties: `.typemd/properties.yaml` (optional, defines reusable property definitions referenced via `use` in type schemas)
- Relations defined as properties in type schemas
- Wiki-links: `[[type/name-ulid]]` syntax in markdown body, with backlink tracking
- SQLite index: `.typemd/index.db`
- TUI session state: `.typemd/tui-state.yaml` (persisted on quit, restored on launch)
- Object files: `objects/<type>/<name>.md`

## Web UI Architecture

- **Shared frontend**: `tmd serve`, try.typemd.io, and desktop app (Wails) share one React + shadcn/ui frontend
- **Storage Interface**: Frontend talks to a `VaultStorage` abstraction
  - `tmd serve` → Go HTTP API (read-write)
  - try.typemd.io → GitHub REST API from browser, no backend (read-only initially, read-write later)
  - Wails → Go bindings (read-write)
- **No SQLite in browser**: try.typemd.io uses in-memory index built from GitHub API responses
- **Design principle**: SQLite is optional acceleration, not a hard dependency — files are always the source of truth

## Language Convention

**English is the primary language** for all project artifacts:

- **Issues** — titles, descriptions, comments
- **Commits** — commit messages and bodies
- **Skills** — skill content in `.claude/skills/`
- **Releases** — release notes and CHANGELOG

Blog posts are the exception: written in Traditional Chinese (zh-tw) first, then synced to English via the `sync-blog` skill.

## Build & Test

```bash
go build ./...
go test ./...
go run ./cmd/tmd [command]
```

## Testing

This project uses two layers of testing:

- **BDD (Godog)** — Define behaviors, establish shared vocabulary, and describe what a feature does from the user's perspective. Gherkin `.feature` files live in `<package>/features/`. BDD scenarios focus on **what**, not implementation details.
- **Unit tests** — Verify precise logic: edge cases, output formats, exact values, error conditions. Traditional Go `testing` style.

When deciding where a test belongs: if it defines a behavior or names a concept, write a BDD scenario. If it validates an implementation detail (e.g. JSON format, lowercase ULID, flag edge cases), write a unit test.

### BDD scope by package

| Package | Testing approach |
|---------|-----------------|
| `core/` | BDD (`core/features/`) + unit tests |
| `tui/`  | BDD (`tui/features/`, planned) + unit tests |
| `web/`  | BDD (`web/features/`, future) |
| `cmd/`  | Minimal — CLI commands delegate to `core/`, covered by core BDD scenarios |
| `mcp/`  | Unit tests — BDD TBD |
