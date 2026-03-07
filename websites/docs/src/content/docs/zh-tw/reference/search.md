---
title: tmd search
description: 在你的 vault 中進行全文搜尋。
sidebar:
  order: 5
---

在檔名、屬性和內文中進行全文搜尋。由 SQLite FTS5 驅動。

```bash
tmd search "concurrency"
tmd search "golang" --json
```

## 選項

| 選項 | 說明 |
|------|------|
| `--json` | 以 JSON 格式輸出結果 |
