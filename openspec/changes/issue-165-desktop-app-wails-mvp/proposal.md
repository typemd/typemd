## Why

TypeMD currently only offers CLI and TUI interfaces. A desktop app provides a more accessible, visual experience for non-terminal users. Wails v3 lets us reuse the Go core library with a web frontend in a native macOS window, avoiding the overhead of Electron. This MVP establishes the `app/` directory architecture referenced in the project roadmap (#6).

## What Changes

- Scaffold a Wails v3 project under `app/` with Go backend and web frontend
- Bind core library functions (vault initialization, object querying) to the frontend via Wails bindings
- Create a minimal web frontend page displaying an object list
- Produce a runnable macOS binary that launches a native window with the web frontend

## Capabilities

### New Capabilities

- `desktop-app-scaffold`: Wails project structure, build configuration, and Go entry point under `app/`
- `desktop-object-list`: Minimal frontend page that lists objects by type, powered by core library bindings

### Modified Capabilities

_None — this is a new interface layer; no existing specs change._

## Impact

- **New directory:** `app/` — Wails project (Go backend + web frontend)
- **New dependency:** `github.com/wailsapp/wails/v3` in `go.mod`
- **Build tooling:** Requires Wails CLI (`wails3`) and Task runner for building
- **No changes** to `core/`, `cmd/`, `tui/`, or `mcp/`
