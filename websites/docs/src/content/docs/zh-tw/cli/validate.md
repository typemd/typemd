---
title: tmd type validate
description: 驗證 Type schema、Object、Relation 和 Wiki-link。
sidebar:
  order: 10
---

驗證 vault 的 Type schema、Object、Relation、Wiki-link 和 name 唯一性。依序執行五個階段：

```bash
tmd type validate
```

## Flags

| Flag | 縮寫 | 說明 |
|------|------|------|
| `--watch` | `-w` | 持續監視檔案變更並重新驗證 |

## Watch 模式

使用 `--watch` 進入持續驗證模式。此指令會監視 `types/`、`properties/` 和 `objects/` 的檔案變更。每次偵測到變更時，會清除終端機、重新同步索引，並重新執行所有五個驗證階段。

```bash
tmd type validate --watch
```

輸出會包含每次驗證的時間戳記：

```
[14:32:05] Validating...

Validation passed.
```

按 `Ctrl+C` 退出 watch 模式。

## 驗證階段

### 階段 1：Schema 驗證

掃描所有 `types/<name>/schema.yaml` 檔案並檢查：

- `name` 欄位為必填
- 每個屬性必須有 `name` 和 `type`
- `select`/`multi_select` 屬性必須定義 `options`
- `relation` 屬性必須定義 `target`

### 階段 2：Object 驗證

依據 Type schema 驗證所有 Object 屬性：

- `select`/`multi_select` 的值必須在允許的 `options` 清單中
- `number` 屬性必須是數值
- `relation` 的目標必須符合預期的 Type

### 階段 3：Relation 驗證

檢查所有已儲存的 Relation，確保來源和目標 Object 都存在。

### 階段 4：Wiki-link 驗證

偵測壞掉的 wiki-link——Object 內文中的 `[[target]]` 參照，但目標 Object 不存在。

### 階段 5：Name 唯一性驗證

檢查所有 schema 中有 `unique: true` 的 type（例如內建的 `tag` type），確保同一 type 下沒有兩個 object 共用相同的 `name` 值。

## 輸出

成功時：

```
Validation passed.
```

失敗時，錯誤依階段分組並附上詳情：

```
Schema errors:
  book.yaml: property "status": select type requires non-empty options

Object errors:
  book/example: property "rating": expected number, got "abc"

Relation errors:
  book/example author person/missing: target object not found

Wiki-link errors:
  book/example: broken wiki-link [[person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj]]

Name uniqueness errors:
  duplicate tag name "golang": tag/golang-01abc and tag/golang-01xyz

found 5 validation error(s)
```
