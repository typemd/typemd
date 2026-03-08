---
title: Object
description: TypeMD 中知識的基本單位。
sidebar:
  order: 2
---

Object 是 TypeMD 知識庫的基本建構元素。你追蹤的一切——書籍、人物、想法、專案——都是 Object。

## 什麼是 Object？

Object 是一個 Markdown 檔案，由兩個部分組成：

- **Frontmatter** — 檔案頂端的 YAML 中繼資料（屬性，例如 title、status、rating）
- **Body** — 自由格式的 Markdown 內文（筆記、想法、參考資料）

```markdown
---
title: Go in Action
status: reading
rating: 4.5
author: person/alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz
---

# Notes

A great book about Go concurrency patterns.

See also: [[note/go-concurrency-patterns]]
```

## Object ID

每個 Object 都有一個格式為 `type/slug-ulid` 的 ID，例如：

```
book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz
```

| 部分 | 說明 |
|------|------|
| `book` | 這個 Object 所屬的 Type |
| `golang-in-action` | Slug——人類可讀的名稱 |
| `01jqr3k5mpbvn8e0f2g7h9txyz` | ULID——唯一識別碼 |

使用 `tmd object create` 時，ULID 會自動附加。如果你手動建立檔案，ULID 是可選的——TypeMD 向後相容不含 ULID 的純 slug。

## Object 不是檔案

在大多數筆記工具中，你把檔案放進資料夾來組織知識。在 TypeMD 中，你透過給 Object 一個 **Type** 並用 **Relation** 連接它們來組織知識。

檔案系統只是儲存層。一個 book Object 和一個 person Object 可以互相引用，不管它們在目錄樹的哪個位置。結構來自你的資料模型，而非資料夾階層。

## Object 的存放位置

Object 儲存在 Vault 的 `objects/` 目錄下，依 Type 分類：

```
objects/
├── book/
│   ├── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
│   └── designing-data-intensive-applications-01jqr4m7npcwo9f3g8h0jtuyab.md
├── person/
│   └── alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz.md
└── note/
    └── go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc.md
```

每個 Type 有自己的子目錄。這讓檔案系統保持整潔，但真正的組織是透過 Type 和 Relation 來實現的。
