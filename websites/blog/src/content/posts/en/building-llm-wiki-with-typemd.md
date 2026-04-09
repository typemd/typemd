---
title: "Building Your LLM Wiki with typemd"
description: "LLM Wiki is a new pattern for AI-maintained personal knowledge bases. typemd was designed for this from day one."
date: 2026-04-07
tags: [philosophy, architecture]
---

Recently, Tobi Lütke shared a concept called [LLM Wiki](https://gist.github.com/karpathy/442a6bf555914893e9891c11519de94f). The core idea: instead of having AI search through your documents from scratch every time you ask a question, have it **continuously maintain a structured knowledge base**. Every time you add a new source, the AI updates summaries, adds cross-references, and flags contradictions. Knowledge accumulates — it's not re-derived on every query.

Reading it felt like looking in a mirror. This is exactly what we've been building toward with typemd. LLM Wiki articulated what we'd been thinking more clearly than we had ourselves, and showed us angles we hadn't considered. In this post, we'll break down the concept, see where typemd stands today, and map out where we're heading next.

## The RAG Bottleneck

Most people's experience with AI and documents looks like RAG (Retrieval-Augmented Generation): upload a pile of files, the AI searches for relevant chunks at query time, and stitches together an answer. NotebookLM, ChatGPT file uploads — they mostly work this way.

The problem: **you start from zero every time.** Ask a question that requires synthesizing five documents, and the AI has to find and piece them together anew. Nothing is built up.

LLM Wiki proposes a different approach: instead of having the AI search on demand, have it **compile knowledge ahead of time** and keep it current. The cross-references are already there. The contradictions have already been flagged. The synthesis already reflects everything you've read.

## The Three-Layer Architecture

LLM Wiki defines three layers:

1. **Raw Sources** — articles, papers, notes you've collected. Immutable source of truth.
2. **Wiki** — structured pages generated and maintained by the AI. Summaries, entity pages, concept pages, comparative analyses.
3. **Schema** — conventions that tell the AI how to organize the wiki: naming rules, page formats, workflows.

Three core operations: **Ingest** (add a new source, AI reads it and updates multiple wiki pages), **Query** (ask questions, file good answers back into the wiki), **Lint** (periodic health checks for contradictions and orphan pages).

Humans curate, ask questions, and think about meaning. The AI handles the grunt work — summarizing, cross-referencing, filing, bookkeeping.

## Where typemd and LLM Wiki Overlap

typemd's core design maps naturally onto the three-layer architecture:

| LLM Wiki | typemd | Description |
|---|---|---|
| Wiki Pages | Objects | Structured Markdown with YAML frontmatter |
| Schema | Type Schemas | `types/<name>/schema.yaml` defines properties, relations, validation |
| Cross-references | Wiki-links | `[[type/name]]` syntax, auto-expanded to full IDs with backlink indexing |
| Typed Relations | Relation properties | Bidirectional, multi-valued, type-constrained relation system |
| Index | SQLite FTS5 | Full-text search + structured property queries |
| Dataview-style queries | Views | Per-type filter/sort/group rules stored as YAML |
| Graph View | `tmd graph` | DOT-format knowledge graph output |
| Git version tracking | Native | A vault is just a git repo of markdown files |

On the AI side, typemd has three core operations: **Describe** (auto-generate summaries from content), **SuggestTags** (analyze content to recommend tags), and **ExploreSchema** (analyze objects to suggest schema improvements). These run through configurable AI providers (Claude CLI or OpenAI-compatible APIs).

Obsidian users will find the syntax familiar — typemd's frontmatter format is Obsidian-compatible, and Markdown captured by Web Clipper can be read directly.

## CLI-First: Not Retrofitted, Native

Practicing the LLM Wiki pattern requires AI to **directly operate** on the knowledge base — read, search, create, update, link. That's why Obsidian released an [official CLI](https://obsidian.md/cli) in February this year, and Logseq shipped a [CLI with a built-in MCP server](https://www.npmjs.com/package/@logseq/cli). These tools are all preparing for the same thing: letting AI agents programmatically operate on knowledge bases.

The difference is that Obsidian and Logseq are GUI-first tools where CLI and MCP are bridge layers added after the fact. typemd has been a CLI tool from day one — `tmd` is the primary interface. The MCP server (`tmd mcp`) isn't a plugin; it's built in. Every operation has a corresponding CLI command, and every output supports JSON.

This doesn't mean GUI doesn't matter — typemd also has a TUI and Web UI. But every core operation can be performed in the terminal, which means AI agents operate on a typemd vault the same way humans do. No translation layer.

## Pattern vs Platform

That said, typemd and LLM Wiki have different positions.

LLM Wiki is a **pattern** — a set of working principles. The wiki is just a pile of Markdown files, the schema is a CLAUDE.md file, and structure is maintained entirely by convention between the LLM and the human. This is flexible, but it also means the LLM has to re-learn the relationships between files every session.

typemd is a **platform** — with a type system, relation engine, reconciler, and index. Structure isn't maintained by convention; it's enforced by schema constraints and automated maintenance. The LLM doesn't have to guess what properties an object has — the schema defines them. It doesn't need to manually track references — the reconciler auto-expands wiki-links and builds backlink indexes. It doesn't need to search every file — SQLite FTS5 has already indexed everything.

Another difference is **who maintains the wiki**. LLM Wiki leans toward "the LLM writes and maintains everything." typemd is designed for **humans and LLMs to co-maintain** — you can edit properties and content directly in the TUI, or let AI operate through MCP. The `locked` property lets you mark objects that shouldn't be automatically modified. The LLM handles the grunt work, but humans retain final control.

## What We're Still Missing

Inspired by LLM Wiki, we took inventory of what's needed to become a complete LLM Wiki platform and opened corresponding issues:

**MCP read/write tools** — The MCP server currently only has `search` and `get_object`, both read-only. The LLM can read but can't write. We're adding full read/write capabilities ([#377](https://github.com/typemd/typemd/issues/377), [#382](https://github.com/typemd/typemd/issues/382)), so any MCP-compatible AI client can directly maintain your vault.

**Onboarding workflow** — The biggest friction point is "how do I turn a pile of existing Markdown into a structured vault?" We're designing a complete workflow ([#381](https://github.com/typemd/typemd/issues/381)): scan sources → LLM classifies and suggests schemas → batch import → verify. Supports multiple sources and incremental additions.

**Ingest and Lint** — The two core LLM Wiki operations. Ingest ([#378](https://github.com/typemd/typemd/issues/378)) lets the LLM read new sources and automatically update related pages. Lint ([#379](https://github.com/typemd/typemd/issues/379)) lets the LLM periodically scan for orphan pages, contradictory descriptions, and missing concept pages. Currently, typemd validation only covers structural checks (schema conformance, reference existence) — semantic health checking is still blank.

**Conflict detection** — When humans and LLMs co-maintain a vault, you need optimistic locking ([#384](https://github.com/typemd/typemd/issues/384)).

**Operation log** — The LLM needs to know what happened in the last session to pick up where it left off ([#383](https://github.com/typemd/typemd/issues/383)).

## Why This Path Is Worth Taking

LLM Wiki puts it well: people abandon wikis because **maintenance costs grow faster than value.** Updating cross-references, keeping summaries in sync, flagging contradictions — nobody wants to do that grunt work. LLMs don't get bored, and they can touch a dozen files in one pass.

But LLMs need **structure** to work efficiently. With a pile of plain Markdown, the LLM has to re-learn the relationships between files every time. typemd's type schemas, relation engine, and reconciler provide that structure — they're the shared language between the LLM and your knowledge.

Most importantly, a typemd vault is just a git repo of markdown files. You can open it with any editor, and git tracks every change. You're never locked in. The LLM is a maintainer, not a gatekeeper.

If you're interested in LLM-driven knowledge management, follow our progress on [GitHub](https://github.com/typemd/typemd).
