## Context

#236 introduced `.typemd/config.yaml` with a `VaultConfig` struct containing `CLI.DefaultType`. The config is loaded during `Vault.Open()` and written via `Vault.WriteConfig()`. Currently, only `tmd init` writes config; there's no CLI management interface.

The existing config system is struct-based with YAML serialization. This design builds on it by adding struct-aware get/set/list operations through a key registry.

## Goals / Non-Goals

**Goals:**
- Provide `tmd config set/get/list` subcommands for CLI config management
- Use dot-notation keys (e.g., `cli.default_type`) mapped to VaultConfig struct fields
- Validate keys against a registry of known config keys
- Create `config.yaml` automatically on first `set` if it doesn't exist
- Graceful behavior when config file is missing (`get` returns empty, `list` shows nothing)

**Non-Goals:**
- Support for arbitrary/unknown keys (use struct-aware validation)
- TUI config management (future)
- Config file migration or versioning
- Nested object values (all values are strings for now)

## Decisions

### 1. Key registry with struct-aware mapping

A `configKeyRegistry` maps dot-notation strings to getter/setter functions on `VaultConfig`:

```go
type configKeyEntry struct {
    Get func(cfg *VaultConfig) string
    Set func(cfg *VaultConfig, value string)
}

var configKeyRegistry = map[string]configKeyEntry{
    "cli.default_type": {
        Get: func(cfg *VaultConfig) string { return cfg.CLI.DefaultType },
        Set: func(cfg *VaultConfig, value string) { cfg.CLI.DefaultType = value },
    },
}
```

**Rationale:** Type-safe, explicit mapping. Adding new config keys requires adding one registry entry. No reflection magic, easy to understand and test.

**Alternatives considered:**
- Reflection-based mapping — complex, hard to debug, overkill for a small key set
- YAML node tree manipulation — loses type safety, allows invalid keys
- Code generation — unnecessary complexity for a handful of keys

### 2. Operations on Vault

New methods on `Vault`:

```go
// GetConfigValue returns the value for a dot-notation key.
// Returns ("", false) if the key is unknown or not set.
func (v *Vault) GetConfigValue(key string) (string, bool)

// SetConfigValue sets a value for a dot-notation key.
// Returns error if the key is unknown.
func (v *Vault) SetConfigValue(key, value string) error

// ConfigKeys returns all known config keys sorted alphabetically.
func (v *Vault) ConfigKeys() []string
```

`SetConfigValue` loads the current config (or creates empty), applies the change, and calls `WriteConfig()`.

**Rationale:** Vault is the facade for all vault operations. Config get/set belongs here, not in a separate service. The operations are simple enough that no service layer is needed (consistent with how type schema CRUD is on Vault directly).

### 3. CLI command structure

```
tmd config set <key> <value>
tmd config get <key>
tmd config list
```

Three subcommands under `tmd config`:

- `set` — requires exactly 2 args (key, value). Validates key, sets value, writes config.
- `get` — requires exactly 1 arg (key). Prints value to stdout (empty line if unset). Exit code 0 always.
- `list` — no args. Prints all set (non-empty) values as `key: value` lines. Empty output if no config.

**Rationale:** Follows `git config` style. Simple, predictable, scriptable.

### 4. Error handling

| Scenario | Behavior |
|----------|----------|
| `set` with unknown key | Error: `unknown config key "<key>". Known keys: cli.default_type` |
| `set` with empty value | Allowed — sets the key to empty string (effectively "unset") |
| `get` with unknown key | Error: `unknown config key "<key>". Known keys: cli.default_type` |
| `get` on unset key | Print empty line, exit 0 |
| `list` with no config file | Print nothing, exit 0 |
| `list` with empty config | Print nothing, exit 0 |

**Rationale:** Unknown keys are always errors (struct-aware validation). Missing values are not errors (graceful degradation). This matches the issue requirements.

## Risks / Trade-offs

- **[Small key set]** — Currently only `cli.default_type` exists. The registry pattern is slightly over-engineered for one key but scales cleanly as keys are added. Worth the small upfront cost.
- **[No unset operation]** — Setting to empty string effectively unsets. A dedicated `unset` subcommand could be added later if needed.
- **[String-only values]** — All values are strings. Future config keys needing booleans or lists would require extending the registry. Acceptable for now since all current and foreseeable keys are strings.
