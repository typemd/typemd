# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/).

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
