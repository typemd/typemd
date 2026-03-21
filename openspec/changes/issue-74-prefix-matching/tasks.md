## 1. Core: SyncResult and syncContext extensions

- [x] 1.1 Write BDD scenarios for SyncResult prefix reporting (expanded count, unresolved list)
- [x] 1.2 Implement step definitions for SyncResult scenarios
- [x] 1.3 Extend `SyncResult` struct with `Expanded int` and `Unresolved []UnresolvedRelation` fields
- [x] 1.4 Extend `syncContext` with `diskObjects map[string]*Object`, schema cache, and `nameIndex map[string]map[string][]string`

## 2. Core: Shared name-to-ID resolution

- [x] 2.1 Write BDD scenarios for name resolution (unique match, no match, ambiguous, full ID passthrough, slugified name lookup)
- [x] 2.2 Implement step definitions for name resolution scenarios
- [x] 2.3 Add `buildNameIndex()` that builds per-type name-to-IDs mapping from walked objects (using both `name` property and slug)
- [x] 2.4 Add `resolveByName(nameIndex, typeName, name)` shared helper that returns (fullID, error)
- [x] 2.5 Add `resolveRelationValue()` that detects ULID suffix presence and calls `resolveByName` for non-ULID values
- [x] 2.6 Add unit tests for edge cases (empty value, type-only reference, duplicate names, slug vs name matching)

## 3. Core: Relation sync in Projector

- [x] 3.1 Write BDD scenarios for relation sync (single-value, multi-value, non-existent target skipped, non-relation property ignored)
- [x] 3.2 Implement step definitions for relation sync scenarios
- [x] 3.3 Add `ObjectIndex.DeleteNonTagRelations()` method to clear non-tag relations during full sync
- [x] 3.4 Add `ObjectIndex.DeleteRelationsByObject(objectID)` method for incremental sync
- [x] 3.5 Implement `syncRelations()` on Projector: iterate objects, read schema, resolve relation values, insert into relations table
- [x] 3.6 Add unit tests for `syncRelations()` with clear-and-rebuild behavior

## 4. Core: Auto-expand write-back

- [x] 4.1 Write BDD scenarios for auto-expand write-back (file updated with full ID, multiple properties expanded, unresolvable left unchanged)
- [x] 4.2 Implement step definitions for write-back scenarios
- [x] 4.3 Add write-back logic in `syncRelations()`: collect modified objects, batch-save after all relations processed
- [x] 4.4 Add unit tests for write-back with mixed resolved/unresolved values in same object

## 5. Core: Integrate into Projector.Sync and SyncFiles

- [x] 5.1 Write BDD scenarios for full sync including relation sync phase
- [x] 5.2 Implement step definitions for full sync scenarios
- [x] 5.3 Wire `syncRelations()` into `Projector.Sync()` between upsertObject loop and syncWikiLinksAndTags
- [x] 5.4 Wire relation sync into `Projector.SyncFiles()` for incremental sync (delete-and-rebuild for changed objects only)
- [x] 5.5 Add integration test: full sync round-trip (write name reference → sync → verify file expanded + relation indexed)

## 6. CLI: Validate integration

- [x] 6.1 Write BDD scenarios for `tmd validate` reporting unresolvable relation references
- [x] 6.2 Implement step definitions for validate scenarios
- [x] 6.3 Update `ValidateRelations` to detect and report name references that were not resolved
