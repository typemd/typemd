## Why

typemd relies entirely on git for version control but provides no integration. Users must manually run `git log` with the correct file path to see an object's history. A `tmd log` command would provide a typemd-aware wrapper that resolves object IDs to file paths and displays formatted commit history.

## What Changes

- Add `tmd log <object-id>` CLI command that wraps `git log` for a specific object file
- Resolve object IDs (including prefix matching with interactive disambiguation) to file paths
- Format output with commit hash, author, date, and message
- Support `--oneline` flag for compact output
- Detect non-git vaults and provide clear error messages

## Capabilities

### New Capabilities

- `object-log`: View git commit history for a specific object via `tmd log <object-id>`

### Modified Capabilities

(none)

## Impact

- **cmd/**: New `log.go` command file registered on `rootCmd`
- **core/**: Uses existing `Vault.ObjectPath()` and `resolveObjectInteractive()` for ID resolution
- **Dependencies**: Relies on `git` being available in PATH (subprocess call)
