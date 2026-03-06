---
title: tmd migrate
description: 更新物件以符合當前的 type schema。
sidebar:
  order: 8.5
---

將指定 type 的所有物件遷移至符合當前 schema。新增缺少的屬性（帶預設值）、移除過時的屬性，並可選擇性地重命名屬性。

```bash
tmd migrate book
tmd migrate book --dry-run
tmd migrate book --rename old_field:new_field
```

## Flags

| Flag | 說明 |
|------|------|
| `--dry-run` | 預覽變更，不修改檔案 |
| `--rename old:new` | 重命名屬性（可重複使用） |

## 運作方式

對指定 type 的每個物件：

1. **重命名** — 若提供 `--rename`，將舊屬性的值搬移到新名稱
2. **新增** — schema 中有但物件中缺少的屬性，以預設值新增
3. **移除** — 物件中有但 schema 中沒有的屬性會被移除

已經符合 schema 的物件會被跳過。

## 範例

更新 schema 後，為所有 book 新增 `isbn` 屬性：

```bash
tmd migrate book
#   book/clean-code: added [isbn]
#   book/effective-go: added [isbn]
# Migration complete: 2 object(s) updated.
```

將 `status` 重命名為 `reading_status`（先更新 schema，再執行遷移）：

```bash
tmd migrate book --rename status:reading_status
```

先預覽變更再套用：

```bash
tmd migrate book --dry-run
```

> **注意：** `--rename` flag 要求新名稱存在於當前 schema，且舊名稱不存在。請先更新 schema 再執行遷移。
