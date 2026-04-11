## Context

`core/doctor.go` currently reports 7 categories (Schemas, Objects, Relations, Wiki-links, Uniqueness, Files, Orphans). The `vault-doctor` OpenSpec spec defines an 8th category — **Index** — which checks index-disk synchronization and auto-rebuilds when out of sync. The requirement is well-defined in the spec, the CLI help text in `cmd/doctor.go` already promises 8 categories including "Index", and the BDD feature file asserts 7 categories. This is a pure implementation gap.

A key constraint: `Vault.Open()` already runs `Reconcile()` + `Project()` on every open (see `core/vault.go:156`). This means in normal operation — where `tmd doctor` opens the vault first and then runs the doctor — the index is always freshly synced by the time `RunDoctor` executes. The Index category still serves a purpose: it provides an observability signal ("how many disk changes were applied on this run") and a safety net if someone runs doctor against a vault whose index was manually wiped or corrupted mid-run.

## Goals / Non-Goals

**Goals:**

- Add an Index category to `RunDoctor()` that appears in the report (position matters: between Files and Orphans, matching `cmd/doctor.go` help text order).
- Detect whether the SQLite index is in sync with disk by comparing indexed IDs against disk-walked IDs.
- When out of sync, auto-rebuild by running the existing `Reconcile()` + `Project()` pipeline, and report the fix count via `DoctorCategory.AutoFixed`.
- Auto-fixed index rebuilds MUST NOT count toward the exit code (already handled by `HasUnresolvedIssues()` since we record them as `AutoFixed`, not `Issues`).
- Update BDD scenarios to assert 8 categories and cover the in-sync / out-of-sync paths.

**Non-Goals:**

- Changing when `Vault.Open()` reconciles. The double-reconcile (once at Open, once at doctor) is acceptable — the second run is cheap when nothing has changed, and this keeps doctor self-contained.
- Deep schema verification (SQLite column types, FTS index integrity, etc.). The Index category only checks that object IDs in the index match IDs on disk.
- Detecting stale property data (e.g., index row exists but has outdated JSON). That would require a full content diff, which is `Reconcile()`'s job — and `Reconcile()` already emits `ObjectUpdated` events when it detects changes.

## Decisions

### Decision 1: Use a fresh `Reconcile()` + `Project()` call for the check and auto-fix

**Rationale:** The cleanest signal is to call `v.Reconcile()` a second time. If it returns zero events (Synced == 0, Deleted == 0), the index was already in sync → category passes with 0 issues, 0 auto-fixed. If it returns events, apply them via `v.Project()` and record the total event count as `AutoFixed`. This leverages all the existing diff logic and keeps the check symmetric with how `Vault.Open()` already behaves.

**Alternatives considered:**

- **Diff `ListIDs()` vs `repo.Walk()`** — lighter-weight, but only catches added/removed objects, not property drift. Rejected: the spec says "matches the current state of files on disk," which is broader than just ID presence.
- **Always rebuild the FTS index** — too expensive and doesn't address the "sync" question. `Rebuild()` on SQLiteObjectIndex only refreshes FTS, not the `objects` table.

### Decision 2: Record fix count as `AutoFixed`, not `Issues`

**Rationale:** The doctor's `HasUnresolvedIssues()` returns `TotalIssues() > 0`. Items placed in `Issues` would make the process exit 1 even after a successful auto-fix — violating the spec's "Exit code when only auto-fixed → 0" scenario. Placing the count in `AutoFixed` keeps the exit semantics correct without further changes to `HasUnresolvedIssues`.

**How:** The new `checkIndexSync` function returns a `DoctorCategory{Name: "Index", AutoFixed: n}` when a rebuild happened, or `DoctorCategory{Name: "Index"}` (empty) when in sync. Errors from `Reconcile()` or `Project()` itself become `Issues` with `SeverityError`.

### Decision 3: Place "Index" between "Files" and "Orphans" in the category list

**Rationale:** `cmd/doctor.go:13-23` already documents the category order as Schemas → Objects → Relations → Wiki-links → Uniqueness → Files → **Index** → Orphans. Matching that order keeps help text and report output aligned.

### Decision 4: BDD scenarios use a minimal out-of-sync trigger

**Rationale:** The cleanest out-of-sync trigger is to add an object file directly to disk *after* the vault has opened but *before* doctor runs. The file watcher is not active in BDD tests, so the index won't automatically pick up the new file. When `RunDoctor` reconciles, it will detect the new ID and auto-fix.

**Test step to add:** `an out-of-band object "<path>" exists` — writes a valid `.md` file with minimal frontmatter directly under `objectsDir/<path>`, bypassing `ObjectService`, so the index doesn't know about it.

## Risks / Trade-offs

- **Risk**: Running `Reconcile()` twice on every `tmd doctor` invocation adds cost. **Mitigation**: The second Reconcile is a no-op in the common case (all files already match index state), so the overhead is a single walk of `objects/` plus an index `ListIDs()` call. Acceptable for a diagnostic command.
- **Risk**: If `Reconcile()` itself returns an error, the Index category becomes a hard error rather than a silent pass. **Mitigation**: This is desired — a broken reconciler is a real problem the user should see.
- **Trade-off**: We don't distinguish between "1 file added, 1 file removed" and "2 files added" in the auto-fix count. The count is the total number of events produced. This matches how the spec describes it ("auto-fixed"), and the detail is visible in debug logs if needed.

## Migration Plan

No migration needed — this is a pure additive feature on an internal check. Existing vaults behave identically; new doctor runs simply show one additional category line.
