## 1. Core: BDD scenarios for GetProperty

- [x] 1.1 Write BDD scenarios for GetProperty resolving object_type (book, person types)
- [x] 1.2 Write BDD scenarios for GetProperty resolving stored properties (present, missing)

## 2. Core: Implement GetProperty

- [x] 2.1 Implement step definitions for GetProperty BDD scenarios
- [x] 2.2 Add `GetProperty(name string) (any, bool)` method on `Object` — resolve `object_type` from `Object.Type`, fall back to `Properties` map for stored properties
- [x] 2.3 Add unit tests for GetProperty edge cases (empty type, nil properties map)

## 3. Core: Verify existing enforcement

- [x] 3.1 Verify existing BDD scenarios for SetProperty rejection and frontmatter stripping still pass (no new code needed, just confirmation)
