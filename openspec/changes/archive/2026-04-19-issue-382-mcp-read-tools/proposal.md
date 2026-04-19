## Why

The MCP server currently exposes only `search` (full-text) and `get_object` (single object by ID) for discovery. An LLM maintaining a wiki needs to navigate and understand the vault's structure efficiently — list objects by type, run structured queries, check backlinks, and get a vault overview — without resorting to dozens of `search` calls or shelling out to CLI. Expanding the MCP read surface unlocks the LLM-as-wiki-maintainer workflow described in #378 (ingest) and #379 (semantic lint).

## What Changes

- Add `vault_overview` MCP tool — one-call summary of all types with counts and recent objects; the LLM's "index.md"
- Add `list_objects` MCP tool — list objects with optional filters (type, properties), sorting, and pagination (`limit`/`offset`)
- Add `query_objects` MCP tool — structured query using the existing `FilterRule` system for rich operators (`contains`, `before`, `after`, etc.)
- Add `list_backlinks` MCP tool — given an object ID, return all wiki-link backlinks and typed incoming relations
- Add `vault_stats` MCP tool — per-type property distribution and fill-rate statistics (wraps existing `TypeStats`)
- All tools delegate to existing Vault facade methods; no changes to the core API
- No breaking changes to existing MCP tools

## Capabilities

### New Capabilities

- `mcp-read-tools`: MCP read tools for vault discovery (overview, list, query, backlinks, stats)

### Modified Capabilities

*(none)*

## Impact

- **Code**: `mcp/tools.go` (register 5 new tools, add 5 handlers), `mcp/tools_test.go` (unit tests for marshalling), BDD scenarios under `mcp/features/` (new file)
- **Docs**: marketplace skill `instructions-guide` and any MCP reference pages that enumerate available tools
- **Dependencies**: none — all tools use existing Vault facade
- **Backwards compatibility**: additive only; existing clients unaffected
