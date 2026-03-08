---
title: --reindex
description: 重建 SQLite 索引和搜尋資料庫。
sidebar:
  order: 8
---

`--reindex` flag 會強制將 `objects/` 目錄完整同步到資料庫，清理孤立的 relation，並重建全文搜尋索引。可以與任何指令組合使用。

> **注意：** 開啟 vault 時，TypeMD 會在索引為空或缺失時自動同步。只有在 vault 未開啟期間編輯了檔案，才需要使用 `--reindex`。

```bash
# 重建索引並啟動 TUI
tmd --reindex

# 重建索引並啟動 MCP server
tmd mcp --reindex

# 重建索引並執行查詢
tmd query "type=book" --reindex
```

## 孤立 relation 清理

當物件從磁碟中刪除時，指向或來自該物件的 relation 會變成孤立的。在 reindex 過程中，這些 dangling reference 會自動被偵測並從索引中移除。系統會顯示受影響的 relation 列表：

```
Warning: Found 2 orphaned relation(s):
  book/golang-in-action -> person/deleted-author (relation: "author")
  person/deleted-author -> book/golang-in-action (relation: "books")
Orphaned relations have been removed from the index.
```

> **注意：** 這只會清理 SQLite 索引。`.md` 檔案中的 frontmatter 不會被修改 — 建議在刪除物件前先使用 `tmd relation unlink --both` 正確移除 relation。
