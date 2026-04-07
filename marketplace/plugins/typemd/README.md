# typemd

Reference guides and workflow skills for working with [typemd](https://github.com/typemd/typemd) vaults.

## vault-guide

Comprehensive typemd reference that teaches AI how to design and manage vaults — CLI commands, type schemas, object format, wiki-links, views, templates, and TUI keybindings.

Activates automatically when working in a typemd vault (`.typemd/` directory present).

## instructions-guide

Guide for using `tmd instructions` to output embedded skills (explore, importer, onboarding) enriched with vault context. Covers integration patterns for feeding vault-aware context to AI tools.

## onboarding

Guided four-phase workflow for importing existing markdown collections into a typemd vault: scan sources, generate a conversion plan, execute the plan, and verify results. Uses `tmd import scan/plan/execute` CLI commands for structured data and AI-driven classification.

## Installation

```
/plugin marketplace add typemd/marketplace
/plugin install typemd@typemd-marketplace
```

## Prerequisites

- A typemd vault (or intent to create one)
- [Claude Code](https://claude.com/code) installed
