---
title: Templates
description: How templates provide default content and auto-generated names for new Objects.
sidebar:
  order: 3
---

Templates give you control over the initial content of new Objects. TypeMD supports two kinds of templates: **object templates** that provide default frontmatter and body content, and **name templates** that auto-generate Object names from patterns.

## Object Templates

Each Type can have one or more object templates stored at `templates/<type>/<name>.md`. Templates are regular Markdown files with optional frontmatter and body content that serve as defaults when creating new Objects.

```
vault/
├── templates/
│   └── book/
│       ├── review.md       # a "review" template for books
│       └── summary.md      # a "summary" template for books
└── objects/
    └── book/
        └── ...
```

### Template content

A template file can include frontmatter properties and body content:

```markdown
---
status: to-read
tags:
  - tag/to-review
---

## Summary

## Key Takeaways

## Notes
```

Template frontmatter properties override the schema's default values. The template body becomes the initial body of the new Object.

### Template resolution

When creating an Object with `tmd object create`, templates are resolved automatically:

- **No templates** — proceeds without a template
- **Single template** — auto-applied without prompting
- **Multiple templates** — interactive prompt for selection (or use `-t` to skip the prompt)

```bash
# If book has one template, it auto-applies
tmd object create book clean-code

# Specify a template explicitly
tmd object create book clean-code -t review

# If multiple templates exist, you'll be prompted to choose
tmd object create book clean-code
```

### System property handling

Templates can override **user-authored** system properties (`name`, `description`, `tags`). **Auto-managed** system properties (`created_at`, `updated_at`) in templates are ignored — they always reflect the actual creation time. Properties in the template that are not defined in the type schema are silently ignored.

## Name Templates

Object names are normally provided manually at creation time. For types with predictable naming patterns (daily journals, meeting notes, weekly reviews), this is repetitive and error-prone. Name templates let type schemas define auto-generated names using placeholders, enabling one-command object creation.

### Defining a name template

Add a `name` entry in the type schema's `properties` array with a `template` field:

```yaml
# types/journal/schema.yaml
name: journal
emoji: 📓
properties:
  - name: name
    template: "日記 {{ date:YYYY-MM-DD }}"
  - name: mood
    type: select
    options:
      - value: good
      - value: okay
      - value: bad
```

The `name` entry in `properties` only allows the `template` field — no `type`, `options`, `pin`, `emoji`, or other property fields.

### Date placeholders

The `{{ date:FORMAT }}` placeholder accepts the following format tokens:

| Token | Description | Example |
|-------|-------------|---------|
| `YYYY` | Four-digit year | `2026` |
| `MM` | Two-digit month | `03` |
| `DD` | Two-digit day | `14` |
| `HH` | Two-digit hour (24h) | `09` |
| `mm` | Two-digit minute | `30` |
| `ss` | Two-digit second | `05` |

Examples:

| Template | Result (on 2026-03-14 09:30) |
|----------|------------------------------|
| `日記 {{ date:YYYY-MM-DD }}` | `日記 2026-03-14` |
| `Meeting {{ date:YYYY-MM-DD HH:mm }}` | `Meeting 2026-03-14 09:30` |
| `Week {{ date:YYYY-MM }}` | `Week 2026-03` |
| `Weekly Review` | `Weekly Review` (no placeholders, used as literal) |

### CLI usage

When a type has a name template, the name argument is optional:

```bash
# Name auto-generated from template
tmd object create journal

# Explicit name overrides the template
tmd object create journal "my-journal"
```

If a type does not have a name template, the name argument is required.

## See also

- [Properties](/basics/properties#system-properties) — system property mutability
- [tmd object create](/tui/create) — CLI reference for creating objects
- [Types](/concepts/types) — type schema structure
