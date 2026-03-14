## Why

Every `tmd object create` produces an empty body with only schema defaults. Users who frequently create structured objects (e.g., book reviews, meeting notes, weekly journals) must manually fill in the same boilerplate content every time. Object templates let types ship pre-built content structures, reducing friction and encouraging consistent formatting.

## What Changes

- New `templates/<type>/` directory at project root for storing per-type templates as regular Markdown files (frontmatter + body)
- `tmd object create` loads and applies templates during object creation:
  - 0 templates → current behavior (empty body)
  - 1 template → auto-apply
  - Multiple templates → interactive selection or explicit `-t <name>` flag
- Template frontmatter properties override schema defaults for the new object
- Template body becomes the initial content of the new object
- `SystemProperty` struct gains an `Immutable` field to distinguish auto-managed properties (`created_at`, `updated_at`) from user-authored ones (`name`, `description`, `tags`). Templates can only override user-authored system properties.
- Placeholder/variable substitution in templates is explicitly out of scope (separate issue)

## Capabilities

### New Capabilities

- `object-templates`: Template discovery, loading, selection, and application during object creation

### Modified Capabilities

- `system-property-registry`: Add `Immutable` field to `SystemProperty` struct to distinguish user-authored vs auto-managed properties

## Impact

- **core/**: `SystemProperty` struct change, new template loading/application logic in `NewObject`, new `templates/` directory convention
- **cmd/**: `tmd object create` gains `-t` / `--template` flag, interactive template selection when multiple templates exist
- **Existing objects**: No impact — templates only affect new object creation
- **Vault structure**: New `templates/` directory at project root (optional, no migration needed)
