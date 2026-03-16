## Context

`tmd object create <type> [name]` requires a type as the first positional argument. Users wanting quick capture must remember type names and format names as slugs. Issue #236 originally proposed a `tmd quick` command, but exploration revealed it would duplicate `tmd object create`. Instead, we enhance the existing command with a vault config default type and automatic slug conversion.

Currently there is no vault-level configuration system. Settings are either file-based (type schemas, shared properties) or runtime-only. TUI session state (`.typemd/tui-state.yaml`) stores UI preferences but is TUI-specific.

## Goals / Non-Goals

**Goals:**
- Allow `tmd object create` to work without specifying a type (read from config)
- Accept natural-language names and auto-convert to valid slugs
- Introduce a vault config system that other features can build on
- Maintain full backward compatibility

**Non-Goals:**
- AI-powered auto-fill (`--ai` flag) — future work
- Migrating TUI preferences from `tui-state.yaml` to config — separate issue
- Interactive name prompt via Bubble Tea when no name is given — deferred to keep scope small
- Config validation or `tmd config` management commands

## Decisions

### 1. Vault config file at `.typemd/config.yaml`

Config lives in the `.typemd/` directory alongside other vault metadata. Structure uses interface-layer namespacing:

```yaml
cli:
  default_type: idea
```

**Rationale:** Grouping by interface (cli, tui, web, mcp) keeps concerns separated and avoids naming collisions. The file is optional — missing file means empty/default config.

**Alternatives considered:**
- Flat keys (`default_type: idea`) — doesn't scale, no namespace
- Root-level `vault:` key — unnecessary nesting for now

### 2. Config struct is minimal

```go
type VaultConfig struct {
    CLI CLIConfig `yaml:"cli"`
}

type CLIConfig struct {
    DefaultType string `yaml:"default_type"`
}
```

Config is loaded during `Vault.Open()` and stored on the Vault struct. Missing file or missing keys result in zero values (no error). The struct can grow organically as new features need configuration.

### 3. Slug conversion in core layer

A `Slugify(name string)` function in `core/` handles conversion:
- Lowercase
- Replace spaces and underscores with hyphens
- Remove non-alphanumeric characters (except hyphens)
- Collapse consecutive hyphens
- Trim leading/trailing hyphens

Applied in `ObjectService.Create()` so all creation paths (CLI, TUI, MCP) benefit. The function is idempotent — already-slugified input passes through unchanged.

**Rationale:** Core-layer conversion means every consumer gets consistent behavior. No risk of one path forgetting to slugify.

**Alternatives considered:**
- CLI-layer only — requires each consumer to remember to slugify
- External library — unnecessary dependency for simple string transforms

### 4. Name property preserves original input

When `ObjectService.Create("idea", "Some Thought", "")` is called:
- Slug for ObjectID/filename: `Slugify("Some Thought")` → `some-thought`
- `name` property in frontmatter: `"Some Thought"` (original input)

This means the display name can differ from the filename slug, which is already possible when users edit the `name` property after creation.

**Rationale:** Users expect their input to appear as the display name, not a slugified version. The slug is an implementation detail of the file system.

### 5. Type argument resolution with smart fallback

The `tmd object create` command changes from `RangeArgs(1, 2)` to `RangeArgs(0, 2)`.

Resolution logic for positional arguments:
- **0 args**: type from `--type` flag or config `cli.default_type`; name empty (requires name template on schema)
- **1 arg**: attempt to resolve as type name → if valid type, treat as type (backward compatible); if not a valid type AND default type is available (flag or config), treat as name
- **2 args**: first is type, second is name (current behavior, unchanged)

**Rationale:** This preserves full backward compatibility while enabling the new `tmd object create "Some Thought"` pattern when a default type is configured.

### 6. `--type` flag without `-t` short form

The `-t` short flag is already taken by `--template`. The new type override uses `--type` only (no short flag).

```bash
tmd object create --type note "Meeting Notes"
```

**Rationale:** Avoiding ambiguity. `--type` is expected to be used less frequently than `--template` since the config default covers the common case.

## Risks / Trade-offs

- **[Ambiguous 1-arg resolution]** → When 1 arg is given, the heuristic checks if it's a valid type. If a user mistypes a type name and has `cli.default_type` set, it would be treated as a name. Mitigation: this is a rare edge case, and the resulting object name would be clearly wrong, prompting the user to notice.
- **[Config file conflicts]** → Multiple processes writing config simultaneously. Mitigation: config is read-only for now; no write path exists. Future config commands should handle this.
- **[Slug collision]** → Different natural-language names could slugify to the same value (e.g., "Some Thought!" and "Some Thought"). Mitigation: ULID suffix ensures filename uniqueness regardless of slug collisions. Types with `unique: true` check the `name` property (original input), not the slug.
