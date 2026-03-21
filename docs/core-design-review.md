# Core Package Design Review

> Reviewed: 2026-03-21

## Executive Summary

The core package follows Clean Architecture with CQRS well. The overall structure is solid — clear separation between ObjectService (commands) and QueryService (queries), domain entities are well-defined, and the file-first architecture is consistently applied. However, there are several areas that could benefit from refinement.

---

## Critical Issues

### 1. ObjectRepository 做了太多 Business Logic

**位置**: `local_object_repository_schema.go` — `GetSchema()`

`GetSchema()` 是一個讀取操作，但它做了三件不該做的事：

1. **抽取 NameTemplate**: 從 `properties` 中找到 `name` property，把 `template` 欄位搬到 `schema.NameTemplate`，然後移除該 property
2. **預設 Version**: 空字串時填入 `DefaultSchemaVersion`
3. **解析 `use` entries**: 呼叫 `resolveSchemaUseEntries()` 展開 shared property 引用

這些是 domain logic，不是 persistence 的職責。

更嚴重的是，`GetSchema()` 還有一個 **隱藏的 side effect**：

```go
if !isDir {
    _ = r.migrateToDirectory(name)  // 讀取操作觸發寫入！
}
```

一個讀取操作不應該修改檔案系統，而且錯誤被靜默忽略。

**建議**: 將 schema 的轉換邏輯移到 service layer（例如新增 `TypeSchemaService`），或至少移到 `Vault.LoadType()` 中處理。Repository 只負責原始的反序列化。

---

### 2. Property Type 驗證邏輯分散

Property type 的定義和驗證散佈在多個檔案中，沒有 single source of truth：

| 位置 | 職責 |
|------|------|
| `type_schema.go` — `validPropertyTypes` | 定義合法的 property type 列表 |
| `type_schema_validate.go` — `validatePropertyType()` | Schema 層級驗證 |
| `type_schema_validate.go` — `ValidateObject()` | Instance 層級驗證 |
| `filter_operator.go` — `OperatorsForType()` | 定義每個 type 的 filter operators |
| `display.go` — `FormatValue()` | 定義每個 type 的顯示格式 |
| `stats.go` — `computePropertyStats()` | 定義每個 type 的統計計算 |

新增一個 property type 需要改 6 個檔案，容易遺漏。

**建議**: 考慮用 registry pattern，將每個 property type 的 validation、operators、formatting、stats 聚合在一起。

---

### 3. Relation Property 驗證不完整

`type_schema_validate.go` 中的 relation property 驗證缺少關鍵檢查：

- **不驗證 `target` 是否指向存在的 type**：schema 中可以寫 `target: nonexistent_type` 而不會報錯
- **不驗證 `inverse` property 是否存在於 target type**：bidirectional relation 設定錯誤時只會在 runtime 出錯
- **`Multiple` 和 `Bidirectional` 的組合沒有一致性驗證**

---

### 4. ViewConfig 不驗證 Schema

`SaveView()` 只驗證 view name，不驗證內容：

- Filter rules 可以引用不存在的 property
- Sort rules 可以引用不存在的 property
- Column names 不檢查是否存在
- Layout enum 不驗證

Schema 演進後 view 會靜默失效。

**建議**: 新增 `ValidateView(typeName string, view *ViewConfig) []error`，在 SaveView 時呼叫。

---

## High Priority Issues

### 5. Vault Facade 太大（66 methods）

Vault 有 66 個方法，其中：
- 12 個是純 delegation 到 ObjectService/QueryService ✓
- 9 個是 path accessors（可以歸入 nested struct）
- `FormatAll/FormatObjects/FormatSchemas` 包含完整的 business logic（直接操作 repo.Walk、os.ReadFile、os.WriteFile）
- `MigrateObjects/MigrateSchemas` 同樣包含完整的 business logic

**建議**:
1. 把 `FormatAll/FormatObjects/FormatSchemas` 移到 `FormattingService`
2. 把 `MigrateObjects/MigrateSchemas` 移到 `MigrationService`
3. 考慮把 path accessors 收斂到 `vault.Paths` struct

### 6. ObjectIndex 介面太大（19 methods），違反 ISP

`ObjectIndex` 混合了 read、write、relation management、wikilink management、lifecycle：

```
Read:      Query, Search, FindRelations, FindBacklinks, ListWikiLinks (5)
Write:     Upsert, Remove, ListIDs, NeedsSync (4)
Relations: InsertRelation, DeleteRelation, DeleteRelationsByName,
           DeleteNonTagRelations, DeleteRelationsByObject,
           CleanOrphanedRelations (6)
WikiLinks: SyncWikiLinks, DeleteWikiLinks (2)
Lifecycle: Rebuild, EnsureSchema (2)
```

`CleanOrphanedRelations()` 尤其問題 — 它同時是 query 又是 mutation，回傳被刪除的 orphaned relations。

**建議**: 即使不拆介面，至少將 `CleanOrphanedRelations` 拆成：
- `FindOrphanedRelations() ([]OrphanedRelation, error)` — 純 query
- `CleanOrphanedRelations() error` — 純 mutation

### 7. TypeSchema 方法散佈 5 個檔案

TypeSchema 相關邏輯分佈在：

| 檔案 | 職責 |
|------|------|
| `type_schema.go` | Domain struct + helpers + Vault facade |
| `type_schema_validate.go` | 驗證 |
| `type_schema_marshal.go` | 序列化 + version + color |
| `shared_properties.go` | `use` entry 解析 |
| `system_property.go` | System property registry |

