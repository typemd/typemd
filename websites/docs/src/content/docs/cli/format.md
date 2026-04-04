---
title: tmd format
description: Normalize object frontmatter and schema YAML formatting.
sidebar:
  order: 10.5
---

Formats all object and schema files with canonical property ordering and YAML style. Similar to `gofmt`, it enforces a consistent file format.

```bash
tmd format                    # format all objects and schemas
tmd format --type book        # format only book objects and schema
tmd format --dry-run          # list files that need formatting (CI mode)
```

## Flags

| Flag | Description |
|------|-------------|
| `--type <name>` | Format only objects and schemas of a specific type |
| `--dry-run` | List files that need formatting without modifying them |

## What it does

### Property ordering

Frontmatter properties are rewritten in canonical order:

1. **System properties** — `name`, `description`, `created_at`, `updated_at`, `tags`
2. **Schema-defined properties** — in the order defined in the type schema
3. **Extra properties** — any properties not in the schema, sorted alphabetically

### YAML normalization

YAML formatting is normalized to `yaml.v3` defaults — consistent quoting, indentation, and value representation.

### Schema formatting

Type schema files (`types/<name>/schema.yaml`) are also reformatted by round-tripping through the canonical serializer.

### What is preserved

- **Body content** is not modified
- **`updated_at`** is not changed — formatting is a pure layout change with no semantic modification

## Dry-run mode

With `--dry-run`, the command lists files that would be changed and exits with code 1 if any files need formatting. This is useful for CI pipelines:

```bash
tmd format --dry-run || echo "Files need formatting"
```

### Exit code

- `0` — all files are already formatted
- `1` — one or more files need formatting (dry-run), or an error occurred
