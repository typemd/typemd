# CLAUDE.md

## Project Overview

typemd is a local-first CLI knowledge management tool. Objects (books, people, ideas) are stored as Markdown files with YAML frontmatter, connected by Relations. SQLite provides indexing.

## Architecture

- **core/** — Core library: objects, types, relations, index
  - **core/ai/** — AI provider abstraction and service layer (Claude CLI integration)
- **cmd/** — CLI commands (Cobra)
- **tui/** — Terminal UI (Bubble Tea)
  - **tui/widget/** — Shared UI primitives (CenteredPopup, OverlayPopup, ToastModel via Layer/Compositor, scroll) used across TUI components
- **mcp/** — MCP server
- **web/** — Web UI: Go HTTP server (`tmd serve`) + Vue 3 frontend (Vite + Tailwind CSS)
  - **web/frontend/** — Vue 3 SPA with vault adapter pattern for API abstraction
- **app/** — Desktop app via Wails + shared Vue 3 frontend (future)
- **websites/** — Non-Go websites (site, docs, blog)
- **marketplace/** — Claude Code marketplace plugins (typemd plugin with vault-guide, instructions-guide, explore, importer, and onboarding skills)

## Core Package Architecture

The `core/` package follows **Clean Architecture** with **CQRS** (Command Query Responsibility Segregation). The design separates concerns into layers with clear dependency rules.

### Layer Diagram

```mermaid
graph TB
    subgraph Consumers
        CMD[cmd/ — CLI]
        TUI[tui/ — Terminal UI]
        MCP[mcp/ — MCP Server]
        WEB[web/ — Web UI]
    end

    subgraph Facade
        V[Vault]
    end

    subgraph Use Cases
        OS[ObjectService — commands]
        QS[QueryService — queries]
        PJ[Projector — file→index sync]
    end

    subgraph Domain
        OBJ[Object — aggregate root]
        TS[TypeSchema — aggregate root]
        OID[ObjectID — value object]
        EVT[DomainEvent — event types]
        ED[EventDispatcher]
    end

    subgraph Infrastructure
        REPO[ObjectRepository — interface]
        IDX[ObjectIndex — interface]
        LR[LocalObjectRepository — files]
        SI[SQLiteObjectIndex — SQLite]
    end

    CMD --> V
    TUI --> V
    MCP --> V
    WEB --> V

    V --> OS
    V --> QS
    V --> PJ
    V --> ED

    OS --> REPO
    OS --> IDX
    OS --> ED
    QS --> REPO
    QS --> IDX
    PJ --> REPO
    PJ --> IDX

    OS --> OBJ
    OS --> TS
    QS --> OBJ

    REPO -.-> LR
    IDX -.-> SI

    OBJ --> OID
    OBJ --> EVT
```

### Key Design Decisions

- **ObjectRepository** returns domain entities (`*Object`, `*TypeSchema`), not raw bytes. Path conventions and serialization are encapsulated in implementations.
- **ObjectIndex** returns `ObjectResult` (lightweight projection) for search results. Full entity retrieval goes through `ObjectRepository.Get(id)`.
- **Vault** is a thin facade / DI container. Object business logic lives in `ObjectService` (commands) and `QueryService` (queries). Type schema CRUD (`SaveType`, `DeleteType`, `CountObjectsByType`) lives directly on Vault since it delegates to `ObjectRepository` without needing a separate service layer.
- **Domain Events** follow "entity produces → use case dispatches" pattern. Entity methods return `DomainEvent`; services collect and dispatch after successful operations.
- **Files are the source of truth**. SQLite index is an acceleration layer maintained by the `Projector`.

### Key Entry Points

- `vault.go` — Vault facade + lifecycle (Open/Close/Init)
- `object_service.go` — command use cases (create/save/link)
- `query_service.go` — query use cases (search/filter/stats)
- `reconciler.go` — file normalization + event emission (full + incremental); `reconciler_relation.go` — relation resolution + tag sync; `reconciler_wikilink.go` — wiki-link sync
- `projector.go` — event-driven SQLite index writer
- `config.go` — VaultConfig (date formats, CLI/TUI/AI config sections)

### TUI Architecture

The TUI uses a three-panel layout (sidebar, body, properties) with **focus mode** (`.` key) for single full-width body. The right panel follows the sidebar cursor via a **right panel mode** system: `panelObject`, `panelTypeEditor`, `panelTemplate`, `panelView` (full-width table/list), `panelStats`, `panelSchemaExplore`, `panelConfig` (config settings page). Keybindings are defined in `tui/keys.go`; see `tui/help.go` for the help popup or the [docs site TUI page](websites/docs/src/content/docs/tui/tui.md) for the full keybinding table.

Key sub-models: `typeEditor` (schema editing + wizard + templates), `viewMode` (table/list views with inline cell editing via `cellEdit`), `propEditor` (inline property editing with type-appropriate widgets), `dateEdit` (segmented input + calendar picker), `configEditor` (config settings editor with two-column category/settings layout), `widget.ToastModel` (transient notifications). File watcher monitors `objects/` and `types/` with debounced incremental sync.

## Data Model

- Objects identified by `type/<slug>-<ulid>` (e.g. `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`)
- All objects have system properties managed by typemd. **Stored** (frontmatter): `name` (preserves original input on creation; auto-populated from slug for pre-slugified names, or from name template if defined), `description` (optional, user-authored), `created_at` (set on creation, immutable), `updated_at` (updated on save, immutable), `tags` (relation to built-in `tag` type, multiple), `locked` (boolean, user-authored, prevents editing when true), `archived` (boolean, user-authored, hides object from default queries when true). Stored properties appear first in frontmatter in that order. **Derived** (stable values inferred from structure/metadata, not stored in frontmatter): `object_type`, `created_by`. **Computed** (dynamic values requiring content parsing or index queries, not stored in frontmatter): `links`, `backlinks`, `updated_by`. Both derived and computed properties are read-only. System properties are either **user-authored** (`name`, `description`, `tags`, `locked`, `archived` — can be overridden by templates) or **auto-managed** (`created_at`, `updated_at` — cannot be overridden).
- Type schemas: `types/<name>/schema.yaml`. Reserved system properties (`name`, `description`, `created_at`, `updated_at`, `tags`, `locked`, `archived`, `object_type`, `links`, `backlinks`, `created_by`, `updated_by`) cannot be redefined; `name` can appear with only a `template` field. Optional schema fields: `plural`, `emoji`, `unique`, `version`, `color`, `description`.
- Views: `types/<name>/views/<view>.yaml` (optional). Two layouts: `list` and `table`. Supports `columns`, `filter`, `sort`, `group_by`. Each type has an implicit default view (list, sort by name asc).
- Built-in types: `tag` (🏷️, plural "tags", unique, backs `tags` system property, has `color` and `icon` string properties), `page` (📄, plural "pages", general-purpose content container), and `source` (📥, plural "sources", tracks ingested raw materials, has `url` string, `author` string, and `ingested_at` date properties). Built-in types exist without YAML files, cannot be deleted, but can be overridden by custom `types/<name>/schema.yaml`.
- Shared properties: `properties/<name>.yaml` (optional per-property files, defines reusable property definitions referenced via `use` in type schemas; property name derived from filename; `use` entries can override `pin`, `emoji`, and `description`)
- Relations defined as properties in type schemas
- Wiki-links: `[[type/name-ulid]]`, `[[type/name]]`, or `[[name]]` syntax in markdown body, with backlink tracking. Shorthand forms are resolved during sync and written back as full IDs.
- SQLite index: `.typemd/index.db`
- TUI session state: `.typemd/tui-state.yaml` (persisted on quit, restored on launch)
- Vault config: `.typemd/config.yaml` — sections: `date_format`/`datetime_format`, `cli.*`, `tui.*` (toast, debounce, theme), `ai.*` (providers, prompts). See `config.go` for full key registry.
- Embedded skills: `core/skills/*/SKILL.md` (via `//go:embed`); vault overrides in `.typemd/instructions/<skill>.md`
- Starter types: `core/starters/*.yaml` (offered during `tmd init`)
- Object templates: `templates/<type>/<name>.md` (applied during `tmd object create`)
- Object files: `objects/<type>/<name>.md`

