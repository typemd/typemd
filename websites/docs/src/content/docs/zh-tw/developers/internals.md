---
title: 內部機制
description: 內部機制 — 快取、同步、查詢管線與格式化。
sidebar:
  order: 4
---

本頁記錄 TypeMD 的內部機制，這些機制對使用者來說是不可見的，但對貢獻者理解系統運作很重要。涵蓋 schema 快取、增量同步、relation 同步、名稱解析、查詢管線和屬性格式化。

## Schema 快取

Vault 維護一個 type schema 的記憶體快取，避免重複讀取磁碟。

### 快取行為

- **首次載入**：`LoadType("book")` 從 `ObjectRepository.GetSchema("book")` 讀取並快取結果。
- **後續載入**：直接回傳快取的 schema，不讀取磁碟。

### 失效規則

| 觸發條件 | 影響範圍 |
|---------|---------|
| `SaveType()` | 失效特定 type 的快取項目 |
| `DeleteType()` | 失效特定 type 的快取項目 |
| `MigrateSchemas()` | 失效整個快取 |
| `.typemd/types/` 的外部檔案變更 | 失效整個快取 |

外部檔案變更由 TUI 的檔案監視器偵測。當 schema 檔案在 Vault API 外部被修改時，監視器會發送 schema 變更訊息，觸發完整快取失效和資料重新整理。

## 增量同步

TUI 檔案監視器支援增量同步，避免每次檔案變更都要完整重建索引。

### Debounce 和路徑收集

監視器在可設定的 debounce 視窗期間（預設：200ms，可透過 `.typemd/config.yaml` 的 `tui.debounce_ms` 設定）收集變更的檔案路徑。同一視窗內的重複路徑會被去重。

### ReconcileFiles 流程

`Reconciler.ReconcileFiles(paths []string)` 增量調和指定的檔案：

1. 對每個仍存在於磁碟的路徑：透過 `ObjectRepository.Get()` 讀取物件，正規化屬性，產出 `ObjectUpserted` 事件。
2. 對每個已不存在的路徑：產出 `ObjectDeleted` 事件。
3. 增量物件調和完成後：執行完整的 wikilink 調和、完整的 tag relation 調和，以及變更物件的 relation 調和 — 全部以 domain events 形式產出。
4. 呼叫端將回傳的事件傳給 `Projector.Apply()`，寫入 SQLite 索引。

### 回退到完整調和

TUI 在以下情況回退到完整的 `Reconciler.Reconcile()`：

- `fileChangedMsg` 的路徑列表為空
- `Reconciler.ReconcileFiles()` 回傳錯誤
- 初始啟動（首次資料載入）

### Schema 檔案監控

監視器也監控 `.typemd/types/`。Schema 檔案變更會產生獨立的訊息，觸發 schema 快取失效和完整資料重新整理 — 增量同步不適用於 schema 變更。

## Relation 同步

調和期間，Reconciler 讀取每個物件的 frontmatter，識別 type schema 中定義的 relation 屬性，並對每個解析成功的 relation 產出 `RelationIndexed` 事件。Projector 接著將對應的記錄插入 SQLite `relations` 資料表。

### 各 relation 類型的同步行為

- **單值 relation**：每個 relation 產出一個 `RelationIndexed` 事件（例如 `author: person/john-doe-01abc...` → 一個事件）。
- **多值 relation**：每個值一個事件（例如 `books: [book/a-01abc..., book/b-01xyz...]` → 兩個事件）。
- **不存在的目標**：跳過 — 如果參照的物件在磁碟上不存在，不產出事件。
- **非 relation 屬性**：忽略 — 只處理 schema 中 `type: relation` 的屬性。

### 完整調和 vs. 增量調和

- **完整調和**（`Reconciler.Reconcile()`）：產出 `RelationsCleared{NonTagOnly: true}` 事件，接著為所有 relations 產出 `RelationIndexed` 事件。Tag relations 透過 `RelationsCleared{TagsOnly: true}` 獨立管理。
- **增量調和**（`Reconciler.ReconcileFiles()`）：為變更物件產出 `RelationsCleared{ObjectID: id}` 事件，接著產出其重建後的 `RelationIndexed` 事件。未變更物件不受影響。

## 名稱解析（relation 前綴解析）

同步期間，沒有 ULID 後綴的 relation 值會被視為 `type/name` 參照，並解析為完整的物件 ID。

### 解析規則

| 輸入 | 行為 |
|------|------|
| `person/john-doe-01abc...`（有 ULID）| 視為完整 ID，不需解析 |
| `person/john-doe`（無 ULID，唯一匹配）| 解析為完整 ID，檔案被更新 |
| `person/nobody`（無匹配）| 維持不變，在 `ReconcileResult` 中報告為未解析 |
| `person/john`（多重匹配）| 維持不變，在 `ReconcileResult` 中報告為歧義 |

