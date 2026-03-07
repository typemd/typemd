---
title: 介紹
description: TypeMD 是什麼，以及它為何存在。
sidebar:
  order: 1
---

TypeMD 是一個本地優先的 CLI 知識管理工具。你的知識庫由 **Object（物件）** 組成——而不是檔案。Markdown 只是儲存格式。

## 理念

大多數筆記工具讓你像電腦一樣思考：檔案、資料夾、階層結構。

TypeMD 讓你用 **Object** 來思考——書籍、人物、想法、會議——透過 **Relation（關聯）** 連結。結構源自你的知識，而非資料夾樹狀結構。

## 核心概念

### Object

Object 是 TypeMD 的基本單位。每個 Object 是一個帶有 YAML frontmatter（屬性）和內文的 Markdown 檔案。

Object ID 的格式為 `type/<slug>-<ulid>`，例如 `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`。透過 CLI 建立時，ULID 會自動附加在檔名後方以確保唯一性。

```markdown
---
title: Go in Action
status: reading
rating: 4.5
---

# Notes

A great book about Go...
```

### Type

每個 Object 都屬於一個 Type（類型）。Type 透過 schema 檔案定義屬性名稱、資料型別和驗證規則。

TypeMD 內建三種 Type：

| Type | 屬性 |
|------|------|
| `book` | title (string)、status (enum: to-read/reading/done)、rating (number) |
| `person` | name (string)、role (string) |
| `note` | title (string)、tags (string) |

自訂 Type schema 放在 `.typemd/types/` 目錄中。

### Relation

Relation 在 Type schema 中定義為 `relation` 類型的屬性。支援：

- **單向 / 雙向** — 雙向 Relation 會自動同步兩端
- **單值 / 多值** — 多值以 YAML 陣列儲存
