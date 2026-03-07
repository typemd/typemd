---
title: tmd mcp
description: 啟動 MCP server 以整合 AI。
sidebar:
  order: 9
---

透過 stdio 啟動 [MCP (Model Context Protocol)](https://modelcontextprotocol.io) server。AI 客戶端（例如 Claude Code）可以透過此協定查詢你的 vault。

```bash
tmd mcp
tmd mcp --vault /path/to/vault
```

## 可用工具

| 工具 | 說明 |
|------|------|
| `search` | 全文搜尋 Object；回傳 `id`、`type` 和 `filename` |
| `get_object` | 依 ID 取得完整 Object 詳情；回傳 `id`、`type`、`filename`、`properties` 和 `body` |
