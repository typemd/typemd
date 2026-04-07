## Context

`ObjectService.SetProperty` handles single property updates with lock check, computed property rejection, and schema validation. The `object-lock` OpenSpec spec requires `SetPropertyMultiple` to exist and respect the lock guard. Currently no batch property setter exists — callers must loop `SetProperty` calls, which means N file writes and N events for N properties.

## Goals / Non-Goals

**Goals:**
- Implement `ObjectService.SetPropertyMultiple` that updates multiple properties in a single operation
- Lock check: reject the entire batch if the object is locked
- Computed property check: reject the batch if any key is a computed system property
- Schema validation: validate each property against the type schema
- Single file write and single `ObjectSaved` event for the batch
- Expose via `Vault.SetPropertyMultiple` facade

**Non-Goals:**
- Partial application (apply some properties if others fail) — the batch is all-or-nothing
- CLI command for batch property setting (no consumer needs it yet)
- TUI integration (TUI uses `SetProperty` for inline editing)

## Decisions

**All-or-nothing semantics:** Validate all properties before writing any. If any property fails validation (computed, schema error), return an error without modifying the object. This prevents partial state.

**Single event dispatch:** Emit one `ObjectSaved` event after the batch write, not one per property. This matches `SaveObject` behavior and avoids N projector updates.

**Method signature:** `SetPropertyMultiple(id string, props map[string]any) error` — mirrors the existing `SetProperty(id, key string, value any) error` pattern but takes a map.

**Empty map:** Calling with an empty map is a no-op that returns nil (no file write, no event).

## Risks / Trade-offs

- **[Risk] Property order in map iteration** → Each property is set via `Object.SetProperty` which handles ordering internally. Map iteration order doesn't affect the persisted frontmatter order.
- **[Risk] Error message ambiguity with multiple failures** → Return error on first failed property. The error message includes the property key, so the caller knows which one failed.
