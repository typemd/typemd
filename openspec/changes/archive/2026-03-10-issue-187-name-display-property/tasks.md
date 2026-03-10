## 1. Core: GetName and name property foundation

- [x] 1.1 Write BDD scenarios for GetName (name present, missing, empty string → fallback)
- [x] 1.2 Implement BDD step definitions for GetName scenarios
- [x] 1.3 Add `GetName()` method to Object (make scenarios pass)
- [x] 1.4 Add unit tests for GetName edge cases (whitespace-only name, special characters)
- [x] 1.5 Modify `OrderedPropKeys()` to always emit `name` as the first key
- [x] 1.6 Add unit test for `name` appearing first in frontmatter key ordering

## 2. Core: Object creation

- [x] 2.1 Write BDD scenario for new object having name populated from slug
- [x] 2.2 Implement BDD step definition for the scenario
- [x] 2.3 Modify `NewObject()` to set `Properties["name"]` from the slug parameter (make scenario pass)

## 3. Core: Schema validation

- [x] 3.1 Write unit test for reserved name rejection
- [x] 3.2 Add validation in `ValidateTypeSchema()` rejecting property named "name" as reserved (make test pass)

## 4. Core: Sync migration

- [x] 4.1 Write BDD scenarios for sync migration (missing name added, existing name preserved)
- [x] 4.2 Implement BDD step definitions for sync migration scenarios
- [x] 4.3 Add migration logic in `Sync()` to populate `name` from `DisplayName()` for objects missing it (make scenarios pass)

## 5. TUI: Display name property

- [x] 5.1 Update TUI list view (`tui/list.go`) to use `GetName()` instead of `DisplayName()`
- [x] 5.2 Update TUI detail title panel (`tui/detail.go`) to use `GetName()` instead of `DisplayName()`

## 6. CLI: Display name property

- [x] 6.1 Update CLI helper (`cmd/helper.go`) to use `GetName()` where `DisplayID()` is used
