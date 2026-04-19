## 1. vault_overview tool

- [x] 1.1 Add table-driven unit tests in `mcp/tools_test.go` for `vault_overview` — covers empty vault, multiple types, recent-list cap, sort order
- [x] 1.2 Implement `vaultOverviewHandler` in `mcp/tools.go` — aggregate `ListTypes` + `LoadType` + `VaultStats` + per-type recent query
- [x] 1.3 Register `vault_overview` tool in `registerTools`
- [x] 1.4 Add `overviewEntry` / `overviewRecent` JSON structs with stable field names

## 2. list_objects tool

- [x] 2.1 Add unit tests in `mcp/tools_test.go` — covers type filter, pagination, unknown type, limit clamping
- [x] 2.2 Implement `listObjectsHandler` — delegate to `vault.QueryObjects` with `TypeFilter` when `type` supplied
- [x] 2.3 Add shared `paginateSummaries` helper for `limit`/`offset`/`total` contract (used by list + query)
- [x] 2.4 Register `list_objects` tool with `type`, `limit`, `offset` parameters

## 3. query_objects tool

- [x] 3.1 Add unit tests for filter parsing, sort, invalid filter, pagination
- [x] 3.2 Add `parseFilters(any)` and `parseSort(any)` helpers that translate JSON inputs into `[]FilterRule` / `[]SortRule` with validation
- [x] 3.3 Implement `queryObjectsHandler` — delegate to `vault.QueryObjects(filters, QuerySort(...))` + pagination
- [x] 3.4 Register `query_objects` tool with `filters`, `sort`, `limit`, `offset`

## 4. list_backlinks tool

- [x] 4.1 Add unit tests — covers wiki-only backlinks, relation-only backlinks, mixed, empty, prefix resolution
- [x] 4.2 Implement `listBacklinksHandler` — combine `vault.ListBacklinks(id)` (wiki) with `vault.Queries.ListRelations(id)` filtered to `ToID == id`
- [x] 4.3 Register `list_backlinks` tool with `id` parameter

## 5. vault_stats tool

- [x] 5.1 Add unit tests — covers partial fill rate, unknown type error, zero-object type
- [x] 5.2 Implement `vaultStatsHandler` — delegate to `vault.TypeStats(type)` and shape response with `name`/`count`/`properties`
- [x] 5.3 Register `vault_stats` tool with required `type` parameter

## 6. Documentation

- [x] 6.1 Update `marketplace/plugins/typemd/skills/instructions-guide/SKILL.md` (or equivalent) to list new MCP read tools
- [x] 6.2 Update any docs site page enumerating MCP tools (`websites/docs/src/content/docs/mcp/*` if it exists)

## 7. Verification

- [x] 7.1 Run `make test` — mcp/core/tui/ai all pass; cmd BDD `Vault_not_in_a_git_repository` / `Object_with_no_commits` fail on main too (pre-existing, unrelated)
- [x] 7.2 Manually test tool registration via `go test ./mcp/...` with table-driven cases
