---
title: Type
description: Type 如何定義 Object 的結構。
sidebar:
  order: 3
---

Type（類型）是定義 Object 種類和屬性的藍圖。

## 什麼是 Type？

把 Type 想成一個帶有結構的分類。`book` Type 知道書有標題、閱讀狀態和評分。`person` Type 知道人有名字和角色。

Type 以 YAML schema 檔案定義，放在 `.typemd/types/`：

```yaml
# .typemd/types/book.yaml
name: book
emoji: 📚
properties:
  - name: title
    type: string
  - name: status
    type: select
    options:
      - value: to-read
      - value: reading
      - value: done
    default: to-read
  - name: rating
    type: number
```

可選的 `emoji` 欄位為 Type 提供視覺圖示，顯示在 CLI 和 TUI 的輸出中。

## 為什麼需要 Type

Type 為你的知識庫帶來**一致性**和**可查詢性**：

- **一致性** — 每個 book Object 都有相同的屬性集合，不會出現一本書用 `status`、另一本用 `reading_state` 的情況。
- **可查詢性** — 因為屬性有型別，你可以精確查詢：「顯示所有 status 是 reading 且 rating > 4 的書」。
- **驗證** — TypeMD 會根據 schema 驗證屬性值（例如 select 值必須在允許的 options 中），及早發現錯誤。

## 內建 Type

TypeMD 內建四種 Type 讓你快速開始：

| Type | 屬性 |
|------|------|
| 📚 `book` | title (string)、status (select: to-read/reading/done)、rating (number) |
| 👤 `person` | role (string) |
| 📝 `note` | title (string) |
| 🏷️ `tag` | color (string)、icon (string) |

你可以修改這些內建 Type，或建立自己的 Type 來符合你的知識領域。

## 標籤

標籤在 TypeMD 中是一級 Object。每個 Object 都有一個 `tags` [系統屬性](/zh-tw/concepts/data-model#系統屬性)，存放對 `tag` Object 的參照。標籤名稱在整個 vault 中必須唯一。

在 frontmatter 中，標籤參照可以使用完整的 Object ID 或標籤名稱：

```yaml
tags:
  - tag/go
  - tag/programming-01jqr3k5mpbvn8e0f2g7h9txyz
```

同步時，TypeMD 會將名稱參照解析為完整 ID，並自動建立不存在的標籤。

## Type vs. 標籤

標籤跨 Type 分類 Object——一本 `book` 和一篇 `note` 可以共用同一個標籤。Type 定義結構——一本 `book` 永遠有 title、status 和 rating。使用 Type 確保結構一致性；使用標籤進行跨類型分類。

## 屬性型別

Type schema 中每個屬性都有一個資料型別：

| 型別 | 說明 | 範例 |
|------|------|------|
| `string` | 文字 | `"Go in Action"` |
| `number` | 整數或浮點數 | `42`、`3.14` |
| `date` | YYYY-MM-DD 格式的日期 | `"2026-01-15"` |
| `datetime` | ISO 8601 日期時間 | `"2026-01-15T10:30:00"` |
| `url` | http/https 開頭的網址 | `"https://example.com"` |
| `checkbox` | 布林值 | `true`、`false` |
| `select` | 固定選項中的單一值 | `"reading"` |
| `multi_select` | 固定選項中的多個值 | `["go", "programming"]` |
| `relation` | 連結到另一個 Object | `"person/alan-donovan"` |

Relation 屬性的 `target`、`bidirectional`、`inverse` 和 `multiple` 欄位，請參閱 [Relation](/zh-tw/concepts/relations) 頁面。

## 驗證

TypeMD 採用寬鬆驗證：

- 只驗證 schema 中定義的屬性
- 允許額外屬性（不在 schema 中的）
- 缺少的屬性不會產生錯誤
- `select`/`multi_select` 的值必須在定義的 `options` 中
- `number` 必須是數值
- `date` 必須是 YYYY-MM-DD 格式
- `url` 必須以 http:// 或 https:// 開頭
- `relation` 的目標會檢查 Type 是否正確
- 屬性名稱 `name`、`description`、`created_at`、`updated_at` 和 `tags` 為[系統屬性](/zh-tw/concepts/data-model#系統屬性)保留，不能在 type schema 中使用
