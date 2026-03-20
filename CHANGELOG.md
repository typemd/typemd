# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

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
