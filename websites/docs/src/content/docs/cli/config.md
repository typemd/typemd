---
title: tmd config
description: Manage vault configuration from the CLI.
sidebar:
  order: 12
---

Manage `.typemd/config.yaml` from the command line. Provides `get`, `set`, and `list` subcommands for reading and writing vault configuration without editing files directly.

## Subcommands

### `tmd config set`

Set a config value. Creates `.typemd/config.yaml` if it doesn't exist.

```bash
tmd config set cli.default_type idea
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
```

## Known Keys

| Key | Description |
|-----|-------------|
| `cli.default_type` | Default type for `tmd object create` when no type argument is given |
| `tui.debounce_ms` | File watcher debounce interval in milliseconds (default: 200) |

Keys use dot-notation that maps to YAML nesting. For example, `cli.default_type` corresponds to:

```yaml
cli:
  default_type: idea
tui:
  debounce_ms: 200
```

## Error Handling

| Scenario | Behavior |
|----------|----------|
| Unknown key on `set` or `get` | Error with list of known keys |
| Unset key on `get` | Empty output, exit code 0 |
| No config file on `list` | Empty output, exit code 0 |

## See Also

- [File structure](/advanced/file-structure/) — where `.typemd/config.yaml` lives
- [tmd init](/cli/init/) — creates initial config when starter types are selected
- [tmd object create](/cli/create/) — uses `cli.default_type` for default type
