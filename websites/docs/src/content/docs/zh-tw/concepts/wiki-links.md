---
title: Wiki-link
description: 使用 [[target]] 語法在 Object 之間建立行內引用。
sidebar:
  order: 5
---

Wiki-link 是在 Markdown 內文中直接引用其他 Object 的方式，是你在寫作時連結想法最簡單的途徑。

## 什麼是 Wiki-link？

當你在筆記中提到另一個 Object 時，可以用雙括號包住它的 ID：

```markdown
---
title: Go in Action
---

# Notes

Great introduction to Go concurrency patterns.
See also [[note/go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc]].
```

你也可以加上顯示文字：

```markdown
See [[note/go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc|Go Concurrency Patterns]] for details.
```

目標必須是完整的 Object ID（包含 ULID 後綴）。

## 運作方式

1. 當索引同步時（`tmd reindex` 或自動同步），indexer 會從每個 Object 的內文解析 `[[...]]` 模式
2. 目標會與資料庫中的現有 Object 進行比對
3. Wiki-link 記錄儲存在 SQLite 索引中，以便快速查詢反向連結
4. 重新同步時，已從內文中移除的 wiki-link 會自動清理——對應的反向連結也會一併移除
5. 無法解析的目標會儲存為壞連結（可透過 `tmd type validate` 偵測）

## Backlink

Wiki-link 會被自動追蹤。如果 Object A 包含 `[[B]]`，那麼 B 就知道 A 引用了它——這叫做 **backlink（反向連結）**。

Backlink 會在查看 Object 時顯示為內建屬性：

```
backlinks: ⟵ A
```

你不需要手動維護 backlink。TypeMD 會從 Markdown 內文中的 wiki-link 語法自動計算。

## Wiki-link vs. Relation

TypeMD 有兩種連接 Object 的方式，它們的用途不同：

| | Wiki-link | Relation |
|---|-----------|----------|
| 定義於 | Markdown 內文 | Type schema（frontmatter） |
| 結構化 | 否——自由格式的行內引用 | 是——具名、帶型別、可查詢 |
| 雙向 | 透過 backlink（唯讀、自動） | 依 schema 設定 |
| 需要 schema | 否 | 是（`type: relation`） |
| 使用場景 | 非正式引用（另見、提及） | 正式連結（作者、專案成員） |

**簡單原則**：用 Relation 處理屬於資料模型的連結，用 wiki-link 處理筆記中的隨意引用。

## 驗證

`tmd type validate` 包含 wiki-link 驗證階段，會偵測壞連結——參照到不存在的 Object。
