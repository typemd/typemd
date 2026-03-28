# CLAUDE.md

## Project Overview

typemd is a local-first CLI knowledge management tool. Objects (books, people, ideas) are stored as Markdown files with YAML frontmatter, connected by Relations. SQLite provides indexing.

## Architecture

- **core/** — Core library: objects, types, relations, index
  - **core/ai/** — AI provider abstraction and service layer (Claude CLI integration)
- **cmd/** — CLI commands (Cobra)
- **tui/** — Terminal UI (Bubble Tea)
  - **tui/widget/** — Shared UI primitives (CenteredPopup, OverlayPopup, ToastModel via Layer/Compositor, scroll) used across TUI components
- **mcp/** — MCP server
- **web/** — Web UI: React + shadcn/ui (future)
- **app/** — Desktop app via Wails + shared React frontend (future)
- **websites/** — Non-Go websites (site, docs, blog)
- **marketplace/** — Claude Code marketplace plugins (typemd plugin with vault-guide, instructions-guide, explore, and importer skills)

## Core Package Architecture

The `core/` package follows **Clean Architecture** with **CQRS** (Command Query Responsibility Segregation). The design separates concerns into layers with clear dependency rules.

### Layer Diagram

```mermaid
graph TB
    subgraph Consumers
        CMD[cmd/ — CLI]
        TUI[tui/ — Terminal UI]
        MCP[mcp/ — MCP Server]
    end

    subgraph Facade
        V[Vault]
    end

    subgraph Use Cases
        OS[ObjectService — commands]
        QS[QueryService — queries]
        PJ[Projector — file→index sync]
    end

    subgraph Domain
        OBJ[Object — aggregate root]
        TS[TypeSchema — aggregate root]
        OID[ObjectID — value object]
        EVT[DomainEvent — event types]
        ED[EventDispatcher]
    end

    subgraph Infrastructure
        REPO[ObjectRepository — interface]
        IDX[ObjectIndex — interface]
        LR[LocalObjectRepository — files]
        SI[SQLiteObjectIndex — SQLite]
    end

    CMD --> V
    TUI --> V
    MCP --> V

    V --> OS
    V --> QS
    V --> PJ
    V --> ED

    OS --> REPO
    OS --> IDX
    OS --> ED
    QS --> REPO
    QS --> IDX
    PJ --> REPO
    PJ --> IDX

    OS --> OBJ
    OS --> TS
    QS --> OBJ

    REPO -.-> LR
    IDX -.-> SI

    OBJ --> OID
    OBJ --> EVT
```

### CQRS Flow

```mermaid
graph LR
    subgraph Command Side
        C[Create / Save / Link] --> OS2[ObjectService]
        OS2 --> R1[ObjectRepository — write file]
        OS2 --> I1[ObjectIndex — upsert index]
        OS2 --> E1[EventDispatcher — emit events]
    end

    subgraph Query Side
        Q[Query / Search / Resolve] --> QS2[QueryService]
        QS2 --> R2[ObjectRepository — read by ID]
        QS2 --> I2[ObjectIndex — search/filter]
    end

    subgraph Projection
        PJ2[Projector.Sync] --> R3[ObjectRepository.Walk]
        R3 --> PJ2
        PJ2 --> I3[ObjectIndex.Upsert]
    end
```

### Key Design Decisions

- **ObjectRepository** returns domain entities (`*Object`, `*TypeSchema`), not raw bytes. Path conventions and serialization are encapsulated in implementations.
- **ObjectIndex** returns `ObjectResult` (lightweight projection) for search results. Full entity retrieval goes through `ObjectRepository.Get(id)`.
- **Vault** is a thin facade / DI container. Object business logic lives in `ObjectService` (commands) and `QueryService` (queries). Type schema CRUD (`SaveType`, `DeleteType`, `CountObjectsByType`) lives directly on Vault since it delegates to `ObjectRepository` without needing a separate service layer.
- **Domain Events** follow "entity produces → use case dispatches" pattern. Entity methods return `DomainEvent`; services collect and dispatch after successful operations.
- **Files are the source of truth**. SQLite index is an acceleration layer maintained by the `Projector`.

### Key Files

| File | Role |
|------|------|
| `display.go` | DisplayProperty struct (IsLocal for non-schema non-system properties, DateFormat/DatetimeFormat for configurable display) + FormatValue/Format rendering + DefaultDateFormat/DefaultDatetimeFormat constants + Vault.BuildDisplayProperties facade (injects format config) |
| `doctor.go` | Doctor health check: RunDoctor orchestrator, DoctorReport, issue categories |
| `doctor_orphan.go` | OrphanDir scanning for objects/ and templates/ without type schemas |
| `domain_event.go` | Domain event types + EventDispatcher |
| `filter_operator.go` | Type-aware filter operator registry (validOperators) + FilterRuleToSQL translation + ValidateFilterOperator + OperatorsForType |
| `format.go` | FormatResult struct + Vault.FormatAll/FormatObjects/FormatSchemas for batch frontmatter normalization |
| `instructions.go` | Embedded skills (`core/skills/*/SKILL.md` via `//go:embed`), SKILL.md parser, ListSkills/GetSkill/GetSkillRaw + vault override via `.typemd/instructions/<skill>.md` + BuildSkillContext for type summaries |
| `list.go` | Vault.ListTypes() facade for listing all available type names |
| `local_object_repository.go` | LocalObjectRepository (struct, path conventions, object CRUD, shared properties) |
| `local_object_repository_schema.go` | Type schema CRUD and migration |
| `local_object_repository_template.go` | Template CRUD |
| `migrate.go` | MigrateObjects (add/remove/rename properties) + MigrateSchemas (enum→select conversion) |
| `name_template.go` | EvaluateNameTemplate for `{{ date:FORMAT }}` placeholder expansion in name templates + ConvertDateFormat (YYYY/MM/DD/HH/mm/ss → Go reference time) |
| `object.go` | Object entity (aggregate root) + IsLocked() + Vault facade methods + GlobObjectIDs (filesystem-only prefix search) |
| `object_id.go` | ObjectID value object |
| `name_resolve.go` | buildNameIndex + resolveByName + resolveRelationValue + isFullObjectID for name-to-ID resolution during sync |
| `object_index.go` | ObjectIndex interface + ObjectResult + SortRule + DeleteNonTagRelations + DeleteRelationsByObject |
| `object_repository.go` | ObjectRepository interface |
| `object_service.go` | ObjectService (command use cases) + ErrObjectLocked + lock guards |
| `projector.go` | Projector (file→index sync): full Sync + incremental SyncFiles (path-based) + objectPathToID helper + relation sync (all schema-defined relations) + name reference auto-expand write-back |
| `query.go` | Vault query facades (QueryObjects/SearchObjects/VaultStats/TypeStats/RebuildIndex) + ObjectResult→Object converters |
| `query_service.go` | QueryService (query use cases) |
| `relation.go` | Relation struct + Vault.LinkObjects/UnlinkObjects facades + relation property helpers (append/remove) |
| `shared_properties.go` | Vault.LoadSharedProperties + SharedPropertiesMap + ValidateSharedProperties + resolveUseEntries |
| `slugify.go` | Slugify() function for converting natural-language names to valid slugs |
| `sqlite_object_index.go` | SQLiteObjectIndex (SQLite queries) |
| `starters.go` | Embedded starter type templates (idea/note/book) + StarterTypes() + Vault.WriteStarterTypes() |
| `stats.go` | VaultStats/TypeSummary/TypeStats/PropertyStats structs + QueryService.VaultStats()/TypeStats() methods + TypeSummary.DisplayName() |
| `sync.go` | SyncResult (Expanded/Unresolved fields for name resolution reporting) + OrphanedRelation/UnresolvedRelation structs + ReasonNotFound/ReasonAmbiguous constants + Vault.SyncIndex() facade for Projector.Sync + Vault.SyncFiles() for incremental sync |
| `system_property.go` | SystemProperty registry (name/description/created_at/updated_at/tags/locked) + IsSystemProperty/IsImmutableSystemProperty helpers |
| `tag.go` | resolveTagReference helper for tag name-to-ID resolution during sync |
| `template.go` | Template entity + Vault facade methods (ListTemplates/LoadTemplate/SaveTemplate/DeleteTemplate) |
| `type_schema.go` | TypeSchema entity + helpers + Vault type CRUD (SaveType/DeleteType/CountObjectsByType) + LoadType with schema cache + InvalidateSchemaCache |
| `type_schema_marshal.go` | YAML serialization (MarshalTypeSchema) + version handling (CompareVersions) + color validation (ValidColorPresets) |
| `type_schema_validate.go` | Schema validation (ValidateSchema) + object validation (ValidateObject) + property type validators + ValidatePropertyValue (string input validation for inline editing) |
| `ulid.go` | GenerateULID + StripULID + ulidSuffixPattern for ULID generation and stripping |
| `validate.go` | Vault-wide validators: ValidateAllObjects, ValidateRelations, ValidateRelationReferences, ValidateWikiLinks, ValidateNameUniqueness, ValidateAllSchemas |
| `vault.go` | Vault facade + lifecycle (Open/Close/Init) |
| `config.go` | VaultConfig struct (DateFormat + DatetimeFormat + CLIConfig + TUIConfig + AIConfig) + YAML loading + WriteConfig + DefaultType() + Config() + GetConfigValue/SetConfigValue/ConfigKeys (key registry) |
| `logging.go` | InitLogging(level, output) — configures global slog.Default() with JSON handler; called once at startup by CLI (`--debug`), TUI (file output), or MCP (stderr) |
| `ai/provider.go` | Provider interface + CompletionRequest/CompletionResponse types |
| `ai/claude_cli.go` | ClaudeCLI provider implementation (subprocess invocation via `claude -p`) + LookupBinary |
| `ai/openai_compatible.go` | OpenAICompatible provider implementation (HTTP POST to `/v1/chat/completions`) for Ollama, LM Studio, vLLM, etc. |
| `ai/service.go` | AIService (Describe, SuggestTags, ExploreSchema) + prompt builders + response types (TagSuggestion, SchemaExploration) |
| `ai/prompts.go` | Default system prompts for describe, tag, and explore operations |
| `view.go` | ViewConfig/FilterRule/GroupRule structs + ViewLayout constants + custom UnmarshalYAML (legacy string→[]GroupRule migration) + Vault view CRUD (ListViews/LoadView/SaveView/DeleteView/DefaultView) |
| `wikilink.go` | WikiLink/StoredWikiLink structs + ParseWikiLinks + RenderWikiLinks + Vault.ListWikiLinks/ListBacklinks facades |
| `date_edit.go` | Date editing widget (dateEdit): segmented input mode (YYYY-MM-DD segments with ↑/↓ increment, digit entry, day-of-week feedback), calendar mode (mini month grid with h/j/k/l navigation, H/L month switch, t today jump), mode toggle via `c`, localToday/daysInMonth helpers |
| `prop_editor.go` | Property editor sub-model (propEditor): cursor navigation, editable property filtering, textinput/select/multi_select/checkbox/date edit widgets, Render with cursor indicator, local property separator (`── Local Properties ──`) with read-only dimmed display |
| `prop_editor_update.go` | Property editor Update: navigation (j/k), Enter to edit, Esc to cancel, type-specific validation via core.ValidatePropertyValue, applyPropertyValue save pipeline, date picker delegation |
| `relation_picker.go` | Relation picker: fuzzy-search candidate loading via QueryService, single/multi-value relation editing, link/unlink tea.Cmd wrappers, inline picker rendering |
| `view_mode_cell_edit.go` | Table view inline cell editing: cellEdit struct, activateCellEdit/cancelCellEdit/applyCellValue, updateCellEdit key handling, isCellReadOnly detection, date picker integration, viewCellToastMsg/viewCellSavedMsg messages |
| `widget/toast.go` | ToastModel (transient notifications): Show/Update/View/Overlay, auto-dismiss via tea.Tick + seq, group aggregation, level filtering (Info/Warning/Error), configurable via ToastConfig |

