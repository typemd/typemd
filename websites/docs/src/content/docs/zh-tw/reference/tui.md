---
title: tmd (TUI)
description: 瀏覽 vault 的互動式終端介面。
sidebar:
  order: 1
---

啟動 TUI 互動式介面，用於瀏覽 Object。

```bash
tmd
tmd --vault /path/to/vault
tmd --readonly
tmd --reindex
```

## 版面配置

TUI 使用多面板版面配置：

| 面板 | 說明 |
|------|------|
| **Object 列表**（左側） | 依 type 分群顯示 object。群組標題會顯示 type emoji（若有定義）、type 複數名稱和 object 數量（如 `▼ 📚 books (3)`）。所有已定義的 type 都會顯示，即使沒有 object。底部有 `+ New Type` 項目。 |
| **標題**（右上方） | 顯示所選 object 的 type emoji、type 名稱和顯示名稱（如 `📖 book · Clean Code`）。選取 type 標題時顯示 type 名稱。無選取時隱藏。 |
| **內文**（右中） | 顯示 object 的 markdown 內文內容。選取 type 標題時會替換為 **type 編輯器**。 |
| **屬性**（右側） | 顯示 schema 屬性、Relation 和 wiki-link 反向連結。預設隱藏，可用 `p` 切換。在窄終端（< 56 欄）上會自動隱藏。Type 編輯器啟用時隱藏。 |

標題面板橫跨整個右側寬度（內文 + 屬性區域）。

## Type 編輯器

將游標移至 type 群組標題時，右側面板會自動開啟 **type 編輯器**，讓你直接在 TUI 中管理 type schema，無需手動編輯 YAML 檔案：

- **Meta 欄位** — 編輯複數名稱、emoji、unique 約束。名稱為唯讀（不支援 type 改名）。
- **屬性** — 分為 Pinned (Header) 和 Properties 兩區。透過多步驟精靈新增屬性，可編輯 emoji、刪除或用移動模式重新排序。
- **Pin 切換** — 在屬性上按 `p` 可在兩區之間移動。
- **屬性詳情** — 在屬性上按 `Enter` 開啟彈出視窗編輯 metadata（emoji）。
- **刪除 type** — 按 `D`（Shift+d）刪除 type（需確認）。內建的 `tag` type 無法刪除。

每次操作後會立即儲存（無需手動存檔）。

| 按鍵 | 動作（Type 編輯器內） |
|------|----------------------|
| `e` | 編輯 meta 欄位（Plural/Emoji：文字輸入；Unique：切換） |
| `Enter` | 開啟屬性詳情彈出視窗 |
| `a` | 新增屬性（多步驟精靈） |
| `d` | 刪除屬性（需確認） |
| `D` | 刪除 type（需確認） |
| `m` | 進入移動模式（`↑`/`↓` 重新排序，`Enter`/`Esc` 退出） |
| `p` | 切換 pin（在 Pinned/Properties 區之間移動屬性） |
| `Esc` | 將焦點移回側邊欄 |

## 模式與快捷鍵

TUI 會依據操作情境切換不同的模式。狀態列底部會以 `[中括號]` 顯示目前模式。每個模式有各自的快捷鍵。

### `[VIEW]` — 一般瀏覽

預設模式，在側邊欄瀏覽 Object 和 Type。

| 按鍵 | 動作 |
|------|------|
| `↑`/`k`、`↓`/`j` | 瀏覽列表 / 捲動內容 |
| `Enter` | 選取 Object / 聚焦 type 編輯器 |
| `Space` | 展開/收合群組 |
| `Tab` | 在面板之間循環焦點 |
| `n` | 新增 Object（在目前的 type 群組中） |
| `e` | 進入編輯模式（內文面板聚焦時） |
| `/` | 進入搜尋模式 |
| `p` | 切換屬性面板 |
| `w` | 切換自動換行 |
| `[`/`]` | 縮小/放大焦點面板 |
| `?`/`h` | 開啟快捷鍵說明 |
| `q`/`Ctrl+C` | 離開 |

### `[TYPE]` — Type 編輯器

Type 編輯器面板取得焦點時啟用。在 type 標題上按 `Tab` 或 `Enter` 進入。

| 按鍵 | 動作 |
|------|------|
| `↑`/`k`、`↓`/`j` | 瀏覽 meta 欄位和屬性 |
| `Enter` | 開啟屬性詳情彈出視窗 / 在 `+ Add Property` 上啟動精靈 |
| `e` | 編輯 meta 欄位（Plural/Emoji：文字輸入，Unique：切換） |
| `a` | 新增屬性（啟動精靈） |
| `d` | 刪除屬性（需確認） |
| `D` | 刪除 type（需確認） |
| `m` | 進入移動模式 |
| `p` | 切換屬性 pin |
| `Tab`/`Esc` | 焦點回到側邊欄 |
| `q`/`Ctrl+C` | 離開 |

### `[EDIT]` — 內文編輯

編輯 Object 內文時啟用。在 VIEW 模式且內文面板聚焦時按 `e`。

| 按鍵 | 動作 |
|------|------|
| 一般文字按鍵 | 編輯內文 |
| `Esc` | 退出編輯模式（有變更時自動儲存） |

