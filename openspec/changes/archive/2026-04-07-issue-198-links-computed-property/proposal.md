## Why

Objects already track outgoing wiki-links in the SQLite `wikilinks` table (for backlink queries and validation), but there is no way to access a given object's outgoing links as a first-class property. Adding `links` as a computed system property lets consumers (TUI property panel, MCP server, CLI) display and query an object's outgoing references without parsing the markdown body themselves. This completes the link/backlink pair — `backlinks` already appears as a display property; `links` is its natural complement.

## What Changes

- Register `links` as a **computed system property** (read-only, not stored in frontmatter)
- Compute the value by querying the `wikilinks` table for outgoing links from the object
- Return a deduplicated list of resolved target object IDs (`type/name-ulid`)
- Expose in display properties so TUI, MCP, and CLI can show outgoing links
- `SetProperty("links", ...)` returns an error (read-only)

## Capabilities

### New Capabilities

- `links-property`: Computed system property that exposes an object's outgoing wiki-link targets as a read-only property

### Modified Capabilities

_(none — existing wiki-link parsing and storage are unchanged)_

## Impact

- **core/**: New computed property registration, resolver implementation, display property integration
- **Existing wiki-link infra**: Read-only consumer — no changes to parsing, sync, or storage
- **TUI/MCP/CLI**: Automatically visible once registered as a display property (no per-consumer changes needed)
