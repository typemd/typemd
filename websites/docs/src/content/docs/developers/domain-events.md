---
title: Domain Events
description: Complete reference for all domain events emitted by TypeMD's core library.
sidebar:
  order: 5
---

TypeMD uses domain events to decouple producers (entities, services) from consumers (index writer, TUI, MCP server). Entity methods return `DomainEvent` values; the use case layer collects and dispatches them after successful operations.

## Subscribing to events

```go
vault.Events.Subscribe(func(e core.DomainEvent) {
    switch e := e.(type) {
    case core.ObjectCreated:
        // handle new object
    case core.PropertyChanged:
        // handle property update
    case core.TypeSaved:
        // handle type schema change
    }
})
```

Events are dispatched synchronously in the order they were collected. Handlers should be lightweight — heavy work should be queued for async processing.

## ObjectService events

Emitted by `ObjectService` during user-initiated operations (create, save, link, unlink).

| Event | Purpose | Payload |
|-------|---------|---------|
| `ObjectCreated` | New object created | `Object` |
| `ObjectSaved` | Existing object saved | `Object` |
| `PropertyChanged` | Single property value changed | `ObjectID`, `Key`, `Old`, `New` |
| `ObjectLinked` | Relation created between objects | `FromID`, `ToID`, `RelName` |
| `ObjectUnlinked` | Relation removed between objects | `FromID`, `ToID`, `RelName` |
| `TagAutoCreated` | Tag object auto-created during sync | `Tag`, `ReferencedBy` |

## Vault type events

Emitted by `Vault.SaveType()` and `Vault.DeleteType()` during type schema operations.

| Event | Purpose | Payload |
|-------|---------|---------|
| `TypeSaved` | Type schema created or updated | `Schema` |
| `TypeDeleted` | Type schema deleted | `Name` |

Events are not emitted when the operation fails (e.g., validation error on `SaveType`, attempting to delete a built-in type).

## Reconciler → Projector events

Emitted by the Reconciler during sync; consumed by the Projector to update the SQLite index.

| Event | Purpose | Payload |
|-------|---------|---------|
| `ObjectUpserted` | Object needs to be written to index | `ID`, `Type`, `Filename`, `PropsJSON`, `Body` |
| `ObjectDeleted` | Object removed from disk (stale cleanup) | `ID` |
| `RelationsCleared` | Clear relations before rebuilding | `ObjectID` (per-object), `NonTagOnly` (full sync), or `TagsOnly` (tag sync) |
| `RelationIndexed` | A single relation needs to be indexed | `Name`, `FromID`, `ToID` |
| `WikiLinksSynced` | Wiki-links resolved for an object | `ObjectID`, `Links []WikiLinkEntry` |

### Sync modes

- **Full reconciliation** (`Reconciler.Reconcile()`): Emits a `RelationsCleared{NonTagOnly: true}` event followed by `RelationIndexed` events for all relations. Tag relations are managed separately via `RelationsCleared{TagsOnly: true}`.
- **Incremental reconciliation** (`Reconciler.ReconcileFiles()`): Emits `RelationsCleared{ObjectID: id}` events for changed objects, then `RelationIndexed` events for their rebuilt relations. Unchanged objects are not affected.

## Design principles

- **Fire-and-forget**: Event dispatch failure does not rollback the operation that produced the event.
- **Entity produces → use case dispatches**: Entity methods return `DomainEvent` values; services collect and dispatch after successful operations.
- **Files are the source of truth**: Events drive the SQLite index, but the index can always be rebuilt from files.
