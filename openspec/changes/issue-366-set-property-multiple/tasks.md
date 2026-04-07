## 1. BDD: SetPropertyMultiple lock guard

- [x] 1.1 Add BDD scenario "SetPropertyMultiple on locked object returns error" to `core/features/object_lock.feature`
- [x] 1.2 Implement step definition for SetPropertyMultiple BDD scenario

## 2. Core: Implement SetPropertyMultiple

- [x] 2.1 Add `ObjectService.SetPropertyMultiple(id string, props map[string]any) error` in `core/object_service.go`
- [x] 2.2 Add `Vault.SetPropertyMultiple` facade method in `core/object.go`

## 3. Unit tests: Edge cases

- [x] 3.1 Add unit tests for SetPropertyMultiple (empty map, computed property rejection, schema validation error, unlocked object success, multiple properties applied atomically)
