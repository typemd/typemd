## 1. Core: Shared Properties Loading

- [x] 1.1 Write BDD scenarios for loading shared properties (file exists, missing, empty, with various property types)
- [x] 1.2 Implement BDD step definitions for shared properties loading scenarios
- [x] 1.3 Add `Use` field to Property struct and `SharedPropertiesFile` struct
- [x] 1.4 Implement `LoadSharedProperties()` on Vault (read, parse, cache)
- [x] 1.5 Add unit tests for LoadSharedProperties edge cases (malformed YAML, caching)

## 2. Core: Shared Properties Validation

- [x] 2.1 Write BDD scenarios for shared properties validation (duplicate names, invalid types, reserved name, select without options)
- [x] 2.2 Implement BDD step definitions for shared properties validation scenarios
- [x] 2.3 Implement `ValidateSharedProperties()` (duplicate names, reserved name, standard property rules)
- [x] 2.4 Add unit tests for ValidateSharedProperties edge cases

## 3. Core: Use Keyword in Type Schema

- [x] 3.1 Write BDD scenarios for `use` keyword parsing and validation (use-only, use+pin, use+emoji, disallowed fields, non-existent ref, name conflict)
- [x] 3.2 Implement BDD step definitions for `use` keyword scenarios
- [x] 3.3 Extend `ValidateSchema()` to validate `use` entries (allowed fields, existence check, name conflicts, post-resolution duplicate check)
- [x] 3.4 Add unit tests for ValidateSchema use-related edge cases (both use and name set, use with type/options/default)

## 4. Core: LoadType Resolution

- [x] 4.1 Write BDD scenarios for `use` resolution in LoadType (no override, pin override, emoji override, mixed properties order)
- [x] 4.2 Implement BDD step definitions for LoadType resolution scenarios
- [x] 4.3 Extend `LoadType()` to resolve `use` entries from shared properties with pin/emoji overrides
- [x] 4.4 Add unit tests for LoadType resolution edge cases (Use field cleared after resolution)

## 5. Integration: ValidateAllSchemas

- [x] 5.1 Extend `ValidateAllSchemas()` to call `ValidateSharedProperties()` and pass shared properties to `ValidateSchema()`
- [x] 5.2 Add unit tests for end-to-end validation with shared properties

## 6. Example Vault

- [x] 6.1 Add `.typemd/properties.yaml` to `examples/book-vault/` with sample shared properties
- [x] 6.2 Update example type schemas to use `use` keyword where appropriate
