## Why

Marketplace skills (explore, importer, vault-guide) currently live as standalone SKILL.md files loaded by Claude Code's plugin system. There is no way for users to programmatically access these skill instructions enriched with vault context — they must manually copy SKILL.md content and gather vault data themselves. A CLI command enables LLM integrations beyond Claude Code to use these skills with proper vault context.

## What Changes

- New `tmd instructions` command that outputs skill instructions enriched with vault context as JSON
- Embed marketplace SKILL.md files into the Go binary via `//go:embed` (following the `core/starters.go` pattern)
- Support vault-level overrides via `.typemd/instructions/<skill>.md`
- `--skill` flag outputs raw SKILL.md content (with frontmatter) for saving into a skills directory
- No-argument mode lists all available skills

## Capabilities

### New Capabilities
- `skill-instructions`: Embed marketplace skills in binary, load with vault override support, parse SKILL.md frontmatter, output enriched JSON with per-skill context injection

### Modified Capabilities

_(none)_

## Impact

- **core/**: New `core/instructions.go` with embedded skills, SKILL.md parsing, context builder, and override loading
- **cmd/**: New `cmd/instructions.go` Cobra command registered on `rootCmd`
- **marketplace/**: Read-only — SKILL.md files are embedded from their current locations
- **Data model**: New optional vault directory `.typemd/instructions/` for skill overrides
