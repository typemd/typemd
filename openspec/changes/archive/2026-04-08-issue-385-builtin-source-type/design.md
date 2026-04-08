## Context

typemd has two built-in types: `tag` (backs the `tags` system property, unique) and `page` (general-purpose container). Built-in types are registered in the `defaultTypes` map in `core/type_schema.go`, exist without YAML files, cannot be deleted, but can be overridden by custom `types/<name>/schema.yaml`.

The ingest workflow (#378) needs a standard type for tracking raw materials. Rather than requiring users to define it manually, `source` becomes the third built-in type.

## Goals / Non-Goals

**Goals:**
- Register `source` as a built-in type following the exact same pattern as `tag` and `page`
- Include properties for tracking provenance: `url`, `author`, `ingested_at`
- Ensure deletion protection, override support, and schema loading work identically to other built-in types

**Non-Goals:**
- Ingest workflow implementation (#378)
- MCP write tools (#377)
- Relation from objects back to sources (can be added later via a system property like `tags`)

## Decisions

**Follow the existing pattern exactly.** No new infrastructure is needed. The `source` type is added to `defaultTypes` with three custom properties (`url`, `author`, `ingested_at`). The existing deletion protection in `DeleteType()` and fallback loading in `GetSchema()` handle it automatically.

**`unique: false`** (the default). Unlike `tag`, sources are not expected to be deduplicated by name — the same URL could be ingested at different times with different results.

**Property types:** `url` and `author` are `string` type (text fields). `ingested_at` is `date` type to leverage the existing date property infrastructure (date picker, format config).

## Risks / Trade-offs

- [Minimal risk] The change is purely additive. No existing behavior is modified. All existing built-in type infrastructure handles it automatically.
- [Trade-off] Not adding a system relation property (like `tags` → `tag`). This keeps the change small; a `sources` system property can be added separately if needed.
