---
title: 設計決策
description: 重要架構決策及其背後的考量。
sidebar:
  order: 3
---

本頁記錄 TypeMD 開發過程中的重要設計決策。每項決策說明了做了什麼、為什麼，以及考慮過哪些替代方案。理解這些決策有助於貢獻者在擴展系統時做出一致的選擇。

## Relations：雙重儲存

**決策**：Relations 同時儲存在 YAML frontmatter 和 SQLite `relations` 資料表中。

Frontmatter 是真實來源（source of truth），確保檔案的可攜性——物件是自包含的 Markdown 檔案，可以做版本控制、同步和手動編輯。資料庫提供快速的反向查詢和關聯查詢，這些操作若從檔案計算會非常昂貴。

當兩個儲存區不一致時（例如手動編輯檔案後），`tmd reindex` 會從 frontmatter 重建資料庫。

**被否決的替代方案**：
- 僅儲存在資料庫——違反 local-first、檔案即真實來源的原則
- 僅儲存在 frontmatter，反向查詢靠掃描檔案——每次查詢的 O(n) 掃描成本過高

### 單值 vs. 多值行為

單值 relation（如 `author`）重新連結時會覆寫。多值 relation（如 `books`）會追加並拒絕重複。這符合使用者的期望：「設定作者」vs.「新增一本書」。

### 雙向透過明確的 inverse

雙向 relation 要求兩端都在 type schema 中宣告 inverse property。當連結 A→B 時，系統自動寫入反向的 B→A。這讓兩個檔案保持一致，使用者不需要手動維護兩端。

## Wiki-links：與 relations 分離

**決策**：Wiki-links (`[[type/name-ulid]]`) 儲存在專用的 `wikilinks` 資料表中，與 schema 定義的 relations 分開。

Relations 是結構化的（在 type schema 中定義，有目標型別和基數）。Wiki-links 是自由形式的（寫在 markdown body 中，任何物件都能連結到任何其他物件）。分開儲存避免混淆兩種本質不同的連結機制。

### 使用完整物件 ID，而非顯示名稱

Wiki-link 目標使用包含 ULID 後綴的完整物件 ID（例如 `[[person/bob-01kk3gqm8zrrbjjwkx90f727y6]]`）。這讓目標解析簡化為直接查找——沒有歧義，不需要模糊匹配。

**取捨**：語法較冗長。用易用性換取實作的簡單和正確性。未來的自動完成功能可以減少摩擦，而不需改變底層格式。

## System property registry

**決策**：System properties（`name`、`description`、`created_at`、`updated_at`）定義在 `core/system_property.go` 的 `[]SystemProperty` slice 中。`IsSystemProperty()` 和 `SystemPropertyNames()` 等 helper function 從這個 registry 衍生所有行為。

在 registry 之前，`name` 是唯一的 system property，其處理邏輯散佈在多個檔案中作為特殊案例。新增 `created_at` 和 `updated_at`（未來可能還有更多）會讓這些特殊案例倍增。

### 為什麼用 slice 而非 map

Slice 保留插入順序，這對 frontmatter 輸出排序很重要（`name` → `description` → `created_at` → `updated_at` → schema properties）。Map 需要額外的排序機制。

### 為什麼不用 callback 處理自動設值

每個 system property 的自動設值行為不同（`name` 來自 slug 或 template，`created_at` 來自 `time.Now()`，`updated_at` 每次儲存時更新，`description` 由使用者撰寫）。將這些編碼為 registry 中的 callback 增加了抽象但沒有降低複雜度。改由 `NewObject` 和 `SaveObject` 直接處理設值邏輯，registry 負責識別和驗證。

### RFC 3339 搭配本地時區

時間戳記使用 `time.Now().Format(time.RFC3339)`，包含本地時區偏移。曾考慮 UTC 但被否決——在 local-first 工具中，人類可讀的本地時間更實用。字串原樣儲存在 YAML frontmatter 中，避免 `time.Time` 值的往返格式化問題。

## Name property：放在 properties map 中

**決策**：`name` 儲存在 `Object.Properties["name"]` 中，而非作為專用的 struct 欄位。

將 `name` 放在 properties map 中意味著它自然地流經現有的 frontmatter 讀寫路徑——不需要修改 SQLite `properties` JSON 欄位的 schema，不需要在 frontmatter 解析中做特殊處理。`GetName()` 從 map 讀取，回退到 `DisplayName()`（從檔案名稱衍生）以保持向後相容。

**被否決的替代方案**：專用的 `Object.Name` struct 欄位需要平行儲存、frontmatter 解析的特殊處理，以及資料庫 schema 變更。

