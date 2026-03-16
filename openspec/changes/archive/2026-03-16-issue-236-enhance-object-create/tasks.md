## 1. Core: Vault Config

- [x] 1.1 Write BDD scenarios for vault config loading (file exists, missing, empty, invalid YAML)
- [x] 1.2 Implement BDD step definitions for vault config scenarios
- [x] 1.3 Add `VaultConfig` and `CLIConfig` structs, config loading in `Vault.Open()`
- [x] 1.4 Add `DefaultType()` accessor method on Vault
- [x] 1.5 Add unit tests for config edge cases (unknown keys ignored, partial config)

## 2. Core: Slug Conversion

- [x] 2.1 Write BDD scenarios for slug conversion (spaces, mixed case, special chars, idempotent)
- [x] 2.2 Implement BDD step definitions for slug conversion scenarios
- [x] 2.3 Add `Slugify()` function in core package
- [x] 2.4 Add unit tests for slug edge cases (empty string, non-ASCII, numbers, consecutive hyphens)

## 3. Core: ObjectService.Create Slug Integration

- [x] 3.1 Write BDD scenarios for Create with natural-language names (slug for filename, original for name property)
- [x] 3.2 Implement BDD step definitions for Create slug scenarios
- [x] 3.3 Modify `ObjectService.Create()` to apply `Slugify()` to filename and preserve original input as name property
- [x] 3.4 Add unit tests for backward compatibility (pre-slugified input unchanged)

## 4. CLI: Type Argument Optional + --type Flag

- [x] 4.1 Add `--type` flag to `createCmd` (no `-t` short form)
- [x] 4.2 Change args validation from `RangeArgs(1, 2)` to `RangeArgs(0, 2)`
- [x] 4.3 Implement type resolution logic (0/1/2 args + flag + config fallback)
- [x] 4.4 Add unit tests for arg resolution (0 args, 1 arg as type, 1 arg as name, 2 args, flag override)

## 5. CLI: Init Config Generation

- [x] 5.1 Modify `tmd init` to create `config.yaml` when idea or note starter is selected
- [x] 5.2 Add unit tests for init config generation (idea selected, note only, neither selected)

## 6. Integration Verification

- [x] 6.1 Run full test suite (`go test ./...`) and verify no regressions
- [x] 6.2 Manual test: `tmd object create "Some Thought"` with config default type
- [x] 6.3 Manual test: `tmd object create --type note "Meeting Notes"`
- [x] 6.4 Manual test: `tmd object create book "Clean Code"` (backward compatible)
