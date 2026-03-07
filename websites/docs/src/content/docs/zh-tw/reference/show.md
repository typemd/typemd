---
title: tmd show
description: 顯示 Object 的完整資訊。
sidebar:
  order: 3
---

顯示 Object 的完整資訊：屬性（包含 Relation）和內文。

```bash
tmd show book/golang-in-action
```

輸出範例：

```
book/golang-in-action

Properties
──────────
  title: Go in Action
  status: reading
  rating: 4.5
  author: → person/alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz

Body
────
  # Notes
  A great book about Go...
```

屬性依 schema 定義的順序顯示。Relation 屬性使用 `→` 表示正向連結，`←` 表示反向連結。
