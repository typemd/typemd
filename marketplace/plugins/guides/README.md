# guides

Reference guides for working with [typemd](https://github.com/typemd/typemd) vaults. Contains two skills:

## vault-guide

Comprehensive typemd reference that teaches AI how to design and manage vaults — CLI commands, type schemas, object format, wiki-links, views, templates, and TUI keybindings.

Activates automatically when working in a typemd vault (`.typemd/` directory present).

## instructions-guide

Guide for using `tmd instructions` to output embedded skills (explore, importer) enriched with vault context. Covers integration patterns for feeding vault-aware context to AI tools.

## Installation

```
/plugin marketplace add typemd/marketplace
/plugin install guides@typemd-marketplace
```

## Prerequisites

- A typemd vault (or intent to create one)
- [Claude Code](https://claude.com/code) installed
