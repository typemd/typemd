---
title: Frontmatter
description: 手動編輯 Object 檔案時的 YAML frontmatter 格式。
sidebar:
  order: 2
---

TypeMD 使用 YAML frontmatter 儲存 Object 的屬性——也就是 Markdown 檔案開頭 `---` 分隔符之間的區塊。本頁說明手動建立或編輯 Object 檔案時的格式規則。

## 什麼是 Frontmatter？

Frontmatter 是 Markdown 檔案最開頭，以三條破折號（`---`）包圍的 YAML 區塊：

```markdown
---
name: Clean Code
description: A handbook of agile software craftsmanship
created_at: 2025-03-10T14:30:00+08:00
updated_at: 2025-03-12T09:15:00+08:00
tags:
  - tag/programming-01jqr3k5mpbvn8e0f2g7h9txyz
status: reading
author: person/robert-martin-01jqr3k5mpbvn8e0f2g7h9txyz
---

你的筆記和內容寫在這裡...
```

第二個 `---` 以上的部分是 YAML 中繼資料（屬性）。以下的部分是 Markdown 內文。

## 屬性排序

Frontmatter 中的屬性必須遵循特定順序。系統屬性永遠依固定序列排在最前面，接著才是 schema 定義的屬性：

1. `name`
2. `description`
3. `created_at`
4. `updated_at`
5. `tags`
6. *（schema 定義的屬性，依 schema 順序排列）*

TypeMD 儲存檔案時會維持這個排序。如果你新增的屬性順序不對，下次透過 CLI 或 TUI 儲存時會自動重新排序。

## 屬性值格式

### 字串

純文字值。除非值包含 YAML 特殊字元，否則引號是選填的：

```yaml
name: Clean Code
description: "A book about writing better code"
```

### 日期

使用 [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) 格式的時間戳記：

```yaml
created_at: 2025-03-10T14:30:00+08:00
updated_at: 2025-03-12T09:15:00+08:00
```

### 陣列

值的列表，以 YAML 陣列格式撰寫：

```yaml
tags:
  - tag/programming-01jqr3k5mpbvn8e0f2g7h9txyz
  - tag/software-01jqr3k5mpbvn8e0f2g7h9txyz
```

### Relation 參照

Relation 透過完整的 Object ID（`type/slug-ulid`）參照其他 Object：

```yaml
author: person/robert-martin-01jqr3k5mpbvn8e0f2g7h9txyz
related_books:
  - book/refactoring-01jqr3k5mpbvn8e0f2g7h9txyz
```

### Select 值

Select 屬性使用與定義的 option 相符的純字串值：

```yaml
status: reading
```

### 數值

不加引號的數字值：

```yaml
rating: 5
pages: 464
```

## 手動建立檔案

你可以不使用 CLI，手動建立 Object 檔案。將 `.md` 檔案放在對應的 `objects/<type>/` 目錄下：

```
objects/book/my-new-book.md
```

手動建立的檔案不需要 ULID 後綴。TypeMD 有沒有 ULID 都能正常辨識。

最簡的手動建立檔案長這樣：

```markdown
---
name: My New Book
---

關於這本書的筆記...
```

## 同步時會發生什麼

TypeMD 處理你的檔案時（每次開啟 vault 時），會套用以下規則：

- **缺少 `name`** — 如果省略，TypeMD 會從檔名的 slug 產生
- **時間戳記** — `created_at` 和 `updated_at` 如果存在會被保留；如果缺少，會根據檔案中繼資料設定
- **未定義的屬性** — 不在 type schema 中定義的屬性會保留在檔案中，但在索引中被過濾掉（不會出現在搜尋或查詢結果中）
- **排序** — 下次儲存時，屬性會重新排序為標準順序

檔案永遠是 source of truth。索引從檔案重建，反過來不會。詳細的索引機制請參閱[資料模型](/zh-tw/developers/data-model)。
