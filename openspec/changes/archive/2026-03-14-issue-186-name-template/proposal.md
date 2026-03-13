## Why

Object names are always manually provided at creation time. For types with predictable naming patterns (e.g., daily journals, meeting notes), this is repetitive and error-prone. Allowing type schemas to define a name template reduces friction and enables one-command object creation.

## What Changes

- Allow the `name` system property to appear in the type schema `properties` array with a `template` field (only `template` is permitted — no type redefinition)
- Add template evaluation at object creation time, supporting `{{ date:FORMAT }}` placeholders with user-friendly format syntax (e.g., `YYYY-MM-DD`)
- Make the `name` CLI argument optional when the type has a name template; still required otherwise
- Users can override the template-generated name by providing an explicit name argument

## Capabilities

### New Capabilities

- `name-template`: Template-based name generation for objects at creation time, including template syntax, date placeholder evaluation, and format conversion

### Modified Capabilities

- `type-schema`: Relax validation to allow `name` in `properties` array with only the `template` field permitted
- `name-property`: Add scenario for template-generated names on object creation

## Impact

- `core/type_schema.go` — Property struct gains `Template` field; `ValidateSchema()` relaxed for `name` with template-only constraint
- `core/object.go` — `NewObject()` evaluates template when name argument is empty
- `cmd/create.go` — `name` argument becomes optional (conditional on type schema)
- Existing type schemas are unaffected (no breaking changes)
