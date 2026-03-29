---
title: tmd template show
description: 顯示 template 的內容。
sidebar:
  order: 9.7
---

顯示 template 的 frontmatter 屬性和正文內容。

```bash
tmd template show book/review
```

輸出範例：

```
book/review

Properties
──────────
  status: draft

Body
────
  ## Review Notes
  Write your review here.
```

參數格式為 `type/name`（例如 `book/review`）。屬性依字母順序排列。若 template 沒有屬性，顯示 `(none)`。若正文為空，顯示 `(empty)`。
