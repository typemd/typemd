---
title: tmd template list
description: 列出可用的 template。
sidebar:
  order: 12
---

列出 vault 中所有可用的 template，可選擇依 type 篩選。

```bash
tmd template list
tmd template list book
tmd template list --json
```

輸出範例：

```
book/review
book/summary
note/meeting
```

## Flags

| Flag | 說明 |
|------|------|
| `--json` | 以 JSON 陣列輸出，包含 `type` 和 `name` 欄位 |

未指定 type 參數時，列出所有 type 的 template。指定 type 時，只顯示該 type 的 template。若沒有 template，輸出為空。