### TUI Architecture

The TUI uses a three-panel layout (sidebar, body, properties) with a **focus mode** (`.` key) that hides sidebar and properties for a single full-width body panel. The layout also has a **right panel mode** system:

- `panelEmpty` — no content selected
- `panelObject` — object detail view (body + properties)
- `panelTypeEditor` — type schema editor (independent sub-model `typeEditor` split across `tui/type_editor.go`, `type_editor_update.go`, `type_editor_render.go`, `type_editor_wizard.go`, `type_editor_prop_detail.go`)
- `panelTemplate` — template detail view (independent sub-model `templateEditor` in `tui/template_editor.go`)
- `panelView` — full-width view mode (independent sub-model `viewMode` in `tui/view_mode.go`, `view_mode_cell_edit.go`) with optional `viewEditor` sub-model (`tui/view_editor.go`) shown as a right split panel

The right panel automatically follows the sidebar cursor: moving to an object shows its detail, moving to a type header shows the type editor. The `typeEditor` sub-model has its own `Update()`/`View()` methods and internal mode state (view, edit, move, add wizard, delete confirmation, property detail popup). The type editor includes a Templates section listing available templates; pressing Enter on a template transitions to `panelTemplate` mode. The `templateEditor` sub-model supports viewing, inline editing (body + properties), creating, and deleting templates. The type editor also includes a Views section listing saved views; pressing Enter on a view or pressing `v` from the sidebar transitions to `panelView` mode — a full-width list that replaces the three-panel layout. In **list layout**, navigation uses ↑/↓ and Enter opens an object. In **table layout**, four-directional cell navigation (↑/↓/←/→ or h/j/k/l) with crosshair highlighting enables inline cell editing: Enter edits the focused cell, `o` opens the object detail view, Tab/Shift+Tab moves between editable cells (skipping read-only columns like `created_at`, `updated_at`, relations). Cell editing uses the `cellEdit` sub-model with type-appropriate widgets: textinput for string/number/datetime/url, date picker for date (segmented input or inline calendar via `dateEdit` sub-model), option picker for select, multi-picker for multi_select, direct toggle for checkbox, and textinput for the NAME column. Validation uses `core.ValidatePropertyValue` with errors shown as toast notifications. Edits auto-save on confirm. The help bar shows `[EDIT]` during textinput editing, `[DATE]` during date segmented input, `[CAL]` during date calendar, and `[PICK]` during picker selection. Esc returns to view list; pressing `e` opens the `viewEditor` as a right split panel (60/40 split, mutually exclusive with preview). The view editor supports inline editing of filter rules, sort rules, and group rules with property picker (text + list), operator picker (scrollable list), and auto-save on every change.

