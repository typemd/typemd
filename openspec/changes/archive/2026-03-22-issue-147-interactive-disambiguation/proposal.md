## Why

When a CLI prefix matches multiple objects, users must manually copy-paste a full ID from the error output and re-run the command. This breaks flow and creates unnecessary friction — especially when objects share a common name prefix (e.g., `book/clean-code` matching two editions). An interactive selection menu lets users pick the intended object immediately, without leaving the command context.

## What Changes

- Add an interactive Bubble Tea selection picker that activates when a CLI command encounters an `AmbiguousMatchError` in an interactive terminal
- The picker displays each candidate's human-readable `name` property (from the index) alongside the full object ID for easy identification
- When stdin is not a terminal (piped input), fall back to the existing error message — no behavior change for non-interactive contexts
- All CLI commands that resolve object IDs (`show`, `link`, `unlink`) gain this disambiguation behavior through a shared helper

## Capabilities

### New Capabilities

- `cli-interactive-disambiguation`: Interactive selection picker for resolving ambiguous object ID prefixes in CLI commands. Covers the Bubble Tea picker UI, terminal detection, fallback behavior, and integration with all ID-accepting commands.

### Modified Capabilities

_(none — this is additive behavior that does not change existing spec-level requirements)_

## Impact

- **cmd/**: New `disambiguate_picker.go` file with Bubble Tea model; updated `helper.go` with shared resolve-or-pick function; minor changes to `show.go`, `link.go`, `unlink.go` to use the new helper
- **core/**: Read-only usage of `AmbiguousMatchError.Matches` — no changes to core logic
- **Dependencies**: Uses existing `bubbletea/v2` and `lipgloss/v2` already in go.mod
- **BDD**: New feature file in `cmd/features/` covering interactive and non-interactive scenarios
