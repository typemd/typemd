---
title: Type Schema
description: 使用屬性、列舉和 Relation 定義自訂 Type schema。
sidebar:
  order: 3
---

Type schema 是 `.typemd/types/` 中的 YAML 檔案，用來定義 Object 的結構。

## 基本格式

```yaml
# .typemd/types/book.yaml
name: book
properties:
  - name: title
    type: string
  - name: status
    type: enum
    values: [to-read, reading, done]
    default: to-read
  - name: rating
    type: number
```

## 屬性型別

| 型別 | 說明 | 範例 |
|------|------|------|
| `string` | 文字 | `"Go in Action"` |
| `number` | 整數或浮點數 | `42`、`3.14` |
| `enum` | 列舉值，需定義 `values` | `"reading"` |
| `relation` | 連結到另一個 Object | `"person/alan"` |

## Relation 屬性

請參閱 [Relation](/zh-tw/guides/relations) 指南，了解如何使用 `target`、`bidirectional`、`inverse` 和 `multiple` 欄位定義 Relation 屬性。

## 驗證

TypeMD 採用寬鬆驗證：

- 只驗證 schema 中定義的屬性
- 允許額外屬性（不在 schema 中的）
- 缺少的屬性不會產生錯誤
- `enum` 的值必須在 `values` 清單中
- `number` 必須是數值
- `relation` 的目標會檢查 Type 是否正確
