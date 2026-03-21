---
title: Configuration
description: How to configure your vault with .typemd/config.yaml.
sidebar:
  order: 7
---

Each vault has an optional configuration file at `.typemd/config.yaml`. It controls CLI defaults, TUI behavior, and AI features. The file is created automatically when you run `tmd init` or `tmd config set`.

## File location

```
my-vault/
└── .typemd/
    └── config.yaml
```

## Structure

The config uses YAML with three top-level namespaces:

```yaml
cli:
  default_type: page
tui:
  debounce_ms: 200
  stats_type_layout: fullscreen
  toast:
    duration_ms: 3000
    dismiss_key: esc
ai:
  enabled: true
  model: claude-sonnet-4-6-20250627
  prompts:
    describe: "Custom prompt for descriptions"
  explore:
    sample_count: 20
    body_truncate: 1000
```

## CLI settings

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `cli.default_type` | string | *(empty)* | Default type for `tmd object create` when no type argument is given |

When set, you can create objects without specifying a type:

```bash
# Without default_type — type is required
tmd object create idea "My Idea"

# With cli.default_type: idea
tmd object create "My Idea"
```

## TUI settings

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `tui.debounce_ms` | int | `200` | File watcher debounce interval in milliseconds |
| `tui.stats_type_layout` | string | `fullscreen` | Stats type detail layout: `fullscreen` or `popup` |
| `tui.toast.position` | string | `bottom-right` | Toast display position: `bottom-right` or `help-bar` |
| `tui.toast.duration_ms` | int | `3000` | Auto-dismiss duration in milliseconds |
| `tui.toast.dismiss_key` | string | `esc` | Key to manually dismiss a toast notification |
| `tui.toast.show_warnings` | bool | `true` | Show warning-level toast notifications |
| `tui.toast.show_success` | bool | `false` | Show info/success-level toast notifications |

The file watcher monitors `objects/` and `.typemd/types/` for changes. The debounce interval controls how quickly the TUI reacts to file modifications — lower values make updates feel more instant, higher values reduce redundant refreshes.

Toast notifications appear in the bottom-right corner of the TUI for transient messages such as sync warnings (unresolved references) and AI operation errors. They auto-dismiss after `duration_ms` and can be manually dismissed by pressing the `dismiss_key`. Error-level toasts always show regardless of configuration.

## AI settings

AI features require the [Claude CLI](https://docs.anthropic.com/en/docs/claude-code) to be installed and authenticated.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `ai.enabled` | bool | `false` | Enable AI features in the TUI |
| `ai.model` | string | *(claude default)* | Override the Claude model (e.g. `claude-haiku-4-5-20251001`) |
| `ai.prompts.describe` | string | *(built-in)* | Custom system prompt for description generation |
| `ai.prompts.tag` | string | *(built-in)* | Custom system prompt for tag suggestions |
| `ai.prompts.explore` | string | *(built-in)* | Custom system prompt for schema exploration |
| `ai.explore.sample_count` | int | `10` | Number of objects sampled for schema exploration |
| `ai.explore.body_truncate` | int | `500` | Max characters of object body included in explore prompt |

When `ai.enabled` is `true`, the TUI adds two AI keybindings:

- **`g`** in object detail — open AI action picker (generate description / suggest tags)
- **`ctrl+e`** from sidebar — enter schema explore mode

## Reading and writing config

You can edit `.typemd/config.yaml` directly, or use the CLI:

```bash
# Set a value
tmd config set ai.enabled true

# Get a value
tmd config get cli.default_type

# List all set values
tmd config list

# List all known keys with defaults
tmd config list --all
```

Only known keys are accepted. Setting an unknown key returns an error listing the valid keys.

## Defaults

If `.typemd/config.yaml` does not exist or a key is not set, TypeMD uses sensible defaults:

- No default type — `tmd object create` requires a type argument
- File watcher debounce at 200ms
- Stats layout in fullscreen mode
- Toast notifications: bottom-right, 3s auto-dismiss, Esc to dismiss, warnings shown, success hidden
- AI features disabled
- Built-in prompts for all AI operations
- 10 sample objects and 500 character body limit for schema exploration

## See also

- [tmd config](/cli/config/) — CLI command reference
- [File Structure](/advanced/file-structure/) — where `.typemd/config.yaml` lives
- [tmd init](/cli/init/) — creates initial config with `default_type: page`
