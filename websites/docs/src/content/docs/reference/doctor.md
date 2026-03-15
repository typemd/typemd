---
title: tmd doctor
description: Comprehensive vault health check with diagnostics and auto-fix.
sidebar:
  order: 11
---

Performs a comprehensive vault health check across 8 categories. A superset of [`tmd type validate`](/reference/validate/) with additional structural integrity checks.

```bash
tmd doctor
```

## Check Categories

### 1. Schemas

Validates all type schemas (same as `tmd type validate` Phase 1).

### 2. Objects

Validates object properties against their type schemas (same as Phase 2).

### 3. Relations

Checks all relation endpoints reference existing objects (same as Phase 3).

### 4. Wiki-links

Detects broken `[[target]]` references (same as Phase 4).

### 5. Uniqueness

Checks types with `unique: true` for duplicate names (same as Phase 5).

### 6. Files

Detects corrupted object files — files under `objects/` with unparseable YAML frontmatter or missing `---` delimiters that are silently skipped during normal operations.

### 7. Index

Checks if the SQLite index is in sync with files on disk. **Auto-fixed**: if the index is stale, it is automatically rebuilt.

### 8. Orphans

Detects directories under `objects/` or `templates/` that don't correspond to any known type schema. Reported as warnings.

## Output

Results are grouped by category with `✓`/`✗` indicators:

```
  ✓ Schemas
  ✓ Objects
  ✓ Relations
  ✓ Wiki-links
  ✓ Uniqueness
  ✗ Files
    [error] book/bad-file.md: yaml: did not find expected key
  ✓ Index (auto-fixed)
  ✗ Orphans
    [warn] object directory objects/ghost has no type schema

2 issue(s) found, 1 auto-fixed.
```

## Exit Code

- `0` — no issues found (auto-fixed items don't count)
- `1` — one or more issues found
