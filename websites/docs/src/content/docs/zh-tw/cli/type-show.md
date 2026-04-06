---
title: tmd type show
description: 顯示 Type schema 的詳細資訊。
sidebar:
  order: 10
---

顯示指定 Type schema 的詳細資訊，包含名稱、複數形式（若有定義）和所有屬性定義。

```bash
tmd type show book
```

輸出範例：

```
Type: 📚 book
Plural: books

Properties
──────────
  title (string)
  status (select) [to-read, reading, done]
  rating (number)
  author (relation) -> person (bidirectional) inverse=books
```

當 Type schema 定義了 `plural` 欄位時，會在 Type 名稱後顯示。若未設定，則省略此行。

如果指定的 Type 不存在，會回傳錯誤。
