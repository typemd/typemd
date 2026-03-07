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

| Tool | Description |
|------|-------------|
| `search` | Full-text search objects; returns `id`, `type`, and `filename` |
| `get_object` | Get full object detail by ID; returns `id`, `type`, `filename`, `properties`, and `body` |
