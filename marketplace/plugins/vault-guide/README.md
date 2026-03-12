# vault-guide

A [typemd](https://github.com/typemd/typemd) reference guide that teaches Claude how to work with typemd vaults.

## What It Does

This skill loads a comprehensive reference of typemd into Claude's context, covering:

- **CLI commands** — all `tmd` subcommands with usage and flags
- **Object format** — markdown + YAML frontmatter structure, system properties
- **Type schemas** — property types, relations, shared properties
- **Wiki-links** — syntax and backlink tracking
- **TUI** — panels, keybindings, modes
- **Vault structure** — directory layout, SQLite index

Once installed, Claude automatically loads this knowledge when working in a typemd vault, so it understands how to create objects, modify schemas, set up relations, and use the CLI correctly.

## Installation

```
/plugin marketplace add typemd/marketplace
/plugin install vault-guide@typemd-marketplace
```

## When It Activates

Claude loads the skill automatically when it detects you're working in a typemd vault (`.typemd/` directory present). You can also reference it explicitly:

```
/vault-guide:vault-guide
```

## Prerequisites

- A typemd vault (or intent to create one)
- [Claude Code](https://claude.com/code) installed
