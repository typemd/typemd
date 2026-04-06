---
title: 領域事件
description: TypeMD 核心程式庫所有領域事件的完整參考。
sidebar:
  order: 5
---

TypeMD 使用領域事件來解耦生產者（實體、服務）與消費者（索引寫入器、TUI、MCP server）。實體方法回傳 `DomainEvent` 值；用例層在操作成功後收集並派發它們。

## 訂閱事件

```go
vault.Events.Subscribe(func(e core.DomainEvent) {
    switch e := e.(type) {
    case core.ObjectCreated:
        // 處理新物件
    case core.PropertyChanged:
        // 處理屬性更新
    case core.TypeSaved:
        // 處理 type schema 變更
    }
})
```

事件按照收集順序同步派發。Handler 應保持輕量 — 較重的工作應排入非同步佇列處理。

## ObjectService 事件

由 `ObjectService` 在使用者發起的操作中產出（建立、儲存、連結、解除連結）。

| 事件 | 用途 | Payload |
|------|------|---------|
| `ObjectCreated` | 新物件被建立 | `Object` |
| `ObjectSaved` | 既有物件被儲存 | `Object` |
| `PropertyChanged` | 單一屬性值變更 | `ObjectID`、`Key`、`Old`、`New` |
| `ObjectLinked` | 兩個物件之間建立 relation | `FromID`、`ToID`、`RelName` |
| `ObjectUnlinked` | 兩個物件之間移除 relation | `FromID`、`ToID`、`RelName` |
| `TagAutoCreated` | 同步期間自動建立 tag 物件 | `Tag`、`ReferencedBy` |

## Vault type 事件

由 `Vault.SaveType()` 和 `Vault.DeleteType()` 在 type schema 操作時產出。

| 事件 | 用途 | Payload |
|------|------|---------|
| `TypeSaved` | Type schema 被建立或更新 | `Schema` |
| `TypeDeleted` | Type schema 被刪除 | `Name` |

操作失敗時不會產出事件（例如 `SaveType` 驗證錯誤、嘗試刪除內建型別）。

## Reconciler → Projector 事件

由 Reconciler 在同步期間產出；由 Projector 消費以更新 SQLite 索引。

| 事件 | 用途 | Payload |
|------|------|---------|
| `ObjectUpserted` | 物件需要寫入索引 | `ID`、`Type`、`Filename`、`PropsJSON`、`Body` |
| `ObjectDeleted` | 物件從磁碟移除（stale 清理） | `ID` |
| `RelationsCleared` | 在重建前清除 relations | `ObjectID`（單一物件）、`NonTagOnly`（完整同步）或 `TagsOnly`（tag 同步） |
| `RelationIndexed` | 單一 relation 需要建立索引 | `Name`、`FromID`、`ToID` |
| `WikiLinksSynced` | 物件的 wiki-links 已解析完成 | `ObjectID`、`Links []WikiLinkEntry` |

### 同步模式

- **完整調和**（`Reconciler.Reconcile()`）：產出 `RelationsCleared{NonTagOnly: true}` 事件，接著為所有 relations 產出 `RelationIndexed` 事件。Tag relations 透過 `RelationsCleared{TagsOnly: true}` 獨立管理。
- **增量調和**（`Reconciler.ReconcileFiles()`）：為變更的物件產出 `RelationsCleared{ObjectID: id}` 事件，然後為重建的 relations 產出 `RelationIndexed` 事件。未變更的物件不受影響。

## 設計原則

- **Fire-and-forget**：事件派發失敗不會回滾產出事件的操作。
- **實體產出 → 用例派發**：實體方法回傳 `DomainEvent` 值；服務在操作成功後收集並派發。
- **檔案是唯一真實來源**：事件驅動 SQLite 索引，但索引永遠可以從檔案重建。
