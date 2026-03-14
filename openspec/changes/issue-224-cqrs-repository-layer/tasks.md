## 1. Define interfaces

- [x] 1.1 Define `ObjectIndex` interface in `core/object_index.go` with query methods (Query, Search, FindRelations, FindBacklinks, ListWikiLinks) and projector write methods (Upsert, Remove, ListIDs, InsertRelation, DeleteRelation, DeleteRelationsByName, CleanOrphanedRelations, SyncWikiLinks, DeleteWikiLinks, Rebuild, EnsureSchema)
- [x] 1.2 Define `ObjectResult` struct (ID, Type, Filename, Properties) as the lightweight query return type
- [x] 1.3 Define `ObjectRepository` interface in `core/object_repository.go` with entity methods (Get, Save, Create, Walk, GlobIDs, ModTime, EnsureDir) and schema/template methods (GetSchema, WriteSchema, ListSchemas, GetTemplate, ListTemplates, GetSharedProperties)

## 2. Extract SQLiteObjectIndex

- [x] 2.1 Create `core/sqlite_object_index.go` struct with `*sql.DB` field and constructor
- [x] 2.2 Move `ensureSchema()` SQL DDL from `vault.go` → `SQLiteObjectIndex.EnsureSchema()`
- [x] 2.3 Move `QueryObjects()` and `SearchObjects()` from `query.go` → `SQLiteObjectIndex.Query()` and `Search()`, returning `[]*ObjectResult` instead of `[]*Object`
- [x] 2.4 Move `RebuildIndex()` from `query.go` → `SQLiteObjectIndex.Rebuild()`
- [x] 2.5 Move `ListRelations()` from `relation.go` → `SQLiteObjectIndex.FindRelations()`
- [x] 2.6 Move relation write operations (INSERT/DELETE) from `relation.go` → `SQLiteObjectIndex.InsertRelation()`, `DeleteRelation()`, `DeleteRelationsByName()`
- [x] 2.7 Move `cleanOrphanedRelations()` from `sync.go` → `SQLiteObjectIndex.CleanOrphanedRelations()`
- [x] 2.8 Move `syncWikiLinks()`, `ListWikiLinks()`, `ListBacklinks()` from `wikilink.go` → `SQLiteObjectIndex` methods
- [x] 2.9 Move object UPSERT/DELETE SQL from `sync.go` and `object.go` → `SQLiteObjectIndex.Upsert()`, `Remove()`, `ListIDs()`
- [x] 2.10 Write unit tests for `SQLiteObjectIndex` — verify all query and write methods
- [x] 2.11 Update `Vault` to hold `ObjectIndex` interface, wire `SQLiteObjectIndex` in `Open()`, verify all existing tests pass

## 3. Extract LocalObjectRepository

- [x] 3.1 Create `core/local_object_repository.go` struct with root path field and constructor, encapsulating path conventions (ObjectsDir, TypesDir, ObjectPath, etc.)
- [x] 3.2 Move `GetObject()` file read + parse from `object.go` → `LocalObjectRepository.Get()`, returning `*Object`
- [x] 3.3 Move `saveObjectFile()` file write from `object.go` → `LocalObjectRepository.Save()`
- [x] 3.4 Move object file creation (O_EXCL write) from `NewObject()` → `LocalObjectRepository.Create()`
- [x] 3.5 Move `filepath.Walk` logic from `walkAndUpsertObjects()` → `LocalObjectRepository.Walk()`, returning `[]*Object`
- [x] 3.6 Move `ResolveID()` glob logic → `LocalObjectRepository.GlobIDs()`
- [x] 3.7 Move `os.Stat` for mtime → `LocalObjectRepository.ModTime()`
- [x] 3.8 Move `LoadType()` from `type_schema.go` → `LocalObjectRepository.GetSchema()`, move `ListTypes()` from `list.go` → `ListSchemas()`
- [x] 3.9 Move `LoadTemplate()` and `ListTemplates()` from `template.go` → `LocalObjectRepository.GetTemplate()` and `ListTemplates()`
- [x] 3.10 Move `LoadSharedProperties()` from `shared_properties.go` → `LocalObjectRepository.GetSharedProperties()`
- [x] 3.11 Move `MigrateSchemas()` schema file write → `LocalObjectRepository.WriteSchema()`
- [x] 3.12 Move `os.MkdirAll` for object dirs → `LocalObjectRepository.EnsureDir()`
- [x] 3.13 Write unit tests for `LocalObjectRepository` — verify all entity methods
- [x] 3.14 Update `Vault` to hold `ObjectRepository` interface, wire `LocalObjectRepository` in `Open()`, verify all existing tests pass

## 4. Extract Projector

- [x] 4.1 Create `core/projector.go` with `Projector` struct holding `ObjectRepository` + `ObjectIndex`
- [x] 4.2 Move `SyncIndex()` orchestration logic → `Projector.Sync()`: walk repo, upsert index, delete stale, clean orphans, sync tags, sync wikilinks, rebuild FTS
- [x] 4.3 Move `syncTagRelations()` and `resolveOrCreateTag()` → `Projector` methods
- [x] 4.4 Move name migration logic (add NameProperty if missing) → `Projector` internal step
- [x] 4.5 Write unit tests for `Projector` — verify sync orchestration (covered by existing SyncIndex BDD + unit tests)
- [x] 4.6 Update `Vault.SyncIndex()` to delegate to `Projector.Sync()`, verify all existing tests pass

## 5. Vault cleanup

- [x] 5.1 Remove all direct `os.*` calls from Vault entity methods — only lifecycle methods (Init, Open, IsInitialized) retain `os` usage
- [x] 5.2 Remove all direct `database/sql` usage from Vault entity methods — only lifecycle methods (Open, Close) retain `sql` usage for DB connection management
- [x] 5.3 Path helper methods remain on Vault for backward compatibility (70+ external call sites in tests, tui, cmd, mcp) — removal deferred to separate follow-up
- [x] 5.4 Run full test suite (`go test ./...`) — all existing BDD scenarios and unit tests pass
- [x] 5.5 Run `go vet ./...` — no warnings
