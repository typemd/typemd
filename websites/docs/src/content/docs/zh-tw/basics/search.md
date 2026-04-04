---
title: 搜尋
description: 在你的 vault 中進行全文搜尋。
sidebar:
  order: 4
---

搜尋讓你透過比對檔名、屬性和內文中的文字來找到 Object。由 SQLite FTS5 驅動，即使在大型 vault 中也能快速回傳結果。如果 SQLite 索引無法使用（遺失或損毀），搜尋會自動退回到不分大小寫的子字串比對，比對對象為 object 名稱、描述和內文。

## 運作方式

開啟 vault 時，TypeMD 會從所有 Object 檔案建立全文搜尋索引。索引涵蓋：

- **檔名** — Object 的 slug 和 ULID
- **屬性** — 所有 frontmatter 值（系統屬性和 schema 定義的屬性）
- **內文** — 每個 Object 的 Markdown 內文

索引儲存在 `.typemd/index.db`，每次開啟 vault 時會自動從檔案同步，使用 CLI 或 TUI 操作時也會即時更新。

## CLI 用法

```bash
tmd search "concurrency"
tmd search "golang" --json
```

| 選項 | 說明 |
|------|------|
| `--json` | 以 JSON 格式輸出結果 |

## TUI 用法

在 TUI 中，按 `/` 進入搜尋模式。結果會隨著你輸入即時篩選。按 `Esc` 清除搜尋結果並回到完整列表。

## 相關頁面

- [查詢](/zh-tw/basics/queries) — 依屬性值進行結構化篩選
- [tmd search](/zh-tw/tui/search) — CLI 參考
- [資料模型](/zh-tw/developers/data-model#sqlite-索引) — 索引的運作方式
