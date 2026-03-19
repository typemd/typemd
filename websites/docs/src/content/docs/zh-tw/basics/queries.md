---
title: 查詢
description: 依屬性值對 Object 進行結構化篩選。
sidebar:
  order: 5
---

查詢（Query）讓你使用結構化篩選規則依屬性值篩選 Object。不同於[搜尋](/zh-tw/basics/search)會在所有內容中進行全文比對，查詢針對特定屬性以型別感知運算子進行篩選。

## 運作方式

查詢使用結構化 `FilterRule` 條件來比對 Object。每條規則指定一個屬性、一個運算子和一個值。提供多個規則時，會以 AND 邏輯組合——只有符合所有條件的 Object 才會被回傳。

`type` 是一個特殊的篩選屬性，比對 Object 的 type 名稱。其他所有屬性都比對 frontmatter 屬性值。

## 列出 Object

使用 CLI 列出所有 Object 或搜尋：

```bash
# 列出所有 Object
tmd object list
tmd object list --json

# 全文搜尋
tmd search "concurrency"
```

## 篩選運算子

View 支援型別感知的篩選運算子。每種屬性類型有一組有效的運算子：

| 屬性類型 | 運算子 |
|----------|--------|
| `string`, `url` | `is`, `is_not`, `contains`, `does_not_contain`, `starts_with`, `ends_with`, `is_empty`, `is_not_empty` |
| `number` | `eq`, `neq`, `gt`, `gte`, `lt`, `lte`, `is_empty`, `is_not_empty` |
| `date`, `datetime` | `eq`, `before`, `after`, `on_or_before`, `on_or_after`, `is_empty`, `is_not_empty` |
| `select` | `is`, `is_not`, `is_empty`, `is_not_empty` |
| `multi_select`, `relation` | `contains`, `does_not_contain`, `is_empty`, `is_not_empty` |
| `checkbox` | `is`, `is_not` |

篩選規則定義在 view 設定中：

```yaml
filter:
  - property: status
    operator: is
    value: reading
  - property: rating
    operator: gt
    value: "3"
```

多個篩選規則以 AND 邏輯組合。

## 排序

查詢支援依屬性值排序。在 TUI [view 模式](/zh-tw/tui/tui#view-模式--表格顯示)中使用時，排序規則定義在 view 設定中：

```yaml
sort:
  - property: rating
    direction: desc
```

排序方向為 `asc`（遞增）或 `desc`（遞減）。多個排序規則依序套用（第一條規則為主排序）。

## 相關頁面

- [搜尋](/zh-tw/basics/search) — 跨所有內容的全文搜尋
- [View](/zh-tw/advanced/file-structure#views) — View 設定格式
- [資料模型](/zh-tw/developers/data-model#查詢架構) — 查詢的詳細說明
