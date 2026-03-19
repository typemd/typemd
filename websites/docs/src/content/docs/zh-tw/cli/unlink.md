---
title: tmd relation unlink
description: 移除兩個 Object 之間的 Relation。
sidebar:
  order: 7
---

從來源 Object 的 frontmatter 移除 Relation。

預設只移除正向端。使用 `--both` 可同時移除目標 Object 上的反向 Relation。此功能僅適用於雙向 Relation（schema 中設定 `bidirectional: true` 且有 `inverse` 欄位）。

Object ID 支援前綴匹配 — 可以省略 ULID 後綴，只要前綴能唯一識別 Object 即可。

```bash
# 只移除正向 Relation（book 的 "author" 屬性）
tmd relation unlink book/golang-in-action author person/alan-donovan

# 同時移除雙向：book 的 "author" 和 person 的反向 "books" 屬性
tmd relation unlink book/golang-in-action author person/alan-donovan --both
```

## 選項

| 選項 | 說明 |
|------|------|
| `--both` | 同時移除目標 Object 上的反向 Relation（需要 schema 中設定 `bidirectional` 和 `inverse` 欄位） |