The TUI file watcher monitors both `objects/` (for object changes) and `.typemd/types/` (for schema changes). Object file changes trigger incremental index sync via `Projector.SyncFiles()`. Schema file changes trigger `InvalidateSchemaCache()` followed by a full refresh. The watcher debounce interval defaults to 200ms and is configurable via `tui.debounce_ms` in `.typemd/config.yaml`.

Type creation uses a **title panel wizard** (`createTypeState` in `tui/create_type.go`): triggered via `+ New Type`, it transforms the title panel into a multi-field form (emoji, name, plural) with Tab cycling and a live type schema preview in the right panel. After creation, the type editor opens automatically.

The TUI supports **AI-powered assistance** when `ai.enabled: true` in config and a provider is available. AI supports multiple providers via `ai.providers` map with `ai.default` selector — `cli` type invokes the `claude` CLI binary, `openai-compatible` type calls any OpenAI-compatible HTTP API (Ollama, LM Studio, vLLM, etc.) at `/v1/chat/completions`. Pressing `g` in object detail view opens an AI action picker popup (Generate description / Suggest tags). Pressing `ctrl+e` from the sidebar enters schema explore mode. AI state is tracked via `aiState` (idle, action picker, loading, preview, showing tags) with corresponding help bar messages and property panel overlays (`tui/ai.go`, `tui/ai_update.go`, `tui/ai_render.go`). AI errors are surfaced as toast notifications.

