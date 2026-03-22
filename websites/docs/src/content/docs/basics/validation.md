---
title: Validation
description: How TypeMD validates type schemas, objects, relations, and links.
sidebar:
  order: 6
---

Validation checks the consistency of your vault — from type schemas to object properties to links between objects. Run it with:

```bash
tmd type validate
```

## Validation phases

Validation runs five phases in order. Each phase checks a different aspect of your vault.

### Phase 1: Schema Validation

Scans all type schemas (`.typemd/types/<name>/schema.yaml`) and `.typemd/properties.yaml` (shared properties) and checks:

- `name` field is required
- Each property must have a `name` and `type`
- `select`/`multi_select` properties must define `options`
- `relation` properties must define `target`
- Property names `description`, `created_at`, `updated_at`, and `tags` are reserved for [system properties](/basics/properties#system-properties)
- `name` in `properties` only allows a `template` field (for [name templates](/basics/templates#name-templates))

### Phase 2: Object Validation

Validates all object properties against their type schema:

- `select`/`multi_select` values must be in the allowed `options` list
- `number` properties must be numeric
- `date` must be in `YYYY-MM-DD` format
- `datetime` must include a time component
- `url` must start with `http://` or `https://`
- `checkbox` must be boolean
- `relation` targets must match the expected type

### Phase 3: Relation Validation

Checks all stored relations to ensure both the source and target objects exist.

### Phase 4: Link Validation

Detects broken links — `[[target]]` references in object body content where the target object does not exist. For shorthand targets (`[[type/name]]` or `[[name]]`) that match multiple objects, the ambiguous target is reported along with the list of matching full IDs.

### Phase 5: Name Uniqueness Validation

Checks all types with `unique: true` in their schema (e.g., the built-in `tag` type) to ensure no two objects of the same type share the same `name` value.

## Lenient validation philosophy

TypeMD uses lenient validation by design:

- **Extra properties are allowed** — Properties not defined in the schema are silently accepted. This lets you add ad-hoc metadata without updating the schema.
- **Missing properties do not cause errors** — An object doesn't need to have every property defined in its type schema.
- **Only schema-defined properties are type-checked** — If a property is in the schema, its value must match the declared type. If it's not in the schema, no validation occurs.

This approach keeps TypeMD flexible while still catching real mistakes — like a `select` value that isn't in the options list, or a `relation` pointing to a non-existent object.

## Output

On success:

```
Validation passed.
```

On failure, errors are grouped by phase:

```
Schema errors:
  book.yaml: property "status": select type requires non-empty options

Object errors:
  book/example: property "rating": expected number, got "abc"

Relation errors:
  book/example author person/missing: target object not found

Wiki-link errors:
  book/example: broken wiki-link [[person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj]]

Name uniqueness errors:
  duplicate tag name "golang": tag/golang-01abc and tag/golang-01xyz

found 5 validation error(s)
```

## See also

- [Properties](/basics/properties) — property types and validation rules
- [Tags](/basics/tags#tag-uniqueness) — tag name uniqueness
- [tmd type validate](/tui/validate) — CLI reference
