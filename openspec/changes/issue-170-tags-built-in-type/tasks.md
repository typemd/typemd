## 1. Tag Built-in Type

- [x] 1.1 Write BDD scenarios for tag built-in type (core/features/tag_type.feature)
- [x] 1.2 Implement step definitions for tag type scenarios
- [x] 1.3 Add `tag` to `defaultTypes` with color and icon properties
- [x] 1.4 Remove `tags` property from built-in `note` type
- [x] 1.5 Add unit tests for tag type edge cases

## 2. SystemProperty Extension

- [x] 2.1 Write BDD scenarios for tags system property (core/features/tags_property.feature)
- [x] 2.2 Implement step definitions for tags system property scenarios
- [x] 2.3 Add `Target` and `Multiple` fields to `SystemProperty` struct
- [x] 2.4 Add `TagsProperty` constant and register `tags` in `systemProperties` slice
- [x] 2.5 Add unit tests for SystemProperty relation fields

## 3. Schema Validation

- [x] 3.1 Write BDD scenarios for tags rejection in type schemas and shared properties
- [x] 3.2 Implement step definitions for schema validation scenarios
- [x] 3.3 Verify existing validation logic rejects `tags` (no code change expected — IsSystemProperty already covers it)
- [x] 3.4 Add unit tests for tags rejection edge cases

## 4. Tag Name Uniqueness

- [x] 4.1 Write BDD scenarios for tag name uniqueness enforcement
- [x] 4.2 Implement step definitions for uniqueness scenarios
- [x] 4.3 Add uniqueness check in `NewObject` for tag type
- [ ] 4.4 Add uniqueness validation in `tmd type validate`
- [x] 4.5 Add unit tests for uniqueness edge cases (case sensitivity, whitespace)

## 5. Tag Reference Resolution in SyncIndex

- [x] 5.1 Write BDD scenarios for tag reference resolution (by ID and by name)
- [x] 5.2 Write BDD scenarios for auto-creation of missing tags
- [x] 5.3 Implement step definitions for resolution and auto-creation scenarios
- [x] 5.4 Implement ULID suffix detection for tag references
- [x] 5.5 Implement tag resolution by name (query tag objects by name property)
- [x] 5.6 Implement auto-creation of missing tag objects during SyncIndex
- [x] 5.7 Write tag relations to `relations` table during SyncIndex
- [x] 5.8 Ensure `tags` property is preserved during SyncIndex property filtering
- [x] 5.9 Add unit tests for ULID detection, name resolution edge cases

## 6. Extend LinkObjects / UnlinkObjects

- [x] 6.1 Write BDD scenarios for linking/unlinking via system property relations
- [x] 6.2 Implement step definitions for link/unlink scenarios
- [x] 6.3 Add `findSystemRelationProperty` function
- [x] 6.4 Modify `LinkObjects` to fall back to system properties
- [x] 6.5 Modify `UnlinkObjects` to fall back to system properties
- [x] 6.6 Add unit tests for system property relation fallback

## 7. Integration Verification

- [x] 7.1 Run full test suite (`go test ./...`) and fix any regressions
- [x] 7.2 Run `go build ./...` to verify clean build
- [ ] 7.3 Manual test: create a vault with tag objects, verify frontmatter ordering, sync, and relation queries