The TUI supports **inline property editing** via the `propEditor` sub-model (`tui/prop_editor.go`, `prop_editor_update.go`, `relation_picker.go`). When the properties panel is focused (Tab from body), a cursor (▸) highlights the current editable property. Navigation uses j/k; Enter activates type-appropriate editing; Esc returns to sidebar. Read-only properties (created_at, updated_at, reverse relations, backlinks, local properties) are displayed but skipped during navigation. Local properties (frontmatter keys not in the type schema or system properties) are rendered after a `── Local Properties ──` separator in dim style. Edit widgets by property type: **string/number/datetime/url** use a textinput with type-specific validation (via `core.ValidatePropertyValue`); **date** opens a dedicated date picker (`dateEdit` sub-model in `tui/date_edit.go`) with two modes — segmented input (three editable YYYY-MM-DD segments with ←→ navigation, ↑↓ increment/decrement with carry, direct digit entry, live day-of-week feedback) and inline calendar (mini month grid with h/j/k/l day navigation, H/L month switch, t today jump, today marker) — toggled with `c`; **checkbox** toggles directly between ☐/☑ on Enter/Space; **select** opens an option list picker; **multi_select** opens a multi-picker with Space toggle and Enter confirm; **relation** (single-value) opens a searchable relation picker with fuzzy filtering and a "(none)" option to clear; **relation** (multi-value, including `tags`) opens a multi-select relation picker with Space toggle, Enter confirm, and fuzzy search. Relation editing uses `Vault.LinkObjects`/`UnlinkObjects` for mutations, handling bidirectional relations automatically. Validation errors appear as toast notifications. Edits auto-save on confirm. The help bar shows `[PROPS]` during navigation, `[EDIT]` during textinput editing, `[DATE]` during date segmented input, `[CAL]` during date calendar, and `[PICK]` during picker selection (select, multi_select, relation). The properties panel border changes to edit color (orange) during active editing. Locked objects (with `locked: true` in frontmatter) display a 🔒 badge right-aligned in the title panel. Editing is blocked for locked objects — pressing Tab (props), `e` (body edit), or `r` (rename) shows a toast notification instead. Pressing `L` (uppercase) toggles the lock state of the selected object.

