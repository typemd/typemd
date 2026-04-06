---
title: tmd stats
description: Vault 整體摘要與指定 Type 的屬性彙總統計。
sidebar:
  order: 20
---

顯示 vault 的統計資訊。不帶參數時顯示整體摘要，搭配 `--type` 時顯示指定 Type 的屬性彙總統計。

## Vault 整體摘要

```bash
tmd stats
```

顯示每個 Type 的 emoji、複數名稱、Object 數量、最後更新日期，以及所有 Type 的總計。

### 輸出範例

```
📚 books       12   last updated 2026-03-15
👤 people       8   last updated 2026-03-14
💡 ideas        5   last updated 2026-03-10
🏷️  tags         3   last updated 2026-02-28
────────────────────────────
  Total         28
```

## 指定 Type 的屬性統計

```bash
tmd stats --type book
```

顯示指定 Type 每個屬性的彙總統計：

| 屬性類型 | 彙總方式 |
|---------|---------|
| `number` | 總和、平均、最小值、最大值 |
| `select` | 各選項出現次數 |
| `checkbox` | true/false 計數 |
| `date` | 最早、最晚 |
| `relation` | 連結數量 |

### 輸出範例

```
📚 book (12 objects)
────────────────────────────────────────

rating (number)
  filled: 12/12
  sum: 46.5  avg: 3.88  min: 2  max: 5

status (select)
  filled: 12/12
  reading         ████ 4
  done            ██████ 6
  to-read         ██ 2

favorite (checkbox)
  filled: 12/12
  true: 3  false: 9

published (date)
  filled: 12/12
  earliest: 2008-08-01  latest: 2024-11-15

author (relation)
  filled: 8/12
  links: 8
```

## JSON 輸出

```bash
tmd stats --json
tmd stats --type book --json
```

兩種模式皆支援 `--json` 以產生機器可讀的輸出。
