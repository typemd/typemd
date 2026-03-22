---
name: instructions-guide
description: Guide for using tmd instructions to output embedded skills enriched with vault context. Use when integrating typemd with AI tools, passing vault context to LLMs, or setting up skill-based workflows.
---

# Instructions Guide

`tmd instructions` outputs embedded skill instructions enriched with your vault's data model. Use it to give AI tools vault awareness without manual context gathering.

## Available Skills

| Skill | Purpose |
|-------|---------|
| `explore` | Analyze existing markdown files and suggest type schemas (read-only) |
| `importer` | Convert existing markdown files into typemd objects |

## Commands

### List skills

```bash
tmd instructions
```

Output:
```
explore   Explore existing markdown files and suggest typemd type schemas...
importer  Convert existing markdown files into typemd objects...
```

### List skills as JSON

```bash
tmd instructions --json
```

Returns a JSON array of `{name, description}` objects.

### Get skill with vault context

```bash
tmd instructions explore
```

Outputs JSON with the skill instructions and vault context:

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
        "description": "Track your reading",
        "properties": [
          {"name": "author", "type": "string", "description": "Book author"}
        ]
      }
    ]
  }
}
```

- `instructions` — full skill body (markdown, no frontmatter)
- `context.types` — all types in the vault with emoji, description, and properties (name, type, description)

If run outside a vault, `context` is omitted.

### Get raw SKILL.md

```bash
tmd instructions explore --skill
```

Outputs the raw SKILL.md content including YAML frontmatter. Use this to save a skill into a Claude Code skills directory or inspect the source. Works without a vault.

## Vault Overrides

Override an embedded skill by placing a file at `.typemd/instructions/<skill>.md`:

```
.typemd/
└── instructions/
    └── explore.md    # overrides the embedded explore skill
```

The override follows the same SKILL.md format — YAML frontmatter with `name` and `description`, followed by markdown body. If no frontmatter is present, the embedded skill's metadata is preserved and only the instructions body is replaced.

## Integration Patterns

### Feed vault context to an AI tool

```bash
# Get explore instructions with vault context
tmd instructions explore | ai-tool --context -

# Or save to a file
tmd instructions importer > /tmp/importer-context.json
```

The JSON output contains everything the AI needs: the skill's instructions plus the vault's type definitions.

### Install a skill into Claude Code

```bash
# Save raw SKILL.md to your project's skills directory
tmd instructions explore --skill > .claude/skills/explore/SKILL.md
```

### Combine with vault-guide

For full vault awareness, use both:
1. Install the `vault-guide` skill (from the guides plugin) for CLI, schema, and object format knowledge
2. Use `tmd instructions` at runtime for vault-specific type context

The vault-guide provides static knowledge (how typemd works), while `tmd instructions` provides dynamic context (what types exist in this vault).