The TUI uses a **toast notification system** (`widget.ToastModel` in `tui/widget/toast.go`) for transient messages. Toasts appear as a floating overlay in the bottom-right corner via lipgloss Layer/Compositor. Three severity levels: Info, Warning, Error. Toasts auto-dismiss after a configurable duration (default 3s) and can be manually dismissed via a configurable key (default Esc, which is consumed by the toast). Multiple messages from a single event are aggregated via group keys (e.g., `⚠ 2 not found`, `⚠ 1 ambiguous`). Toast is initialized from `tui.toast.*` config via `newToastFromConfig()`. Current use cases: sync unresolved reference warnings, AI operation errors, property validation errors, and table cell editing errors.

## Data Model

- Objects identified by `type/<slug>-<ulid>` (e.g. `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`)
- All objects have system properties managed by typemd: `name` (preserves original input on creation; auto-populated from slug for pre-slugified names, or from name template if defined), `description` (optional, user-authored), `created_at` (set on creation, immutable), `updated_at` (updated on save, immutable), `tags` (relation to built-in `tag` type, multiple), `locked` (boolean, user-authored, prevents editing when true). These appear first in frontmatter in that order. System properties are either **user-authored** (`name`, `description`, `tags`, `locked` — can be overridden by templates) or **auto-managed** (`created_at`, `updated_at` — cannot be overridden).
- Type schemas: `.typemd/types/<name>/schema.yaml`. Cannot define properties named `description`, `created_at`, `updated_at`, or `tags` — they're reserved system properties; `name` can appear in `properties` with only a `template` field for auto-generated names. Type schemas support optional `plural` (for display in collection contexts), `unique` (to enforce name uniqueness), `version` (semver-style `"major.minor"` string for schema migration tracking, default `"0.0"`), `color` (preset name or `#RGB`/`#RRGGBB` hex for visual theming), and `description` (free-text type documentation) fields. Properties also support an optional `description` field for documenting their purpose.
- Views: `.typemd/types/<name>/views/<view>.yaml` (optional, defines layout + columns + filter + sort + group_by for presenting objects of a type). Two layouts: `list` (name + optional inline values) and `table` (columnar NAME + property columns with headers). Optional `columns` field (`[]string`) specifies which properties to display; defaults: list = none (name only), table = all properties. `group_by` is `[]GroupRule` (array of `{property: string}`) supporting multi-level grouping; legacy single-string format (`group_by: "genre"`) is auto-migrated on load. Each type has an implicit default view (list layout, sort by name asc) that materializes as `views/default.yaml` when customized. ViewConfig is in core; layout rendering is in TUI.
- Built-in types: `tag` (🏷️, plural "tags", unique, backs `tags` system property, has `color` and `icon` string properties) and `page` (📄, plural "pages", general-purpose content container). Built-in types exist without YAML files, cannot be deleted, but can be overridden by custom `.typemd/types/<name>/schema.yaml`.
- Shared properties: `.typemd/properties.yaml` (optional, defines reusable property definitions referenced via `use` in type schemas; `use` entries can override `pin`, `emoji`, and `description`)
- Relations defined as properties in type schemas
- Wiki-links: `[[type/name-ulid]]`, `[[type/name]]`, or `[[name]]` syntax in markdown body, with backlink tracking. Shorthand forms are resolved during sync and written back as full IDs.
- SQLite index: `.typemd/index.db`
- TUI session state: `.typemd/tui-state.yaml` (persisted on quit, restored on launch; stores `selected_object_id` or `selected_type_name`, expanded groups, scroll offset, panel widths, props visibility, and optionally view mode state — `view_type_name`, `view_name`, `view_cursor`, `view_scroll`, `view_expanded_groups` — when the TUI was in view mode on exit)
- Vault config: `.typemd/config.yaml` (top-level `date_format`/`datetime_format` control date display formatting with tokens `YYYY`/`MM`/`DD`/`HH`/`mm`/`ss`, defaults `YYYY-MM-DD` and `YYYY-MM-DD HH:mm:ss`; interface-layer namespacing; `cli.default_type` sets the default type for `tmd object create`; `tmd init` always creates this with `default_type: page`; `tui.toast.*` configures toast notifications — `position`, `duration_ms`, `dismiss_key`, `show_warnings`, `show_success`; `ai.enabled` enables AI features, `ai.default` selects active provider, `ai.providers` maps named provider configs with `type` (`cli`/`openai-compatible`), `base_url`, `model`, `api_key`; `ai.prompts.*` customizes system prompts, `ai.explore.sample_count`/`ai.explore.body_truncate` configures schema explore sampling)
- Embedded skills: `core/skills/*/SKILL.md` (embedded in binary via `//go:embed`; available via `tmd instructions` — explore, importer)
- Skill overrides: `.typemd/instructions/<skill>.md` (optional, vault-level overrides for embedded skills; follows same SKILL.md frontmatter format with `name` and `description` fields)
- Starter type templates: `core/starters/*.yaml` (embedded in binary via `//go:embed`; offered during `tmd init` as opt-in type schemas — idea, note, book)
- Object templates: `templates/<type>/<name>.md` (optional, Markdown files with frontmatter property overrides and body content applied during `tmd object create`; single template auto-applies, multiple templates prompt for selection or use `-t` flag)
- Object files: `objects/<type>/<name>.md`

