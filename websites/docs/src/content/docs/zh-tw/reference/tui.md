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
| **Object 列表**（左側） | 依 type 分群顯示 object。群組標題會顯示 type emoji（若有定義）、type 複數名稱和 object 數量（如 `▼ 📚 books (3)`）。 |
| **標題**（右上方） | 顯示所選 object 的 type emoji、type 名稱和顯示名稱（如 `📖 book · Clean Code`）。未選取 object 時隱藏。 |
| **內文**（右中） | 顯示 object 的 markdown 內文內容。 |
| **屬性**（右側） | 顯示 schema 屬性、Relation 和 wiki-link 反向連結。預設隱藏，可用 `p` 切換。在窄終端（< 56 欄）上會自動隱藏。 |

標題面板橫跨整個右側寬度（內文 + 屬性區域）。

## 鍵盤快捷鍵

| 按鍵 | 動作 |
|------|------|
| `↑` / `k` | 向上移動（瀏覽列表 / 捲動詳情） |
| `↓` / `j` | 向下移動（瀏覽列表 / 捲動詳情） |
| `Enter` / `Space` | 選取 Object / 展開/收合群組 |
| `Tab` | 在面板之間循環焦點 |
| `e` | 進入編輯模式（聚焦在內文或屬性面板時） |
| `/` | 進入搜尋模式 |
| `Esc` | 退出編輯模式 / 離開搜尋 / 清除結果 |
| `p` | 切換屬性面板 |
| `w` | 切換自動換行 |
| `[` / `]` | 縮小 / 放大焦點面板 |
| `?` / `h` | 開啟快捷鍵說明 |
| `q` / `Ctrl+C` | 離開 |

## 編輯模式

當**內文面板**取得焦點時，按下 `e` 即可進入原地編輯模式：唯讀的 viewport 會被替換為 textarea，並預填目前 object 的內文內容。你可以自由輸入、刪除及使用標準文字編輯快捷鍵導航。狀態列會顯示 `[EDIT]`，面板邊框也會變為橙色。

編輯期間，所有全域導航按鍵（`j`/`k`、`Tab`）都會被攔截，不會切換面板。按下 `Esc` 退出編輯模式時，變更會**自動儲存**至磁碟，並返回 `[VIEW]` 模式。

若檔案在你載入後已被外部程式修改，狀態列會顯示**衝突警告**。按 `y` 覆蓋寫入、`n` 從磁碟重新載入，或 `Esc` 取消。

在屬性面板按 `e` 同樣會進入編輯模式（僅顯示視覺指示；屬性編輯功能尚未支援）。

## 唯讀模式

使用 `--readonly` 啟動可防止任何編輯操作。在此模式下：

- `e` 鍵被停用——無法進入編輯模式
- 不會對 Object 檔案或 SQLite index 執行任何寫入操作
- 狀態列顯示 `[READONLY]` 而非 `[VIEW]`
- 快捷鍵說明中會隱藏編輯相關的快捷鍵

適用於在共用終端機、簡報或自動化情境中安全瀏覽 vault。

## 會話狀態

TUI 在離開時會自動將會話狀態儲存到 `.typemd/tui-state.yaml`。下次啟動時會還原：

- **選取的 Object** — 游標回到同一個 Object（以 Object ID 識別）
- **展開的群組** — Type 群組維持你離開時的展開/收合狀態
- **面板尺寸** — 左側面板和屬性面板的寬度
- **屬性面板可見性** — 屬性面板是否顯示
- **焦點面板** — 哪個面板擁有焦點（左側、內文或屬性）
- **捲動偏移** — Object 列表的垂直捲動位置

若之前選取的 Object 已被刪除，TUI 會 fallback 到同 type 群組的第一個 Object，再到整體第一個 Object。

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
