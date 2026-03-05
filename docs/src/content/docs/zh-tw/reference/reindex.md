---
title: tmd reindex
description: 重建 SQLite 索引和搜尋資料庫。
sidebar:
  order: 8
---

掃描 `objects/` 目錄，將所有檔案同步到資料庫，並重建全文搜尋索引。在 TypeMD 外部手動編輯檔案後使用。

> **注意：** 開啟 vault 時，TypeMD 會在索引為空或缺失時自動同步。只有在 vault 未開啟期間編輯了檔案，才需要執行 `tmd reindex`。

```bash
tmd reindex
```
