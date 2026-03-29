---
title: tmd template create
description: 建立新的 template 檔案。
sidebar:
  order: 9.8
---

建立新的 template 檔案並在編輯器中開啟。

```bash
tmd template create book/review
tmd template create note/meeting
```

參數格式為 `type/name`。此指令會：

1. 在 `templates/<type>/<name>.md` 建立檔案（如果 type 目錄不存在會自動建立）
2. 使用你的編輯器開啟檔案（`$EDITOR`、`$VISUAL`，或 `vi`）

如果 template 已存在，指令會回傳錯誤。
