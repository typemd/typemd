---
title: 關聯
description: Object 之間如何透過結構化的關聯互相連接。
sidebar:
  order: 4
---

Relation（關聯）是兩個 Object 之間的結構化關聯。它用來表達一本書有作者、一個專案有成員、或一個想法關聯到一場會議。

## 什麼是 Relation？

Relation 是在 Type schema 中定義的 `relation` 類型屬性。不同於普通的超連結或 wiki-link，Relation 具備以下特性：

- **具名** — 它有一個屬性名稱，例如 `author` 或 `books`，而不只是一個 URL
- **帶型別** — 它知道目標 Type（例如 `author` 指向 `person`）
- **可選雙向** — 將書關聯到作者時，可以自動將作者關聯回這本書

## 單向 vs. 雙向

**單向** Relation 是單方面的關聯。在書上設定 `author` 會指向一個人，但那個 person Object 不知道這個關聯。

**雙向** Relation 讓兩端保持同步。當你把書的 `author` 關聯到一個人，那個人的 `books` 屬性也會自動更新——反之亦然。

```yaml
# book.yaml
- name: author
  type: relation
  target: person
  bidirectional: true
  inverse: books

# person.yaml
- name: books
  type: relation
  target: book
  multiple: true
  bidirectional: true
  inverse: author
```

## 單值 vs. 多值

- **單值** — 一本書有一個 `author`。重新關聯會覆蓋先前的值。
- **多值** — 一個人可以有很多 `books`。關聯會附加到清單中。

這由 schema 中的 `multiple` 欄位控制。

## Relation 欄位

| 欄位 | 說明 |
|------|------|
| `target` | 目標 Object 的 Type 名稱 |
| `multiple` | 該屬性是否儲存多個值（陣列）。單值 relation 重新關聯時會覆蓋；多值 relation 會附加。 |
| `bidirectional` | 建立關聯時自動同步反向端 |
| `inverse` | 目標 Type schema 上的屬性名稱 |

## 使用 Relation

### 建立關聯

```bash
tmd relation link book/golang-in-action author person/alan-donovan
```

當 `bidirectional: true` 時，這會自動更新書的 `author` 和人物的 `books` 屬性。

### 移除關聯

```bash
tmd relation unlink book/golang-in-action author person/alan-donovan --both
```

使用 `--both` 來同時移除反向端。這只在 schema 中定義了 `bidirectional: true` 和 `inverse` 欄位時才有效。

不加 `--both` 時，只會移除指定的那一端。目標 Object 的 inverse 屬性仍會保留參照。

## 顯示指示符

檢視 Object 的屬性時，Relation 使用方向指示符來表示關聯方向：

| 指示符 | 意義 | 範例 |
|--------|------|------|
| `→` | 正向 relation（此物件關聯到另一個物件） | `→ person/alan-donovan` |
| `←` | 反向 relation（透過 inverse 屬性關聯） | `← book/clean-code` |
| `⟵` | Backlink（來自 wiki-link，非 schema relation） | `⟵ note/my-note` |

## Relation vs. 連結

TypeMD 也支援[連結](/zh-tw/concepts/links)（`[[type/slug]]`）來做非正式的行內引用。Relation 是結構化、由 schema 定義的；連結是自由格式、寫在 Markdown 內文中的。詳細比較請見[連結](/zh-tw/concepts/links)頁面。
