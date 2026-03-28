## 1. Go HTTP Server

- [x] 1.1 Create `web/server.go` with REST API handlers wrapping `core.Vault`
- [x] 1.2 Implement DTO types for JSON serialization (`typeItem`, `objectItem`, `objectDetail`, `displayProp`)
- [x] 1.3 Implement property value parsing (`parsePropertyValue`) with `strconv` for type safety
- [x] 1.4 Implement time serialization (`serializeValue`) for `time.Time` → RFC3339

## 2. CLI Command

- [x] 2.1 Create `cmd/serve.go` with `tmd serve` cobra command and `--port` flag
- [x] 2.2 Wire server to `openVault()` helper and pass embedded frontend

## 3. Embedded Frontend

- [x] 3.1 Create `web/frontend.go` with `//go:embed frontend/dist` and `FrontendFS()` helper
- [x] 3.2 Add SPA fallback in server routes (try exact file, fall back to `index.html`)

## 4. React Frontend Setup

- [x] 4.1 Scaffold React + Tailwind CSS 4 + Vite project in `web/frontend/`
- [x] 4.2 Create vault adapter (`src/lib/vault.js`) abstracting all API calls
- [x] 4.3 Configure Vite proxy for dev mode (`/api` → `localhost:3000`)

## 5. Three-Panel Layout

- [x] 5.1 Implement `App.jsx` with sidebar / body / properties layout and keyboard shortcuts
- [x] 5.2 Implement `Sidebar.jsx` with expandable type groups and object list
- [x] 5.3 Implement `Body.jsx` with view/edit modes and save (⌘+Enter)
- [x] 5.4 Implement `Properties.jsx` with pinned/unpinned sections in card layout
- [x] 5.5 Implement `PropertyRow.jsx` with inline editing (text, number, date, checkbox)

## 6. Object Creation

- [x] 6.1 Implement `CreateDialog.jsx` with type selector, name input, and template picker

## 7. Build Infrastructure

- [x] 7.1 Add `build-frontend` target to Makefile; make `test` depend on it
- [x] 7.2 Update `.github/workflows/ci.yml` to setup Node.js and build frontend before Go build
- [x] 7.3 Update `.goreleaser.yml` with `before.hooks` to build frontend

## 8. Documentation

- [x] 8.1 Update `CLAUDE.md` architecture, layer diagram, Web UI section, and build instructions
- [x] 8.2 Update `README.md` and `README.zh-TW.md` with Web UI feature, `tmd serve` usage, tech stack
- [x] 8.3 Update `CONTRIBUTING.md` with Node.js prerequisite and frontend build workflow
- [x] 8.4 Create `websites/docs/src/content/docs/cli/serve.md` command reference
- [x] 8.5 Update `developers/architecture.md` diagram and multi-platform table
- [x] 8.6 Update `getting-started/introduction.md` and `quick-start.md` with Web UI
