---
title: tmd type validate
description: 驗證 Type schema、Object、Relation 和 Wiki-link。
sidebar:
  order: 10
---

驗證 vault 的 Type schema、Object、Relation 和 Wiki-link。依序執行四個階段：

```bash
tmd type validate
```

## 驗證階段

### 階段 1：Schema 驗證

掃描所有 `.typemd/types/*.yaml` 檔案並檢查：

- `name` 欄位為必填
- 每個屬性必須有 `name` 和 `type`
- `enum` 屬性必須定義 `values`
- `relation` 屬性必須定義 `target`

### 階段 2：Object 驗證

依據 Type schema 驗證所有 Object 屬性：

- `enum` 的值必須在允許的 `values` 清單中
- `number` 屬性必須是數值
- `relation` 的目標必須符合預期的 Type

### 階段 3：Relation 驗證

檢查所有已儲存的 Relation，確保來源和目標 Object 都存在。

### 階段 4：Wiki-link 驗證

偵測壞掉的 wiki-link——Object 內文中的 `[[target]]` 參照，但目標 Object 不存在。

## 輸出

成功時：

```
Validation passed.
```

失敗時，錯誤依階段分組並附上詳情：

```
Schema errors:
  book.yaml: property "status": enum must define values

Object errors:
  book/example: property "rating": expected number, got "abc"

Relation errors:
  book/example author person/missing: target object not found

Wiki-link errors:
  book/example: broken wiki-link [[person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj]]

found 4 validation error(s)
```
