## Context

The computed system property infrastructure (#371) registers `object_type` as a computed property with schema validation and write rejection. However, there is no runtime resolution — no way to _read_ `object_type` through the property API. The value is already available as `Object.Type`, but consumers must know to access the struct field directly instead of using the property system.

## Goals / Non-Goals

**Goals:**
- Provide a unified `GetProperty(name)` method on `Object` that resolves both stored (frontmatter) and computed properties
- `object_type` returns `Object.Type` — the type derived from the file path convention `objects/<type>/<name>.md`
- Establish the resolver pattern for future computed properties (`links`, `backlinks`, `created_by`, `updated_by`)

**Non-Goals:**
- Implementing other computed properties (separate issues)
- Adding `object_type` to query result projections (already available via `ObjectResult.Type`)
- Displaying `object_type` in TUI property panel (already shown via sidebar grouping)

## Decisions

### 1. Add `GetProperty()` to `Object` instead of a standalone resolver

**Decision:** Add `func (o *Object) GetProperty(name string) (any, bool)` that checks computed properties first, then falls back to `o.Properties[name]`.

**Rationale:** Property access should be uniform. Consumers shouldn't need to know whether a property is stored or computed — they call `GetProperty()` and get the value. This follows the existing pattern where `GetName()` already encapsulates name resolution logic.

**Alternatives considered:**
- A separate `ComputedPropertyResolver` service — adds unnecessary indirection for values already on the entity.
- A map-based registry of resolver functions — over-engineering for the current scope. Can be introduced later if computed properties need external data (e.g., `links` needing index queries).

### 2. Computed properties resolved within `Object` entity

**Decision:** `object_type` is resolved from `Object.Type` directly. No external dependencies needed.

**Rationale:** `Object.Type` is always populated when an object is loaded (derived from file path by the repository). This is a pure data derivation with zero cost, so it belongs on the entity, not a service.

### 3. Return type is `(any, bool)` following Go map access pattern

**Decision:** `GetProperty` returns `(value any, exists bool)` to distinguish between "property not set" and "property set to zero value".

**Rationale:** Consistent with `map[string]any` access pattern. For computed properties, `exists` is always `true` since they're always derivable.

## Risks / Trade-offs

- **[Risk] Future computed properties may need external data** (e.g., `links` needs the index) → `GetProperty()` handles `object_type` inline; future properties can extend the method or introduce a resolver interface when needed. No premature abstraction.
- **[Risk] Consumers bypassing `GetProperty()` and reading `Properties` map directly** → This is acceptable for stored properties. Document `GetProperty()` as the canonical way to access any property including computed ones.
