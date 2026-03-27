## Why

The TUI sidebar displays objects grouped by type with no sorting, filtering, or alternative presentation options. As vaults grow, users need a way to see their objects through different lenses — sorted by rating, filtered by status, grouped by category. This is the foundation for the v0.5.0 "See Your Vault" release theme.

## What Changes

- Introduce the **View** concept: a type-scoped configuration of layout + filter + sort + group_by
- Each type can have multiple saved views stored as YAML files
- Every type has an implicit default view (list layout, sorted by name) that can be customized
- **BREAKING**: Type schema storage migrates from single file (`.typemd/types/book.yaml`) to directory structure (`.typemd/types/book/schema.yaml`) with auto-migration on read
- Core layer gains `ViewConfig` struct and Vault facade CRUD methods for views
- `QueryService` gains sort support via `QueryOptions` struct
- `SQLiteObjectIndex` gains `ORDER BY` support for property-based sorting
- TUI gains a full-width View mode with table display (property columns), toggleable preview panel, and Enter/Esc navigation
- Filter operators are type-aware (string, number, date, select, multi_select, relation, checkbox each have appropriate operators)
- View and Query are distinct concepts: View = single-type presentation, Query = cross-type search (Query is out of scope, separate issue)

## Capabilities

### New Capabilities
- `view-config`: ViewConfig struct, YAML serialization, Vault facade CRUD (ListViews, LoadView, SaveView, DeleteView, DefaultView)
- `view-storage`: Type directory structure migration (single file → directory), view file discovery under `.typemd/types/<type>/views/`
- `query-sort`: QueryOptions struct with sort support, SQLiteObjectIndex ORDER BY generation
- `query-filter-operators`: Type-aware filter operators (is, contains, gt, lt, before, after, is_empty, etc.) mapped to SQL conditions
- `tui-view-mode`: Full-width View mode in TUI with table display (property columns), toggleable preview panel (`p`), view creation (`+ Add View`), view selection popup (huh v2 Select), and Enter/Esc navigation

### Modified Capabilities
- `type-schema`: Type schema loading migrates from single-file to directory structure with auto-migration; read path checks `types/<name>/schema.yaml` first, then falls back to `types/<name>.yaml` and auto-upgrades
- `tui-layout`: TUI gains a new full-width View mode that replaces the three-panel layout when active

## Impact

- **core/**: New files `view.go` (ViewConfig struct + Vault methods), modifications to `query_service.go` (QueryOptions), `sqlite_object_index.go` (ORDER BY + filter operators), `local_object_repository.go` (directory-based type schema loading), `type_schema.go` (auto-migration)
- **tui/**: New files for view mode rendering and navigation, modifications to type editor (views section), keyboard shortcut registration
- **Data migration**: Existing `.typemd/types/*.yaml` files auto-migrate to directory structure on first read; old format support will be removed in a future version
- **No breaking API changes** for CLI commands or MCP server (views are a new capability, not a modification of existing commands)
