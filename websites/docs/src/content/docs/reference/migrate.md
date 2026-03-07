---
title: tmd migrate
description: Update objects to match the current type schema.
sidebar:
  order: 8.5
---

Migrates all objects of a given type to match the current schema. Adds missing properties (with defaults), removes obsolete ones, and optionally renames properties.

```bash
tmd migrate book
tmd migrate book --dry-run
tmd migrate book --rename old_field:new_field
```

## Flags

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview changes without modifying files |
| `--rename old:new` | Rename a property (repeatable) |

## What it does

For each object of the specified type:

1. **Rename** — if `--rename` is provided, moves the old property value to the new name
2. **Add** — properties in the schema but missing from the object are added with their default value
3. **Remove** — properties in the object but not in the schema are removed

Objects that already match the schema are skipped.

## Examples

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
