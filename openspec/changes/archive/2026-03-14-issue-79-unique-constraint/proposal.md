## Why

Multiple objects of the same type can currently share identical `name` values. For types like `tag` or `person`, duplicates cause confusion and broken references. Today, only the built-in `tag` type has hardcoded uniqueness enforcement — there is no way for user-defined types to opt in.

## What Changes

- Add `unique: true` option to type schema YAML files
- Generalize the existing tag-only name uniqueness check into a schema-driven mechanism that works for any type
- Reject object creation when a same-type object with an identical `name` property already exists (when `unique: true`)
- Add uniqueness violation detection to `tmd type validate` for all unique types
- Remove hardcoded tag uniqueness logic and replace with the generalized mechanism
- Built-in `tag` type schema gains `unique: true` (preserving existing behavior)

## Capabilities

### New Capabilities

- `unique-constraint`: Schema-level unique constraint on object names — declaration in type schema, enforcement at creation time, and validation of existing objects

### Modified Capabilities

- `type-schema`: TypeSchema struct gains a `Unique` field; schema YAML supports `unique: true/false`

## Impact

- **core/type_schema.go** — `TypeSchema` struct, YAML parsing, built-in `tag` default
- **core/object.go** — `NewObject()` creation-time uniqueness check
- **core/tag.go** — `checkTagNameUnique()` generalized or removed
- **core/validate.go** — `ValidateTagNameUniqueness()` generalized to `ValidateNameUniqueness()`
- **core/features/** — New BDD scenarios for unique constraint behavior
- No breaking changes — `unique` defaults to `false`, preserving current behavior for all existing types
