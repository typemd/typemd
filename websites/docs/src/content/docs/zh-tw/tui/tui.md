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
| **內文**（右中） | 頂部顯示釘選屬性（若有），接著顯示 object 的 markdown 內文內容。選取 type 標題時會替換為 **type 編輯器**。 |
| **屬性**（右側） | 顯示 schema 屬性、Relation 和 wiki-link 反向連結。預設隱藏，可用 `p` 切換。在窄終端（< 56 欄）上會自動隱藏。Type 編輯器或樣板編輯器啟用時隱藏。 |

標題面板橫跨整個右側寬度（內文 + 屬性區域）。

## Type 編輯器

將游標移至 type 群組標題時，右側面板會自動開啟 **type 編輯器**，讓你直接在 TUI 中管理 type schema，無需手動編輯 YAML 檔案：

- **Meta 欄位** — 編輯複數名稱、emoji、unique 約束。名稱為唯讀（不支援 type 改名）。
- **屬性** — 分為 Pinned (Header) 和 Properties 兩區。透過多步驟精靈新增屬性，可編輯 emoji、刪除或用移動模式重新排序。
- **Pin 切換** — 在屬性上按 `p` 可在兩區之間移動。
- **屬性詳情** — 在屬性上按 `Enter` 開啟彈出視窗編輯 metadata（emoji）。
- **刪除 type** — 按 `D`（Shift+d）刪除 type（需確認）。內建的 `tag` type 無法刪除。
- **樣板** — Templates 區塊列出該 type 可用的樣板。在樣板上按 `Enter` 開啟**樣板編輯器**。在 `+ Add Template` 上按 `Enter` 可輸入名稱建立新樣板。

每次操作後會立即儲存（無需手動存檔）。

| 按鍵 | 動作（Type 編輯器內） |
|------|----------------------|
| `e` | 編輯 meta 欄位（Plural/Emoji：文字輸入；Unique：切換） |
| `Enter` | 開啟屬性詳情彈出視窗 / 開啟樣板 / 啟動新增精靈 |
| `a` | 新增屬性（多步驟精靈） |
| `d` | 刪除屬性（需確認） |
| `D` | 刪除 type（需確認） |
| `m` | 進入移動模式（`↑`/`↓` 重新排序，`Enter`/`Esc` 退出） |
| `p` | 切換 pin（在 Pinned/Properties 區之間移動屬性） |
| `Esc` | 將焦點移回側邊欄 |

## 釘選屬性

Type schema 可以為屬性設定 `pin` 值（正整數），將其顯示在內文面板的頂部。釘選屬性會依 pin 值排序（數值越小優先順序越高），以 key-value 的形式顯示，並在內文內容前加上分隔線。設有 emoji 的屬性會在值旁邊顯示 emoji。

釘選屬性會從屬性面板中**排除**，以避免重複顯示。

```yaml
# .typemd/types/book/schema.yaml
properties:
  - name: status
    type: select
    emoji: 📋
    pin: 1        # 在內文面板中最先顯示
  - name: rating
    type: number
    emoji: ⭐
    pin: 2        # 第二個顯示
  - name: title
    type: string  # 沒有 pin — 顯示在屬性面板中
```

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
| `v` | 開啟目前 type 的 view 模式 |
| `n` | 新增 Object 並編輯內文（在目前的 type 群組中） |
| `N` | 快速建立 — 批次模式（在目前的 type 群組中） |
| `e` | 進入編輯模式（內文面板聚焦時） |
| `r` | 重新命名 Object（在標題面板中直接編輯） |
| `/` | 進入搜尋模式 |
| `p` | 切換屬性面板 |
| `w` | 切換自動換行 |
| `[`/`]` | 縮小/放大焦點面板 |
| `?`/`h` | 開啟快捷鍵說明 |
| `q`/`Ctrl+C` | 離開 |

### View 模式

在 type 群組標題或物件上按 `v` 進入 **view 模式**。三欄佈局會被全寬顯示取代。View 模式支援兩種 layout：

- **List layout** — 每列顯示物件的 emoji 和名稱，後方可選擇性附加屬性值，以 " · " 分隔。這是新 view 的預設 layout。
- **Table layout** — 欄位式顯示，包含 NAME 欄和屬性欄位。欄位標題會顯示排序指示符號（↑/↓）。欄位數量會依終端機寬度自動調整。

