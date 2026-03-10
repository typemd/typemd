## Why

Every TUI launch starts from scratch — all type groups collapsed, cursor at the top, default panel widths. Users must manually re-expand groups, navigate back to their previous object, and adjust the layout each time. This friction compounds with frequent use and larger vaults.

## What Changes

- TUI saves session state to `.typemd/tui-state.yaml` on exit
- TUI restores session state on next launch, with graceful fallback if saved state is stale
- State storage follows a two-tier model (global + vault-level), similar to git config, where vault-level overrides global
- Search state is explicitly excluded from persistence (ephemeral by nature)

## Capabilities

### New Capabilities
- `tui-session-state`: Persist and restore TUI session state across restarts, including selected object, expanded groups, panel dimensions, focus panel, and scroll offset

### Modified Capabilities
- `tui-layout`: Startup behavior changes — instead of hardcoded defaults, the TUI attempts to restore previous session state before falling back to defaults

## Impact

- **tui/**: New state save/load logic in `app.go`, changes to `Start()` initialization and quit handler
- **core/**: May need a utility for reading/writing JSON state files, or this can live entirely in `tui/`
- **.typemd/**: New file `tui-state.yaml` (user-local, not necessarily gitignored — user decides)
- **~/.config/tmd/**: Future global config location (vault-level only for now)
