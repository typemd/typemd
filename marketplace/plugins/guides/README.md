# guide

A [typemd](https://github.com/typemd/typemd) reference guide that teaches AI how to work with typemd vaults and use `tmd instructions` for enriched skill output.

## What It Does

This skill loads a comprehensive reference of typemd into the AI's context, covering:

- **CLI commands** — all `tmd` subcommands with usage and flags
- **Object format** — markdown + YAML frontmatter structure, system properties
- **Type schemas** — property types, relations, shared properties
- **Wiki-links** — syntax and backlink tracking
- **TUI** — panels, keybindings, modes
- **Vault structure** — directory layout, SQLite index
- **Skill instructions** — how to use `tmd instructions` to output embedded skills (explore, importer) enriched with vault context for AI integrations

## Installation

```
/plugin marketplace add typemd/marketplace
/plugin install guide@typemd-marketplace
```

## When It Activates

Claude loads the skill automatically when it detects you're working in a typemd vault (`.typemd/` directory present). You can also reference it explicitly:

```
/guide:instructions-guide
```

## Prerequisites

- A typemd vault (or intent to create one)
- [Claude Code](https://claude.com/code) installed
