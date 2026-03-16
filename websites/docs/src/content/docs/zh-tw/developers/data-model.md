---
title: 資料模型
description: TypeMD 如何索引資料及管理 SQLite 加速層。
sidebar:
  order: 1
---

TypeMD 採用雙儲存設計：Markdown 檔案是 source of truth，SQLite 資料庫提供快速索引和搜尋。本頁說明索引層的技術細節。

使用者導向的檔案結構說明，請參閱[檔案結構](/zh-tw/advanced/file-structure)。

## 檔案即 Source of Truth

最根本的設計原則是**檔案永遠是 source of truth**。SQLite 索引是加速層，隨時可以從檔案重建。如果索引被刪除或損壞，不會有任何資料遺失——執行 `tmd reindex` 或開啟 vault 就會重建。

這個設計帶來以下好處：

- 使用 Git 進行版本控制（只有 Markdown 和 YAML 檔案重要）
- 用任何文字編輯器手動編輯
- 透過檔案同步工具（Dropbox、iCloud、Syncthing）同步
- 不需要資料庫遷移就能在不同機器間搬移

## SQLite 索引

索引儲存在 `.typemd/index.db`，包含：

- **Object 中繼資料** — type、檔名和所有 frontmatter 屬性
- **Wiki-link 記錄** — 從 Object 內文提取，用於反向連結追蹤
- **全文搜尋索引** — 由 [FTS5](https://www.sqlite.org/fts5.html) 驅動，涵蓋檔名、屬性和內文

索引檔案不應直接編輯。它完全由 TypeMD 管理。

## 自動同步行為

索引在以下情況會自動同步：

- **開啟 Vault 時** — 開啟 vault 且資料庫為空或缺失時（例如 `git clone` 後首次開啟），TypeMD 會走訪所有 Object 檔案並填入索引
- **CLI/TUI 操作時** — 建立、儲存或刪除 Object 時會立即更新索引
- **手動重建** — 在 TypeMD 外部編輯檔案後，執行 `tmd reindex` 重建索引

同步流程由 **Projector** 元件處理，它透過 Repository 走訪所有 Object 檔案並 upsert 到索引中。Projector 在系統中的角色請參閱[系統架構](/zh-tw/developers/architecture)。

## 查詢架構

TypeMD 提供兩種查詢路徑：

### 結構化查詢

使用 `tmd query` 依屬性篩選 Object。條件使用 `key=value` 格式，以空格分隔（AND 邏輯）：

```bash
tmd query "type=book"
tmd query "type=book status=reading"
tmd query "type=book" --json
```

結構化查詢針對 SQLite 索引執行以提升效能，回傳輕量的 `ObjectResult` 投影而非完整的 Object 實體。

### 全文搜尋

使用 `tmd search` 搜尋檔名、屬性和內文。由 SQLite FTS5 驅動：

```bash
tmd search "concurrency"
tmd search "golang" --json
```

### TUI 搜尋

在 TUI 中，按 `/` 進入搜尋模式。結果會即時篩選。按 `Esc` 清除結果並回到完整列表。

## 唯一性約束機制

Type schema 可透過設定 `unique: true` 來啟用 name 唯一性檢查。啟用後，TypeMD 會阻止建立同一 type 下擁有相同 `name` 值的多個 Object。

```yaml
# .typemd/types/person.yaml
name: person
unique: true  # 同一個人名稱只能有一個
properties:
  - name: role
    type: string
```

唯一性在建立時透過檢查索引中是否已存在相同 type 和 name 的 Object 來強制執行。也可透過 `tmd type validate` 驗證。內建的 `tag` type 預設啟用 `unique: true`。

## Vault 設定

`.typemd/config.yaml` 是選填的 vault 層級設定檔，使用介面層命名空間：

```yaml
# .typemd/config.yaml
cli:
  default_type: idea
```

目前支援的設定：

| Key | 說明 |
|-----|------|
| `cli.default_type` | `tmd object create` 省略 type 參數時使用的預設 object type |

Config 檔在 vault 開啟時載入。若檔案不存在或為空，所有設定使用零值（不會出錯）。無效的 YAML 會產生錯誤。
