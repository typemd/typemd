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
plural: books
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

可選的 `plural` 欄位提供 Type 在集合情境中的正確複數顯示名稱（如 TUI 群組標題中的「▼ 📚 books (3)」）。若未設定，則使用 Type 名稱作為預設值。

可選的 `emoji` 欄位為 Type 提供視覺圖示，顯示在 CLI 和 TUI 的輸出中。

## 為什麼需要 Type

Type 為你的知識庫帶來**一致性**和**可查詢性**：

- **一致性** — 每個 book Object 都有相同的屬性集合，不會出現一本書用 `status`、另一本用 `reading_state` 的情況。
- **可查詢性** — 因為屬性有型別，你可以精確查詢：「顯示所有 status 是 reading 且 rating > 4 的書」。
- **驗證** — TypeMD 會根據 schema 驗證屬性值（例如 select 值必須在允許的 options 中），及早發現錯誤。

## 內建 Type

TypeMD 有兩個內建 Type：

| Type | 屬性 | 用途 |
|------|------|------|
| 🏷️ `tag` | color (string)、icon (string) | 支援 `tags` 系統屬性；具有 `unique: true` 以強制 name 唯一性 |
| 📄 `page` | _（無）_ | 通用內容容器，用於自由形式的寫作 |

內建 type 存在於每個 vault 中，無法刪除。你可以透過建立自訂的 `.typemd/types/<name>.yaml` 檔案來覆蓋它們並新增屬性。所有其他 type 都由使用者透過 `.typemd/types/*.yaml` 檔案自行定義。TypeMD 刻意不預設附帶 `book` 或 `note` 等帶有主觀用途的 type——你可以設計符合自己領域的 type。

關於標籤的詳細說明，請參閱[標籤](/zh-tw/basics/tags)。

關於 object template 的詳細說明，請參閱[模板](/zh-tw/basics/templates)。

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

如果某個屬性出現在多個 type 中，你可以定義一次然後重複使用。詳情請參閱 [Shared Properties](/zh-tw/basics/properties#shared-properties)。

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
- 屬性名稱 `description`、`created_at`、`updated_at` 和 `tags` 為[系統屬性](/zh-tw/advanced/file-structure#系統屬性)保留，不能在 type schema 中使用。`name` 可以出現在 `properties` 中，但只允許搭配 `template` 欄位（用於 [name template](/zh-tw/basics/templates#name-template)）。
