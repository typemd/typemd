---
title: tmd object create
description: Create a new object from a type schema.
sidebar:
  order: 2.5
---

Creates a new object file (Markdown + YAML frontmatter) based on the type schema.

```bash
tmd object create [type] [name]
tmd object create book "Clean Code"
tmd object create book "Clean Code" -t review
tmd object create --type idea "Some Thought"
tmd object create "Some Thought"              # uses default type from config
tmd object create book
```

The `type` argument can be omitted if `--type` flag is set or `cli.default_type` is configured in `.typemd/config.yaml`. When a single argument matches a known type name, it is treated as the type. To use a name that collides with a type name, use the `--type` flag explicitly.

Names are automatically converted to slugs for the filename (e.g., "Clean Code" becomes `clean-code` in the filename), while the original input is preserved as the `name` property in frontmatter.

If the type has a name template, the `name` argument is optional — the name will be auto-generated from the template. When no name is provided and no name template is defined, an error is returned.

The command generates a `.md` file under `objects/<type>/` with all schema-defined properties set to their default values (or `null` if no default is specified). The object is also added to the SQLite index.

Each object is assigned a unique ULID suffix, so multiple objects of the same type can share the same name (e.g. two `book/clean-code-<ulid>` files with different ULIDs). If the type does not exist, an error is returned.

## Flags

| Flag | Description |
|------|-------------|
| `--type` | Object type (overrides config default, no short form) |
| `-t`, `--template` | Template name to use (from `templates/<type>/`) |

## Default type

If `cli.default_type` is configured in `.typemd/config.yaml`, you can omit the type argument:

```yaml
# .typemd/config.yaml
cli:
  default_type: idea
```

```bash
# These are equivalent when default_type is "idea":
tmd object create idea "Some Thought"
tmd object create "Some Thought"
```

The `--type` flag always takes precedence over the config default.

## Templates

If the type has templates defined in `templates/<type>/`, they are resolved automatically:

- **No templates** — proceeds without a template
- **Single template** — auto-applied without prompting
- **Multiple templates** — interactive prompt for selection (or use `-t` to skip the prompt)

Template frontmatter properties override schema defaults. Template body becomes the initial object body. Auto-managed system properties (`created_at`, `updated_at`) in templates are ignored.