沒有統一的 `TypeSchemaService`。Vault 直接代理到 repository 做 CRUD（SaveType、DeleteType），繞過 service layer。

### 8. `name` Property 的雙重身份

`name` 同時是：
1. System property（`system_property.go` 中定義）
2. Schema property（`ValidateSchema` 允許 schema 中出現 `name`，但只能有 `template` 欄位）

`GetSchema()` 把 schema 中的 `name` property 抽出來變成 `schema.NameTemplate`，然後把 property 移除。這個轉換邏輯隱藏在 repository 中，不明顯。

---

## Medium Priority Issues

### 9. Error Handling 不一致

**deprecated os.IsNotExist vs modern errors.Is**:

| 使用 deprecated 形式 | 使用 modern 形式 |
|---|---|
| `Create()` — `os.IsExist(err)` | `GetSharedProperties()` — `errors.Is` |
| `Walk()` — `os.IsNotExist(err)` | `DeleteSchema()` — `errors.Is` |
| `GetTemplate()` — `os.IsNotExist(err)` | `DeleteTemplate()` — `errors.Is` |

**建議**: 統一使用 `errors.Is(err, os.ErrNotExist)`。

### 10. Validator 回傳型別不一致

| Validator | 回傳型別 |
|-----------|---------|
| `ValidateAllSchemas` | `map[string][]error` |
| `ValidateAllObjects` | `[]error` |
| `ValidateRelations` | `[]error` |

Doctor 因此需要兩種不同的轉換函式（`errorsToCategory` 和 `mapErrorsToCategory`）。

### 11. Silent Skip Pattern

多處在遇到問題時靜默跳過，不回報：

| 位置 | 行為 |
|------|------|
| Projector `upsertObject` | JSON marshal 失敗時 `return nil` |
| Projector `syncTagRelations` | Unresolvable tag `continue` |
| SQLiteObjectIndex `Query` | Unsafe sort property 被靜默忽略 |
| `GetSchema()` | Directory migration 錯誤 `_ = ...` |

**建議**: 至少收集到 SyncResult 或回傳 warning list。

### 12. Concurrency Safety 未文件化

ObjectService 沒有任何 locking：
- `Create` 的 uniqueness check 和實際建立之間有 race condition
- `Link` 的 bidirectional 操作涉及兩次 Save，可能互相覆蓋

目前是 safe 的（CLI/TUI/MCP 都是 single-threaded），但應在 Vault 的 godoc 中明確記錄此假設。

### 13. Property.Default 從未被驗證

Property 結構有 `Default` 欄位，但：
- 不驗證 default value 是否符合 property type
- 不驗證 select property 的 default 是否在 options 中

### 14. `use` Override 的 Pin 零值問題

`shared_properties.go` 的 `resolveUseEntries()`：

```go
if resolved.Pin == 0 {
    resolved.Pin = shared.Pin
}
```

`pin: 0` 被當作「未設定」，但 0 也可能是有意義的值。應改用 `*int`（pointer）來區分零值與未設定。

---

## Low Priority Issues

### 15. ID 格式化有三種方式

- `display.go` — `displayObjectID()` 使用 `StripULID`
- `object_id.go` — `ParseObjectID().DisplayID()`
- Projector 中直接呼叫 `StripULID()`

應統一使用一種方式。

### 16. `BacklinksDisplayKey` Magic String

`display.go` 中硬編碼 `"backlinks"` 作為 display key。如果某個 type 有一個名為 `backlinks` 的 property，會產生衝突。

### 17. Stats 中的日期比較用字串

`stats.go` 中 `DateStats` 的 min/max 用字串的 lexicographic comparison：

```go
if earliest == "" || s < earliest { earliest = s }
```

只在日期格式為 `YYYY-MM-DD` 時正確。應 parse 為 `time.Time` 再比較。

### 18. `queryWikiLinks` 的 Column 參數未驗證

`sqlite_object_index.go` 中 `queryWikiLinks(column, objectID)` 將 `column` 直接嵌入 SQL。目前是 safe 的（只有 hardcoded 的 `"from_id"` 和 `"to_id"`），但沒有防禦性驗證。

### 19. PropertyChanged Event 定義了但未使用

`domain_event.go` 定義了 `PropertyChanged`，但 `ObjectService.SetProperty()` 只 dispatch `ObjectSaved`，沒有 dispatch `PropertyChanged`。

### 20. SyncResult 資訊不足

`SyncResult` 只回報 expanded 和 unresolved relations，不回報 added/updated/deleted 的數量。消費者無法區分 full rebuild 和 incremental update。

---

## Design Strengths（值得保留的設計）

1. **Clean CQRS 分離**: ObjectService (write) / QueryService (read) 界線清晰
2. **DI-friendly**: 所有 service 透過 constructor injection，無 global state
3. **Domain Event Pattern**: entity produces → service dispatches，一致且乾淨
4. **File-first Architecture**: 檔案是 source of truth，SQLite 是加速層
5. **Incremental Sync**: `SyncFiles()` 支援 file watcher 的增量同步
6. **Schema Caching**: Projector per-sync cache 避免重複載入
7. **Interface-based Dependencies**: ObjectRepository 和 ObjectIndex 都是 interface
8. **Vault as Thin Facade**: 大部分方法是純 delegation（除了 format/migrate）
