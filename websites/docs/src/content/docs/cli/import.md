---
title: tmd import
description: Import external markdown files into the vault.
sidebar:
  order: 25
---

Import external markdown files into the vault using a three-step workflow: scan, plan, execute.

## Scan sources

```bash
tmd import scan <paths...>
```

Scan source directories or files for markdown content, extracting frontmatter patterns and collecting file statistics.

### Example

```bash
tmd import scan ~/notes ~/docs/blog
```

Outputs a JSON `ScanResult` including:

- `sources` — each markdown file with path, size, and frontmatter keys
- `file_count` — total markdown files found
- `directories` — directory structure with per-directory file counts
- `patterns` — aggregate frontmatter key statistics (frequency and sample values)
- `existing_types` — vault type schemas with their properties
- `no_frontmatter_count` — files without YAML frontmatter

## Generate a plan

```bash
tmd import plan <classifications-file>
tmd import plan classifications.json --output plan.json
```

Generate an import plan from a JSON file containing object classifications. The classifications file is a JSON array of `ObjectPlan` entries, typically produced by an AI analyzing a scan result.

### Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Write plan to file instead of stdout |

### ObjectPlan format

```json
[
  {
    "source_path": "books/clean-code.md",
    "type_name": "book",
    "name": "Clean Code",
    "properties": {"author": "Robert C. Martin"},
    "body": "Book content here...",
    "conflict": "none",
    "depends_on": []
  }
]
```

The plan command detects new types that need to be created, checks for conflicts with existing objects, and computes dependency-ordered import sequence (tags first, then independent objects, then dependent objects).

## Execute a plan

```bash
tmd import execute <plan-file>
```

Execute an import plan to create types and objects in the vault.

### Execution phases

1. **Create types** — new type schemas listed in the plan
2. **Create objects** — in dependency order, respecting conflict flags (`skip` or `none`)
3. **Reconcile** — resolve wiki-links across all imported objects

### Output

Returns a JSON `ImportReport`:

```json
{
  "types_created": 1,
  "objects_created": 8,
  "objects_skipped": 2,
  "objects_failed": 0,
  "details": [...],
  "unresolved_refs": [],
  "suggestions": []
}
```

## AI-assisted workflow

The import commands are designed for AI orchestration via the `onboarding` skill:

1. AI runs `tmd import scan` to analyze source files
2. AI classifies files into types and maps properties
3. AI constructs the classifications JSON and generates a plan
4. User reviews and approves the plan
5. AI runs `tmd import execute` with the approved plan

Use `tmd instructions onboarding` to get the full onboarding skill instructions with vault context.
