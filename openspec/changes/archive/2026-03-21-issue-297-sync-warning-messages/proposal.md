## Why

When Projector sync encounters unresolved relation references, the TUI displays a generic aggregated toast like "⚠ 3 unresolved refs". Users cannot tell whether references were not found or ambiguous, making it difficult to diagnose and fix the issue without running `tmd validate`. The `UnresolvedRelation` struct already carries `Reason` ("not_found" / "ambiguous") and `Matches` data, but the TUI's toast conversion discards this information.

## What Changes

- Improve the `refreshData()` toast conversion logic to use separate group keys based on `UnresolvedRelation.Reason`, producing distinct messages like "⚠ 2 not found" and "⚠ 1 ambiguous" instead of a single "⚠ 3 unresolved refs"
- Update existing tests to verify the new group-based message format

## Capabilities

### New Capabilities

_(none)_

### Modified Capabilities

- `toast-notification`: The "Sync warnings display as toast" requirement changes from a single "unresolved refs" group to reason-specific groups ("not found" / "ambiguous")

## Impact

- `tui/app.go` — `refreshData()` toast conversion logic (~10 lines changed)
- `tui/sync_toast_test.go` — test expectations updated
- `openspec/specs/toast-notification/spec.md` — requirement updated
