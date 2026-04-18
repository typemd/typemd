---
title: tmd config
description: Manage vault configuration from the CLI.
sidebar:
  order: 23
---

Manage `.typemd/config.yaml` from the command line. Provides `get`, `set`, and `list` subcommands for reading and writing vault configuration without editing files directly.

## Subcommands

### `tmd config set`

Set a config value. Creates `.typemd/config.yaml` if it doesn't exist.

```bash
tmd config set cli.default_type idea
tmd config set ai.enabled true
```

Only known keys are accepted. Setting an unknown key returns an error listing the valid keys.

### `tmd config get`

Get a config value. Prints the value to stdout. If the key is not set, prints nothing.

```bash
tmd config get cli.default_type
# → idea
```

### `tmd config list`

List all set (non-empty) config values in `key: value` format. Prints nothing if no config exists or all values are empty.

```bash
tmd config list
# cli.default_type: idea
# ai.enabled: true
```

Use `--all` to show all known keys, including unset keys with their defaults.

## Known Keys

### CLI

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `cli.default_type` | string | *(empty)* | Default type for `tmd object create` when no type argument is given |

### TUI

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `tui.debounce_ms` | int | `200` | File watcher debounce interval in milliseconds |
| `tui.stats_type_layout` | string | `fullscreen` | Stats type detail layout: `fullscreen` or `popup` |
| `tui.keybindings.<action>` | string | *(per-action default)* | Override a global TUI keybinding by action name. See [TUI customization](/tui/tui/#customizing-keybindings) for the full action list. |

### AI

AI features require the `claude` CLI to be installed and authenticated.

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `ai.enabled` | bool | `false` | Enable AI features in the TUI (`g` and `ctrl+e` keybindings) |
| `ai.model` | string | *(claude default)* | Override the Claude model (e.g. `claude-haiku-4-5-20251001`) |
| `ai.prompts.describe` | string | *(built-in)* | Custom system prompt for AI description generation |
| `ai.prompts.tag` | string | *(built-in)* | Custom system prompt for AI tag suggestions |
| `ai.prompts.explore` | string | *(built-in)* | Custom system prompt for AI schema exploration |
| `ai.explore.sample_count` | int | `10` | Number of objects sampled for schema exploration |
| `ai.explore.body_truncate` | int | `500` | Max characters of object body included in schema explore prompt |

Keys use dot-notation that maps to YAML nesting:

```yaml
cli:
  default_type: idea
tui:
  debounce_ms: 200
  stats_type_layout: fullscreen
ai:
  enabled: true
  model: claude-sonnet-4-6-20250627
  prompts:
    describe: "Custom prompt for descriptions"
  explore:
    sample_count: 20
    body_truncate: 1000
```

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Unknown key on `set` or `get` | Error with list of known keys |
| Unset key on `get` | Empty output, exit code 0 |
| No config file on `list` | Empty output, exit code 0 |

## See Also

- [Configuration](/basics/configuration/) — configuration overview and reference
- [File structure](/advanced/file-structure/) — where `.typemd/config.yaml` lives
- [tmd init](/cli/init/) — creates initial config with `default_type: page`
- [TUI AI Assist](/tui/tui/#ai-assist) — AI features enabled by `ai.enabled`
