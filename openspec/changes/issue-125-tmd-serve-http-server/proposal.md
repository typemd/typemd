## Why

typemd currently only supports TUI and CLI interactions. Users need a browser-based interface to browse, edit, and create objects — especially useful for users who prefer GUI over terminal, or when sharing a vault view on a local network. The `tmd serve` command starts a local HTTP server that wraps the same `core.Vault` API, enabling read-write access through a React frontend.

## What Changes

- Add `tmd serve` CLI command that starts a local HTTP server (default port 3000)
- Add REST API endpoints in `web/server.go` for types, objects, properties, and templates
- Add React frontend in `web/frontend/` with three-panel layout (sidebar, body, properties) mirroring the TUI
- Embed built frontend into the Go binary via `//go:embed` for single-binary deployment
- Add vault adapter pattern (`web/frontend/src/lib/vault.js`) abstracting API calls for future backend swaps
- Support inline property editing, body editing, and object creation in the browser
- Add `build-frontend` target to Makefile; update CI and GoReleaser to build frontend before Go binary

## Capabilities

### New Capabilities

- `web-serve`: HTTP server lifecycle, REST API endpoints, embedded frontend serving with SPA fallback
- `web-ui-layout`: Browser-based three-panel layout (sidebar, body, properties) with keyboard shortcuts and focus mode
- `web-object-crud`: Object browsing, property editing, body editing, and object creation via the web interface

### Modified Capabilities

## Impact

- **New files**: `web/server.go`, `web/frontend.go`, `cmd/serve.go`, `web/frontend/` (React app)
- **Modified files**: `Makefile`, `.github/workflows/ci.yml`, `.goreleaser.yml`, `CLAUDE.md`, `README.md`, `README.zh-TW.md`, `CONTRIBUTING.md`
- **New dependencies**: React 19, Tailwind CSS 4, Vite 6 (frontend only, not Go dependencies)
- **Build process**: `go build` now requires `web/frontend/dist/` to exist (embedded via `//go:embed`)
