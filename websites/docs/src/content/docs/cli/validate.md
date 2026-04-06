---
title: tmd type validate
description: Validate type schemas, objects, relations, and wiki-links.
sidebar:
  order: 16
---

Validates the vault's type schemas, objects, relations, wiki-links, and name uniqueness. Runs five phases in order:

```bash
tmd type validate
```

## Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--watch` | `-w` | Continuously watch for changes and re-validate |

## Watch Mode

Use `--watch` to enter continuous validation mode. The command monitors `types/`, `properties/`, and `objects/` for file changes. On each change it clears the terminal, re-syncs the index, and re-runs all five validation phases. Rapid changes are debounced (200ms) into a single validation run.

```bash
tmd type validate --watch
```

Output includes a timestamp for each validation cycle:

```
[14:32:05] Validating...

Validation passed.
```

Press `Ctrl+C` to exit watch mode.

## Validation Phases

### Phase 1: Schema Validation

Scans all `types/<name>/schema.yaml` files and checks:

- `name` field is required
- Each property must have a `name` and `type`
- `select`/`multi_select` properties must define `options`
- `relation` properties must define `target`

### Phase 2: Object Validation

Validates all object properties against their type schema:

- `select`/`multi_select` values must be in the allowed `options` list
- `number` properties must be numeric
- `relation` targets must match the expected type

### Phase 3: Relation Validation

Checks all stored relations to ensure both the source and target objects exist.

### Phase 4: Wiki-link Validation

Detects broken wiki-links — `[[target]]` references in object body content where the target object does not exist.

### Phase 5: Name Uniqueness Validation

Checks all types with `unique: true` in their schema (e.g., the built-in `tag` type) to ensure no two objects of the same type share the same `name` value.

## Output

On success:

```
Validation passed.
```

On failure, errors are grouped by phase with details:

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
