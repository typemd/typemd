## 1. System Property Registry

- [x] 1.1 Write BDD scenarios for system property registry (IsSystemProperty, SystemPropertyNames, registry contents)
- [x] 1.2 Implement BDD step definitions for registry scenarios
- [x] 1.3 Create `core/system_property.go` with SystemProperty type, registry slice, IsSystemProperty(), SystemPropertyNames()
- [x] 1.4 Add unit tests for registry edge cases (empty string, case sensitivity)

## 2. Schema Validation via Registry

- [x] 2.1 Write BDD scenarios for schema validation rejecting all system property names (created_at, updated_at)
- [x] 2.2 Implement BDD step definitions for schema validation scenarios
- [x] 2.3 Refactor `ValidateSchema` to use `IsSystemProperty()` instead of hardcoded `NameProperty` check
- [x] 2.4 Refactor `ValidateSharedProperties` to use `IsSystemProperty()` instead of hardcoded `NameProperty` check
- [x] 2.5 Add unit tests for validation edge cases

## 3. Timestamps on Object Creation

- [x] 3.1 Write BDD scenarios for new objects having created_at and updated_at
- [x] 3.2 Implement BDD step definitions for timestamp creation scenarios
- [x] 3.3 Update `NewObject` to set `created_at` and `updated_at` using `time.Now().Format(time.RFC3339)`
- [x] 3.4 Add unit tests for timestamp format validation (RFC 3339 with timezone offset)

## 4. Timestamps on Save

- [x] 4.1 Write BDD scenarios for SaveObject and SetProperty updating updated_at
- [x] 4.2 Implement BDD step definitions for save timestamp scenarios
- [x] 4.3 Update `saveObjectFile` to set `updated_at` to current time, preserve `created_at`
- [x] 4.4 Add unit tests for save timestamp behavior (created_at unchanged, updated_at refreshed)

## 5. Frontmatter Ordering and Sync

- [x] 5.1 Write BDD scenarios for frontmatter ordering with system properties
- [x] 5.2 Implement BDD step definitions for ordering scenarios
- [x] 5.3 Refactor `OrderedPropKeys` to use `SystemPropertyNames()` for ordering instead of hardcoded `NameProperty`
- [x] 5.4 Refactor `SyncIndex` property filtering to preserve all system properties via registry
- [x] 5.5 Add unit tests for ordering edge cases (missing timestamps, no schema)

## 6. Graceful Absence for Existing Objects

- [x] 6.1 Write BDD scenarios for existing objects without timestamps
- [x] 6.2 Implement BDD step definitions for graceful absence scenarios
- [x] 6.3 Verify SyncIndex does not add timestamps to existing objects (no migration)
- [x] 6.4 Add unit tests for GetObject with missing timestamps
