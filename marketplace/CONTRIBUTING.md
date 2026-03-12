# Contributing to typemd Marketplace

Thank you for contributing a plugin to the typemd marketplace! This guide covers everything you need to submit a plugin.

## Submission Process

1. Fork the [typemd/typemd](https://github.com/typemd/typemd) repository
2. Create your plugin under `marketplace/plugins/<your-plugin-name>/`
3. Add your plugin entry to `marketplace/.claude-plugin/marketplace.json`
4. Submit a pull request — CI will automatically validate your plugin

## Naming Rules

- Plugin names use **kebab-case** (lowercase letters, numbers, hyphens): `reading-notes`, `auto-relations`
- Names prefixed with `tmd-` or `typemd-` are **reserved for official plugins** and will be rejected by CI

## Plugin Structure

Every plugin must follow Claude Code's plugin format:

```
plugins/<your-plugin-name>/
├── .claude-plugin/
│   └── plugin.json          # Plugin metadata (required)
├── skills/                   # Claude skills (at least one of skills/ or commands/)
│   └── <skill-name>/
│       └── SKILL.md
├── commands/                 # Slash commands (optional)
│   └── <command>.md
└── README.md                 # Usage documentation (required)
```

### plugin.json

```json
{
  "name": "your-plugin-name",
  "description": "Brief description of what your plugin does",
  "version": "1.0.0",
  "author": {
    "name": "Your Name"
  },
  "keywords": ["relevant", "tags"]
}
```

### SKILL.md

Every `SKILL.md` must have YAML frontmatter with at least a `description` field:

```markdown
---
description: What this skill does and when Claude should use it
---

Skill instructions here...
```

### marketplace.json Entry

Add your plugin to the `plugins` array in `marketplace/.claude-plugin/marketplace.json`:

```json
{
  "name": "your-plugin-name",
  "source": "./plugins/your-plugin-name",
  "description": "Brief description",
  "version": "1.0.0",
  "category": "workflow",
  "tags": ["relevant", "tags"]
}
```

## Quality Requirements

- **One plugin, one purpose** — keep plugins focused on a single workflow or capability
- **README.md is required** — explain what the plugin does, how to use it, and any prerequisites
- **SKILL.md must have a description** — Claude uses this to decide when to load the skill
- **Test locally before submitting** — use `claude --plugin-dir ./marketplace/plugins/your-plugin` to verify your plugin works

## Local Testing

Test your plugin before submitting:

```bash
claude --plugin-dir ./marketplace/plugins/your-plugin-name
```

Then try your skill:

```
/your-plugin-name:skill-name
```
