## Context

The TUI handles transient messages ad-hoc: `createState.flash` (auto-dismiss via `tea.Tick` + seq), `saveErr` (manual dismiss), `aiError` (any-key dismiss), and `type_editor.saveErr` / `template_editor.saveErr` (never dismiss). Each has different behavior and rendering location.

The existing `tui/widget/popup.go` provides `OverlayPopup` using lipgloss Layer/Compositor for centered popups. Toast will follow the same Layer/Compositor pattern but position in the bottom-right corner instead of center.

## Goals / Non-Goals

**Goals:**

- Provide a reusable `ToastModel` in `tui/widget/` that any TUI component can use for transient notifications
- Support three severity levels (Info, Warning, Error) with distinct visual styling
- Auto-dismiss via `tea.Tick` with configurable duration, manual dismiss via configurable key
- Toast consumes the dismiss key (does not propagate to other components)
- Group-based message aggregation: multiple messages with the same group key merge into a summary (e.g., "⚠ 3 unresolved refs")
- Replace `aiError` state with Toast notifications in this PR
- Wire up sync warning toasts from `SyncResult.Unresolved`
- Add all 5 `tui.toast.*` config keys

**Non-Goals:**

- Replacing `saveErr`, `createState.flash`, or `type_editor.saveErr` (future PRs)
- help-bar inline position (future; `position` config key is added but only `bottom-right` is implemented)
- Toast stacking (multiple simultaneous toasts) — only one toast visible at a time; new toast replaces current
- Toast queue or history

## Decisions

### 1. ToastModel as a Bubble Tea sub-model in `tui/widget/`

The `ToastModel` follows the Bubble Tea `Model` pattern with its own `Update()` and `View()` methods. It lives in `tui/widget/toast.go` alongside `popup.go`.

**Why:** Consistent with the existing widget pattern. Sub-model encapsulation keeps toast logic isolated from the main app model.

**Alternative considered:** A simple render function (like `CenteredPopup`). Rejected because toast has state (timer, messages, group aggregation) that doesn't fit a pure render function.

### 2. Single active toast with replacement semantics

Only one toast is visible at a time. Calling `Show()` while a toast is active replaces it and resets the timer. Group messages within the same `Show()` call are aggregated, but a new `Show()` replaces the previous toast entirely.

**Why:** Keeps the widget simple. Multiple stacked toasts add visual complexity (positioning, z-ordering, individual timers) without clear benefit for current use cases.

### 3. Group aggregation via `Show()` accepting multiple items

```go
type ToastItem struct {
    Message string
    Group   string // messages with same group are counted, not listed
}

func (t *ToastModel) Show(level ToastLevel, items []ToastItem) tea.Cmd
```

When multiple items share the same group, they render as `"⚠ 3 unresolved refs"` (using the group string as the summary label). Items without a group render individually. This keeps aggregation logic inside Toast while callers control what gets grouped.

**Why:** Issue specifies "multiple messages from a single event are summarized." The caller knows the items; Toast knows how to summarize.

**Alternative considered:** Toast accumulates messages over time and auto-groups. Rejected — adds complexity (when to reset accumulation?) and makes behavior less predictable.

### 4. Bottom-right positioning via lipgloss Layer

Toast renders using `lipgloss.NewLayer().X(x).Y(y)` with coordinates calculated to place the popup in the bottom-right corner, similar to how `OverlayPopup` places popups at center. The toast is composited on top of the existing screen content.

**Why:** Consistent with existing overlay pattern in `popup.go`. Bottom-right avoids occluding primary content (sidebar, body).

### 5. Timer mechanism: `tea.Tick` with monotonic sequence

Each `Show()` increments a sequence counter. The `tea.Tick` callback includes this sequence. When the tick fires, it only dismisses if the sequence matches (prevents stale ticks from dismissing newer toasts). This is the same pattern used by `createState.flash`.

### 6. Esc consumption via Update() return

`ToastModel.Update()` returns `(ToastModel, tea.Cmd, bool)` where the third value `consumed` indicates whether the key was handled. The main app checks `consumed` before passing the key to other handlers.

**Why:** Clean separation — Toast doesn't need to know about other components, and the app doesn't need to inspect Toast's internal state to decide key routing.

### 7. aiError migration

Remove `aiError` from `aiState` enum. In `updateAIDescribeResult` and `updateAITagResult`, error paths call `m.toast.Show(widget.ToastError, ...)` instead of setting `aiState = aiError`. Remove `updateAIError()`, remove aiError rendering from help bar, remove aiError branch from `app.go` Update switch.

Also update `applySelectedTags` which sets `aiError` on link/resolve failures — these become error toasts.

### 8. Config structure

```go
type ToastConfig struct {
    Position     string `yaml:"position,omitempty"`      // "bottom-right" (default) or "help-bar"
    DurationMs   int    `yaml:"duration_ms,omitempty"`   // default 3000
    DismissKey   string `yaml:"dismiss_key,omitempty"`   // default "esc"
    ShowWarnings *bool  `yaml:"show_warnings,omitempty"` // default true (pointer for zero-value detection)
    ShowSuccess  *bool  `yaml:"show_success,omitempty"`  // default false
}
```

Added as `Toast ToastConfig` field on `TUIConfig`. All 5 keys registered in `configKeyRegistry` with `tui.toast.*` prefix.

**Why `*bool`:** Distinguishes "not set" (nil → use default) from "explicitly set to false". Without pointer, `false` and "not set" are indistinguishable in YAML.

## Risks / Trade-offs

- **[Esc consumption changes UX]** → Users accustomed to pressing Esc to exit panels will first dismiss any visible toast. Since toasts auto-dismiss in 3s, this is a minor friction. Mitigation: short default duration means toasts rarely block Esc.
- **[Single toast limitation]** → If sync warnings and AI errors fire simultaneously, only the latest shows. Mitigation: acceptable for v1; can add stacking later if needed.
- **[Config key `position` added but only one value works]** → `help-bar` value is accepted but behaves identically to `bottom-right` in v1. Mitigation: document that `help-bar` is reserved for future use.
