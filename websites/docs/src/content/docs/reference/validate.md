---
title: tmd validate
description: Validate type schemas, objects, and relations.
sidebar:
  order: 10
---

Validates the vault's type schemas, objects, and relations. Runs three phases in order:

```bash
tmd validate
```

## Validation Phases

### Phase 1: Schema Validation

Scans all `.typemd/types/*.yaml` files and checks:

- `name` field is required
- Each property must have a `name` and `type`
- `enum` properties must define `values`
- `relation` properties must define `target`

### Phase 2: Object Validation

Validates all object properties against their type schema:

- `enum` values must be in the allowed `values` list
- `number` properties must be numeric
- `relation` targets must match the expected type

### Phase 3: Relation Validation

Checks all stored relations to ensure both the source and target objects exist.

## Output

On success:

```
Validation passed.
```

On failure, errors are grouped by phase with details:

```
Schema errors:
  book.yaml: property "status": enum must define values

Object errors:
  book/example: property "rating": expected number, got "abc"

Relation errors:
  book/example author person/missing: target object not found

found 3 validation error(s)
```
