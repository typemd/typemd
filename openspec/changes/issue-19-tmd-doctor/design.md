## Context

typemd has an existing `tmd type validate` command that covers schema conformance (5 phases). However, it does not detect corrupted files (silently skipped by `Walk()`), index-disk desync, or orphan files/directories/templates. There is no repair capability.

The `Vault` already auto-syncs on `Open()` when `NeedsSync()` is true, and `Projector.Sync()` already cleans orphaned relations. The doctor command builds on these foundations, consolidating and extending them into a comprehensive diagnostic tool.

## Goals / Non-Goals

**Goals:**
- Provide a single `tmd doctor` command that checks all aspects of vault health
- Reuse existing validation functions from `core/validate.go` (no duplication)
- Add new checks: corrupted frontmatter, index-disk sync, orphan files/dirs/templates
- Implement tiered auto-fix: safe ops auto-fix, dangerous ops report only (no `--fix` flag in v1)
- Output grouped summary (git status style) with final summary line
- Exit code 0 = healthy, 1 = issues found

**Non-Goals:**
- Interactive repair mode (prompt user to confirm dangerous fixes) — future work
- `--fix` flag for dangerous operations (deleting orphan files) — future work
- Deprecating or modifying `tmd type validate` — it stays as-is
- JSON output format — future work

## Decisions

### 1. Core diagnostic engine as `DoctorReport` struct

The doctor logic lives in `core/doctor.go` as a `RunDoctor(v *Vault)` function returning a `DoctorReport` struct. The report contains per-category results with issues typed by severity (error vs warning vs auto-fixed).

**Why:** Separating the diagnostic engine from CLI presentation follows the existing pattern (e.g. `ValidateAllSchemas` returns data, `cmd/validate.go` formats output). This keeps `core/` testable and allows TUI/MCP to reuse the diagnostics.

**Alternative considered:** Adding checks directly in `cmd/doctor.go`. Rejected because it breaks the clean architecture boundary.

### 2. Eight check categories in fixed order

| # | Category | Source | Auto-fix |
|---|----------|--------|----------|
| 1 | Schemas | `ValidateAllSchemas()` | No |
| 2 | Objects | `ValidateAllObjects()` | No |
| 3 | Relations | `ValidateRelations()` | No |
| 4 | Wiki-links | `ValidateWikiLinks()` | No |
| 5 | Name uniqueness | `ValidateNameUniqueness()` | No |
| 6 | Corrupted files | New: `ScanCorruptedFiles()` | No |
| 7 | Index sync | `NeedsSync()` + `SyncIndex()` | Yes (auto-rebuild) |
| 8 | Orphans | New: scan for orphan dirs/templates | No (report only) |

Checks 1-5 delegate directly to existing functions. Checks 6-8 are new.

**Why this order:** Schema/object/relation/wikilink/uniqueness are the "conformance" checks (matching `tmd type validate` order). Corrupted files, index sync, and orphans are "structural integrity" checks that come after.

### 3. Corrupted file scanning via new `WalkAll()` method

`LocalObjectRepository.Walk()` silently skips unparseable files. Instead of modifying `Walk()` (which would break Projector assumptions), add a new `WalkAll()` method on `LocalObjectRepository` that returns both parsed objects and a list of `CorruptedFile{Path, Error}` entries.

**Why:** `Walk()` is used by Projector which relies on silent skipping. A separate method avoids breaking existing behavior.

**Alternative considered:** Adding an error callback to `Walk()`. Rejected because it changes the interface signature.

### 4. Orphan detection via filesystem scan

Scan `objects/` directories against known type schemas (from `ListSchemas()`). Any directory not matching a known type is an orphan. Also scan `templates/` directories the same way.

**Why:** Simple filesystem comparison. No index needed — types come from `.typemd/types/*.yaml` + built-in defaults.

### 5. Report structure

```go
type DoctorReport struct {
    Categories []DoctorCategory
}

type DoctorCategory struct {
    Name     string
    Issues   []DoctorIssue
    AutoFixed int
}

type DoctorIssue struct {
    Severity IssueSeverity // Error, Warning
    Message  string
}
```

The CLI renders this as grouped summary with `✓`/`✗` per category and a final summary line.

## Risks / Trade-offs

- **[Performance on large vaults]** → Doctor runs all checks sequentially. For large vaults this could be slow. Mitigation: checks 1-5 already query the index (fast). Check 6 (`WalkAll`) requires full disk scan but is the same cost as `SyncIndex`. Acceptable for a diagnostic command.
- **[WalkAll duplicates Walk logic]** → Some code duplication between `Walk()` and `WalkAll()`. Mitigation: `WalkAll()` can internally call common helpers. The duplication is minimal and justified by keeping `Walk()` stable.
- **[No dangerous auto-fix]** → Orphan files/dirs are reported but not deleted. Users must manually clean up. Mitigation: clear error messages with file paths. Interactive fix can be added later.
