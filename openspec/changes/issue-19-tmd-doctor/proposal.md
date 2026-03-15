## Why

There is no way to check vault health status. Issues like orphaned relations, corrupted frontmatter, index-disk desync, and orphan files can go undetected. `tmd type validate` covers schema conformance but not structural integrity. Users need a comprehensive diagnostic tool that can both detect and fix vault issues.

## What Changes

- Add `tmd doctor` CLI command as a comprehensive vault health check (superset of `tmd type validate`)
- Add 8 check categories: schema validation, object property validation, relation target validation, wikilink validation, name uniqueness, corrupted frontmatter detection, index-disk sync, and orphan file/directory/template detection
- Add tiered auto-fix: safe operations (reindex, clean orphaned relations) run automatically; dangerous operations (delete orphan files) require user confirmation
- Add `core.Doctor` diagnostic engine with structured health report
- Output uses grouped summary style (like `git status`), showing results per category with a final summary line

## Capabilities

### New Capabilities
- `vault-doctor`: Comprehensive vault health check with 8 diagnostic categories and tiered auto-fix. Covers all checks from `tmd type validate` plus corrupted frontmatter detection, index-disk sync verification, and orphan file/directory/template detection.

### Modified Capabilities

_(none — `tmd type validate` remains unchanged, doctor reuses its validation functions internally)_

## Impact

- **core/**: New `doctor.go` with diagnostic engine, new `doctor_test.go`; minor changes to `local_object_repository.go` to expose corrupted-file scanning (currently Walk silently skips)
- **cmd/**: New `doctor.go` command registration
- **core/validate.go**: No changes — doctor calls existing validation functions
- **Dependencies**: No new external dependencies
