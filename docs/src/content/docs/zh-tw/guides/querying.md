---
title: 查詢
description: 在你的 vault 中篩選和搜尋 Object。
sidebar:
  order: 4
---

TypeMD 提供兩種尋找 Object 的方式：結構化查詢和全文搜尋。

## 結構化查詢

使用 `tmd query` 依屬性篩選 Object。條件使用 `key=value` 格式，以空格分隔（AND 邏輯）。

```bash
tmd query "type=book"
tmd query "type=book status=reading"
tmd query "type=book" --json
```

## 全文搜尋

使用 `tmd search` 搜尋檔名、屬性和內文。由 SQLite FTS5 驅動。

```bash
tmd search "concurrency"
tmd search "golang" --json
```

## TUI 搜尋

在 TUI 中，按 `/` 進入搜尋模式。結果會即時篩選。按 `Esc` 清除結果並回到完整列表。
