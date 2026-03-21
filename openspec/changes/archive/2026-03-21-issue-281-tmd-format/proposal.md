## Why

Objects are formatted correctly only when created or saved through typemd (`ObjectService.Create`/`Save`), which calls `OrderedPropKeys` to sort properties. However, files edited externally or files written before a schema reorder retain their original property order and YAML style. There is no way to batch-reformat existing files. Similar to `go fmt` / `gofmt`, typemd needs a canonical formatting command.

## What Changes

- Add `tmd format` CLI command that rewrites object Markdown files with normalized frontmatter
- Support `--type <name>` flag to format only objects of a specific type
- Support `--dry-run` flag that lists files needing formatting without writing (exit code 1 for CI)
- Also format type schema YAML files (`.typemd/types/`) via `MarshalTypeSchema` round-trip
- Property ordering follows existing `OrderedPropKeys` logic: system properties first, then schema-defined order, then extras alphabetically
- YAML normalization via `yaml.v3` re-serialization (consistent quoting, indentation, null handling)
- Body content is preserved unchanged
- `updated_at` is NOT modified during formatting (pure formatting change, no semantic modification)

## Capabilities

### New Capabilities
- `object-format`: Batch formatting of object frontmatter with property reordering and YAML normalization, including dry-run mode for CI verification

### Modified Capabilities

(none)

## Impact

- **New files**: `cmd/format.go` (CLI command), `core/format.go` (core logic)
- **Core package**: New `FormatObjects` / `FormatSchemas` methods on Vault, reusing existing `OrderedPropKeys`, `writeFrontmatter`, `MarshalTypeSchema`
- **No breaking changes**: Purely additive command
- **No dependency changes**: Uses existing `gopkg.in/yaml.v3`
