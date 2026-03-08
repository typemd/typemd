## Context

TypeMD's `app/` directory is reserved for a desktop app (currently empty with `.gitkeep`). The project uses Go 1.25, which exceeds Wails v3's requirement of Go 1.23+. The core library already provides `Vault`, `QueryObjects`, `LoadType`, and other functions that can be directly bound to a web frontend.

Wails v3 renders web content in a native window using the system's WebKit (macOS), producing a lightweight binary (~15MB) with no bundled browser.

## Goals / Non-Goals

**Goals:**
- Scaffold a Wails v3 project under `app/` that integrates with the existing Go module
- Bind core library functions so the frontend can query objects
- Render a minimal object list page in a native macOS window
- Produce a runnable `tmd-app` binary

**Non-Goals:**
- Full-featured UI (editing, search, relations) — this is MVP only
- Cross-platform builds (Linux, Windows) — macOS first
- Custom theming or design system — basic functional UI
- Wails v2 support — targeting v3 only
- Changing the existing `core/`, `tui/`, or `cmd/` packages

## Decisions

### Use Wails v3 (alpha) over v2

Wails v3 offers a transparent build system, better binding generation, and is the actively developed version. The alpha status is acceptable for an MVP that may iterate.

**Alternative considered:** Wails v2 — more stable but will be superseded; building on v3 avoids a future migration.

### Vanilla TypeScript frontend (no framework)

For an MVP object list, a heavy framework is unnecessary. Vanilla TS keeps dependencies minimal and build fast. Can migrate to React/Svelte later if needed.

**Alternative considered:** React/Svelte — adds complexity and build tooling for a single-page MVP.

### Separate `app/` binary entry point

The desktop app lives under `app/main.go` with its own `main` function, separate from `cmd/tmd/main.go`. This follows the project's architecture (each interface has its own entry point) and avoids coupling the CLI with desktop-specific code.

### Bind a service struct wrapping core.Vault

Create an `AppService` struct in `app/` that wraps `core.Vault` and exposes methods for the frontend (e.g., `ListObjects()`, `GetObject(id)`). This provides a clean API boundary and allows adding desktop-specific logic without modifying core.

## Risks / Trade-offs

- **Wails v3 alpha stability**: API may change before stable release → Pin a specific v3 version in `go.mod`; MVP scope is small enough to adapt
- **macOS-only initially**: No Windows/Linux support → Acceptable for MVP; Wails v3 supports all three platforms when ready
- **Build tooling dependency**: Requires `wails3` CLI and `task` runner → Document prerequisites clearly in README
- **Module structure**: `app/` shares the same Go module as the rest of the project → Wails v3 supports this; no separate module needed
