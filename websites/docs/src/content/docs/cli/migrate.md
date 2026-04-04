---
title: tmd migrate
description: Migrate type schemas and objects.
sidebar:
  order: 8.5
---

Migrates type schemas and objects. The command has two modes depending on whether a type argument is provided.

```bash
tmd migrate                          # migrate schemas (enum → select)
tmd migrate --dry-run                # preview schema migrations
tmd migrate book                     # migrate book objects to match schema
tmd migrate book --dry-run
tmd migrate book --rename old_field:new_field
```

## Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview changes without modifying files |
| `--rename old:new` | Rename a property (repeatable, only with type argument) |

## Schema migration (no type argument)

When run without arguments, `tmd migrate` scans all type schemas (`types/<name>/schema.yaml`) and converts the legacy `enum` property type (with `values`) to the current `select` type (with `options`).

```bash
tmd migrate
#   book: converted enum → select for [status]
# Schema migration complete: 1 type(s) updated.
```

If all schemas are already up to date, no changes are made.

## Object migration (with type argument)

When a type name is provided, `tmd migrate <type>` updates all objects of that type to match the current schema:

1. **Rename** — if `--rename` is provided, moves the old property value to the new name
2. **Add** — properties in the schema but missing from the object are added with their default value
3. **Remove** — properties in the object but not in the schema are removed

Objects that already match the schema are skipped.

### Examples

Add a new `isbn` property to all books after updating the schema:

```bash
tmd migrate book
#   book/clean-code: added [isbn]
#   book/effective-go: added [isbn]
# Migration complete: 2 object(s) updated.
```

Rename `status` to `reading_status` (update the schema first, then migrate):

```bash
tmd migrate book --rename status:reading_status
```

Preview changes before applying:

```bash
tmd migrate book --dry-run
```

> **Note:** The `--rename` flag requires that the new name exists in the current schema and the old name does not. Update the schema before running the migration.
