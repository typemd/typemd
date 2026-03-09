## Context

typemd is a local-first knowledge management tool where Markdown files with YAML frontmatter are the source of truth. The TUI provides a three-panel interface (list, body, properties). We want to bring this same experience to the browser via try.typemd.io — a pure frontend SPA that reads vault data from GitHub repos.

The frontend is designed to be shared across three surfaces: try.typemd.io (GitHub API backend), `tmd serve` (Go HTTP backend), and Wails desktop app (Go bindings). This change builds the first surface.

## Goals / Non-Goals

**Goals:**
- Ship a working read-only vault browser at try.typemd.io
- Establish the VaultStorage interface that future backends will implement
- Mirror the TUI layout faithfully in the browser
- Zero backend — static site deployment only

**Non-Goals:**
- Write support (Phase C, separate change)
- `tmd serve` Go HTTP backend (separate change, #3)
- Wails integration (separate change, #6)
- User accounts, analytics, or server-side state
- Offline/PWA support

## Decisions

### 1. Project location: `websites/try/`

The project lives in `websites/try/` alongside the existing `websites/` directory (site, docs, blog). It is a standalone React project with its own `package.json`, independent of the Go codebase.

**Alternatives considered:**
- `web/` directory (reserved for shared frontend code in the future)
- Separate repo (harder to keep VaultStorage interface in sync)

### 2. VaultStorage interface in TypeScript

```typescript
interface VaultStorage {
  listTypes(): Promise<TypeSchema[]>
  getTypeSchema(typeName: string): Promise<TypeSchema>
  listObjects(filter?: ObjectFilter): Promise<ObjectSummary[]>
  getObject(id: string): Promise<Object>
  searchObjects(keyword: string): Promise<ObjectSummary[]>
}
```

The GitHubBackend implements this by:
1. Fetching `.typemd/types/*.yaml` via GitHub Contents API
2. Fetching `objects/**/*.md` via GitHub Trees API (single call for all files)
3. Parsing frontmatter with `gray-matter`, building an in-memory index
4. Search implemented as case-insensitive substring match over name + body

**Why Trees API**: A single `GET /repos/{owner}/{repo}/git/trees/{sha}?recursive=1` call returns the entire file tree. Combined with blob fetches, this minimizes API calls vs. iterating directories with Contents API.

**Alternatives considered:**
- GraphQL API (more efficient but more complex, harder to debug)
- Contents API only (too many round-trips for large vaults)

### 3. In-memory index with lazy loading

On connect:
1. Fetch repo tree → build file list
2. Fetch all type schemas (small files, few calls)
3. Fetch object file list but NOT full content yet
4. For sidebar: parse only filenames to get type/name grouping
5. On object select: fetch full content (frontmatter + body) on demand

This avoids fetching all object bodies upfront, which could be slow for large vaults.

**Alternatives considered:**
- Fetch everything upfront (simple but slow for 100+ objects)
- IndexedDB cache (adds complexity, defer to later)

### 4. Layout: CSS Grid three-panel

```
grid-template-columns: [left] auto [body] 1fr [props] auto;
grid-template-rows: [title] auto [content] 1fr;
```

- Left panel: resizable (drag handle or fixed width), collapsible on mobile
- Title panel: spans body + props columns
- Properties panel: toggleable (hidden by default, matches TUI behavior)
- Responsive: on narrow screens, collapse to single-column with navigation

### 5. Token handling

- PAT stored in React state by default (memory only)
- Optional "Remember on this device" checkbox → stores in localStorage
- Token sent via `Authorization: token <PAT>` header on GitHub API calls
- No token needed for public repos
- Clear token button in UI

### 6. Deployment

Static build deployed to GitHub Pages via GitHub Actions. Custom domain `try.typemd.io` via CNAME.

## Risks / Trade-offs

- **GitHub API rate limits** → Unauthenticated: 60 req/hr. With PAT: 5,000 req/hr. Lazy loading mitigates this. Show rate limit status in UI.
- **Large vaults** → Trees API returns all files but individual blob fetches could be many. Mitigate with pagination in sidebar and lazy content loading.
- **CORS** → GitHub REST API supports CORS for browser requests. No proxy needed.
- **Frontmatter parsing in browser** → `gray-matter` works in browser but adds ~15KB. Acceptable trade-off for compatibility with Go parser output.
- **Shared frontend drift** → VaultStorage interface is the contract. Future backends must implement the same interface. Risk of drift is low if interface is well-defined from the start.
