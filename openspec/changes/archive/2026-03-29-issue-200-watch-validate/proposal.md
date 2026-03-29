## Why

`tmd type validate` runs once and exits, requiring users to manually re-run the command after each file edit to check for new errors. This creates a slow feedback loop during schema and object authoring. A `--watch` mode — similar to `tsc --watch` — would continuously monitor vault files and re-run validation on changes, giving instant feedback.

## What Changes

- Add `--watch` / `-w` flag to `tmd type validate` command
- When `--watch` is active, the command enters a persistent loop: watches `.typemd/types/`, `.typemd/properties.yaml`, and `objects/` for file changes
- On change: debounce (200ms), clear terminal, re-index, re-run all validation phases, display results
- Graceful exit on Ctrl+C via signal handling
- Reuse existing `fsnotify` + debounce patterns from `tui/watcher.go`

## Capabilities

### New Capabilities

- `watch-validate`: Continuous file watching with debounced re-validation for `tmd type validate --watch`

### Modified Capabilities

_(none — this is a new additive feature with no changes to existing validation behavior)_

## Impact

- **Code**: `cmd/validate.go` (add flag + watch loop), possibly a new shared watcher utility extracted from `tui/watcher.go`
- **Dependencies**: Already uses `fsnotify` (existing dependency)
- **APIs**: No changes to core validation APIs — reuses existing `ValidateAll*` functions
- **Systems**: CLI-only change, no impact on TUI, MCP, or Web
