## Context

`tmd object create` currently produces objects with an empty body and only schema-default property values. Users who create structured objects (book reviews, meeting notes) must manually add boilerplate content every time. The `NewObject` method in `core/object.go` hardcodes an empty body (`""`) and initializes properties from schema defaults only.

Templates are stored as regular Markdown files at `templates/<type>/` in the vault root. Each template has frontmatter (property overrides) and body (initial content). The system property registry (`SystemProperty` struct) currently lacks the ability to distinguish user-authored from auto-managed properties, which is needed to determine which system properties a template can override.

## Goals / Non-Goals

**Goals:**

- Allow types to define one or more templates as Markdown files
- Apply templates during `tmd object create` (auto, interactive, or explicit)
- Extend `SystemProperty` with an `Immutable` field to distinguish auto-managed properties
- Template frontmatter overrides schema defaults for mutable system properties and schema properties

**Non-Goals:**

- Placeholder/variable substitution in templates (separate issue)
- Template inheritance or composition
- Template validation via `tmd type validate` (future enhancement)
- TUI template selection (future enhancement)

## Decisions

### 1. Template storage: `templates/<type>/` at vault root

**Choice:** Store templates at `templates/<type>/<name>.md` in the vault root directory.

**Alternatives considered:**
- `.typemd/templates/<type>/` — keeps templates in the config area, but templates are user-facing content, not internal config
- `objects/<type>/_templates/` — co-locates with objects, but pollutes object scanning and requires filtering logic
- Inline in type schema YAML — limits templates to property defaults, cannot include body content

**Rationale:** Templates are user-created content (like objects), not internal metadata (like type schemas). Placing them at the vault root makes them discoverable. The `templates/` directory mirrors `objects/` as a top-level content directory. No impact on existing object scanning logic.

### 2. Template file format: standard Markdown

**Choice:** Templates use the same format as object files — YAML frontmatter + Markdown body.

**Rationale:** No new format to learn. Users can copy an existing object as a starting template. The existing `parseFrontmatter` function can parse templates.

### 3. Template selection: three modes

**Choice:**
- 0 templates → current behavior (empty body, schema defaults only)
- 1 template → auto-apply without prompting
- Multiple templates + no `-t` flag → interactive selection (list template names)
- `-t <name>` flag → explicit template by name (filename without `.md`)

**Rationale:** Auto-apply for single templates eliminates friction. Interactive selection for multiple templates avoids requiring users to remember template names. The `-t` flag supports scripting and power users.

### 4. Property merge order

**Choice:** Properties are merged in this order (later wins):
1. Schema defaults (`Property.Default`)
2. Template frontmatter values
3. Auto-managed system properties (`name`, `created_at`, `updated_at`)

**Rationale:** Template values override schema defaults (that's the point of templates). Auto-managed system properties always win because they must reflect actual creation state. User-authored system properties (`name`, `description`, `tags`) from templates are applied at step 2 and preserved (not overwritten at step 3).

### 5. SystemProperty.Immutable field

**Choice:** Add `Immutable bool` to `SystemProperty` struct. Properties with `Immutable: true` (`created_at`, `updated_at`) cannot be overridden by templates. Properties with `Immutable: false` (`name`, `description`, `tags`) can be.

**Rationale:** A single boolean is sufficient — the distinction is binary (auto-managed vs user-authored). Using a field on the registry struct makes the rule declarative and easy to query when applying templates.

### 6. Core API: NewObject gains templateName parameter

**Choice:** Extend `NewObject(typeName, name string)` to `NewObject(typeName, name, templateName string)`. When `templateName` is non-empty, load and apply the template. When empty, behave as today.

**Rationale:** Adding a parameter keeps the API simple. The CLI layer handles template discovery and selection, then passes the chosen template name to core. This keeps template selection logic (interactive prompts, `-t` flag) in `cmd/` where it belongs.

### 7. Template path helpers on Vault

**Choice:** Add `TemplatesDir()`, `TypeTemplatesDir(typeName)`, and `TemplatePath(typeName, templateName)` methods to `Vault`, following the existing pattern of `ObjectsDir()`, `ObjectDir()`, `ObjectPath()`.

**Rationale:** Consistent with existing Vault path conventions. Centralizes path logic.

## Risks / Trade-offs

- **Template frontmatter with unknown properties** → Silently ignored. Template properties not in the schema are dropped during property initialization, same as how extra frontmatter in objects is handled by `SyncIndex`. This is acceptable for MVP.
- **Template body references non-existent objects** → Not validated. Wiki-links in template body are treated as plain text until the object is synced. Acceptable for MVP.
- **Interactive selection blocks scripting** → Mitigated by the `-t` flag which skips interactive selection. If neither `-t` nor a single template is available, the command errors in non-TTY contexts rather than hanging.
