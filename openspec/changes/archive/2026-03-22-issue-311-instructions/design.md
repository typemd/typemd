## Context

Marketplace skills (explore, importer, vault-guide) are SKILL.md files under `marketplace/plugins/*/skills/*/SKILL.md`. They are currently loaded by Claude Code's plugin system at runtime. The issue asks for a `tmd instructions` command that embeds these skills in the binary and outputs them enriched with vault context as JSON, enabling LLM integrations beyond Claude Code.

The project already has a pattern for embedding files: `core/starters.go` uses `//go:embed starters/*.yaml` to embed starter type templates. The instructions feature follows this same pattern.

SKILL.md files have YAML frontmatter with `name` and `description` fields, followed by markdown content.

## Goals / Non-Goals

**Goals:**
- Embed marketplace skills in the Go binary
- Parse SKILL.md frontmatter (name, description) and body
- Output JSON with instructions + vault context (type summaries)
- Support vault-level overrides via `.typemd/instructions/<skill>.md`
- Support `--skill` flag for raw SKILL.md output
- List available skills when no argument given

**Non-Goals:**
- Marketplace plugin install/management
- Dynamic skill discovery
- Per-skill context customization beyond type summaries (future)

## Decisions

### 1. Embed in `core/` following the starters pattern

Place embedded SKILL.md files and the instructions API in `core/instructions.go`, mirroring `core/starters.go`. Use `//go:embed` with a `skills/` directory under `core/` that contains copies of the marketplace SKILL.md files.

**Alternative considered:** Embed directly from `marketplace/plugins/*/skills/*/SKILL.md` paths. Rejected because `//go:embed` requires the embedded files to be under the package directory (or a subdirectory), and `marketplace/` is outside `core/`. Copying skills to `core/skills/` keeps the embed simple and the files self-contained.

### 2. SKILL.md frontmatter parsing with `gopkg.in/yaml.v3`

Parse YAML frontmatter between `---` delimiters to extract `name` and `description`. The body is everything after the closing `---`. This uses the existing `gopkg.in/yaml.v3` dependency.

### 3. Per-skill context injection via Go functions

Each skill registers a context-builder function that takes a `*Vault` and returns structured context. For the initial implementation, all skills receive the same context: type summaries (name, emoji, description, properties with name/type/description). This is built from `Vault.ListTypes()` + `Vault.LoadType()`.

**Alternative considered:** No context injection (raw instructions only). Rejected because the primary value of `tmd instructions` over reading SKILL.md directly is the enriched vault context.

### 4. Override via `.typemd/instructions/<skill>.md`

Check `.typemd/instructions/<skill>.md` before falling back to embedded. The override file follows the same SKILL.md format. If the override has no frontmatter, treat the entire file as instructions body and use embedded metadata.

### 5. JSON output structure

```json
{
  "name": "explore",
  "description": "...",
  "instructions": "# Explore\n\n...",
  "context": {
    "types": [
      {
        "name": "book",
        "emoji": "📚",
        "description": "...",
        "properties": [
          {"name": "author", "type": "string", "description": "..."}
        ]
      }
    ]
  }
}
```

The `instructions` field contains the body (no frontmatter). The `context` field contains vault data. For `--skill` mode, raw SKILL.md bytes are printed to stdout (no JSON wrapping).

### 6. List mode output

`tmd instructions` with no argument outputs a simple list:

```
explore     Explore existing markdown files and suggest typemd type schemas...
importer    Convert existing markdown files into typemd objects...
vault-guide typemd reference guide...
```

With `--json`, outputs a JSON array of `{name, description}` objects.

## Risks / Trade-offs

- **Stale embedded copies:** Skills in `core/skills/` could drift from `marketplace/` source files. → Mitigation: document in CONTRIBUTING.md that marketplace skill changes must be synced to `core/skills/`. A future CI check could enforce this.
- **Context size:** Large vaults with many types could produce verbose JSON. → Mitigation: type summaries are compact (name, emoji, description, property list). Body content is not included.
