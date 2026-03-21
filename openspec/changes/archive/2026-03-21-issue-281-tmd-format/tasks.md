## 1. Core: Object Formatting

- [x] 1.1 Write BDD scenarios for object formatting (property reordering, YAML normalization, body preservation, updated_at unchanged, already-formatted skip)
- [x] 1.2 Implement BDD step definitions for object formatting scenarios
- [x] 1.3 Add `FormatObjects(typeName string, dryRun bool) (*FormatResult, error)` to Vault — reads each object, re-serializes with `OrderedPropKeys` + `writeFrontmatter`, compares bytes, writes if different
- [x] 1.4 Add unit tests for edge cases (empty frontmatter, object without schema, nil properties)

## 2. Core: Schema Formatting

- [x] 2.1 Write BDD scenarios for schema formatting (non-canonical YAML reformat, already-formatted skip, built-in type skip)
- [x] 2.2 Implement BDD step definitions for schema formatting scenarios
- [x] 2.3 Add `ListSchemaNames()` to `LocalObjectRepository` — walks `.typemd/types/` and returns file-backed schema names
- [x] 2.4 Add `FormatSchemas(typeName string, dryRun bool) (*FormatResult, error)` to Vault — loads each schema, re-marshals with `MarshalTypeSchema`, compares bytes, writes if different
- [x] 2.5 Add unit tests for schema formatting edge cases (schema with name template, schema with use entries)

## 3. Core: Combined Format Entry Point

- [x] 3.1 Add `FormatResult` struct and `FormatAll(typeName string, dryRun bool) (*FormatResult, error)` to Vault — combines object + schema formatting results
- [x] 3.2 Write BDD scenarios for type filter validation (valid type, invalid type error)
- [x] 3.3 Implement step definitions for type filter scenarios

## 4. CLI: `tmd format` Command

- [x] 4.1 Write BDD scenarios for CLI command (no args, --type flag, --dry-run flag, --dry-run exit code)
- [x] 4.2 Implement BDD step definitions for CLI scenarios
- [x] 4.3 Add `cmd/format.go` with Cobra command — `--type` string flag, `--dry-run` bool flag, calls `vault.FormatAll()`, prints summary, returns exit code 1 on dry-run with changes
- [x] 4.4 Add unit tests for CLI output formatting (summary messages, dry-run output)
