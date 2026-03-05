---
title: Object 與 Type
description: 瞭解 TypeMD 中的 Object 和 Type。
sidebar:
  order: 1
---

## Object

Object（物件）是 TypeMD 的基本單位。每個 Object 以 Markdown 檔案儲存，包含 YAML frontmatter（屬性）和內文。

### Object ID

Object ID 的格式為 `type/filename`，例如 `book/golang-in-action`。`objects/` 底下的每個目錄是一個 **Type 命名空間**——不同 Type 可以共用相同的檔名。

### 資料結構

```
vault/
├── .typemd/
│   ├── types/              # Type schema 定義（YAML）
│   │   ├── book.yaml
│   │   └── person.yaml
│   └── index.db            # SQLite 索引（自動更新）
└── objects/
    ├── book/
    │   └── golang-in-action.md
    └── person/
        └── alan-donovan.md
```

## Type

每個 Object 都屬於一個 Type（類型）。Type 透過儲存在 `.typemd/types/` 的 schema 檔案定義屬性名稱、資料型別和驗證規則。

### 內建 Type

| Type | 屬性 |
|------|------|
| `book` | title (string)、status (enum: to-read/reading/done)、rating (number) |
| `person` | name (string)、role (string) |
| `note` | title (string)、tags (string) |

### 屬性型別

| 型別 | 說明 | 範例 |
|------|------|------|
| `string` | 文字 | `"Go in Action"` |
| `number` | 整數或浮點數 | `42`、`3.14` |
| `enum` | 列舉值，需定義 `values` | `"reading"` |
| `relation` | 連結到另一個 Object | `"person/alan"` |
