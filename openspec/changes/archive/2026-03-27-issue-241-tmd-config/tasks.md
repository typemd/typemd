## 1. Core: Config Key Registry

- [x] 1.1 Write BDD scenarios for config key registry (known key accepted, unknown key rejected, list known keys)
- [x] 1.2 Implement BDD step definitions for config key registry scenarios
- [x] 1.3 Add `configKeyEntry` struct and `configKeyRegistry` map in `vault_config.go`
- [x] 1.4 Add `GetConfigValue()`, `SetConfigValue()`, `ConfigKeys()` methods on Vault
- [x] 1.5 Add unit tests for edge cases (empty value set, get on unset known key, sorted key list)

## 2. CLI: Config Subcommands

- [x] 2.1 Add `config.go` with `configCmd` parent and `configSetCmd`, `configGetCmd`, `configListCmd` subcommands
- [x] 2.2 Add unit tests for CLI commands (set valid key, set unknown key, get set key, get unset key, get unknown key, list with values, list empty)
- [x] 2.3 Register `configCmd` on root command

## 3. Integration Verification

- [x] 3.1 Run full test suite (`go test ./...`) and verify no regressions
- [x] 3.2 Manual test: `tmd config set cli.default_type idea` creates/updates config
- [x] 3.3 Manual test: `tmd config get cli.default_type` returns value
- [x] 3.4 Manual test: `tmd config list` shows all set values
- [x] 3.5 Manual test: `tmd config set unknown.key value` shows error with known keys
