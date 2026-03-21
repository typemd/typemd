## 1. Core: ToastConfig in VaultConfig

- [x] 1.1 Add ToastConfig struct and `Toast` field to TUIConfig in `core/config.go`
- [x] 1.2 Register all 5 `tui.toast.*` keys in configKeyRegistry
- [x] 1.3 Add unit tests for toast config loading and registry (existing BDD config scenarios cover structural loading; unit tests for new keys)

## 2. Widget: ToastModel

- [x] 2.1 Write unit tests for ToastModel (show, dismiss, auto-dismiss tick, group aggregation, level filtering, key consumption)
- [x] 2.2 Implement ToastModel in `tui/widget/toast.go` (types, Show, Update, View, Overlay rendering)

## 3. TUI: Integrate Toast into app model

- [x] 3.1 Add `toast widget.ToastModel` field to `model` struct in `tui/app.go`
- [x] 3.2 Wire ToastModel.Update() in main Update loop with Esc consumption logic
- [x] 3.3 Wire ToastModel overlay rendering in View() (after all other overlays, before help)
- [x] 3.4 Pass ToastConfig from vault config to ToastModel during initialization

## 4. TUI: Migrate aiError to Toast

- [x] 4.1 Remove `aiError` from `aiState` enum and `aiError` field from model struct
- [x] 4.2 Update `updateAIDescribeResult` and `updateAITagResult` error paths to use toast.Show()
- [x] 4.3 Update `applySelectedTags` error paths to use toast.Show()
- [x] 4.4 Remove `updateAIError()` function
- [x] 4.5 Remove aiError help bar branch from app_render.go
- [x] 4.6 Update unit tests to verify AI errors produce toast instead of aiError state

## 5. TUI: Wire sync warning toasts

- [x] 5.1 After sync operations that return SyncResult, convert Unresolved items to warning toast with group key
- [x] 5.2 Add unit test verifying sync unresolved items produce grouped warning toast
