## Why

Getting existing content into typemd is the biggest friction point in adoption. The current `explore` skill (read-only analysis) and `importer` skill (single-file conversion) require manual glue work between them — users must create type schemas by hand after exploring and before importing. This skill closes that gap with a unified four-phase workflow: scan, plan, execute, verify.

## What Changes

- Add a new `onboarding` embedded skill (`core/skills/onboarding/SKILL.md`) that orchestrates the full import pipeline
- Add a new `onboarding` marketplace skill (`marketplace/plugins/typemd/skills/onboarding/SKILL.md`) for AI discovery
- Add `tmd import scan <paths...>` CLI command — scan sources and output structured analysis (JSON)
- Add `tmd import plan <paths...>` CLI command — generate a conversion plan based on scan + existing schemas
- Add `tmd import execute <plan-file>` CLI command — execute a confirmed plan (create types, batch-create objects, resolve relations)
- Add scan/plan/execute logic in `core/` — file scanning, type inference, plan generation, batch import with dependency ordering
- Support incremental additions — detect previously imported sources to avoid duplication

## Capabilities

### New Capabilities

- `source-scanning`: Scan source directories/files and produce structured analysis — file count, directory structure, frontmatter patterns, naming conventions, size distribution
- `conversion-planning`: Generate a structured conversion plan — map each source file to a target type + property mapping, determine import order by dependency, present for user review
- `batch-import`: Execute a confirmed plan — create type schemas, batch-create objects in dependency order, resolve wiki-links and relations in a second pass, handle conflicts (skip/overwrite/merge)
- `import-verification`: Report import results — created/failed/skipped counts, unresolved references, schema adjustment suggestions, follow-up action recommendations

### Modified Capabilities

_(none — this is a new workflow; existing explore/importer skills coexist unchanged)_

## Impact

- **core/**: New scan, plan, and execute logic; new embedded skill
- **cmd/**: New `tmd import` command group with `scan`, `plan`, `execute` subcommands
- **marketplace/**: New onboarding skill referencing the embedded skill
- **Dependencies**: No new external dependencies — uses existing `core/` object creation, type schema, and relation infrastructure
- **Existing skills**: `explore` and `importer` remain unchanged and coexist for simpler use cases
