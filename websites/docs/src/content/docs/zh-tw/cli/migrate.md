---
title: tmd migrate
description: 遷移 type schemas 和物件。
sidebar:
  order: 9
---

遷移 type schemas 和物件。根據是否提供 type 參數，指令有兩種模式。

```bash
tmd migrate                          # 遷移 schemas (enum → select)
tmd migrate --dry-run                # 預覽 schema 遷移
tmd migrate book                     # 遷移 book 物件以符合 schema
tmd migrate book --dry-run
tmd migrate book --rename old_field:new_field
```

## Flags

| Flag | 說明 |
|------|------|
| `--dry-run` | 預覽變更，不修改檔案 |
| `--rename old:new` | 重命名屬性（可重複使用，僅適用於有 type 參數時） |

## Schema 遷移（無 type 參數）

不帶參數執行時，`tmd migrate` 會掃描所有 type schemas（`types/<name>/schema.yaml`），將舊版 `enum` 屬性類型（含 `values`）轉換為當前的 `select` 類型（含 `options`）。

```bash
tmd migrate
#   book: converted enum → select for [status]
# Schema migration complete: 1 type(s) updated.
```

若所有 schemas 已是最新版本，則不會進行任何變更。

## 物件遷移（有 type 參數）

提供 type 名稱時，`tmd migrate <type>` 會更新該 type 的所有物件以符合當前 schema。對每個物件：

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
