---
title: tmd template delete
description: 刪除 template 檔案。
sidebar:
  order: 15
---

從 vault 中刪除 template 檔案。

```bash
tmd template delete book/review
tmd template delete book/review --force
```

## Flags

| Flag | 縮寫 | 說明 |
|------|------|------|
| `--force` | `-f` | 跳過確認提示 |

在互動式終端機中，刪除前會提示確認。使用 `--force` 可跳過提示（適用於腳本）。參數格式為 `type/name`。
