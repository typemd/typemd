## Why

`tmd object create <type> <name>` requires users to provide both type and a slug-formatted name every time, creating too much friction for quick idea capture. Users need a lower-friction creation flow: omit type to use a configured default, and enter natural-language names that auto-convert to slugs.

## What Changes

- Add **vault config system** (`.typemd/config.yaml`) with interface-layer namespacing (`cli`, `tui`)
- Make `tmd object create`'s **type argument optional** — falls back to `cli.default_type` from config
- Add **automatic slug conversion** for names — spaces to hyphens, lowercasing, etc.
- Add `--type` / `-t` flag to override the configured default type

## Capabilities

### New Capabilities
- `vault-config`: Vault-level configuration file system (`.typemd/config.yaml`) with interface-layer namespacing (cli, tui). Initially implements `cli.default_type`.
- `slug-conversion`: Automatic conversion of natural-language names to valid slug format (spaces → hyphens, lowercasing, etc.)

### Modified Capabilities
- `object-templates`: `tmd object create`'s type argument changes from required to optional; adds `-t` flag for type override

## Impact

- **core/**: New vault config loading logic, slug conversion utility
- **cmd/**: Modify `create.go` args handling (type becomes optional), add `-t` flag
- **File format**: New `.typemd/config.yaml`
- **Backward compatibility**: Fully backward-compatible — behavior unchanged without config
