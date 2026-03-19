## Why

TUI session state (`.typemd/tui-state.yaml`) does not save view mode context. When a user exits the TUI while in view mode and reopens it, the view is lost and the TUI starts in normal sidebar mode. This breaks the "resume where I left off" expectation, especially for users who work primarily in a specific view.

## What Changes

- Add view mode fields to `SessionState`: active view type name, view name, cursor position, scroll offset, and expanded group labels
- `captureState()` captures view mode state when `rightPanel == panelView`
- `applySessionState()` restores view mode on startup when view state fields are present
- Graceful fallback when the saved type or view no longer exists

## Capabilities

### New Capabilities

- `view-mode-persistence`: Persisting and restoring TUI view mode state across sessions, including cursor position and expanded groups within the view

### Modified Capabilities

- `tui-session-state`: SessionState gains view mode fields; capture and restore logic extended to handle view mode context

## Impact

- `tui/state.go` — SessionState struct, captureState(), applySessionState(), loadSessionState()
- `tui/app.go` — Start() function to trigger view mode restoration on launch
- `tui/view_mode.go` — May need method to expose group expanded state for capture
- `.typemd/tui-state.yaml` — New optional fields (backward compatible: old files without view fields still work)