所有導航按鍵（`j`/`k`、`Tab`）被攔截，不會切換面板。面板邊框會變為橙色。

若檔案在編輯期間被外部程式修改，會出現 **`[CONFLICT]`** 警告。按 `y` 覆蓋寫入、`n` 從磁碟重新載入、`Esc` 關閉警告。

### `[MOVE]` — 屬性重新排序

在 type 編輯器中重新排序屬性時啟用。在 TYPE 模式的屬性上按 `m`。

| 按鍵 | 動作 |
|------|------|
| `↑`/`k`、`↓`/`j` | 上下移動屬性（與鄰居交換） |
| `Enter`/`Esc` | 確認並儲存新順序 |

跨越 Pinned/Properties 區域邊界時，會自動切換 pin 值。

### `[ADD PROPERTY]` — 新增屬性精靈

新增 type 屬性的多步驟精靈。在 TYPE 模式按 `a` 或在 `+ Add Property` 上按 `Enter`。

| 步驟 | 輸入 | 按鍵 |
|------|------|------|
| 1. 名稱 | 屬性名稱文字輸入 | `Enter`：下一步、`Esc`：取消 |
| 2. 型別 | 從列表選擇（string、number、date⋯） | `↑`/`↓`：選擇、`Enter`：下一步、`Esc`：返回 |
| 2b. 選項 | 逗號分隔的值（select/multi_select 用） | `Enter`：建立、`Esc`：返回 |
| 3. Relation | 目標 type、multiple、bidirectional、inverse 名稱 | `Tab`：下一欄位、`Enter`/`Space`：切換/確認、`Esc`：返回 |

### `[PROPERTY]` — 屬性詳情彈出視窗

編輯屬性 metadata 的彈出視窗。在 TYPE 模式的屬性上按 `Enter`。

| 按鍵 | 動作 |
|------|------|
| `Enter`/`e` | 編輯選取的欄位 |
| `Esc` | 關閉彈出視窗 |

編輯欄位時：`Enter` 儲存、`Esc` 取消（還原）。

### `[DELETE]` / `[DELETE TYPE]` — 刪除確認

刪除屬性（`d`）或 type（`D`）時顯示。

| 按鍵 | 動作 |
|------|------|
| `y` | 確認刪除 |
| `n`/`Esc` | 取消 |

### `[NEW TYPE]` / `[NEW OBJECT]` — 建立輸入

建立新 type（`+ New Type`）或 Object（`n`）時顯示。

| 按鍵 | 動作 |
|------|------|
| 文字按鍵 | 輸入 type 名稱 / Object 名稱 |
| `Enter` | 建立 |
| `Esc` | 取消 |

### `[READONLY]` — 唯讀模式

以 `--readonly` 啟動時啟用。`e` 和 `n` 鍵被停用，不執行任何寫入操作，快捷鍵說明中會隱藏編輯相關的快捷鍵。

## 會話狀態

TUI 在離開時會自動將會話狀態儲存到 `.typemd/tui-state.yaml`。下次啟動時會還原：

- **選取的 Object 或 Type** — 游標回到同一個 Object（以 ID 識別）或 type 標題（以名稱識別）
- **展開的群組** — Type 群組維持你離開時的展開/收合狀態
- **面板尺寸** — 左側面板和屬性面板的寬度
- **屬性面板可見性** — 屬性面板是否顯示
- **捲動偏移** — Object 列表的垂直捲動位置

若之前選取的 Object 已被刪除，TUI 會 fallback 到同 type 群組的第一個 Object，再到整體第一個 Object。焦點啟動時一律在側邊欄。

搜尋狀態不會被保留——每次啟動都從全新的畫面開始。

若狀態檔案遺失或損毀，TUI 會靜默 fallback 到預設啟動行為。

## 重建索引

`--reindex` flag 會強制將 `objects/` 目錄完整同步到資料庫，清理孤立的 relation，並重建全文搜尋索引。這是一個全域 flag，可以與任何指令組合使用。

> **注意：** 開啟 vault 時，TypeMD 會在索引為空或缺失時自動同步。只有在 vault 未開啟期間編輯了檔案，才需要使用 `--reindex`。

```bash
# 重建索引並啟動 TUI
tmd --reindex

# 重建索引並啟動 MCP server
tmd mcp --reindex

# 重建索引並執行查詢
tmd query "type=book" --reindex
```

### 孤立 relation 清理

當物件從磁碟中刪除時，指向或來自該物件的 relation 會變成孤立的。在 reindex 過程中，這些 dangling reference 會自動被偵測並從索引中移除。系統會顯示受影響的 relation 列表：

```
Warning: Found 2 orphaned relation(s):
  book/golang-in-action -> person/deleted-author (relation: "author")
  person/deleted-author -> book/golang-in-action (relation: "books")
Orphaned relations have been removed from the index.
```

> **注意：** 這只會清理 SQLite 索引。`.md` 檔案中的 frontmatter 不會被修改 — 建議在刪除物件前先使用 `tmd relation unlink --both` 正確移除 relation。

## 自動重新整理

TUI 透過 fsnotify 監控 `objects/` 目錄。當檔案被建立、修改或刪除時，會自動同步資料庫並重新整理畫面（200ms 防抖），並盡可能保持目前的選取狀態。
