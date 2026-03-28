---
title: tmd serve
description: Start the web UI server for browser-based vault access.
sidebar:
  order: 12
---

Starts an HTTP server that provides a browser-based UI and REST API for your vault.

## Usage

```bash
tmd serve
tmd serve -p 8080
```

Opens the web UI at `http://localhost:3000` (default port). The interface mirrors the TUI's three-panel layout: sidebar (type groups + object list), body (markdown content), and properties panel.

## Flags

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | `-p` | `3000` | Port to listen on |
| `--host` | | `127.0.0.1` | Host to bind to (use `0.0.0.0` for all interfaces) |
| `--vault` | | `.` | Path to vault directory (global flag) |

## Web UI Features

- **Three-panel layout** — sidebar, body, and properties panel
- **Object browsing** — expand type groups, click to view objects
- **Property editing** — double-click a property value to edit inline
- **Body editing** — click "Edit" to switch to a markdown editor (⌘+Enter to save, Esc to cancel)
- **Object creation** — press `N` or click "+ New" in the sidebar

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `.` | Toggle focus mode (full-width body) |
| `p` | Toggle properties panel |
| `n` | Open create object dialog |

## REST API

The server exposes a REST API under `/api/` for programmatic access.

### Types

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/types` | List all types with metadata and object counts |
| GET | `/api/types/{name}` | Get type schema with property definitions |

### Objects

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/objects` | List all objects (optionally filter with `?type=book`) |
| GET | `/api/objects/{type}/{slug}` | Get single object with properties and body |
| PUT | `/api/objects/{type}/{slug}` | Update object body (`{"body": "..."}`) |
| POST | `/api/objects` | Create object (`{"type": "book", "name": "...", "template": "..."}`) |

### Properties

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/properties/{type}/{slug}` | Get display-ready properties for an object |
| PUT | `/api/properties/{type}/{slug}/{key}` | Update a single property (`{"value": "..."}`) |

### Templates

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/templates/{type}` | List available templates for a type |

## Development

For frontend development with hot reload:

```bash
# Terminal 1: Go server
tmd serve

# Terminal 2: Vite dev server (proxies /api to localhost:3000)
cd web/frontend && npm run dev
```

The Vite dev server runs at `http://localhost:5173` with hot module replacement. API requests are automatically proxied to the Go server.
