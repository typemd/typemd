## 1. TypeSchema: Add Unique field

- [x] 1.1 Write BDD scenarios for unique field parsing (unique true, false, omitted)
- [x] 1.2 Implement step definitions for unique field scenarios
- [x] 1.3 Add `Unique bool` field to TypeSchema struct with `yaml:"unique"` tag
- [x] 1.4 Add unit test: built-in tag schema has Unique true

## 2. Core: Generalize name uniqueness check

- [x] 2.1 Write BDD scenarios for creation-time uniqueness enforcement (duplicate rejected, first succeeds, different types allowed, non-unique type allows duplicates)
- [x] 2.2 Implement step definitions for uniqueness scenarios
- [x] 2.3 Add `checkNameUnique(typeName, name string) error` method to Vault
- [x] 2.4 Update `NewObject()` to use schema-driven uniqueness check instead of hardcoded tag check
- [x] 2.5 Remove `checkTagNameUnique()` from tag.go
- [x] 2.6 Update built-in tag default schema to include `Unique: true`

## 3. Validation: Generalize uniqueness validation

- [x] 3.1 Write BDD scenarios for validation (no duplicates passes, duplicates reported, non-unique types skipped)
- [x] 3.2 Implement step definitions for validation scenarios
- [x] 3.3 Add `ValidateNameUniqueness(vault *Vault) []ValidationError` that scans all unique types
- [x] 3.4 Replace `ValidateTagNameUniqueness()` call in validate command with `ValidateNameUniqueness()`
- [x] 3.5 Remove `ValidateTagNameUniqueness()` from validate.go

## 4. Cleanup and verify

- [x] 4.1 Update existing tag_uniqueness.feature scenarios to work with generalized mechanism
- [x] 4.2 Run full test suite and fix any regressions
