## Context

The Projector currently handles two types of cross-object references during sync: wikilinks (in markdown body) and tag relations (the `tags` system property). General relation properties defined in type schemas (e.g., `author`, `books`) are only indexed when created through `ObjectService.Link()`. Hand-edited frontmatter relations are stored in the properties JSON but never reach the `relations` SQLite table.

Additionally, all relation values must be full ULID-suffixed object IDs. Users editing frontmatter manually must know and type the complete ID, which is tedious and error-prone.

## Goals / Non-Goals

**Goals:**
- Resolve relation value prefixes (e.g., `person/john-doe`) to full IDs during Projector sync
- Auto-expand resolved prefixes by writing the full ID back to the object file
- Sync all schema-defined relation properties to the SQLite `relations` table during sync
- Report resolution results (expanded, ambiguous, not found) through `SyncResult`

**Non-Goals:**
- Runtime prefix resolution (keeping short prefixes in files) — files must always converge to full IDs
- Changing `ObjectService.Link()`/`Unlink()` behavior — they already work with full IDs
- Interactive disambiguation during sync — that belongs to CLI commands (#147)

## Decisions

### Decision 1: Resolve-then-write-back in Projector, not in ObjectService

**Choice**: Add prefix resolution and relation indexing to `Projector.Sync()` and `SyncFiles()`.

**Rationale**: The Projector is already responsible for normalizing object state during sync (e.g., adding missing `name` property). Prefix resolution is another normalization step. `ObjectService` operates on fully resolved IDs and should not need to handle prefixes.

**Alternative considered**: Adding resolution to `ObjectService.Save()`. Rejected because Save is used by many paths (link, unlink, template apply) where IDs are already resolved.

### Decision 2: Process relations after upsert, before wikilinks/tags

**Choice**: Add a new `syncRelations()` phase between `upsertObject` and `syncWikiLinksAndTags`.

**Flow**:
1. Walk all objects → upsert each (existing)
2. **syncRelations** → resolve prefixes, write-back, index relations (new)
3. syncWikiLinksAndTags (existing)

**Rationale**: Relations need the full object set available for prefix resolution (like tags need it for name lookup). Processing after upsert ensures all objects are in the index. Processing before wikilinks/tags maintains the existing flow.

### Decision 3: Clear-and-rebuild for relation sync (like tags)

**Choice**: During full sync, delete all non-tag relations from the `relations` table, then rebuild from frontmatter. During incremental sync, delete relations for changed objects only, then rebuild those.

**Rationale**: This matches the existing `syncTagRelations` pattern (clear-then-rebuild). It's simpler than diffing and handles property renames, removed relations, and manual edits correctly.

**Alternative considered**: Incremental diff-based sync. Rejected as overly complex for the current scale and the fact that tags already use clear-and-rebuild.

### Decision 4: Shared name-to-ID resolution method

**Choice**: Create a shared `ResolveByName(typeName, name string) (string, error)` method that resolves a human-readable name (e.g., `john-doe`) to a full object ID (e.g., `person/john-doe-01abc...`). This is used by both relation prefix resolution and wiki-link resolution.

**Resolution strategy**: During sync, the Projector has all objects in memory. Build a name index (`map[type+name]objectID`) from the walked objects. For a relation value like `person/john-doe`, strip the type prefix, look up `john-doe` in the name index for type `person`. This is O(1) lookup, much faster than filesystem glob.

**Rationale**: The name-to-ID resolution is a general capability needed by:
1. Relation prefix resolution (this change)
2. Wiki-link shorthand resolution (#176, future)
3. Tag resolution (already exists as `resolveTagReference`, can be unified)

Designing it as a shared method now avoids duplication later. The existing `resolveTagReference` is a specialized version of this pattern — it uses `tagNameIndex` for name-based lookup. The new method generalizes this to all types.

**Alternative considered**: Reusing `GlobIDs` (filesystem glob). Rejected because name resolution is semantically about the `name` property, not filename prefixes. A name index built from walked objects is both faster and more correct (handles slugified names properly).

### Decision 5: Collect write-backs and batch-save after all relations are processed

**Choice**: During relation sync, collect all objects that need prefix expansion. After processing all objects, save them in a batch.

**Rationale**: An object may have multiple relation properties, some needing expansion. Batching avoids multiple writes per object. It also prevents partial writes if a later resolution fails.

### Decision 6: Add syncContext fields for relation tracking

**Choice**: Extend `syncContext` with:
- `diskObjects map[string]*Object` — all objects by ID (for write-back)
- `diskSchemas map[string]*TypeSchema` — cached schemas (reuse from objectSyncer)
- `nameIndex map[string]map[string]string` — per-type name-to-ID index (`nameIndex[type][name] = objectID`)

**Rationale**: Relation sync needs access to both the object (for frontmatter values) and its schema (for relation property definitions). The name index enables O(1) name-to-ID resolution during sync. These are already loaded during upsert and can be cached.

## Risks / Trade-offs

- **Write-back during sync modifies files** → The TUI file watcher could re-trigger sync. Mitigation: The watcher is single-shot (delivers one message then stops), so a write-back during sync won't cause re-entry. The next watcher cycle will see the already-expanded file and do nothing.
- **Ambiguous prefixes silently skip** → Users may not notice unresolved references. Mitigation: Report in `SyncResult` and surface via `tmd validate`.
- **Performance of GlobIDs per relation value** → For large vaults with many relations, filesystem glob per value could be slow. Mitigation: Most relation values will already be full IDs (written by `Link()`). Only hand-edited prefixes trigger glob. Future optimization: batch resolve via index query.
- **Bidirectional relation sync complexity** → When syncing from frontmatter, should inverse relations be auto-created? Decision: No. The Projector only indexes what's in the file. If user writes `author: person/X` on a book, the Projector creates the forward relation. The inverse (`books` on person) should only exist if the person's file also has it. `Link()` handles bidirectional auto-creation; sync respects file contents literally.
