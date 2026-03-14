## Context

The `core/` package centers on a `Vault` struct with 50+ methods that directly calls `os.*` for file I/O (22 call sites) and `database/sql` for SQLite operations (20 call sites). This tight coupling prevents alternative storage backends needed for the Web UI (`try.typemd.io` via GitHub API, Wails desktop via Go bindings). The current architecture works but doesn't scale to multi-platform delivery.

The codebase already has a natural CQRS shape: commands (NewObject, SaveObject, LinkObjects) write to both files and SQLite, while most queries (QueryObjects, SearchObjects) read from SQLite. This refactoring formalizes that separation.

## Goals / Non-Goals

**Goals:**
- Extract `ObjectRepository` interface for entity persistence by known ID (file-backed for local)
- Extract `ObjectIndex` interface for search/query/discovery (SQLite-backed for local)
- Extract `Projector` for file→index synchronization (currently `SyncIndex`)
- Vault becomes a Use Case coordinator depending only on interfaces
- Interfaces designed for Web UI compatibility (GitHub API + in-memory index)
- All existing tests pass at every step — zero behavioral changes

**Non-Goals:**
- Implementing GitHubObjectRepository or InMemoryObjectIndex (future, when Web UI is built)
- Splitting Vault into separate CommandService/QueryService structs
- Changing Vault's public method signatures (consumers are unaffected)
- Refactoring domain logic (validation, frontmatter parsing) — stays in core
- Async event-driven CQRS — dual writes remain synchronous

## Decisions

### 1. Repository returns domain entities, not raw bytes

Repository methods return `*Object`, `*TypeSchema`, `*Template`, etc. Serialization (frontmatter parsing, YAML decoding) is the repository implementation's responsibility.

**Rationale:** Follows DDD — Repository is the boundary between domain and infrastructure. The domain layer works with entities, never raw storage formats. Shared parsing utilities (`parseFrontmatter`, `writeFrontmatter`) become internal helpers used by repository implementations.

**Alternative considered:** Returning `[]byte` and parsing in the Use Case layer. Rejected because it leaks storage format concerns upward and forces every consumer to handle parsing.

### 2. ObjectRepository scope includes schemas, templates, and shared properties

A single `ObjectRepository` interface covers all file-backed entities: objects, type schemas, templates, and shared properties. These are all "things with a file-based ID that are the source of truth on disk."

**Rationale:** All file-backed entities share the same storage pattern (read by ID from filesystem). A single interface keeps the implementation count manageable. The GitHub API implementation would back all of these from the same API surface.

**Alternative considered:** Separate interfaces per entity type (`ObjectStore`, `SchemaStore`, `TemplateStore`). Rejected because it creates interface proliferation without meaningful behavioral differences — they all read/write files by ID.

### 3. ObjectIndex returns ObjectResult, not full entities

Index queries return `ObjectResult` (ID, Type, Filename, Properties) — a lightweight projection, not the full `*Object` with body. When the full entity is needed, the caller uses `ObjectRepository.Get(id)`.

**Rationale:** The index is a read model optimized for search and listing. Returning full entities would either require storing the body in SQLite (wasteful duplication) or fetching from disk (defeating the purpose of the index). The two-step pattern (search → fetch) matches the CQRS read model concept.

### 4. Projector write methods live on the ObjectIndex interface

Methods like `Upsert`, `Remove`, `SyncRelations`, `SyncWikiLinks`, and `Rebuild` are on `ObjectIndex` rather than a separate `ProjectorStore` interface.

**Rationale:** These methods are the write side of the read model — they're how the index gets populated. The `Projector` component calls them, but they're defined on `ObjectIndex` because the SQLite implementation owns the schema and write logic. Separating them would create a 1:1 interface-to-implementation mapping with no practical benefit.

### 5. Incremental extraction order: Index first, then Repository, then Projector

Phase 1: Extract SQLite operations → `SQLiteObjectIndex`
Phase 2: Extract file I/O → `LocalObjectRepository`
Phase 3: Extract SyncIndex → `Projector`
Phase 4: Clean up Vault to depend only on interfaces

**Rationale:** Index extraction is the safest first step because SQLite queries are already well-isolated in query.go, wikilink.go, and relation.go (read side). File I/O is more deeply intertwined with domain logic, making it a harder second step. The Projector naturally falls out once both sides are extracted.

### 6. Vault retains domain logic and orchestration

Validation (`ValidateObject`), frontmatter formatting (`OrderedPropKeys`), name generation (`EvaluateNameTemplate`, `GenerateULID`), and ULID handling remain in the `core` package as domain utilities. Vault orchestrates these with the repositories.

**Rationale:** These are domain concerns, not infrastructure. Moving them into repositories would violate DDD by putting business rules in the infrastructure layer.

## Risks / Trade-offs

- **[Risk] Large refactoring surface** → Mitigated by incremental phases. Each phase is a separate PR that maintains all tests. If any phase causes issues, it can be reverted independently.

- **[Risk] Dual-write consistency** → Commands write to both ObjectRepository and ObjectIndex synchronously. If one fails, state is inconsistent. → Mitigated by writing to ObjectRepository first (source of truth), then ObjectIndex. SyncIndex/Projector can reconcile on next run. This matches current behavior.

- **[Risk] Performance regression** → Adding interface indirection and entity construction in repository methods could slow hot paths. → Mitigated by the fact that file I/O and SQLite queries already dominate latency. One extra struct allocation per operation is negligible.

- **[Trade-off] ObjectRepository is a large interface** → It covers objects, schemas, templates, and shared properties. This is intentional — splitting would create interface proliferation without benefit. If the interface grows unwieldy during implementation, it can be decomposed later.
