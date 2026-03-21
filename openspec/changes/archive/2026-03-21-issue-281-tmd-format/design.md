## Context

typemd objects are Markdown files with YAML frontmatter. When created or saved via `ObjectService`, properties are ordered using `OrderedPropKeys` (system → schema → extras alphabetically) and serialized via `writeFrontmatter` (per-key `yaml.Marshal`). However, externally-edited files or files predating a schema reorder retain their original formatting. There is no batch reformat command.

Type schemas are serialized via `MarshalTypeSchema` which uses a Go struct with `yaml` tags, producing consistent output. But existing schema files may have been hand-edited with different formatting.

## Goals / Non-Goals

**Goals:**
- Canonical formatting for all object files (frontmatter property order + YAML style)
- Canonical formatting for all type schema files (YAML re-serialization)
- CI-friendly dry-run mode (exit code 1 when files need formatting)
- Filter by type name for targeted formatting
- Reuse existing serialization infrastructure (`OrderedPropKeys`, `writeFrontmatter`, `MarshalTypeSchema`)

**Non-Goals:**
- Formatting Markdown body content (body is preserved byte-for-byte)
- Formatting view YAML files, shared properties, or config files
- Custom YAML style options (always uses `yaml.v3` defaults)
- Updating `updated_at` timestamp during formatting

## Decisions

### 1. Core-layer `FormatResult` pattern

**Decision:** Add `FormatObjects(typeName string, dryRun bool)` and `FormatSchemas(dryRun bool)` methods to Vault, returning a `FormatResult` with changed file paths.

**Rationale:** Follows the same pattern as `MigrateObjects` / `MigrateSchemas` — core logic in the domain layer, CLI is a thin presentation wrapper. Enables future use from MCP/TUI.

**Alternative considered:** Putting logic directly in `cmd/format.go` RunE. Rejected because it violates the existing layered architecture and prevents reuse.

### 2. Byte comparison for change detection

**Decision:** Compare the serialized output bytes against the original file bytes. Only write if different.

**Rationale:** Simple, reliable, and avoids false positives. Since `yaml.v3` is deterministic for the same input, identical files won't be rewritten. This also naturally handles the YAML normalization — any quoting/indentation differences will show as byte differences.

### 3. Object formatting via read-reserialize round-trip

**Decision:** For each object: `parseFrontmatter` → `OrderedPropKeys(props, schema)` → `writeFrontmatter(props, body, keyOrder)` → compare with original bytes → write if different.

**Rationale:** Reuses 100% of existing infrastructure. No new serialization logic needed.

### 4. Schema formatting via load-remarshal round-trip

**Decision:** For each schema: `GetSchema(name)` → `MarshalTypeSchema(schema)` → compare with file bytes → write if different.

**Rationale:** `MarshalTypeSchema` already produces canonical output via struct-tagged serialization. Round-tripping through it normalizes any formatting differences.

**Caveat:** `GetSchema` resolves `use` entries and strips the `name` property (converting it to `NameTemplate`). `MarshalTypeSchema` reverses this. The round-trip should be lossless, but schema files with comments will lose them (acceptable trade-off since YAML comments are not part of the data model).

### 5. Schema iteration via `ListSchemaNames` + `schemaPath`

**Decision:** Add a `ListSchemaNames()` method to `ObjectRepository` that returns all schema names by walking `.typemd/types/`. Use existing `schemaPath` to locate each file.

**Rationale:** No existing method walks all schema files. `ListTypes` on Vault returns names but includes built-in types without files. Need a repo-level method that only returns file-backed schemas.

### 6. `--type` flag applies only to objects, not schemas

**Decision:** When `--type book` is specified, only format objects of type `book`. All schemas are always formatted (or none, if `--type` is given — since schema formatting is already fast and orthogonal).

**Revised:** `--type` filters both objects and schemas. When `--type book` is specified, format only `book` objects and the `book` schema file.

**Rationale:** Users filtering by type likely want targeted formatting. Formatting all schemas when they asked for one type would be surprising.

## Risks / Trade-offs

- **YAML comment loss** → Schema files with hand-written YAML comments will lose them after formatting. Mitigation: This is expected behavior for a formatter (same as `gofmt` dropping non-doc comments). Document in command help text.
- **Frontmatter library parsing differences** → The `adrg/frontmatter` parser and `yaml.v3` marshaler may handle edge cases differently (e.g., multi-line strings, special characters). Mitigation: The same libraries are used throughout typemd, so any object that was created/saved by typemd will round-trip cleanly. Edge cases only affect hand-crafted files.
- **Large vault performance** → Walking all files is O(n). Mitigation: This is a batch operation run infrequently; same as `gofmt ./...`. No index needed.
