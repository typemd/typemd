---
title: tmd query
description: 依屬性篩選 Object。
sidebar:
  order: 4
---

依屬性篩選 Object。條件使用 `key=value` 格式，以空格分隔（AND 邏輯）。

```bash
tmd query "type=book"
tmd query "type=book status=reading"
tmd query "type=book" --json
```

## 選項

| 選項 | 說明 |
|------|------|
| `--json` | 以 JSON 格式輸出結果 |
