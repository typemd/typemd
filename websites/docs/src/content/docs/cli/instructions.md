---
title: tmd instructions
description: Output embedded skill instructions enriched with vault context.
sidebar:
  order: 24
---

Output skill instructions enriched with vault context. Skills are embedded in the `tmd` binary and can be overridden per-vault.

## List available skills

```bash
tmd instructions
```

Lists all embedded skills with their descriptions.

### Example output

```
explore     Explore existing markdown files and suggest typemd type schemas...
importer    Convert existing markdown files into typemd objects...
onboarding  Import existing markdown collections into a structured typemd vault...
```

### JSON output

```bash
tmd instructions --json
```

Returns a JSON array of `{name, description}` objects.

## Get skill with vault context

```bash
tmd instructions explore
```

Outputs the skill instructions as JSON, enriched with vault context (type summaries including properties).

### Example output

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
          {"name": "author", "type": "string"}
        ]
      }
    ]
  }
}
```

If run outside a vault, the `context` field is omitted.

## Raw skill output

```bash
tmd instructions explore --skill
```

Outputs the raw SKILL.md content including YAML frontmatter, suitable for saving directly into a skills directory. Works without a vault.

## Vault overrides

Place a file at `.typemd/instructions/<skill>.md` to override an embedded skill. The override file follows the same SKILL.md format with optional `name` and `description` frontmatter fields. If no frontmatter is present, the override body replaces the instructions while the embedded skill's metadata is preserved.
