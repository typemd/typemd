---
title: tmd object show
description: 顯示 Object 的完整資訊。
sidebar:
  order: 3
---

顯示 Object 的完整資訊：屬性（包含 Relation 和反向連結）和內文。

支援前綴匹配 — 可以省略 ULID 後綴，只要前綴能唯一識別 Object 即可。如果在互動式終端機中有多個 Object 匹配，會顯示選擇器讓你挑選。在非互動模式（piped input）中，指令會回傳錯誤並列出所有匹配的 ID。

```bash
tmd object show book/golang-in-action
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
  backlinks: ⟵ note/reading-list-01jqr4a2bcdef0123456789xyz

Local Properties
────────────────
  mood: reflective

Body
────
  # Notes
  A great book about Go...
```

屬性依 schema 定義的順序顯示。Relation 屬性使用 `→` 表示正向連結，`←` 表示反向連結。Wiki-link 反向連結使用 `⟵` 表示其他 Object 的引用。

如果 Object 包含不在 type schema 中定義的 frontmatter 屬性（「local properties」），會顯示在獨立的 **Local Properties** 區段。此區段僅在有 local properties 時出現。
