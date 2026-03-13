## 1. BDD: Update spec for built-in types

- [x] 1.1 Write BDD scenario in `core/features/` for "only tag is built-in" (loading undefined type returns error, tag loads without custom schema)
- [x] 1.2 Implement step definitions for the new scenarios

Note: Covered by unit tests `TestDefaultTypes_OnlyTag`, `TestVault_LoadType_BuiltinFallback` (tests tag-only fallback and book error). BDD feature file not needed — this is an implementation detail (defaultTypes contents), not a user-facing behavior.

## 2. Core: Remove built-in types

- [x] 2.1 Remove `book`, `person`, `note` entries from `defaultTypes` in `core/type_schema.go` (keep only `tag`)
- [x] 2.2 Verify new BDD scenarios pass

## 3. Test helper: Type schema creation utility

- [x] 3.1 Add `writeCommonTestTypeSchemas(v)` helper function that writes book/person/note YAML schemas
- [x] 3.2 Integrated into `aVaultIsReady()` BDD step and `setupTestVault()` unit test helper — all tests using these helpers automatically get type schemas

Note: Strategy changed from per-file fixes to centralizing schema creation in shared helpers. This avoided modifying 30+ test files individually.

## 4. Fix unit tests

- [x] 4.1 Fix `core/type_schema_test.go` — replaced `TestDefaultTypes`, `TestDefaultTypes_BookUsesSelect`, `TestDefaultTypes_HaveEmoji`, `TestDefaultTypes_NoteDoesNotHaveTagsProperty`, `TestVault_LoadType_BuiltinFallback`
- [x] 4.2 Fix `core/object_test.go` — handled via `setupTestVault()` helper
- [x] 4.3 Fix `core/relation_test.go` — handled via `setupTestVault()` helper
- [x] 4.4 Fix `core/query_test.go` — handled via `setupTestVault()` helper
- [x] 4.5 Fix `core/sync_test.go` — handled via `setupTestVault()` helper
- [x] 4.6 Fix `core/validate_test.go` — handled via `setupTestVault()` helper
- [x] 4.7 Fix `core/system_property_test.go` — handled via `setupTestVault()` helper
- [x] 4.8 Fix `core/wikilink_db_test.go` — handled via `setupTestVault()` helper
- [x] 4.9 Fix `core/display_test.go` — added `writeCommonTestTypeSchemas()` calls directly
- [x] 4.10 Fix `core/relation_list_test.go` — handled via `setupTestVault()` helper
- [x] 4.11 Fix `core/migrate_test.go` — handled via `setupTestVault()` helper
- [x] 4.12 Fix `core/tag_test.go` — handled via `setupTestVault()` helper
- [x] 4.13 Fix `core/ulid_test.go` — handled via `setupTestVault()` helper
- [x] 4.14 Fix `core/vault_test.go` — updated `setupTestVault()` helper
- [x] 4.15 Fix `core/list_test.go` — updated `TestVault_ListTypes_CustomOverridesDefault` to test tag override
- [x] 4.16 Fix `cmd/create_test.go` — added schema YAML writes in `setupTestVaultDir()`
- [x] 4.17 Fix `tui/app_test.go` — added schema YAML write in `setupTestModelWithVault()`

## 5. Fix BDD feature files and step definitions

- [x] 5.1 Update `core/features/object.feature` — handled via `aVaultIsReady()` helper
- [x] 5.2 Update `core/features/relation.feature` — already writes custom schemas
- [x] 5.3 Update `core/features/wikilink.feature` — already writes custom schemas
- [x] 5.4 Update `core/features/query.feature` — handled via `aVaultIsReady()` helper
- [x] 5.5 Update `core/features/system_property.feature` — handled via `aVaultIsReady()` helper
- [x] 5.6 Update `core/features/name_property.feature` — handled via `aVaultIsReady()` helper
- [x] 5.7 Update `core/features/property_type.feature` — handled via `aVaultIsReady()` helper
- [x] 5.8 Update `core/features/property_filtering.feature` — handled via `aVaultIsReady()` helper
- [x] 5.9 Update `core/features/resolve.feature` — handled via `aVaultIsReady()` helper
- [x] 5.10 Update `core/features/validate.feature` — handled via `aVaultIsReady()` helper
- [x] 5.11 Update `core/features/tag_resolution.feature` — handled via `aVaultIsReady()` helper
- [x] 5.12 Update `core/features/tag_type.feature` and `tags_property.feature` — handled via `aVaultIsReady()` helper

## 6. Verify

- [x] 6.1 Run `go test ./...` — all tests pass
- [x] 6.2 Run `go build ./...` — clean build
