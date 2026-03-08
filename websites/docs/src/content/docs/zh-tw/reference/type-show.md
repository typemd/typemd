---
title: tmd type show
description: 顯示 Type schema 的詳細資訊。
sidebar:
  order: 9
---

顯示指定 Type schema 的詳細資訊，包含所有屬性定義。

```bash
tmd type show book
```

輸出範例：

```
Type: 📚 book

Properties
──────────
  title: string
  status: enum (reading, read, to-read)
  rating: number
  author: relation → person
```

如果指定的 Type 不存在，會回傳錯誤。
