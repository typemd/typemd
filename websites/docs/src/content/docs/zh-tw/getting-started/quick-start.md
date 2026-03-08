---
title: 快速開始
description: 不到一分鐘就能開始使用 TypeMD。
sidebar:
  order: 3
---

## 1. 初始化 Vault

```bash
tmd init
```

這會在目前目錄建立 `.typemd/` 目錄結構和 SQLite 資料庫。

## 2. 開啟 TUI

```bash
tmd
tmd --readonly  # 唯讀模式（停用編輯功能）
```

這會啟動三欄介面來瀏覽你的 vault。

## 3. 建立你的第一個 Object

透過 CLI 建立 Object：

```bash
tmd object create book golang-in-action
# 建立 book/golang-in-action
```

CLI 會自動在檔名後附加 ULID 以確保唯一性，產生的檔案路徑為 `objects/book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md`。

你也可以直接手動建立檔案（不需要 ULID），例如在 `objects/book/golang-in-action.md` 建立：

```markdown
---
title: Go in Action
status: reading
rating: 4.5
---

# Notes

A great book about Go...
```

TUI 會自動偵測新檔案並顯示它。如果你 clone 了一個現有的 vault，資料庫會在首次開啟時自動建立——不需要手動執行 `tmd reindex`。

## 4. 查詢與搜尋

```bash
# 依 Type 和屬性篩選
tmd query "type=book status=reading"

# 全文搜尋
tmd search "concurrency"
```
