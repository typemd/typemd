## Context

typemd's core library (`core.Vault`) is consumed by CLI, TUI, and MCP. The web UI adds a fourth consumer: a Go HTTP server that exposes `Vault` operations as REST endpoints, with a React SPA served from the same binary.

The existing architecture already supports this — `Vault` is a facade with no UI coupling. The web server wraps it in HTTP handlers, and the frontend calls the API through an adapter module.

## Goals / Non-Goals

**Goals:**
- `tmd serve` starts an HTTP server serving both REST API and embedded frontend
- Three-panel layout (sidebar, body, properties) matching TUI's information architecture
- Read-write CRUD: browse objects, edit properties inline, edit body, create objects
- Single binary deployment — frontend embedded via `//go:embed`
- Adapter pattern for API calls — enable future backend swaps (static, Wails)

**Non-Goals:**
- Authentication / multi-user access
- Real-time sync (file watcher → websocket push)
- Markdown rendering (body displayed as raw text)
- Type schema editing, view configuration, AI features
- Static site export (`tmd export --static`)
- Share frontend code with Wails desktop app (deferred)

## Decisions

### REST API over RPC-style

REST (`GET /api/objects/{type}/{slug}`) over RPC (`POST /api/getObject`). REST is more natural for the resource-oriented data model, enables browser caching, and is easier to explore with curl.

### Standard library `net/http` over frameworks

Go 1.22+ `http.ServeMux` supports method-based routing (`GET /api/types/{name}`), eliminating the need for gin/chi. Fewer dependencies, simpler code.

### Object ID in URL as `{type}/{slug}` path segments

Object IDs contain a slash (`book/clean-code-01xxx`). Go's `{id...}` catch-all wildcard can only appear at the end of a pattern, so `/api/objects/{id...}/properties` is invalid. Solution: split into two path parameters `{type}/{slug}`. Properties get their own top-level route prefix (`/api/properties/{type}/{slug}`) instead of nesting under objects.

### Embedded frontend with SPA fallback

`//go:embed frontend/dist` bakes the built React app into the binary. The server tries exact file matches first, then falls back to `index.html` for client-side routing. In dev mode, Vite runs separately with proxy to the Go server.

### Vault adapter pattern in frontend

All API calls go through `src/lib/vault.js`. Components never call `fetch` directly. This enables swapping the HTTP adapter for Wails bindings (desktop) or preloaded JSON (static export) without changing any component code.

## Risks / Trade-offs

- **No file watcher** — Changes made outside the web UI (e.g., in a text editor) won't appear until the page is refreshed. Acceptable for MVP; websocket push can be added later.
- **N+1 in ListTypes** — `handleListTypes` calls `LoadType()` and `CountObjectsByType()` per type in a loop. Acceptable for small vaults (< 50 types); optimize with batch query if needed.
- **Frontend embedded in binary** — `go build` requires `web/frontend/dist/` to exist. CI and Makefile handle this, but developers must run `make build-frontend` before `go build`.
