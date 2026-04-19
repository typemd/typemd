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

### 讀取工具

| 工具 | 說明 |
|------|------|
| `search` | 全文搜尋 Object；回傳 `id`、`type` 和 `filename` |
| `get_object` | 依 ID 取得完整 Object 詳情；回傳 `id`、`type`、`filename`、`properties` 和 `body` |
| `list_types` | 列出所有可用的 Type schema 與基本資訊（name、plural、emoji、properties） |
| `vault_overview` | 一次取得 vault 結構摘要：每個 type 的數量、emoji、描述與最近更新的 Objects |
| `list_objects` | 列出 Object 摘要，可選 `type` 過濾與 `limit`/`offset` 分頁 |
| `query_objects` | 透過 `filters`（FilterRule 陣列）進行結構化查詢，支援 `sort`、`limit`、`offset` |
| `list_backlinks` | 回傳指定 Object 的 wiki-link backlinks 與 typed relation backlinks |
| `vault_stats` | 取得指定 `type` 的屬性填入分佈（filled 數量與 fill rate） |

### 寫入工具

| 工具 | 說明 |
|------|------|
| `create_object` | 建立新的 Object；參數 `type`、`name`，可選 `template`、`properties`、`body` |
| `update_object` | 更新 Object 的 properties（合併）與／或 body（整體替換）；參數 `id`、`properties`、`body` |
| `link_objects` | 建立兩個 Object 之間的 relation；參數 `from_id`、`relation`、`to_id` |
| `unlink_objects` | 移除 relation；參數 `from_id`、`relation`、`to_id`，可選 `both` |
