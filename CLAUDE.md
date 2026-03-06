# CLAUDE.md

## Project Overview

typemd is a local-first CLI knowledge management tool. Objects (books, people, ideas) are stored as Markdown files with YAML frontmatter, connected by Relations. SQLite provides indexing.

## Architecture

- **core/** — Core library: objects, types, relations, index
- **cmd/** — CLI commands (Cobra)
- **tui/** — Terminal UI (Bubble Tea)
- **mcp/** — MCP server
- **web/** — Web UI (future)
- **app/** — Desktop app via Wails (future)

## Data Model

- Objects identified by `type/filename` (e.g. `book/golang-in-action`)
- Type schemas: `.typemd/types/*.yaml`
- Relations defined as properties in type schemas
- SQLite index: `.typemd/index.db`
- Object files: `objects/<type>/<name>.md`

## Build & Test

```bash
go build ./...
go test ./...
go run ./cmd/tmd [command]
```

## Testing

- Tests live in the same package as implementation
- Use `t.TempDir()` for isolation
- Cover happy path and error/edge cases
