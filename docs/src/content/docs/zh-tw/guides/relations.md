---
title: Relation
description: 使用具型別的雙向連結來連接 Object。
sidebar:
  order: 2
---

Relation（關聯）在 Type schema 中定義為 `relation` 類型的屬性，讓你用具名的連結來連接 Object。

## 定義 Relation

```yaml
# .typemd/types/book.yaml
name: book
properties:
  - name: title
    type: string
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
```

```yaml
# .typemd/types/person.yaml
name: person
properties:
  - name: name
    type: string
  - name: books
    type: relation
    target: book
    multiple: true
    bidirectional: true
    inverse: author
```

## Relation 欄位

| 欄位 | 說明 |
|------|------|
| `target` | 目標 Object 的 Type 名稱 |
| `multiple` | 該屬性是否儲存多個值（陣列） |
| `bidirectional` | 連結時自動同步反向端 |
| `inverse` | 目標 Type schema 上的屬性名稱 |

## 使用 Relation

### 建立連結

```bash
tmd link book/golang-in-action-01jqr3k5mp... author person/alan-donovan-01jqr3k5mp...
```

當 `bidirectional: true` 時，這會自動更新書的 `author` 和人物的 `books` 屬性。

### 移除連結

```bash
tmd unlink book/golang-in-action-01jqr3k5mp... author person/alan-donovan-01jqr3k5mp... --both
```

使用 `--both` 來同時移除反向端。這只在 schema 中定義了 `bidirectional: true` 和 `inverse` 欄位時才有效。
