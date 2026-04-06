## Why

The `object-lock` OpenSpec spec requires that `SetPropertyMultiple` on a locked object returns an error, but `SetPropertyMultiple` does not exist in the codebase. Only `SetProperty` (single property) is implemented. Adding a batch property setter enables atomic multi-property updates and fulfills the spec requirement.

## What Changes

- Add `ObjectService.SetPropertyMultiple(id string, props map[string]any)` method with lock check, computed property rejection, and schema validation
- Add `Vault.SetPropertyMultiple()` facade method delegating to `ObjectService`
- Dispatch a single `ObjectSaved` event for the batch update (not one per property)

## Capabilities

### New Capabilities

(none — this extends the existing object-lock capability)

### Modified Capabilities

- `object-lock`: The existing spec already defines `SetPropertyMultiple` behavior. Implementation now fulfills that requirement — no spec change needed.

## Impact

- **core/object_service.go** — new `SetPropertyMultiple` method
- **core/object.go** (Vault facade) — new delegation method
- **core/features/object_lock.feature** — new BDD scenario for `SetPropertyMultiple` on locked object
- **core/object_service_test.go** or **core/system_property_test.go** — unit tests for edge cases
- No breaking changes to existing API
