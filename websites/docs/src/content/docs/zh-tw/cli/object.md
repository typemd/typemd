---
title: Object 指令
description: 建立、檢視和列出 object 的指令。
---

建立、檢視、列出和搜尋 vault 中的 object。

| 指令 | 說明 |
|------|------|
| [`tmd object create`](/zh-tw/cli/create/) | 從 type schema 建立新 object |
| [`tmd object show`](/zh-tw/cli/show/) | 顯示 object 的完整資訊 |
| [`tmd object list`](/zh-tw/cli/object-list/) | 列出 vault 中所有 object |
| `tmd object lock <id>` | 鎖定 object，防止編輯 |
| `tmd object unlock <id>` | 解鎖已鎖定的 object |
| `tmd object archive <id>` | 封存 object，從預設查詢中隱藏 |
| `tmd object unarchive <id>` | 取消封存 object，使其重新出現在查詢結果中 |
| [`tmd search`](/zh-tw/cli/search/) | 在 vault 中進行全文搜尋 |
