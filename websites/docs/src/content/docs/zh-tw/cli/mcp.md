---
title: tmd mcp
description: 啟動 MCP server，用於 AI 整合。
sidebar:
  order: 9
---

啟動透過 stdio 的 [MCP (Model Context Protocol)](https://modelcontextprotocol.io) server。AI 客戶端（例如 Claude Code）可以透過此協定讀寫你的 vault。

```bash
tmd mcp
tmd mcp --vault /path/to/vault
```

## 可用工具

| 工具 | 說明 |
|------|------|
| `search` | 全文搜尋 Object；回傳 `id`、`type` 和 `filename` |
| `get_object` | 依 ID 取得完整 Object 詳情；回傳 `id`、`type`、`filename`、`properties` 和 `body` |
| `list_types` | 列出所有可用的 Type Schema 及中繼資料（name、plural、emoji、properties） |
| `create_object` | 建立新 Object；參數：`type`、`name`，可選 `template`、`properties`、`body` |
| `update_object` | 更新 Object 的屬性（合併）和/或內文（取代）；參數：`id`、`properties`、`body` |
| `link_objects` | 建立兩個 Object 之間的 Relation；參數：`from_id`、`relation`、`to_id` |
| `unlink_objects` | 移除 Relation；參數：`from_id`、`relation`、`to_id`，可選 `both` |
