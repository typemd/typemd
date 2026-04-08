## Context

The MCP server (`mcp/`) currently registers two read-only tools (`search`, `get_object`) in `tools.go`. All tools follow the same pattern: `mcplib.NewTool()` for registration, a handler function returning `server.ToolHandlerFunc`, and delegation to `core.Vault` facade methods.

The `core.Vault` facade already exposes all required write operations: `NewObject`, `SaveObject`, `LinkObjects`, `UnlinkObjects`, and `ListTypes`. No core changes are needed.

## Goals / Non-Goals

**Goals:**

- Add five MCP tools: `create_object`, `update_object`, `link_objects`, `unlink_objects`, `list_types`
- Follow the existing tool registration pattern in `tools.go`
- Provide clear error messages for invalid inputs and locked objects
- Merge semantics for `update_object` properties (partial update, not full replace)

**Non-Goals:**

- WebMCP tools (#134, #136) — separate scope
- Type schema CRUD via MCP — out of scope
- Batch/bulk operations — per-object overhead is acceptable at typical MCP call rates
- Optimistic locking (#384) — extends write tools in a future issue

## Decisions

### 1. Tool parameter design

`create_object` accepts `type` (required), `name` (required), `template` (optional), `properties` (optional JSON object), and `body` (optional). The `name` parameter maps to the `filename` parameter in `Vault.NewObject()`. After creation, if `properties` or `body` are provided, the object is updated via `SaveObject`.

**Rationale:** Two-step create-then-update keeps the handler simple and reuses existing Vault methods without needing a new combined facade method.

### 2. update_object merge semantics

`update_object` accepts `id` (required), `properties` (optional JSON object), and `body` (optional). Properties are merged: only provided keys are updated, existing keys not in the input are preserved. If `body` is provided, it replaces the existing body entirely.

**Rationale:** Merge semantics match the issue requirement and are more useful for LLM agents that typically set specific properties without knowing all existing ones.

### 3. All tools in a single file

All new tools are added to `mcp/tools.go` alongside existing tools. No new files needed.

**Rationale:** The file is small and the pattern is uniform. Splitting would add indirection without benefit.

### 4. list_types returns schema details

`list_types` returns type names with their schema metadata (plural, emoji, properties) rather than just names. This gives LLM agents enough context to construct valid `create_object` calls.

**Rationale:** A name-only list would require a follow-up `get_type` call for every type, defeating the purpose.

## Risks / Trade-offs

- **[Per-object reconcile overhead]** → Each `SaveObject` triggers reconcile + index update. Acceptable at typical MCP call rates (single-digit per minute). If batch ingestion becomes a use case, deferred reconciliation can be added later.
- **[Locked object writes]** → `SaveObject` already checks `locked` and returns an error. The MCP handler surfaces this as a tool error. No additional checking needed.
- **[Concurrent writes]** → File watcher in TUI handles sync. MCP writes go through the same Vault facade, so consistency is maintained through the file system as source of truth.
