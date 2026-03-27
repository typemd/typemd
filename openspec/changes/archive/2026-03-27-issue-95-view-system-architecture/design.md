## Context

The TUI sidebar displays objects in a flat, type-grouped list with no sorting, filtering, or alternative presentation. The `QueryService` supports basic `key=value` filtering but no sorting. As vaults grow, users need configurable ways to browse objects within a type — sorted by rating, filtered by status, grouped by category.

The codebase follows Clean Architecture with CQRS. `Vault` is a thin facade for simple CRUD; `ObjectService` handles commands; `QueryService` handles queries. Type schemas are stored as `.typemd/types/<name>.yaml` single files. The TUI uses a three-panel layout (sidebar | body | properties) with right panel modes (`panelObject`, `panelTypeEditor`, `panelTemplate`).

Related issues: #15 (View system epic), #12 (Query system — separate concept, out of scope), #256 (basic list view implementation — blocked by this design).

## Goals / Non-Goals

**Goals:**
- Define the View concept as a type-scoped configuration of layout + filter + sort + group_by
- Migrate type schema storage from single-file to directory structure to accommodate view files
- Extend `QueryService` with sort support
- Define type-aware filter operators for all property types
- Enable a full-width View mode in TUI with navigation stack

**Non-Goals:**
- Query system (cross-type dynamic search) — separate issue
- View layouts beyond `list` — future work (table, kanban, gallery, calendar)
- Inline views embedded in markdown body — future work
- AND/OR filter group nesting — future work (current filters are AND-only)
- Relative date filters (past_week, next_month) — future work
- CLI commands for view management — future work

## Decisions

### Decision 1: View is type-scoped, Query is cross-type

**Choice:** View binds to a single type. Query (cross-type search) is a separate concept for a future issue.

**Alternatives considered:**
- View handles both single-type and cross-type → Complex, cross-type views can only filter/sort system properties since custom properties differ per type
- Only Query, no View → "Query" implies search/retrieval; "View" better conveys presentation of a known collection

**Rationale:** Different types have different property schemas. A type-scoped View can leverage all properties for filter/sort/group. Cross-type queries are limited to system properties and serve a different use case.

### Decision 2: Type schema migrates to directory structure

**Choice:** Types migrate from `.typemd/types/book.yaml` to `.typemd/types/book/schema.yaml` with views stored under `.typemd/types/book/views/`. Auto-migration on read: if old format is detected, automatically upgrade to directory structure. Old format support will be removed in a future version.

**Alternatives considered:**
- Views stored in separate `.typemd/views/` directory → Separates views from their type, less cohesive
- Views embedded in type schema YAML → Bloats schema file, mixes data definition with presentation

**Rationale:** Directory structure keeps type and its views co-located. Auto-migration provides seamless upgrade path.

### Decision 3: ViewConfig struct lives in core, rendering lives in UI

**Choice:** `ViewConfig` struct and Vault facade CRUD methods (ListViews, LoadView, SaveView, DeleteView, DefaultView) live in core. Layout rendering, group_by logic, and View mode UX live in TUI/Web layers.

**Alternatives considered:**
- ViewConfig only in TUI → Web UI would need to duplicate YAML reading and ViewConfig struct
- Full ViewService in core → Over-engineering for v0.5.0 where Vault facade CRUD is sufficient

**Rationale:** Core provides cross-platform shared config; UI layers handle platform-specific rendering. Follows existing pattern where Vault facade does CRUD for type schemas and templates.

### Decision 4: QueryService extended with sort, not a new service

**Choice:** Extend `QueryService.Query()` to accept `QueryOptions` struct with sort rules. `SQLiteObjectIndex` generates `ORDER BY json_extract(...)` clauses.

**Alternatives considered:**
- New ViewQueryService → Unnecessary; sorting is a natural extension of querying
- Sort in UI layer after fetching all results → Inefficient for large datasets; SQLite can sort efficiently

**Rationale:** Sorting belongs in the query/index layer for efficiency. Adding sort to the existing query path is minimal and consistent.

### Decision 5: Default view is implicit, materializes on customization

**Choice:** Every type has an implicit default view (layout: list, sort by name asc). When the user customizes the default view, it materializes as `views/default.yaml`. Without this file, the system generates the default in memory.

**Alternatives considered:**
- Always write default.yaml on type creation → Creates unnecessary files
- No default view, require explicit creation → Poor UX for first-time use

**Rationale:** Zero-config experience: views work out of the box. File only appears when the user intentionally customizes.

### Decision 6: Filter operators are type-aware

**Choice:** Each property type has a defined set of valid operators:

| Property Type | Operators |
|---|---|
| string | is, is_not, contains, does_not_contain, starts_with, ends_with, is_empty, is_not_empty |
| number | eq, neq, gt, gte, lt, lte, is_empty, is_not_empty |
| date / datetime | eq, before, after, on_or_before, on_or_after, is_empty, is_not_empty |
| select | is, is_not, is_empty, is_not_empty |
| multi_select | contains, does_not_contain, is_empty, is_not_empty |
| relation | contains, does_not_contain, is_empty, is_not_empty |
| checkbox | is, is_not |

**Alternatives considered:**
- Uniform operators for all types → Confusing (what does "gt" mean for a string?)
- Minimal operators only (eq/neq/empty) → Too limiting for meaningful filtering

**Rationale:** Matches user expectations from Notion/Anytype. Type-aware validation prevents nonsensical filters.

### Decision 7: TUI View mode is full-width table with preview panel

**Choice:** Entering a View replaces the entire TUI with a full-width table displaying object names and property columns (pinned first). A toggleable preview panel (`p`) splits the layout: table on the left, object preview on the right. Selecting an object (Enter) opens the standard three-panel detail view within the View context. Esc returns to the View list, then to the sidebar.

**Alternatives considered:**
- Full-width list with names only → No added value over sidebar; users reported it felt empty
- View renders in body panel, sidebar stays → Limited horizontal space for table content
- Toggle between sidebar and View with no nesting → Can't browse objects within a View context

**Rationale:** Table display with property columns justifies the full-width mode by showing information the sidebar cannot. Preview panel provides quick inspection without leaving the list. Navigation stack (sidebar → view list → object detail) provides natural drill-down and drill-up flow.

### Decision 8: View selection popup uses huh v2 Select

**Choice:** When pressing `v` with multiple saved views, a centered popup using `charm.land/huh/v2` Select field appears for selection. Single view (or no saved views) enters directly.

**Alternatives considered:**
- Hand-coded popup with manual cursor → More code (~50 lines), no built-in filtering
- bubbles/list → Too heavy for 3-5 items (has pagination, spinner, etc.)

**Rationale:** huh v2 is in the same `charm.land` namespace as bubbletea/bubbles, zero dependency conflict. Select provides keyboard navigation and potential filtering out of the box.

## Risks / Trade-offs

- **Type directory migration is a breaking change** → Mitigated by auto-migration on read. Old format is detected and upgraded transparently. Future version removes backward compatibility.
- **Filter operator validation requires type schema lookup** → FilterRule validation must resolve property types from the schema. This couples filter validation to schema loading, but this is acceptable since views are always type-scoped.
- **Full-width View mode is a significant TUI change** → Mitigated by implementing as a new panel mode alongside existing modes. Sidebar browsing is unchanged; View mode is additive.
- **group_by in UI layer may lead to duplication** → Both TUI and Web need grouping logic. Acceptable for now; can extract to a shared package if patterns converge.
