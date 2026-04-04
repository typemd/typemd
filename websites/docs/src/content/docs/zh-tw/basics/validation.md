---
title: 驗證
description: TypeMD 如何驗證 type schema、object、relation 和 link。
sidebar:
  order: 6
---

驗證會檢查你的 vault 一致性——從 type schema 到 object 屬性再到 object 之間的連結。使用以下指令執行：

```bash
tmd type validate
```

若要在開發過程中持續驗證，加上 `--watch` 來監控檔案變更並自動重新驗證：

```bash
tmd type validate --watch
```

## 驗證階段

驗證依序執行五個階段，每個階段檢查 vault 的不同面向。

### 階段 1：Schema 驗證

掃描所有 type schema（`types/<name>/schema.yaml`）和 `properties/properties.yaml`（共用屬性）並檢查：

- `name` 欄位為必填
- 每個屬性必須有 `name` 和 `type`
- `select`/`multi_select` 屬性必須定義 `options`
- `relation` 屬性必須定義 `target`
- 屬性名稱 `name`、`description`、`created_at`、`updated_at`、`tags`、`locked` 和 `archived` 為[系統屬性](/zh-tw/basics/properties#系統屬性)保留
- `properties` 中的 `name` 只允許 `template` 欄位（用於 [name template](/zh-tw/basics/templates#name-template)）

### 階段 2：Object 驗證

依據 type schema 驗證所有 object 屬性：

- `select`/`multi_select` 的值必須在允許的 `options` 清單中
- `number` 屬性必須是數值
- `date` 必須是 `YYYY-MM-DD` 格式
- `datetime` 必須包含時間部分
- `url` 必須以 `http://` 或 `https://` 開頭
- `checkbox` 必須是布林值
- `relation` 的目標必須符合預期的 type

### 階段 3：Relation 驗證

檢查所有已儲存的 relation，確保來源和目標 object 都存在。

### 階段 4：Link 驗證

偵測壞掉的連結——object 內文中的 `[[target]]` 參照，但目標 object 不存在。對於匹配到多個物件的簡寫目標（`[[type/name]]` 或 `[[name]]`），會回報有歧義的目標，並列出匹配的完整 ID 清單。

### 階段 5：Name 唯一性驗證

檢查所有 schema 中有 `unique: true` 的 type（例如內建的 `tag` type），確保同一 type 下沒有兩個 object 共用相同的 `name` 值。

## 寬鬆驗證哲學

TypeMD 刻意採用寬鬆驗證：

- **允許額外屬性** — 不在 schema 中定義的屬性會被靜默接受。這讓你可以在不更新 schema 的情況下添加臨時的中繼資料。
- **缺少的屬性不會產生錯誤** — Object 不需要擁有 type schema 中定義的每個屬性。
- **只對 schema 定義的屬性做型別檢查** — 如果屬性在 schema 中，它的值必須符合宣告的型別。如果不在 schema 中，則不會進行驗證。

這種方式讓 TypeMD 保持彈性，同時仍然能捕捉真正的錯誤——例如 `select` 的值不在 options 清單中，或 `relation` 指向不存在的 object。

## 輸出

成功時：

```
Validation passed.
```

失敗時，錯誤依階段分組：

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

## 相關頁面

- [屬性](/zh-tw/basics/properties) — 屬性型別和驗證規則
- [標籤](/zh-tw/basics/tags#標籤唯一性) — 標籤 name 唯一性
- [tmd type validate](/zh-tw/tui/validate) — CLI 參考