| 按鍵 | 動作 |
|------|------|
| `↑`/`k`、`↓`/`j` | 瀏覽列表 |
| `Enter`/`Space` | 開啟物件詳情 / 展開收合群組 |
| `e` | 開啟 view 編輯器（右側分割面板） |
| `p` | 切換預覽面板（右側） |
| `Esc` | 返回側邊欄 |
| `q`/`Ctrl+C` | 離開 |

開啟**預覽**（`p`）時，右側面板顯示游標所在物件的屬性和內文摘要。上下移動時預覽會自動跟隨更新。預覽和 view 編輯器互斥，同時只能開啟一個。

若 type 有多個已儲存的 [view](/zh-tw/advanced/file-structure#view)，按 `v` 時會彈出選擇選單。若只有預設 view，則直接進入。

Type 編輯器也會顯示 **Views** 區塊，可以瀏覽已儲存的 view 並透過 `+ Add View` 新增。

#### View 編輯器

在 view 模式中按 `e` 開啟 **view 編輯器**，以右側分割面板顯示。編輯器包含五個區塊：

| 區塊 | 說明 |
|------|------|
| **Layout** | 按 `Enter` 在 `list` 和 `table` 之間切換。 |
| **Columns** | 新增、移除和重新排序要顯示的屬性。在 list layout 中，選取的屬性會顯示為行內值。在 table layout 中，會成為欄位標題。 |
| **Filter** | 新增篩選規則（屬性 + 運算子 + 值）來限縮顯示的物件。 |
| **Sort** | 定義排序規則（屬性 + 方向）控制列的排列順序。 |
| **Group By** | 依一或多個屬性將物件分群。 |

| 按鍵 | 動作 |
|------|------|
| `Tab` | 下一個區塊 |
| `Shift+Tab` | 上一個區塊 |
| `Enter` | 編輯欄位 / 切換 layout / 新增規則 |
| `x`/`d` | 刪除選取的規則 |
| `K` | 將選取的規則向上移動 |
| `J` | 將選取的規則向下移動 |
| `D` | 刪除整個 view（需確認） |
| `Esc` | 關閉編輯器並返回 view 列表 |

變更會在每次編輯後自動儲存。

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

### `[TEMPLATE]` — 樣板編輯器

檢視或編輯樣板時啟用。在 type 編輯器中的樣板上按 `Enter` 進入。

| 按鍵 | 動作 |
|------|------|
| `e` | 編輯樣板內文（進入文字區域編輯模式） |
| `d` | 刪除樣板（需確認） |
| `Tab` | 在內文和屬性面板之間切換焦點 |
| `↑`/`k`、`↓`/`j` | 捲動內文（內文聚焦時） / 瀏覽屬性（屬性聚焦時） |
| `Enter` | 編輯屬性值（屬性聚焦時） |
| `Esc` | 返回 type 編輯器 |
| `q`/`Ctrl+C` | 離開 |

編輯內文時：`Esc` 儲存、`Ctrl+C` 取消。編輯屬性時：`Enter` 確認、`Esc` 取消。清空屬性值會將該屬性從樣板 frontmatter 中移除。

### `[NEW TEMPLATE]` — 新增樣板輸入

在 type 編輯器中透過 `+ Add Template` 建立樣板時顯示。

| 按鍵 | 動作 |
|------|------|
| 文字按鍵 | 樣板名稱 |
| `Enter` | 建立空白樣板 |
| `Esc` | 取消 |

### `[DELETE]` / `[DELETE TYPE]` — 刪除確認

刪除屬性（`d`）、type（`D`）或樣板（在樣板編輯器中按 `d`）時顯示。

| 按鍵 | 動作 |
|------|------|
| `y` | 確認刪除 |
| `n`/`Esc` | 取消 |

### `[NEW TYPE]` — 建立新 Type

透過 `+ New Type` 觸發。**標題面板**會變成行內建立表單，包含三個欄位：emoji（選填）、name（必填）、plural（選填）。右側面板會顯示 type schema 的**即時預覽**。

| 按鍵 | 動作 |
|------|------|
| 文字按鍵 | 在聚焦的欄位中輸入 |
| `Tab` | 循環聚焦：name → plural → emoji |
| `Enter` | 建立 type 並開啟 type editor |
| `Esc` | 取消 |

### `[NEW OBJECT]` — 建立並編輯模式

按 `n` 觸發。**標題面板**會變成行內建立表單，包含名稱輸入和（有樣板時）樣板循環選擇器。內文和屬性面板會顯示所選樣板內容的**即時預覽**。

表單配置：`📚 book · [name█] 📝 review`

| 按鍵 | 動作 |
|------|------|
| 文字按鍵 | 輸入 Object 名稱 |
| `Tab` | 切換名稱與樣板欄位的焦點 |
| `←`/`→` | 循環切換樣板（樣板欄位聚焦時） |
| `Enter` | 建立 Object 並進入內文編輯模式 |
| `Esc` | 取消 |

- 若 type 有**多個樣板**，`Tab` 可切換至樣板選擇器，`←`/`→` 在各樣板和 `(none)` 選項之間循環。只有單一樣板時會自動選取，顯示為靜態標籤。
- 若 type 定義了 **name template**（如 `{{ date:YYYY-MM-DD }}`），名稱會預先填入計算後的值。可以編輯或直接按 Enter 接受。
- 切換樣板會即時更新內文和屬性面板的預覽。
- 重複名稱錯誤（適用於 `unique: true` 的 type）會在標題面板中行內顯示，修改文字後自動清除。

### `[QUICK CREATE]` — 批次建立模式

按 `N` 觸發。與「建立並編輯」模式相同的標題面板表單，但建立後會留在輸入模式，方便快速連續建立多個 Object。

選取的樣板在整個批次中持續有效 — 所有 Object 使用同一個樣板。Name template 會預先填入但可隨時編輯。

| 按鍵 | 動作 |
|------|------|
| 文字按鍵 | 輸入 Object 名稱 |
| `Tab` | 切換名稱與樣板欄位的焦點 |
| `←`/`→` | 循環切換樣板（樣板欄位聚焦時） |
| `Enter` | 建立 Object、清空輸入、準備建立下一個 |
| `Esc` | 退出批次模式（選取最後建立的 Object） |

每次建立成功後，標題面板會短暫顯示成功訊息（如 `✓ Created: my-book`）。

### `[READONLY]` — 唯讀模式

以 `--readonly` 啟動時啟用。`e`、`n` 和 `N` 鍵被停用，不執行任何寫入操作，快捷鍵說明中會隱藏編輯相關的快捷鍵。

## 會話狀態

TUI 在離開時會自動將會話狀態儲存到 `.typemd/tui-state.yaml`。下次啟動時會還原：

- **選取的 Object 或 Type** — 游標回到同一個 Object（以 ID 識別）或 type 標題（以名稱識別）
- **展開的群組** — Type 群組維持你離開時的展開/收合狀態
- **面板尺寸** — 左側面板和屬性面板的寬度
- **屬性面板可見性** — 屬性面板是否顯示
- **捲動偏移** — Object 列表的垂直捲動位置
- **View mode** — 若離開時處於 view mode，下次啟動時會自動還原該 view，包含游標位置、捲動偏移和 view 內展開的群組

若之前選取的 Object 已被刪除，TUI 會 fallback 到同 type 群組的第一個 Object，再到整體第一個 Object。焦點一律從側邊欄開始，確保一致的操作體驗（還原 view mode 時除外，此時焦點會設定在內文面板）。

若儲存的 view 所屬的 type 或 view 已被刪除，TUI 會 fallback 到該 type 的 default view，若無任何 view 則回到正常的 sidebar 模式。

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

# 重建索引並列出 Object
tmd object list --reindex
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

TUI 透過 fsnotify 監控 `objects/` 目錄。當檔案被建立、修改或刪除時，會對變更的檔案進行增量索引同步並重新整理畫面（200ms 防抖），並盡可能保持目前的選取狀態。防抖間隔可透過 `.typemd/config.yaml` 中的 `tui.debounce_ms` 自訂。

TUI 也會監控 `.typemd/types/` 目錄的 schema 變更。當 type schema 被外部修改時，schema cache 會被清除並觸發完整重新整理。