## Web UI Architecture

- **`tmd serve`** starts a Go HTTP server with REST API + embedded Vue 3 SPA
- **REST API** (`web/server.go`): endpoints under `/api/` for types, objects, properties, templates (read-write)
- **Vue 3 frontend** (`web/frontend/`): Vite + Tailwind CSS, three-panel layout mirroring the TUI (sidebar, body, properties), three-theme system (warm/dark/light) with CSS custom properties
- **Vault adapter** (`web/frontend/src/lib/vault.js`): all API calls go through this single module; swap implementation for different backends
- **Embedded frontend** (`web/frontend.go`): `//go:embed frontend/dist` bakes the built SPA into the Go binary; `tmd serve` serves it for non-API routes with SPA fallback
- **Dev mode**: run `tmd serve` (port 3000) + `cd web/frontend && npm run dev` (port 5173, Vite proxies `/api` to 3000)
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
make test                  # builds frontend + go build + go test + go vet
make build-frontend        # only build web/frontend/dist
```

The Go binary embeds `web/frontend/dist` via `//go:embed`. Frontend must be built before `go build`.

## Debugging

- **CLI**: `tmd --debug <command>` enables DEBUG-level JSON logging to stderr. Useful for inspecting sync flow, query execution, and AI provider calls.
- **TUI**: Always writes DEBUG-level JSON logs to `.typemd/logs/{YYYY-MM-DD}.log`. Check this file to diagnose TUI issues (sync, file watcher, property editing) without corrupting terminal output.

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
| `cmd/`  | BDD (`cmd/features/`) + unit tests |
| `mcp/`  | Unit tests — BDD TBD |
