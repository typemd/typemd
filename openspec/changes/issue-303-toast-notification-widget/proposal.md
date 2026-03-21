## Why

The TUI currently has no unified notification system. Transient messages are handled ad-hoc across different components (`createState.flash`, `saveErr`, `aiError`), each with different dismiss behaviors and rendering locations. This inconsistency makes it hard to add new notification sources and creates a fragmented user experience.

## What Changes

- Add a reusable `ToastModel` widget in `tui/widget/toast.go` that surfaces transient notifications as a floating overlay in the bottom-right corner (lipgloss Layer/Compositor)
- Support three severity levels: Info, Warning, Error
- Auto-dismiss after configurable duration (default 3s), manually dismissable via configurable key (default Esc)
- Toast consumes the dismiss key — does not propagate to other components
- Support message aggregation via group key (e.g., multiple unresolved refs → `⚠ 3 unresolved refs`)
- Add `tui.toast.*` configuration keys to `VaultConfig` and `configKeyRegistry`
- **Migrate `aiError` to Toast** — remove `aiError` state from `aiState` enum and replace with toast notifications
- Wire up initial use cases: sync warnings (`SyncResult.Unresolved`) and AI error results

## Capabilities

### New Capabilities

- `toast-notification`: Toast notification widget — ToastModel API, rendering, auto-dismiss timer, Esc handling, group-based message aggregation, and configuration

### Modified Capabilities

- `vault-config`: Adding `tui.toast.*` keys (position, duration_ms, dismiss_key, show_warnings, show_success) to TUIConfig and configKeyRegistry

## Impact

- **New files**: `tui/widget/toast.go`, `tui/widget/toast_test.go`
- **Modified**: `core/config.go` (TUIConfig struct, configKeyRegistry), `tui/ai.go` (remove aiError state), `tui/ai_update.go` (error → toast), `tui/ai_render.go` (remove aiError rendering), `tui/app.go` (integrate ToastModel), `tui/app_render.go` (overlay compositing), `tui/update.go` (Esc routing, sync warning toast)
- **No breaking changes** — existing config files remain valid, new keys are optional with defaults
