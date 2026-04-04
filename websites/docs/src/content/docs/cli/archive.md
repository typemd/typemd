---
title: tmd object archive / unarchive
description: Archive or unarchive objects (soft delete).
sidebar:
  order: 3.6
---

Archive an object to hide it from default queries without deleting the file. The object remains on disk and can be unarchived at any time.

Supports prefix matching — you can omit the ULID suffix if the prefix uniquely identifies an object.

## Archive

```bash
tmd object archive book/clean-code
```

Sets `archived: true` in the object's frontmatter. The object will no longer appear in `tmd object list` or other queries by default.

If the object is already archived, the command prints a message and exits successfully (no-op).

Archiving a locked object is allowed — the archive flag bypasses the lock guard.

## Unarchive

```bash
tmd object unarchive book/clean-code
```

Removes the `archived` property from the object's frontmatter. The object will appear in queries again.

If the object is not archived, the command prints a message and exits successfully (no-op).

## Behavior

- Archived objects are excluded from `tmd object list` and all default queries
- Use `tmd object list --include-archived` to include archived objects
- `tmd object show <id>` still works on archived objects (direct access by ID is unaffected)
- Wiki-links and relations to archived objects still resolve
- Files stay in `objects/` — archiving is purely a metadata flag
