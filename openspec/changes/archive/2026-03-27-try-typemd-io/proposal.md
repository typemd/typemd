## Why

typemd currently requires local installation to use. A hosted web app at try.typemd.io would let users connect their GitHub repo and browse their vault in the browser with zero setup — lowering the barrier to trying typemd and eventually serving as a full web product.

## What Changes

- Add a new `websites/try/` directory with a React + shadcn/ui SPA
- Implement a `VaultStorage` interface abstraction with a `GitHubBackend` that reads vault data via GitHub REST API
- Parse YAML frontmatter, type schemas, and wiki-links in the browser using an in-memory index (no SQLite)
- Replicate the TUI three-panel layout (sidebar with grouped type list, markdown body, optional properties panel, title panel) in the web UI
- Support connecting to public repos without auth, and private repos via Personal Access Token (PAT)
- PAT optionally persisted to localStorage with user consent
- Read-only in this phase; read-write deferred to a future change

## Capabilities

### New Capabilities

- `web-vault-storage`: VaultStorage interface abstraction with GitHubBackend implementation — fetches repo content via GitHub API, parses frontmatter, builds in-memory object index
- `web-ui-layout`: Three-panel web layout mirroring TUI — sidebar (grouped type list with expand/collapse, search), body (markdown rendering with clickable wiki-links), optional properties panel, title panel
- `web-github-connect`: Landing page with repo URL input, optional PAT field, localStorage opt-in, connection flow

### Modified Capabilities

(none)

## Impact

- New directory: `websites/try/` (React project, independent of Go codebase)
- New dependencies: React, shadcn/ui, react-markdown, gray-matter (or similar YAML parser)
- Deployment: static site to GitHub Pages or Vercel
- No changes to existing Go code — this is a standalone frontend project
- The `VaultStorage` interface design will inform future `tmd serve` (#3) and Wails desktop app (#6) frontend work