### 透過 sync 遷移

沒有 `name` property 的現有物件會在 `Vault.Sync()` 期間使用從檔案名稱衍生的顯示名稱回填。這搭載在現有的 sync 機制上——使用者零額外步驟。

## Name templates：放在 properties 陣列中

**決策**：Name template 定義為 type schema `properties` 陣列中的 `- name: name` 項目，僅允許 `template` 欄位。

```yaml
properties:
  - name: name
    template: "Journal {{ date:YYYY-MM-DD }}"
  - name: content
    type: string
```

`properties` 陣列是 property 設定的自然位置。曾考慮頂層的 `name_template` 欄位，但被否決——這會建立將 property 設定散佈在 schema 各處的先例。

### 僅在建立時求值

`NewObject()` 在 name 引數為空時求值 template，將結果作為靜態字串寫入 `name` property。Template 字串不儲存在物件中。這簡單、可預測，且允許使用者事後編輯名稱。

### 人性化的日期格式

Template 使用 `{{ date:YYYY-MM-DD }}` 語法，而非 Go 的參考時間格式（`2006-01-02`）。Template 引擎在內部轉換 token：

| Token | Go 對應 | 範例 |
|-------|---------|------|
| YYYY  | 2006    | 2026 |
| MM    | 01      | 03   |
| DD    | 02      | 14   |
| HH    | 15      | 09   |
| mm    | 04      | 30   |
| ss    | 05      | 00   |

## Property type system：明確型別而非自動偵測

**決策**：`date` 和 `datetime` 是獨立的 property 型別，而非單一自動偵測的型別。

`date` property 始終儲存 `YYYY-MM-DD`；`datetime` 始終儲存 `YYYY-MM-DDTHH:MM:SS`。使用者在 schema 中宣告意圖。ISO 8601 字串在 SQLite 中天然可排序。

**被否決的替代方案**：接受兩種格式的單一 `date` 型別——會讓驗證和顯示不一致。

### Options 使用物件陣列

`select` 和 `multi_select` 型別使用 `options: [{value: x, label: X}]` 而非舊的 `values: [x]` 格式。`label` 欄位讓顯示名稱可以和儲存值不同（例如 `value: in-progress`, `label: In Progress`）。

**被否決的替代方案**：平行的 `values` 和 `labels` 陣列——容易出錯的耦合。

### 延遲 typed SQLite 儲存

Properties 作為 JSON blob 儲存在 `objects.properties` 中。曾考慮 typed `object_properties` 資料表但延遲實作——只有搭配能利用它的查詢語法（如 `rating>4`、日期範圍）時才有價值。現在建構會迫使在沒有消費者的情況下做出關於 multi-select 儲存的過早決策。

## Shared properties：單層 `use` 參照

**決策**：Shared properties 放在 `.typemd/properties.yaml`。Type schema 透過 `use: <name>` 參照它們，僅允許 `pin` 和 `emoji` 作為覆寫。

```yaml
# .typemd/properties.yaml
properties:
  - name: due_date
    type: date
    emoji: 📅

# .typemd/types/task/schema.yaml
properties:
  - use: due_date
    pin: 1
```

解析發生在 `LoadType()` 中——解析後，每個 `use` 項目被替換為從 shared definition 複製並套用覆寫的完整解析 `Property`。下游程式碼永遠看不到 `use` 項目。

**被否決的替代方案**：`ref` 關鍵字——避免與 JSON Schema `$ref` 混淆。不支援繼承或多層組合——`use` 是單層查找，不可能遞迴。

## TUI session state：物件 ID 而非游標索引

**決策**：TUI 持久化 `selectedObjectID`（例如 `book/clean-code-01jqr...`）而非游標索引。

物件 ID 在跨 session 間是穩定的，即使物件被新增或刪除。游標索引 3 在變更後可能指向完全不同的物件。

### Fallback 策略

當儲存的物件不再存在時，TUI 回退到同一 type group 的第一個物件，再回退到整體第一個物件。這讓使用者留在 vault 的同一「鄰近區域」。

### 僅在離開時儲存

State 僅在使用者離開時寫入 `.typemd/tui-state.yaml`——不會持續寫入。崩潰會遺失 state，但對便利功能來說是可接受的。更簡單的方式避免了檔案系統的額外開銷。

### 載入錯誤時靜默失敗

如果 state 檔案遺失、損壞或包含無效資料，TUI 靜默回退到預設的啟動行為。State 持久化是便利功能，不是關鍵功能——使用者不應該被損壞的 state 檔案阻擋。
