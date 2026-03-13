## Why

typemd ships with opinionated built-in types (`book`, `person`, `note`) in `defaultTypes`, making the tool appear tailored to specific use cases. Since `tag` is the only type that backs a system property (`tags`), it should be the only built-in type. Users should define all other types via `.typemd/types/*.yaml`.

## What Changes

- **BREAKING**: Remove `book`, `person`, and `note` from `defaultTypes` in `core/type_schema.go`
- Retain `tag` as the sole built-in type (it backs the `tags` system property)
- `tmd new <type>` for an undefined type now returns an error prompting the user to create a type schema first (this already works via `LoadType` returning "unknown type" — no code change needed)
- Update all tests and BDD features that implicitly depend on built-in `book`/`person`/`note` types to explicitly create type schema YAML files in test vaults

## Capabilities

### New Capabilities

_(none — this is a removal, not an addition)_

### Modified Capabilities

- `type-schema`: Built-in types reduced from 4 to 1 (`tag` only). `LoadType` behavior unchanged but fewer fallback types available.

## Impact

- **core/type_schema.go**: Remove 3 entries from `defaultTypes`
- **Tests (~197 occurrences across 18 Go test files)**: Tests using `book`/`person`/`note` types need explicit type schema YAML files in their test vaults
- **BDD features (~81 occurrences across 13 feature files)**: Same — test setup must create type schemas
- **Users**: Existing vaults using `book`/`person`/`note` without custom schemas will break — users must create `.typemd/types/book.yaml` etc. This is intentional.
