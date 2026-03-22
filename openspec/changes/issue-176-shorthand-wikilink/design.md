## Context

Wiki-links currently require the full `type/name-ulid` format. The relation system already has name-based resolution via `buildNameIndex()` and `resolveByName()` in `name_resolve.go`, with write-back of expanded IDs to frontmatter. This design extends the same pattern to wiki-links in markdown body content.

The `syncWikiLinks()` function in `projector.go` currently does a simple lookup against `diskIDs` (full ID → bool). The `syncContext` already includes `nameIndex` built by `buildNameIndex()`, which is used for relation resolution but not yet for wiki-links.

## Goals / Non-Goals

**Goals:**

- Support three wiki-link formats: `[[type/name-ulid]]` (exact), `[[type/name]]` (type-qualified), `[[name]]` (same-type)
- Resolve shorthand links during sync and write back full IDs to source files
- Report ambiguous/unresolved shorthand links via SyncResult and validation
- Provide `tmd fix wikilinks` to batch-expand shorthand links

**Non-Goals:**

- Autocomplete for shorthand links in TUI/CLI (separate issue)
- Cross-type fallback for `[[name]]` (same-type only)
- Changes to DB schema or wiki-link storage format
- Changes to wiki-link parsing regex (already accepts any target string)

## Decisions

### 1. Resolution order in syncWikiLinks()

**Decision:** Three-step resolution: full ID → type-qualified name → same-type name.

```
Target string
    │
    ├── contains "/" AND has ULID suffix?
    │   → Full ID: lookup in diskIDs (existing behavior)
    │
    ├── contains "/"  AND no ULID suffix?
    │   → Type-qualified: split into type + name, resolveByName(nameIndex, type, name)
    │
    └── no "/"?
        → Same-type: resolveByName(nameIndex, sourceObjectType, target)
```

**Rationale:** Reuses existing `isFullObjectID()` and `resolveByName()`. The check order is unambiguous — the presence of `/` and ULID suffix fully determines which resolution path to take.

### 2. Write-back mechanism

**Decision:** Replace shorthand targets in-place within the object body string, then save the modified object to disk.

**Approach:**
- Track which bodies are modified during `syncWikiLinksAndTags()`
- Use string replacement: `[[shorthand]]` → `[[full-id]]` and `[[shorthand|text]]` → `[[full-id|text]]`
- Save modified objects via `repo.Save()` after sync completes (similar to `writeBackModified()` for relations)

**Rationale:** Consistent with how relation name expansion works. Body write-back is slightly more complex than frontmatter write-back (regex replacement vs. map update) but the same principle applies.

### 3. Ambiguity handling

**Decision:** Ambiguous shorthand links are left unresolved (not written back) and reported in SyncResult.

- `syncWikiLinks()` returns unresolved wiki-links alongside expanded count
- Toast notifications aggregate: `⚠ 2 ambiguous wiki-links`
- `ValidateWikiLinks()` reports each ambiguous link with the list of matching full IDs

**Rationale:** Consistent with relation ambiguity handling. The user must disambiguate manually.

### 4. SyncResult extension

**Decision:** Add `WikiLinksExpanded int` and `UnresolvedWikiLinks []UnresolvedWikiLink` fields to `SyncResult`.

```go
type UnresolvedWikiLink struct {
    ObjectID string
    Target   string
    Reason   string   // "not_found" or "ambiguous"
    Matches  []string // populated when ambiguous
}
```

**Rationale:** Mirrors `Expanded`/`Unresolved` fields already used for relation resolution. Enables toast notifications and detailed reporting.

### 5. tmd fix wikilinks command

**Decision:** New subcommand under `tmd fix` that walks all objects, resolves shorthand wiki-links, and writes back full IDs.

**Approach:**
- Build name index from all objects (reuse `buildNameIndex()`)
- For each object, parse wiki-links and resolve shorthand targets
- Replace in body and save
- Report summary: N expanded, M unresolved

**Rationale:** Provides a way to bulk-normalize wiki-links without running full sync. Useful after imports or manual edits.

### 6. Validate changes

**Decision:** `ValidateWikiLinks()` extended to resolve shorthand targets. Ambiguous links reported with suggestion of matching full IDs.

```
book/my-notes-01abc: ambiguous wiki-link [[clean-code]] — matches:
  book/clean-code-01xyz
  book/clean-code-principles-01def
```

**Rationale:** Helps users fix ambiguous links by showing available options.

## Risks / Trade-offs

- **Body modification during sync** — Write-back modifies source files, which could surprise users editing in external editors. → Mitigation: Same pattern as relation expansion; file watcher debounce handles concurrent edits.
- **Name collision across types** — `[[golang]]` could mean different things in different contexts. → Mitigation: `[[name]]` is strictly same-type; use `[[type/name]]` for cross-type references.
- **Performance** — Name index is already built during sync; no additional walk required. Wiki-link resolution adds per-link overhead but this is negligible for typical vault sizes.
- **Regex replacement in body** — Must handle edge cases like links inside code blocks or escaped brackets. → Mitigation: Current parser already handles these; replacement uses the same regex pattern.
