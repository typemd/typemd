---
title: Wiki-links
description: 在內文中用 wiki 語法連結 Object，並自動追蹤反向連結。
sidebar:
  order: 7.5
---

Wiki-links 讓你在 Markdown 內文中用 `[[target]]` 語法參照其他 Object。系統會自動計算反向連結（backlinks），讓被連結的 Object 知道誰引用了它。

## 語法

```markdown
參見 [[book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz]]。
```

可加上顯示文字：

```markdown
參見 [[book/clean-code-01jqr3k5mpbvn8e0f2g7h9txyz|Clean Code]]。
```

目標必須是完整的 Object ID（包含 ULID 後綴）。

## 運作方式

1. 當索引同步時（`tmd reindex` 或自動同步），indexer 會從每個 Object 的內文解析 `[[...]]` 模式
2. 目標會與資料庫中的現有 Object 進行比對
3. Wiki-link 記錄儲存在 SQLite 索引中，以便快速查詢反向連結
4. 無法解析的目標會儲存為壞連結（可透過 `tmd validate` 偵測）

## 反向連結

反向連結會以內建的 `backlinks` 屬性顯示。如果 Object A 包含 `[[B]]`，則 B 的屬性會顯示：

```
backlinks: ⟵ A
```

反向連結在 `tmd show` 和 TUI 屬性面板中都可見。

## Wiki-links 與 Relation 的差異

| | Wiki-links | Relation |
|--|-----------|----------|
| **定義位置** | Markdown 內文 | Type schema 屬性 |
| **語法** | `[[target]]` | `tmd link` 指令 |
| **需要 Schema** | 否 | 是（`type: relation`） |
| **雙向** | 透過反向連結（唯讀） | 可設定（`bidirectional: true`） |
| **儲存位置** | `wikilinks` 表 | `relations` 表 + frontmatter |

## 驗證

`tmd validate` 包含 wiki-link 驗證階段，會偵測壞連結——參照到不存在的 Object。
