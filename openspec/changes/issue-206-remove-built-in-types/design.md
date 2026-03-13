## Context

`defaultTypes` in `core/type_schema.go` contains 4 built-in types: `book`, `person`, `note`, `tag`. The first three are opinionated examples that should be user-defined. Only `tag` is structurally required (it backs the `tags` system property).

Tests and BDD features extensively use `book`/`person`/`note` as test fixtures, relying on `LoadType` falling back to `defaultTypes`. After removal, tests must create explicit type schema YAML files.

## Goals / Non-Goals

**Goals:**
- Remove `book`, `person`, `note` from `defaultTypes`
- Keep all existing tests passing by providing explicit type schemas in test fixtures
- Maintain `tag` as the sole built-in type

**Non-Goals:**
- Changing `LoadType` behavior or error messages
- Adding a `tmd type init` scaffolding command (separate issue)
- Migration tooling for existing user vaults

## Decisions

### 1. Test fixture strategy: helper function creates type schema YAML files

**Decision**: Add a test helper (e.g., `writeTestTypeSchema`) that writes `.typemd/types/<name>.yaml` files into the test vault directory. Each test that uses `book`/`person`/`note` calls this helper during setup.

**Alternatives considered**:
- Replace all test types with `tag` — impractical, many tests need non-tag types with specific properties
- Keep a `testTypes` map for tests only — would mask the real `LoadType` behavior, tests wouldn't catch regressions

### 2. Preserve existing type schemas in test fixtures

**Decision**: The test helper writes schemas matching the current `defaultTypes` definitions (same properties, emoji, options). This ensures tests don't change behavior, only how the type is loaded.

### 3. BDD step: reusable step for type schema creation

**Decision**: Add a BDD step like `Given a type "<name>" exists` or similar that creates the YAML file. This keeps feature files clean and avoids duplicating schema setup in every scenario.

## Risks / Trade-offs

- **Risk**: Missing a test that relies on built-in types → CI catches it immediately since tests fail
- **Risk**: Breaking existing user vaults → Intentional and documented in release notes. Users must create `.typemd/types/book.yaml` etc.
- **Trade-off**: Test helper duplicates the old `defaultTypes` schemas → acceptable, these are test fixtures not production code
