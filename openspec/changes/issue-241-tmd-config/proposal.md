## Why

After #236 introduced `.typemd/config.yaml` with the `cli.default_type` setting, the only way to create or modify vault configuration is manual file editing. Existing vaults (or vaults initialized with `--no-starters`) have no `config.yaml` at all. Users need a CLI command to manage vault config without touching files directly.

## What Changes

- Add `tmd config set <key> <value>` — set a config value (creates `config.yaml` if missing)
- Add `tmd config get <key>` — get a config value (empty output if unset)
- Add `tmd config list` — list all config values in `key: value` format
- Dot-notation keys map to YAML nesting (e.g., `cli.default_type` → `cli:\n  default_type:`)
- Struct-aware key validation — only known keys are accepted

## Capabilities

### New Capabilities
- `config-cli`: CLI subcommands (`set`, `get`, `list`) for managing `.typemd/config.yaml` with dot-notation keys and struct-aware validation

### Modified Capabilities
- `vault-config`: Add `SetConfigValue()` and `GetConfigValue()` methods to Vault for programmatic config access; add `ConfigKeys()` for listing known keys

## Impact

- **core/**: New methods on Vault and VaultConfig for get/set/list operations, key registry
- **cmd/**: New `config.go` with `set`, `get`, `list` subcommands
- **Backward compatibility**: Fully backward-compatible — no existing behavior changes