### 名稱索引

Reconciler 從走訪的物件建立每個 type 的名稱索引。每個物件的 slug 和原始名稱都會被索引。同一 type 內的重複名稱會產生歧義項目。

### 自動展開回寫

當前綴成功解析後，Reconciler 會將展開的完整 ID 回寫到物件的 frontmatter 檔案。多個屬性可以在一次檔案寫入中展開。無法解析的前綴維持不變。

### 多值解析

對於 `multiple: true` 的 relation，陣列中的每個值獨立解析。完整 ID 和前綴的混合可以正確處理 — 完整 ID 維持原樣，前綴各自獨立解析。

### ReconcileResult 報告

`ReconcileResult` 包含：

- `Expanded` — 成功解析的前綴數量
- `Unresolved` — 未解析參照的列表，含前綴和原因

### 共用的解析方法

名稱解析方法在 relation 前綴解析和未來的 wiki-link 簡寫解析之間共用。它接受一個每 type 的名稱索引，將 `type/name` 參照解析為完整的物件 ID。

## 結構化查詢篩選

查詢管線使用結構化的 `[]FilterRule` 參數，取代原始的篩選字串。

### FilterRule 結構

每條規則指定一個屬性名稱、運算子和值。空的 slice 回傳所有物件（不篩選）。多條規則以 AND 邏輯組合。

### SQL 對應

| 屬性 | SQL 對應 |
|------|---------|
| `type` | 直接欄位參照：`type = ?` |
| 其他屬性 | JSON 提取：`json_extract(properties, '$.property') = ?` |

### TypeFilter 便利函式

`TypeFilter(typeName string) []FilterRule` 回傳單元素的篩選條件，用於最常見的 type 查詢模式。

## 顯示屬性格式化

`DisplayProperty` 提供兩個格式化方法，用於在 TUI 和 CLI 中呈現屬性值。

### FormatValue()

回傳不含 key 前綴的格式化值：

| 屬性型別 | 範例輸入 | 輸出 |
|---------|---------|------|
| `string` | `"Robert Martin"` | `Robert Martin` |
| `date` | `2024-01-15` | `2024-01-15` |
| `datetime` | `2024-01-15T10:30:00` | `2024-01-15T10:30:00` |
| `multi_select` | `["go", "cli"]` | `[go, cli]` |
| `checkbox`（true）| `true` | `✓` |
| `checkbox`（false）| `false` | （空字串）|
| `relation` | `person/robert-martin-01abc` | `→ person/robert-martin` |
| `backlink` | `note/my-note-01abc` | `⟵ note/my-note` |
| `reverse relation` | `book/clean-code-01abc` | `← book/clean-code` |
| `nil` | `nil` | （空字串）|

### Format()

組合輸出為 `key + ": " + FormatValue()`。用於屬性詳細面板。值的部分委派給 `FormatValue()`。

### 在 view mode 中的使用

View mode 的表格列使用 `FormatValue()` 來格式化屬性欄位和預覽面板，確保所有顯示情境中的格式一致。

## Domain events

Reconciler 和 Projector 透過 domain events 溝通。Reconciler 產出描述變更的事件；Projector 消費事件來更新 SQLite 索引。

### Reconciler → Projector 事件

| 事件 | 用途 | Payload |
|------|------|---------|
| `ObjectUpserted` | 物件需要寫入索引 | `ID`、`Type`、`Filename`、`PropsJSON`、`Body` |
| `ObjectDeleted` | 物件從磁碟移除（stale 清理） | `ID` |
| `RelationsCleared` | 在重建前清除 relations | `ObjectID`（單一物件）、`NonTagOnly`（完整同步）或 `TagsOnly`（tag 同步） |
| `RelationIndexed` | 單一 relation 需要建立索引 | `Name`、`FromID`、`ToID` |
| `WikiLinksSynced` | 物件的 wiki-links 已解析完成 | `ObjectID`、`Links []WikiLinkEntry` |

### ObjectService 事件

這些事件由 `ObjectService` 在使用者發起的操作中產出：

| 事件 | 用途 | Payload |
|------|------|---------|
| `ObjectCreated` | 新物件被建立 | `Object` |
| `ObjectSaved` | 既有物件被儲存 | `Object` |
| `PropertyChanged` | 單一屬性值變更 | `ObjectID`、`Key`、`Old`、`New` |
| `ObjectLinked` | 兩個物件之間建立 relation | `FromID`、`ToID`、`RelName` |
| `ObjectUnlinked` | 兩個物件之間移除 relation | `FromID`、`ToID`、`RelName` |
| `TagAutoCreated` | 同步期間自動建立 tag 物件 | `Tag`、`ReferencedBy` |
