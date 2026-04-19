## Context

The MCP server (`mcp/`) is the typemd MCP integration surface. Today it exposes `search`, `get_object`, `list_types`, `create_object`, `update_object`, `link_objects`, and `unlink_objects`. Discovery is limited to full-text search; the LLM must guess or rely on prior context to understand vault structure.

The core layer already provides every primitive needed:

- `Vault.QueryObjects(filter, opts...)` with `FilterRule` and `QuerySort`
- `Vault.VaultStats()` and `Vault.TypeStats(typeName)`
- `Vault.ListBacklinks(objectID)` (wiki-links)
- `QueryService.ListRelations(objectID)` (typed relations)
- `Vault.ListTypes()` + `Vault.LoadType(name)`

The task is to expose these in an LLM-friendly shape through MCP tools.

## Goals / Non-Goals

**Goals:**

- Add five MCP read tools (`vault_overview`, `list_objects`, `query_objects`, `list_backlinks`, `vault_stats`)
- Reuse existing core facades; no changes to the Vault API
- JSON payloads optimised for LLM consumption — concise, consistent, stable field names
- Pagination on any tool that can return more than a handful of rows (`list_objects`, `query_objects`)
- Follow the existing `mcp/tools.go` registration pattern

**Non-Goals:**

- Semantic or vector search (future work)
- Write tools (already covered by #377)
- Streaming / subscription-style tools
- CLI / Web API changes

## Decisions

### 1. JSON filter syntax for `query_objects`

`FilterRule` exposes `{property, operator, value}`. Expose the same shape directly as a JSON array — no DSL string. This avoids round-tripping and lets the LLM reuse the exact vocabulary already documented in type-schema files and query specs.

```json
{"filters": [{"property": "status", "operator": "is", "value": "reading"}]}
```

**Alternative considered**: DSL string (e.g. `"status = 'reading' AND rating > 3"`). Rejected — requires a parser and duplicates validation logic that the index already performs.

### 2. Pagination contract

`list_objects` and `query_objects` both accept `limit`/`offset`. Defaults: `limit=50`, `offset=0`. Hard ceiling: `limit=500` (clamped server-side, not rejected). Each response includes `total` so the LLM can plan further calls.

**Alternative considered**: cursor-based pagination. Rejected for now — offset is simpler, matches the existing index query path, and is sufficient for vaults with <10K objects. Can upgrade later without breaking callers if we keep the response shape extensible.

### 3. `vault_overview` shape

One array entry per type with: `name`, `plural`, `emoji`, `description`, `count`, `recent` (array of `{id, name, updated_at}` capped at 5). Count comes from `VaultStats().Counts[type]`; recent comes from `QueryObjects(TypeFilter(type), QuerySort(SortRule{Property: "updated_at", Direction: "desc"}))` clipped to 5. This is the LLM's "index.md" — one tool call returns everything needed to start navigating.

### 4. `list_backlinks` separates wiki vs relation backlinks

The core model treats wiki-links and typed relations as distinct edges. Merging them in the response would lose the distinction. Return two separate arrays: `wiki_backlinks` (bare source IDs/names) and `relation_backlinks` (source plus relation name).

### 5. Summary vs. detail responses

All list-style tools return lightweight summaries — ID, type, name, `updated_at` — not full property maps. Callers needing full objects use `get_object`. This mirrors the existing `searchHandler` shape and keeps payloads small for LLM context windows.

### 6. BDD coverage location

Per CLAUDE.md, `mcp/` historically has been unit-tested only ("BDD TBD"). For this change we add a lightweight BDD feature file under `mcp/features/mcp-read-tools.feature` to document behaviour in Gherkin, and supporting step definitions. If setup cost is too high, we fall back to table-driven unit tests in `mcp/tools_test.go` — decision is deferred to tasks.

## Risks / Trade-offs

- **[Risk]** Large vaults make `list_objects` without pagination slow → **Mitigation**: enforce hard limit of 500 and require explicit `offset` paging.
- **[Risk]** `query_objects` exposes the full `FilterRule` operator set, some of which may not be documented externally → **Mitigation**: tool description references `list_types` and existing view/query docs; invalid operators produce actionable errors.
- **[Risk]** `vault_overview` fans out per-type queries for recent objects, which on a vault with many types could be O(N types × index roundtrip) → **Mitigation**: types are usually <20; if this becomes an issue we can add an index-level shortcut later without changing the tool contract.
- **[Trade-off]** We return summaries, not full objects — clients doing "list then fetch" pay an extra round trip. Accepted because it keeps list responses small enough for LLM context windows.

## Migration Plan

Additive change. No existing tool contracts change. Deployment is "add new tools"; rollback is "remove tools" — no data migrations.

## Open Questions

- Should `query_objects` accept a `type` shortcut field alongside `filters`, or require callers to express type via `{property: "type", operator: "is", value: "<type>"}`? Lean toward the explicit `FilterRule` form for consistency with the index, but worth revisiting after first LLM usage.
- Should `list_backlinks` collapse duplicate sources (same object referenced multiple times via both wiki-link and relation)? Current proposal keeps them distinct in two arrays — callers can dedupe if needed.
