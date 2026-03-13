## 1. Core: Name template evaluation

- [x] 1.1 Write BDD scenarios for name template evaluation (date placeholder, static text, no placeholders)
- [x] 1.2 Implement step definitions for template evaluation scenarios
- [x] 1.3 Add `NameTemplate` field to `TypeSchema` struct and template extraction in `LoadType()`
- [x] 1.4 Implement `EvaluateNameTemplate()` function with `{{ date:FORMAT }}` support and YYYY→Go format conversion
- [x] 1.5 Add unit tests for date format conversion edge cases (YYYY-MM, HH:mm:ss, mixed tokens)

## 2. Core: Schema validation for name entry

- [x] 2.1 Write BDD scenarios for name entry validation (template accepted, type rejected, other fields rejected)
- [x] 2.2 Implement step definitions for name validation scenarios
- [x] 2.3 Modify `ValidateSchema()` to allow `name` in properties with template-only constraint
- [x] 2.4 Add unit tests for validation edge cases (name with no fields, name with template+emoji, etc.)

## 3. Core: Object creation with template

- [x] 3.1 Write BDD scenarios for object creation with template (auto name, override, no template no name)
- [x] 3.2 Implement step definitions for object creation scenarios
- [x] 3.3 Modify `NewObject()` to evaluate template when name argument is empty
- [x] 3.4 Add unit tests for slug generation from template output

## 4. CLI: Optional name argument

- [x] 4.1 Change `cmd/create.go` from `ExactArgs(2)` to `RangeArgs(1, 2)`
- [x] 4.2 Pass empty string to `NewObject()` when name arg is omitted
- [x] 4.3 Update command usage and examples
