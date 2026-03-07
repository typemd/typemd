---
title: tmd create
description: 從 Type schema 建立新 Object。
sidebar:
  order: 2.5
---

根據 Type schema 建立新的 Object 檔案（Markdown + YAML frontmatter）。

```bash
tmd create book clean-code
tmd create person robert-martin
```

指令會在 `objects/<type>/` 下產生 `.md` 檔案，所有 schema 定義的屬性設為預設值（若無預設值則為 `null`）。Object 同時會被加入 SQLite 索引。

如果 Type 不存在或同名的 Object 已存在，會回傳錯誤。
