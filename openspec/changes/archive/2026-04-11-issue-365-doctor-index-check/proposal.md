## Why

The `vault-doctor` spec defines an **Index category** requirement — the doctor should verify that the SQLite index is in sync with files on disk, and auto-rebuild if not. The spec currently describes 8 categories, but `core/doctor.go` only implements 7 (Schemas, Objects, Relations, Wiki-links, Uniqueness, Files, Orphans). The Index category was never wired up, leaving a documented-but-unimplemented requirement.

This gap was surfaced by a consistency check cross-referencing BDD, OpenSpec, and code (see issue #365).

## What Changes

- Add an **Index category** to `RunDoctor()` that detects whether the SQLite index is in sync with files on disk.
- When the index is out of sync, the doctor automatically rebuilds it (via `Reconcile()` + `Project()`) and records the count as auto-fixed under the Index category.
- When the index is in sync, the Index category passes silently with 0 issues and 0 auto-fixed.
- Update `core/features/doctor.feature` to assert **8** categories (up from 7) and include scenarios for the Index category.
- Auto-fixed index rebuilds do **not** count toward the process exit code (already honored by `HasUnresolvedIssues`).

## Capabilities

### New Capabilities

None.

### Modified Capabilities

- `vault-doctor`: The "Doctor checks index-disk synchronization" requirement is already defined in the spec but not implemented. This change implements it — no spec text changes, but we need a delta file to link the change to the capability.

## Impact

- **Affected code**: `core/doctor.go` (new `checkIndexSync` function, new category in `RunDoctor`), `core/features/doctor.feature` (category count update + new scenarios), `core/bdd_steps_doctor_test.go` (if a new step is needed).
- **No CLI / cmd changes**: `cmd/doctor.go` already documents 8 categories including "Index" — it was ahead of the implementation.
- **No API or dependency changes**.
