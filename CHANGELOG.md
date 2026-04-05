# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

## [v0.8.0] - 2026-04-05

### Breaking Changes

- Move types/ and properties/ to Vault Root — type schemas and shared property files now live at `types/` and `properties/` in the vault root instead of inside `.typemd/`. Automatic migration on vault open. Update `.gitignore` if needed. (#362)

### Added

- Object Archival — soft-delete objects with `tmd object archive`; archived objects hide from default queries but files remain; restore with `tmd object unarchive` (#34)
- Graph Export — `tmd graph` exports relations and wiki-links as DOT format for Graphviz visualization; supports `--type` filtering (#25)
- Template CLI — `tmd template list/show/create/delete` commands for managing object templates without editing files (#248)
- `tmd log` — show git commit history for a specific object (#24)
- `tmd type validate --watch` — continuous validation that re-runs on file changes (#200)
- Markdown Rendering — TUI body panel renders headings, bold, italic, code blocks, blockquotes, and lists with syntax highlighting (#103)
- Config Settings Page — press `,` in TUI to open an interactive config editor (#355)
- Per-Property Files — shared properties split from single `properties.yaml` into individual `properties/<name>.yaml` files; automatic migration (#363)
- SQLite Fallback — graceful degradation to filesystem queries when the index is unavailable (#180)

### Changed

- Reconciler/Projector Split — internal refactor separating file normalization (Reconciler) from index writing (Projector) for clearer responsibilities (#339)
- Validation Consolidation — deduplicated validation logic into a single `type_schema_validate.go` (#340)

[v0.8.0]: https://github.com/typemd/typemd/releases/tag/v0.8.0

## [v0.7.0] - 2026-03-28

### Breaking Changes

- Remove `--reindex` Flag — the global `--reindex` flag has been removed; the index now syncs automatically on every vault open, so manual reindexing is no longer needed (#325)
- Remove Legacy Single-File Type Schema — the old `.typemd/types/<name>.yaml` format is no longer supported; only directory format `.typemd/types/<name>/schema.yaml` is accepted. Run `tmd migrate` on v0.6.0 first if you still have single-file schemas (#273)
- Panel Resize Keys — remapped from `[`/`]` to `-`/`=`

### Added

- Inline Property Editing — press Enter on any property in the TUI to edit it in-place; supports all property types including string, number, date, select, checkbox, url, and text (#87)
- Relation Picker — relation properties open a fuzzy-search picker for selecting target objects; supports both single and multiple relations (#88)
- Table Cell Editing — press Enter on a cell in table view to edit the property value inline, same editing widgets as the properties panel (#316)
- Date Picker — date and datetime properties use a segmented input with an inline calendar popup for precise date selection (#324)
- Object Locking — set `locked: true` in frontmatter to prevent accidental edits; the TUI shows a lock indicator and blocks editing on locked objects (#157)
- Configurable Date Format — set `date_format` and `datetime_format` in `.typemd/config.yaml` to customize how dates are displayed throughout the TUI (#323)
- Local Properties Section — properties not defined in the type schema are displayed in a visually separated "Local Properties" section in the TUI and preserved during sync (#285)
- Local LLM Support — configure OpenAI-compatible providers (Ollama, LM Studio, vLLM) via `ai.providers` map in config; multi-provider switching with `ai.default` selector; existing `ai.enabled`/`ai.model` configs auto-migrate (#319)
- `tmd serve` — new command that starts an HTTP server with a REST API and Vue 3 frontend for browser-based vault access; includes three-theme system (warm/dark/light) (#125)
- Structured Logging — all packages use `slog` for structured JSON logging; TUI writes to `.typemd/logs/`, CLI uses `--debug` for stderr output (#326)
- Focus Mode — press `.` in TUI to toggle single full-width body panel

### Fixed

- Template Deletion — the `d` key now correctly deletes templates in the type editor (#317)
- Property Indentation — property values now align with panel header in the TUI

[v0.7.0]: https://github.com/typemd/typemd/releases/tag/v0.7.0

## [v0.6.0] - 2026-03-22

### Added

- AI MVP — auto-describe objects, suggest tags, and explore schema improvements in TUI; press `g` on an object for AI actions, `Ctrl+E` on a type for schema explore (#294)
- `tmd instructions` Command — output embedded skill instructions enriched with vault context (type summaries) as JSON; supports `--skill` for raw SKILL.md output and `--json` for list mode (#311)
- Skill Instructions Override — place `.typemd/instructions/<skill>.md` to override embedded skills per-vault (#311)
- Marketplace Guides Plugin — `vault-guide` skill teaches AI how to manage vaults; `instructions-guide` skill teaches AI how to use `tmd instructions` for context feeding
- Shell Completion — tab-complete object IDs, type names, and relation names in all CLI commands (#119)
- Interactive Disambiguation — fuzzy picker when a prefix matches multiple objects (#147)
- Shorthand Wiki-Links — `[[name]]` and `[[type/name]]` syntax auto-resolved to full IDs during sync (#176)
- Stats Dashboard — TUI stats mode showing vault-wide and per-type statistics (#280)
- Toast Notifications — transient overlay messages in TUI for sync warnings, AI errors, and other events (#303)
- `tmd format` Command — normalize object frontmatter ordering and formatting (#281)
- Sync Warnings — TUI surfaces not-found and ambiguous relation reference warnings via toast (#297)
- Object Rename — press `r` in TUI to rename an object

### Changed

- Prefix Matching in Relations — relation property values now support prefix matching during sync (#74)
- Wiki-Link Validation — `ValidateWikiLinks` resolves from files instead of DB for accuracy (#306)

### Fixed

- Brew formula name corrected to `typemd-cli` across all documentation

[v0.6.0]: https://github.com/typemd/typemd/releases/tag/v0.6.0

## [v0.5.0] - 2026-03-20

### Breaking Changes

- Remove `tmd query` — replaced by view filters (structured `FilterRule` objects) and `tmd search`; migrate scripts that use `tmd query` accordingly (#267)

### Added

- View System — define saved views per type in `.typemd/types/<name>/views/<view>.yaml` with layout, filter, sort, group_by, and columns configuration (#95, #256)
- List View Layout — `layout: list` displays objects as a name list with optional inline values; the default view for every type (#256)
- Table View Layout — `layout: table` displays objects in a columnar table with configurable `columns` for property selection (#97, #282)
- View Editor — press `e` in view mode to open a side-panel editor for filter rules, sort rules, and group_by with property/operator pickers and auto-save (#258, #279)
- View Filters — type-aware filter operators (`is`, `contains`, `before`, `gt`, etc.) applied at the query pipeline level; each property type has its own valid operator set (#263, #264)
- Multi-level Group By — `group_by` supports an array of `{property: string}` rules for nested grouping; legacy single-string format auto-migrates (#279)
- View Mode Session Persistence — TUI remembers which view was open (type, view name, cursor, scroll, expanded groups) across restarts (#265)
- View Selection Filtering — popup for selecting views supports text filtering (#259)
- `tmd stats` Command — aggregate statistics for total objects, objects per type, and property usage (#30)
- Incremental Index Sync — TUI uses `Projector.SyncFiles()` for path-based incremental sync instead of full reindex (#231)

### Changed

- Type Schema Directory Format — `MigrateSchemas` now handles directory-format type schemas alongside legacy single-file format (#260)
- Unified Property Formatting — `FormatValue` provides consistent property value rendering across view mode table and object detail (#262)
- Structured Filter Model — legacy filter strings replaced with typed `FilterRule` objects supporting type-aware operator validation (#267)
- TUI Overlay System — help overlay and popups use lipgloss `Layer`/`Compositor` for correct layering (#261, #276)

### Fixed

- View editor forwards changes to view mode and supports Shift+Tab navigation
- View editor no longer shows duplicate columns
- View mode layout, CJK text alignment, and object detail navigation corrected
- Sidebar displays object name instead of slug

[v0.5.0]: https://github.com/typemd/typemd/releases/tag/v0.5.0

## [v0.4.0] - 2026-03-18

### Added

- Built-in Page Type — `page` (📄) is a new built-in type for free-form content, always available without a YAML file (#245)
- Vault Health Check — `tmd doctor` scans for orphan directories, missing types, and other vault integrity issues (#19)
- Vault Configuration — `tmd config get/set/list` manages persistent settings in `.typemd/config.yaml`; `cli.default_type` sets the default type for object creation (#241)
- Starter Type Templates — `tmd init` offers optional starter schemas (idea 💡, note 📝, book 📚) to bootstrap new vaults (#235)
- Type Schema Versioning — `version` field (semver `"major.minor"`) on type schemas for migration tracking (#45)
- Type Schema Colors — `color` field supports preset names and hex codes for visual theming (#228)
- Type & Property Descriptions — `description` field on type schemas and properties for documentation (#228)
- TUI Template Management — view, edit, create, and delete object templates in the type editor (#250)
- TUI Object Creation Wizard — inline title-panel input with live template preview and smart slug conversion (#229, #239)
- TUI Type Creation Wizard — multi-field wizard (emoji, name, plural) in the title panel with live schema preview (#230)
- Flexible Object Creation — `tmd object create` type is now optional (falls back to `cli.default_type`), names auto-convert from natural language to slugs (#236, #240)
- Frontmatter Identity-First Ordering — system properties appear in consistent order: `name`, `description`, `created_at`, `updated_at`, `tags` (#199)

[v0.4.0]: https://github.com/typemd/typemd/releases/tag/v0.4.0

## [v0.3.0] - 2026-03-14

### Breaking Changes

- Built-in Types Removed — `book`, `person`, `note` are no longer created by `tmd init`; define your own types instead (#208)
- Reserved System Properties — `description`, `created_at`, `updated_at`, `tags` are now reserved names; type schemas that define properties with these names will fail validation. Remove them before upgrading (#193, #201, #204)

### Added

- Object Templates — place Markdown files in `templates/<type>/` to pre-fill frontmatter and body content on object creation; single template auto-applies, multiple templates prompt for selection (#173)
- Name Templates — auto-generate object names from templates (e.g., `日記 {{ date:YYYY-MM-DD }}`) by defining a `template` on the `name` property in type schemas (#186)
- Plural Display Names — `plural` field on type schemas for grammatically correct collection labels in the TUI (#205)
- Unique Constraint — `unique: true` on type schemas to prevent duplicate object names within a type (#79)
- Tag Name Validation — `tmd type validate` checks for duplicate tag names across the vault (#215)
- System Properties — `description`, `created_at`, `updated_at`, `tags` are now built-in system properties managed by typemd on every object (#193, #201, #204)
- Built-in Tag Type — `tag` is a built-in type with auto-creation during sync when objects reference non-existent tags (#204)
- TUI Type Editor — full CRUD for type schemas directly in the TUI: view, edit, add/remove properties, reorder (#207)
- Domain Events — entity operations emit domain events (`ObjectCreated`, `ObjectSaved`, `PropertyChanged`, `ObjectLinked`, `TagAutoCreated`) for extensibility (#226)
- CQRS Architecture — core refactored to separate command (`ObjectService`) and query (`QueryService`) responsibilities with `ObjectRepository` and `ObjectIndex` interfaces (#224)

### Fixed

- TUI Emoji Alignment — consistent width handling for emoji with variation selectors

[v0.3.0]: https://github.com/typemd/typemd/releases/tag/v0.3.0

## [v0.2.0] - 2026-03-11

### Breaking Changes

- `name` Property — now a reserved system property; type schemas that manually define a `name` property will fail validation. Remove any `name` entries from your type schemas before upgrading (#187)

### Added

- Property Type System — define 9 property types (`string`, `text`, `number`, `bool`, `date`, `datetime`, `url`, `enum`, `relation`) in type schemas (#8)
- Shared Properties — define reusable property definitions in `.typemd/properties.yaml` and reference them via `use` in type schemas (#188)
- Emoji on Types — add an `emoji` field to type schemas for visual identification in the TUI (#145)
- Emoji on Properties — add an `emoji` field to property schemas for compact display (#144)
- TUI Title Panel — dedicated header showing the type emoji and object name when viewing an object (#169)
- TUI Pinned Properties — mark properties as `pinned: true` in schema for prominent display in the TUI detail view (#168)
- TUI Session Persistence — cursor position, selected object, and panel state are restored across TUI restarts (#82)
- `--readonly` flag — launch the TUI in read-only mode to disable all editing (#107)
- `--reindex` flag — global flag to force rebuild the SQLite index on startup, replacing the `tmd reindex` subcommand (#159)
- Prefix Matching — resolve objects by a short prefix of their ULID suffix instead of the full ID (#72)
- Homebrew Installation — install via `brew install typemd/tap/typemd-cli` (#140)

### Changed

- `name` Property — now a required system property automatically populated from the object slug; type schemas cannot define a property named `name` (#187)
- TUI Object List — type emoji shown in group headers alongside type name (#163)
- Undefined Properties — properties not declared in the type schema are silently filtered during sync (#174)

### Fixed

- Relation Display — ULID suffix stripped from relation property display values

[v0.2.0]: https://github.com/typemd/typemd/releases/tag/v0.2.0

## [v0.1.0] - 2026-03-08

### Added

- Objects & Types — define typed schemas in YAML, create objects as Markdown files with `tmd object create` (#18)
- ULID filenames — unique suffix for conflict-free object naming (#48)
- Relations — bidirectional links via `tmd relation link` / `tmd relation unlink`, single-value overwrite and multi-value append
- Wiki-links & Backlinks — `[[target]]` syntax in markdown body with automatic backlink tracking (#10)
- Querying — `tmd query` for type/property filtering, `tmd search` for full-text search, both with `--json` output
- Validation — `tmd type validate` checks schema integrity, property types, orphaned relations, and broken wiki-links (#20)
- Migration — `tmd migrate` updates existing objects when schemas evolve (#22)
- Auto-reindex — SQLite index is automatically rebuilt when empty or missing (#41)
- Orphan cleanup — stale relations detected and removed during reindex (#21)
- CLI reorganization — commands grouped by resource type: `tmd object`, `tmd type`, `tmd relation` (#141)
- TUI — three-panel layout (#47), in-place body editing (#85), edit mode with visual indicator (#84), auto-save on exit (#86), help popup (#104)
- TUI display — ULID stripped from display names (#75), reduced indentation (#57), grouped object list (#43)
- MCP Server — `tmd mcp` exposes vault to AI assistants
- `.gitignore` on init — `tmd init` creates `.typemd/.gitignore` to exclude `index.db` (#1)
- `tmd` binary — `go install` produces `tmd` binary (#61)
- Documentation site with English and zh-TW support (#50, #54)
- BDD testing framework with Godog and Gherkin feature files (#111, #112)
- GitHub Actions release workflow for multi-platform binaries (#39)
- Codebase refactoring — unified naming conventions, extracted helpers, improved error handling (#56)
- Vault structure refactoring — remove `objects/` directory layer (#117)

[v0.1.0]: https://github.com/typemd/typemd/releases/tag/v0.1.0
