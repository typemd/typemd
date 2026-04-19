---
title: tmd mcp
description: Start the MCP server for AI integration.
sidebar:
  order: 9
---

Starts an [MCP (Model Context Protocol)](https://modelcontextprotocol.io) server over stdio. AI clients (e.g. Claude Code) can query your vault through this protocol.

```bash
tmd mcp
tmd mcp --vault /path/to/vault
```

## Available Tools

### Read tools

| Tool | Description |
|------|-------------|
| `search` | Full-text search objects; returns `id`, `type`, and `filename` |
| `get_object` | Get full object detail by ID; returns `id`, `type`, `filename`, `properties`, and `body` |
| `list_types` | List all available type schemas with metadata (name, plural, emoji, properties) |
| `vault_overview` | One-call vault summary: per-type count, emoji, description, and recent objects |
| `list_objects` | List object summaries with optional `type` filter and `limit`/`offset` pagination |
| `query_objects` | Structured query with `filters` (FilterRule array), optional `sort`, `limit`, `offset` |
| `list_backlinks` | Return wiki-link and typed-relation backlinks for an object id |
| `vault_stats` | Per-property distribution stats for a single `type` (fill counts, fill rate) |

### Write tools

| Tool | Description |
|------|-------------|
| `create_object` | Create a new object; params: `type`, `name`, optional `template`, `properties`, `body` |
| `update_object` | Update an object's properties (merge) and/or body (replace); params: `id`, `properties`, `body` |
| `link_objects` | Create a relation between two objects; params: `from_id`, `relation`, `to_id` |
| `unlink_objects` | Remove a relation; params: `from_id`, `relation`, `to_id`, optional `both` |
