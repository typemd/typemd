---
title: 連結
description: 使用 [[target]] 語法在 Object 之間建立行內引用。
sidebar:
  order: 5
---

連結是在 Markdown 內文中使用 wiki-link 語法（`[[...]]`）直接引用其他 Object 的方式，是你在寫作時連結想法最簡單的途徑。

## 語法

當你在筆記中提到另一個 Object 時，可以用雙括號包住它的 ID：

```markdown
---
name: Go in Action
---

# Notes

Great introduction to Go concurrency patterns.
See also [[note/go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc]].
```

你也可以加上顯示文字：

```markdown
See [[note/go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc|Go Concurrency Patterns]] for details.
```

### 簡寫語法

你不需要每次都打完整 ID。TypeMD 支援更短的格式，在同步時自動解析：

| 格式 | 範例 | 解析方式 |
|------|------|---------|
| `[[type/name-ulid]]` | `[[book/clean-code-01abc...]]` | 精確匹配（完整 ID） |
| `[[type/name]]` | `[[book/clean-code]]` | 在指定 type 內依名稱解析 |
| `[[name]]` | `[[clean-code]]` | 在**來源物件的同 type** 內依名稱解析 |

```markdown
See [[book/clean-code]] for more.
Also check [[my-other-note]].
```

當簡寫連結只匹配到一個物件時，同步時會自動展開為完整 ID 並寫回你的檔案。如果簡寫匹配到多個物件，則視為壞連結——使用 `tmd type validate` 查看有歧義的匹配。

## 運作方式

1. 當索引同步時（vault 開啟時自動執行），indexer 會從每個 Object 的內文解析 `[[...]]` 模式
2. 目標會與現有 Object 進行比對——完整 ID 精確匹配，簡寫目標在對應 type 內依名稱解析
3. 解析成功的簡寫目標會以完整 ID 寫回來源檔案（如同 relation 名稱展開）
4. 連結記錄儲存在 SQLite 索引中，以便快速查詢反向連結
5. 重新同步時，已從內文中移除的連結會自動清理——對應的反向連結也會一併移除
6. 無法解析或有歧義的目標會儲存為壞連結（可透過 `tmd type validate` 偵測）

## 反向連結（Backlink）

每個連結都會自動建立反向連結。如果 Object A 包含 `[[B]]`，那麼 B 就知道 A 引用了它——這叫做 **backlink（反向連結）**。

反向連結是發現相關內容的強大方式，完全不需要手動操作。你只要在一個方向建立連結，TypeMD 就會自動追蹤反向關係。

Backlink 會在查看 Object 時顯示為內建屬性：

```
backlinks: ⟵ A
```

你不需要手動維護 backlink。TypeMD 會從 Markdown 內文中的 wiki-link 語法自動計算。當來源 Object 內文中的連結被移除時，對應的反向連結會在下次同步時自動清理。

## 連結 vs. Relation

TypeMD 有兩種連接 Object 的方式，它們的用途不同：

| | 連結 | Relation |
|---|------|----------|
| 定義於 | Markdown 內文 | Type schema（frontmatter） |
| 結構化 | 否——自由格式的行內引用 | 是——具名、帶型別、可查詢 |
| 雙向 | 透過 backlink（唯讀、自動） | 依 schema 設定 |
| 需要 schema | 否 | 是（`type: relation`） |
| 使用場景 | 非正式引用（另見、提及） | 正式連結（作者、專案成員） |

**簡單原則**：用 Relation 處理屬於資料模型的連結，用連結處理筆記中的隨意引用。

## 修正簡寫連結

使用 `tmd fix wikilinks` 可以一次將 vault 中所有簡寫 wiki-link 展開為完整 ID。這在匯入內容或批次編輯檔案後特別有用。

## 驗證

`tmd type validate` 包含連結驗證階段，會偵測壞連結——參照到不存在的 Object。對於匹配到多個物件的簡寫連結（有歧義），驗證會回報有歧義的目標，並列出匹配的完整 ID 清單。
