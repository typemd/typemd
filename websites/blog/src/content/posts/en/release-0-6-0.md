---
title: "TypeMD 0.6.0 — AI Meets Your Vault"
description: "v0.6.0 brings AI features: auto-describe objects, suggest tags, explore schemas, and tmd instructions to give any AI tool full vault awareness."
date: 2026-03-22
tags: [release, devlog]
---

0.5.0 laid the foundation with the View system, letting you see your data your way. This time we're going bigger — **letting AI understand your vault directly**.

## AI Is Here: Describe, Tag, Schema Explore

Three new AI actions in the TUI, all powered by the `claude` CLI:

- Press `g` on an object, choose "Generate description" — AI reads your content and writes a description
- Press `g` on an object, choose "Suggest tags" — AI recommends tags based on the content
- Press `Ctrl+E` on a type — AI samples objects of that type and suggests schema improvements (missing properties, better type choices)

To enable AI features:

```bash
tmd config set ai.enabled true
```

## tmd instructions: Give Any AI Tool Vault Awareness

Beyond the TUI's built-in AI features, we went a step further — making your vault's knowledge accessible to **any** AI tool.

```bash
tmd instructions explore
```

This outputs the explore skill's full instructions, plus your vault's type summaries (names, emojis, descriptions, properties). Feed this JSON to any AI and it immediately understands your data model.

Two embedded skills:
- **explore** — analyze markdown files and suggest type schemas
- **importer** — convert existing markdown files into vault objects

Want to customize a skill? Drop a `.typemd/instructions/<skill>.md` file to override it.

## Marketplace: typemd Plugin

We also launched the first plugin on the [typemd marketplace](https://github.com/typemd/typemd/tree/main/marketplace) — **typemd**. Once installed, Claude Code automatically learns how to work with your vault:

```bash
/plugin marketplace add typemd/marketplace
/plugin install typemd@typemd-marketplace
```

For the full setup walkthrough, see the [AI Setup guide](https://docs.typemd.io/getting-started/ai-setup).

## Smarter Navigation

This release also levels up the CLI navigation experience:

- **Shell completion** — tab-complete object IDs, type names, and relation names
- **Interactive disambiguation** — prefix matches multiple objects? A fuzzy picker pops up
- **Shorthand wiki-links** — write `[[Clean Code]]` or `[[book/clean-code]]`, auto-resolved to full IDs during sync

## What Else

- **Stats dashboard** — TUI stats mode for a vault-wide overview
- **Toast notifications** — sync warnings and AI errors as floating overlay messages
- **`tmd format`** — one command to normalize all object frontmatter
- **Object rename** — press `r` in the TUI to rename an object

## What's Next

The AI foundation is in place. Next up in 0.7.0, we're making the TUI editing experience more complete — inline property editing, relation pickers, so you can do everything without leaving the TUI.

Upgrade:

```bash
brew upgrade typemd/tap/typemd-cli
```

Or grab pre-built binaries from [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.6.0).
