---
title: AI Setup
description: Set up AI-powered features with Claude Code and the typemd marketplace plugin.
sidebar:
  order: 4
---

typemd integrates with AI through two channels:

- **TUI AI features** — auto-describe, auto-tag, and schema explore powered by the `claude` CLI
- **Skill instructions** — `tmd instructions` outputs embedded skills enriched with vault context, usable by any AI tool

## Prerequisites

- [typemd](/getting-started/installation) installed with a vault initialized
- [Claude Code](https://claude.com/code) installed

## 1. Install the typemd plugin

The **typemd** plugin teaches Claude Code how to work with your vault. Add the typemd marketplace and install it:

```bash
claude
# Inside Claude Code:
/plugin marketplace add typemd/marketplace
/plugin install typemd@typemd-marketplace
```

This installs two skills:

| Skill | What it does |
|-------|-------------|
| **vault-guide** | Teaches AI how to design and manage vaults — CLI commands, type schemas, object format, wiki-links, views, templates |
| **instructions-guide** | Teaches AI how to use `tmd instructions` for vault-aware context feeding |

The vault-guide skill activates automatically when Claude Code detects a `.typemd/` directory.

## 2. Enable TUI AI features

Add the following to your vault config:

```yaml
# .typemd/config.yaml
ai:
  enabled: true
```

Or use the CLI:

```bash
tmd config set ai.enabled true
```

This enables AI actions in the TUI:

| Key | Action |
|-----|--------|
| `g` | Open AI action picker (on an object) — Generate description / Suggest tags |
| `Ctrl+E` | Schema explore mode (on a type) — AI analyzes objects and suggests schema improvements |

These features require the `claude` CLI binary to be installed and accessible in your PATH.

## 3. Use skill instructions

`tmd instructions` outputs embedded skill instructions enriched with your vault's type definitions. This lets any AI tool understand your vault without manual context gathering.

```bash
# List available skills
tmd instructions

# Get explore skill with vault context as JSON
tmd instructions explore

# Save a skill for use in a Claude Code skills directory
tmd instructions explore --skill > .claude/skills/explore/SKILL.md
```

The JSON output includes your vault's type summaries (names, emojis, descriptions, properties), so the AI immediately knows your data model. See [`tmd instructions`](/cli/instructions) for full details.

## What's next

With the typemd plugin installed and AI enabled, you can:

- Ask Claude Code to **analyze your markdown files** and suggest type schemas (explore skill)
- Ask Claude Code to **import existing files** into your vault (importer skill)
- Use **TUI AI features** to auto-generate descriptions and tags for objects
- Pipe `tmd instructions` output to **any AI tool** for vault-aware workflows
