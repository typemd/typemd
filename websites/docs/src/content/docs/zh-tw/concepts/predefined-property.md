---
title: 預定義屬性型別
description: Type schema 中可用的內建屬性型別。
sidebar:
  order: 3.5
---

Type schema 中的每個屬性都有一個 `type`，決定它接受什麼值以及如何驗證。

:::note
本頁介紹的是 **schema 定義的屬性型別**——你在 `.typemd/types/*.yaml` 中宣告的屬性。關於 TypeMD 自動管理的系統屬性（`name`、`created_at`、`updated_at`），請參閱[系統屬性](/zh-tw/concepts/data-model#系統屬性)。
:::

## 可用型別

| 型別 | 說明 | 值範例 |
|------|------|--------|
| `string` | 純文字 | `"Go in Action"` |
| `number` | 整數或浮點數 | `42`、`3.14` |
| `date` | `YYYY-MM-DD` 格式的日期 | `"2026-03-09"` |
| `datetime` | ISO 8601 日期時間 | `"2026-03-09T10:30:00+08:00"` |
| `url` | 網址（僅限 http/https） | `"https://example.com"` |
| `checkbox` | 布林值 | `true`、`false` |
| `select` | 從選項中單選 | `"reading"` |
| `multi_select` | 從選項中多選 | `["fiction", "sci-fi"]` |
| `relation` | 連結到另一個 Object | `"person/alan-donovan"` |

## 基本型別

### string

接受任何文字值。

```yaml
- name: title
  type: string
```

### number

接受整數和浮點數。

```yaml
- name: rating
  type: number
```

### date

接受 `YYYY-MM-DD` 格式的日期。YAML 自動解析的日期（例如不加引號的 `2026-03-09`）也會被接受。

```yaml
- name: published
  type: date
```

### datetime

接受 ISO 8601 日期時間值，必須包含時間部分（至少包含時和分）。

```yaml
- name: due_at
  type: datetime
```

### url

接受以 `http://` 或 `https://` 開頭的網址。不支援其他協定（ftp、ssh 等）。

```yaml
- name: website
  type: url
```

### checkbox

僅接受布林值 `true` 或 `false`。字串值如 `"true"` 會被拒絕。

```yaml
- name: favorite
  type: checkbox
```

## 選擇型別

### select

從預定義清單中單選。需要 `options` 陣列，每個項目有 `value`（必填）和可選的 `label` 顯示名稱。

```yaml
- name: status
  type: select
  options:
    - value: to-read
      label: To Read
    - value: reading
      label: Reading
    - value: done
      label: Done
```

省略 `label` 時，預設等於 `value`。

### multi_select

從預定義清單中多選。`options` 格式與 `select` 相同。屬性值為字串列表。

```yaml
- name: tags
  type: multi_select
  options:
    - value: fiction
      label: Fiction
    - value: non-fiction
      label: Non-Fiction
    - value: classic
      label: Classic
```

單一字串值（例如 `tags: fiction`）會自動視為單元素列表（`["fiction"]`）。

## Relation 型別

將 Object 連結到另一個 Object。`target`、`bidirectional`、`inverse` 和 `multiple` 欄位的詳細說明，請參閱 [Relation](/zh-tw/concepts/relations) 頁面。

```yaml
- name: author
  type: relation
  target: person
  bidirectional: true
  inverse: books
```

## 驗證

TypeMD 採用寬鬆驗證——只驗證 schema 中定義的屬性：

- 允許額外屬性（不在 schema 中的）
- 缺少的屬性不會產生錯誤
- `select` / `multi_select` 的值必須在 `options` 清單中
- `number` 必須是數值
- `date` 必須是 `YYYY-MM-DD` 格式
- `datetime` 必須包含時間部分
- `url` 必須以 `http://` 或 `https://` 開頭
- `checkbox` 必須是布林值
- `relation` 的目標會檢查 Type 是否正確