## Web UI Architecture

- **Shared frontend**: `tmd serve`, try.typemd.io, and desktop app (Wails) share one React + shadcn/ui frontend
- **Storage Interface**: Frontend talks to a `VaultStorage` abstraction
  - `tmd serve` → Go HTTP API (read-write)
  - try.typemd.io → GitHub REST API from browser, no backend (read-only initially, read-write later)
  - Wails → Go bindings (read-write)
- **No SQLite in browser**: try.typemd.io uses in-memory index built from GitHub API responses
- **Design principle**: SQLite is optional acceleration, not a hard dependency — files are always the source of truth

## Language Convention

**English is the primary language** for all project artifacts:

- **Issues** — titles, descriptions, comments
- **Commits** — commit messages and bodies
- **Skills** — skill content in `.claude/skills/`
- **Releases** — release notes and CHANGELOG

Blog posts are the exception: written in Traditional Chinese (zh-tw) first, then synced to English via the `sync-blog` skill.

## Build & Test

```bash
make test
```

## Debugging

- **CLI**: `tmd --debug <command>` enables DEBUG-level JSON logging to stderr. Useful for inspecting sync flow, query execution, and AI provider calls.
- **TUI**: Always writes DEBUG-level JSON logs to `.typemd/logs/{YYYY-MM-DD}.log`. Check this file to diagnose TUI issues (sync, file watcher, property editing) without corrupting terminal output.

## Testing

This project uses two layers of testing:

- **BDD (Godog)** — Define behaviors, establish shared vocabulary, and describe what a feature does from the user's perspective. Gherkin `.feature` files live in `<package>/features/`. BDD scenarios focus on **what**, not implementation details.
- **Unit tests** — Verify precise logic: edge cases, output formats, exact values, error conditions. Traditional Go `testing` style.

When deciding where a test belongs: if it defines a behavior or names a concept, write a BDD scenario. If it validates an implementation detail (e.g. JSON format, lowercase ULID, flag edge cases), write a unit test.

### BDD scope by package

| Package | Testing approach |
|---------|-----------------|
| `core/` | BDD (`core/features/`) + unit tests |
| `tui/`  | BDD (`tui/features/`, planned) + unit tests |
| `web/`  | BDD (`web/features/`, future) |
| `cmd/`  | BDD (`cmd/features/`) + unit tests |
| `mcp/`  | Unit tests — BDD TBD |
