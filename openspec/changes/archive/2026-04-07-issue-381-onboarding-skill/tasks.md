## 1. Core: Source Scanning

- [x] 1.1 Write BDD scenarios for source scanning (scan directory, multiple sources, no markdown files, frontmatter extraction, existing vault types)
- [x] 1.2 Implement BDD step definitions for source scanning scenarios
- [x] 1.3 Add `ScanResult`, `SourceInfo`, `DirInfo`, `FrontmatterStats` structs in `core/import.go`
- [x] 1.4 Implement `ScanSources(paths []string) (*ScanResult, error)` method on Vault — scan files, extract frontmatter patterns, include existing types
- [x] 1.5 Add unit tests for scan edge cases (non-existent path, empty directory, files without frontmatter, mixed file types)

## 2. Core: Conversion Planning

- [x] 2.1 Write BDD scenarios for conversion planning (generate plan, dependency ordering, incremental import detection, new type schemas in plan)
- [x] 2.2 Implement BDD step definitions for conversion planning scenarios
- [x] 2.3 Add `ImportPlan`, `TypePlan`, `ObjectPlan` structs in `core/import.go`
- [x] 2.4 Implement `GeneratePlan(scan *ScanResult, classifications []ObjectPlan) (*ImportPlan, error)` — accept AI-classified objects, order by dependency, detect existing objects
- [x] 2.5 Add unit tests for dependency ordering (tags first, relation targets first, circular dependency handling)

## 3. Core: Batch Import Execution

- [x] 3.1 Write BDD scenarios for plan execution (create types, create objects in order, handle conflicts, resolve wiki-links, report progress)
- [x] 3.2 Implement BDD step definitions for batch import scenarios
- [x] 3.3 Implement `ExecutePlan(plan *ImportPlan) (*ImportReport, error)` — create types, create objects in dependency order, run reconciliation
- [x] 3.4 Add `ImportReport` struct with created/skipped/failed counts and details
- [x] 3.5 Add unit tests for conflict handling (skip, overwrite) and error recovery (continue on single failure)

## 4. Core: Import Verification

- [x] 4.1 Write BDD scenarios for import verification (report with counts, unresolved references, follow-up suggestions)
- [x] 4.2 Implement BDD step definitions for verification scenarios
- [x] 4.3 Implement unresolved reference detection in `ImportReport` — check wiki-links and relations after reconciliation
- [x] 4.4 Add follow-up suggestion generation based on report data
- [x] 4.5 Add unit tests for report formatting and suggestion logic

## 5. CLI: Import Command Group

- [x] 5.1 Write BDD scenarios for CLI commands (`tmd import scan`, `tmd import plan`, `tmd import execute` — argument parsing, JSON output format, error messages)
- [x] 5.2 Implement BDD step definitions for CLI import scenarios
- [x] 5.3 Add `cmd/import.go` with `importCmd` parent command and `scanCmd`, `planCmd`, `executeCmd` subcommands
- [x] 5.4 Implement `tmd import scan <paths...>` — call `Vault.ScanSources()`, output JSON
- [x] 5.5 Implement `tmd import plan <paths...>` — call `Vault.ScanSources()` then `Vault.GeneratePlan()`, output JSON plan file
- [x] 5.6 Implement `tmd import execute <plan-file>` — read plan file, call `Vault.ExecutePlan()`, output progress and report
- [x] 5.7 Add unit tests for CLI argument parsing edge cases

## 6. Skills: Embedded + Marketplace

- [x] 6.1 Create `core/skills/onboarding/SKILL.md` — four-phase workflow instructions (scan → plan → execute → verify), property type inference heuristics, relation discovery patterns, conflict resolution guidance
- [x] 6.2 Register the onboarding skill in `core/instructions.go` embedded skill list
- [x] 6.3 Create `marketplace/plugins/typemd/skills/onboarding/SKILL.md` — discovery layer referencing vault-guide and `tmd instructions onboarding`
- [x] 6.4 Add unit test verifying the onboarding skill is discoverable via `core.ListSkills()` and `core.GetSkill()`
