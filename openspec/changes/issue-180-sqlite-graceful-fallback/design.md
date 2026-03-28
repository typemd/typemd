## Context

typemd's query system follows CQRS: `QueryService` coordinates an `ObjectRepository` (filesystem, source of truth) and an `ObjectIndex` (SQLite, acceleration layer). Currently, when `ObjectIndex` methods fail (missing DB, corrupt file, connection error), the error propagates directly to the caller, causing the entire operation to fail.

The `ObjectRepository.Walk()` method already provides reliable filesystem scanning that silently skips corrupted files. This is the natural fallback path.

## Goals / Non-Goals

**Goals:**

- `QueryService.Query()` falls back to filesystem scanning + in-memory filtering when the index is unavailable
- `QueryService.Search()` falls back to filesystem scanning + substring matching when the index is unavailable
- `VaultStats()` and `TypeStats()` also fall back (they delegate to `Query()` internally)
- Emit a `slog.Warn` when entering fallback mode
- All filter operators supported by `FilterRuleToSQL` are also supported by in-memory matching

**Non-Goals:**

- Fallback for `ListRelations()`, `ListBacklinks()`, `ListWikiLinks()` — these require relation/wikilink tables that only exist in the index; no filesystem equivalent
- Fallback for write operations (Upsert, Remove, etc.)
- Auto-rebuilding the index on fallback detection
- DB-becomes-available-mid-session re-enablement (requires vault restart)
- Performance optimization of fallback path (O(n) is acceptable)

## Decisions

### D1: Fallback at QueryService layer, not ObjectIndex

**Decision:** Implement fallback logic in `QueryService.Query()` and `QueryService.Search()`, catching errors from `ObjectIndex` and falling back to `ObjectRepository.Walk()`.

**Alternatives considered:**
- *Wrap ObjectIndex in a fallback decorator*: Would require implementing the full ObjectIndex interface for the fallback, including write methods that don't make sense in fallback mode. Over-engineered.
- *Add fallback at Vault layer*: Would bypass QueryService, breaking the layered architecture. QueryService is the right place since it already coordinates repo and index.

**Rationale:** QueryService already holds both `repo` and `index` references. The fallback is a simple try-index-then-walk pattern in just two methods.

### D2: In-memory filter matching via `MatchFilter()` function

**Decision:** Create a `MatchFilter(obj *Object, rule FilterRule) bool` function in a new file `core/filter_match.go` that mirrors `FilterRuleToSQL` operators but works against in-memory `Object` fields.

**Alternatives considered:**
- *Reuse FilterRuleToSQL by creating an in-memory SQLite DB*: Clever but fragile and slow — defeats the purpose of a lightweight fallback.
- *Only support "type" and "is" operators in fallback*: Too limited; views use all operators.

**Rationale:** Direct mirror of `FilterRuleToSQL` operators ensures consistent behavior between indexed and fallback paths. The function is straightforward — it's a switch on `rule.Operator` reading from `obj.Properties`.

### D3: In-memory sort via `SortObjects()` function

**Decision:** Create a `SortObjects(objects []*Object, sort []SortRule)` function that applies sort rules to an `[]*Object` slice using `sort.SliceStable`.

**Rationale:** Sorting is needed for both `Query()` fallback and potentially future in-memory operations. Separate from filtering for single-responsibility.

### D4: Substring search for Search() fallback

**Decision:** Fallback search uses case-insensitive substring matching against object name (from `Properties["name"]`), description (`Properties["description"]`), and body text.

**Alternatives considered:**
- *Regex matching*: More powerful but users don't expect regex in a search box; also risk of invalid patterns.
- *No search fallback (return error)*: Inconsistent with the Query fallback story.

**Rationale:** Substring match provides a reasonable degraded experience. No ranking or relevance scoring — results are returned in walk order.

### D5: Warn once per fallback invocation

**Decision:** Each fallback invocation logs `slog.Warn("index unavailable, using filesystem fallback")` with the operation name. No persistent state tracking (e.g., "already warned in this session").

**Rationale:** Simple and sufficient. Log-based warning is observable by CLI (`--debug`) and TUI (log file). Per-invocation warning is acceptable since fallback should be rare.

## Risks / Trade-offs

- **[Performance]** Fallback is O(n) file reads per query. → Acceptable for small-to-medium vaults (hundreds of objects). Document in issue/PR that large vaults (thousands of objects) will experience noticeable latency.
- **[Incomplete fallback]** Relations and backlinks have no fallback — `BuildDisplayProperties()` will fail if the index is down. → Explicitly out of scope. The right panel will show an error for relations/backlinks, but the core list/search/sidebar will work.
- **[Filter parity]** In-memory filter must match all SQLite operators exactly. → Mitigated by testing both paths with the same BDD scenarios.
- **[Walk() cost]** Each fallback call does a full Walk(). No caching between calls. → Acceptable given fallback is a degraded mode, not a performance target.
