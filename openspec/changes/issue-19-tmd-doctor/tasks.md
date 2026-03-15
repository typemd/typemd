## 1. Core: DoctorReport data model

- [x] 1.1 Write BDD scenarios for doctor report structure (categories, issues, severity, auto-fixed count)
- [x] 1.2 Implement BDD step definitions for doctor report scenarios
- [x] 1.3 Define `DoctorReport`, `DoctorCategory`, `DoctorIssue`, `IssueSeverity` types in `core/doctor.go`
- [x] 1.4 Add unit tests for report helper methods (total issues, total auto-fixed, has errors)

## 2. Core: Corrupted file scanning

- [x] 2.1 Write BDD scenarios for corrupted file detection (malformed YAML, missing frontmatter, all valid)
- [x] 2.2 Implement BDD step definitions for corrupted file scenarios
- [x] 2.3 Add `WalkAll()` method to `LocalObjectRepository` returning both parsed objects and `CorruptedFile` entries
- [x] 2.4 Add `ScanCorruptedFiles()` function in `core/doctor.go` using `WalkAll()`
- [x] 2.5 Add unit tests for `WalkAll()` edge cases (empty dir, nested dirs, mixed valid/invalid files)

## 3. Core: Orphan detection

- [x] 3.1 Write BDD scenarios for orphan detection (orphan object dirs, orphan template dirs, no orphans)
- [x] 3.2 Implement BDD step definitions for orphan detection scenarios
- [x] 3.3 Add `ScanOrphanDirs()` function in `core/doctor.go` scanning objects/ and templates/ against known types
- [x] 3.4 Add unit tests for orphan detection edge cases (no templates dir, built-in types, empty dirs)

## 4. Core: RunDoctor orchestrator

- [x] 4.1 Write BDD scenarios for full doctor run (healthy vault, vault with mixed issues, index auto-fix)
- [x] 4.2 Implement BDD step definitions for full doctor run scenarios
- [x] 4.3 Implement `RunDoctor(v *Vault)` function orchestrating all 8 checks and returning `DoctorReport`
- [x] 4.4 Add unit tests for RunDoctor (exit code logic: no issues = 0, issues = 1, only auto-fixed = 0)

## 5. CLI: tmd doctor command

- [x] 5.1 Add `cmd/doctor.go` with Cobra command registration, calling `core.RunDoctor()` and rendering grouped summary output
- [x] 5.2 Add unit tests for output formatting (✓/✗ indicators, summary line, auto-fixed display)
