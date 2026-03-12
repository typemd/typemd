---
title: 資料模型
description: TypeMD 如何儲存和索引資料。
sidebar:
  order: 6
---

## 儲存

Object 以 Markdown 檔案搭配 YAML frontmatter 儲存在 `objects/<type>/` 底下。完整的 Object ID 格式為 `type/<slug>-<ulid>`，例如 `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`。透過 CLI 建立的 Object，檔名會自動附加 ULID（26 位小寫字元）以確保唯一性。

### 系統屬性

所有 Object 都有四個由 TypeMD 管理的系統屬性：

| 屬性 | 說明 |
|------|------|
| `name` | 顯示名稱，建立時從 slug 自動填入 |
| `description` | 選填的單行摘要，用於列表顯示和搜尋結果 |
| `created_at` | 建立時間戳記，RFC 3339 格式（僅設定一次，不會修改） |
| `updated_at` | 最後修改時間戳記，RFC 3339 格式（每次儲存時更新） |

系統屬性在 frontmatter 中永遠依上述順序排在最前面，接著才是 schema 定義的屬性。這些名稱是保留的，不能在 type schema 或 shared properties 中使用。

```
vault/
├── .typemd/
│   ├── types/              # Type schema 定義（YAML）
│   │   ├── book.yaml
│   │   └── person.yaml
│   ├── index.db            # SQLite 索引（自動更新）
│   └── tui-state.yaml      # TUI 會話狀態（自動儲存）
└── objects/
    ├── book/
    │   └── golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md
    └── person/
        └── alan-donovan-01jqr3k5mpbvn8e0f2g7h9txyz.md
```

## 索引

TypeMD 使用 SQLite 搭配 FTS5 進行索引。索引儲存在 `.typemd/index.db`，包含：

- Object 中繼資料（Type、檔名、屬性）
- 從 Object 內文提取的 Wiki-link 記錄（用於反向連結追蹤）
- 涵蓋檔名、屬性和內文的全文搜尋索引

開啟 vault 時，若資料庫為空或缺失（例如 clone 後首次開啟），索引會自動同步。使用 TUI 或 CLI 指令時也會保持更新。在 TypeMD 外部手動編輯檔案後，使用 `tmd --reindex` 重建。

## 查詢

TypeMD 提供兩種尋找 Object 的方式：結構化查詢和全文搜尋。

### 結構化查詢

使用 `tmd query` 依屬性篩選 Object。條件使用 `key=value` 格式，以空格分隔（AND 邏輯）。

```bash
tmd query "type=book"
tmd query "type=book status=reading"
tmd query "type=book" --json
```

### 全文搜尋

使用 `tmd search` 搜尋檔名、屬性和內文。由 SQLite FTS5 驅動。

```bash
tmd search "concurrency"
tmd search "golang" --json
```

### TUI 搜尋

在 TUI 中，按 `/` 進入搜尋模式。結果會即時篩選。按 `Esc` 清除結果並回到完整列表。
