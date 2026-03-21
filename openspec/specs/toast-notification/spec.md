### Requirement: ToastModel displays transient notifications as bottom-right overlay

The `ToastModel` SHALL render toast notifications as a floating overlay in the bottom-right corner of the terminal using lipgloss Layer/Compositor. The toast SHALL appear on top of existing screen content without displacing layout.

#### Scenario: Toast appears in bottom-right corner

- **WHEN** a toast is shown with level Info and message "Saved successfully"
- **THEN** the toast SHALL render as an overlay positioned in the bottom-right corner
- **AND** the background content SHALL remain visible outside the toast area

#### Scenario: Toast renders with no active message

- **WHEN** no toast is active
- **THEN** `ToastModel.View()` SHALL return an empty string

### Requirement: Toast supports three severity levels with distinct styling

The `ToastModel` SHALL support three severity levels: Info, Warning, and Error. Each level SHALL have a distinct visual prefix and color.

#### Scenario: Info toast

- **WHEN** a toast is shown with level Info
- **THEN** the toast SHALL display with an info prefix ("ℹ")

#### Scenario: Warning toast

- **WHEN** a toast is shown with level Warning
- **THEN** the toast SHALL display with a warning prefix ("⚠")

#### Scenario: Error toast

- **WHEN** a toast is shown with level Error
- **THEN** the toast SHALL display with an error prefix ("✗")

### Requirement: Toast auto-dismisses after configurable duration

The `ToastModel` SHALL auto-dismiss after a configurable duration. The default duration SHALL be 3000 milliseconds. The `Show()` method SHALL return a `tea.Cmd` that schedules the dismissal.

#### Scenario: Toast auto-dismisses after default duration

- **WHEN** a toast is shown with no custom duration configured
- **THEN** the toast SHALL auto-dismiss after 3000 milliseconds

#### Scenario: Toast auto-dismisses after custom duration

- **WHEN** `tui.toast.duration_ms` is set to 5000
- **AND** a toast is shown
- **THEN** the toast SHALL auto-dismiss after 5000 milliseconds

#### Scenario: New toast resets the timer

- **WHEN** a toast is active with 1 second remaining
- **AND** a new toast is shown
- **THEN** the timer SHALL reset to the full duration for the new toast
- **AND** the previous toast's tick SHALL NOT dismiss the new toast

### Requirement: Toast is manually dismissable via configurable key

The `ToastModel` SHALL be manually dismissable by pressing a configurable key. The default key SHALL be Esc. When a toast is visible, the dismiss key SHALL be consumed by the toast and NOT propagated to other components.

#### Scenario: Esc dismisses active toast

- **WHEN** a toast is active
- **AND** the user presses Esc
- **THEN** the toast SHALL be dismissed
- **AND** the Esc key SHALL NOT be propagated to other components

#### Scenario: No toast active — Esc passes through

- **WHEN** no toast is active
- **AND** the user presses Esc
- **THEN** the key SHALL NOT be consumed by ToastModel
- **AND** it SHALL propagate to other components normally

#### Scenario: Custom dismiss key

- **WHEN** `tui.toast.dismiss_key` is set to "q"
- **AND** a toast is active
- **AND** the user presses "q"
- **THEN** the toast SHALL be dismissed

### Requirement: Toast aggregates messages by group key

When `Show()` is called with multiple `ToastItem` values sharing the same group key, the `ToastModel` SHALL aggregate them into a single summary line showing the count and group label.

#### Scenario: Multiple items with same group

- **WHEN** `Show()` is called with 3 items all having group "unresolved refs"
- **THEN** the toast SHALL display "⚠ 3 unresolved refs" (one summary line)

#### Scenario: Single item with group

- **WHEN** `Show()` is called with 1 item having group "unresolved refs"
- **THEN** the toast SHALL display "⚠ 1 unresolved ref" (singular form not required; "1 unresolved refs" is acceptable)

#### Scenario: Items without group

- **WHEN** `Show()` is called with items that have no group key
- **THEN** each item's message SHALL be displayed individually

#### Scenario: Mixed grouped and ungrouped items

- **WHEN** `Show()` is called with 2 items in group "unresolved refs" and 1 ungrouped item with message "sync complete"
- **THEN** the toast SHALL display both the group summary and the ungrouped message

### Requirement: Toast respects show_warnings and show_success config

The `ToastModel` SHALL respect the `show_warnings` and `show_success` configuration flags to filter which toasts are displayed.

#### Scenario: Warnings suppressed

- **WHEN** `tui.toast.show_warnings` is set to false
- **AND** a Warning-level toast is shown
- **THEN** the toast SHALL NOT be displayed

#### Scenario: Warnings enabled by default

- **WHEN** `tui.toast.show_warnings` is not configured
- **AND** a Warning-level toast is shown
- **THEN** the toast SHALL be displayed

#### Scenario: Success suppressed by default

- **WHEN** `tui.toast.show_success` is not configured
- **AND** an Info-level toast is shown
- **THEN** the toast SHALL NOT be displayed

#### Scenario: Success enabled

- **WHEN** `tui.toast.show_success` is set to true
- **AND** an Info-level toast is shown
- **THEN** the toast SHALL be displayed

#### Scenario: Error toasts always shown

- **WHEN** regardless of configuration
- **AND** an Error-level toast is shown
- **THEN** the toast SHALL always be displayed

### Requirement: AI errors use toast instead of aiError state

The TUI SHALL use Toast for AI operation errors instead of the `aiError` state. When an AI describe or tag operation fails, the error SHALL be shown as an Error-level toast. The `aiError` state and `aiError` field SHALL be removed from the model.

#### Scenario: AI describe error shows toast

- **WHEN** an AI describe operation returns an error
- **THEN** an Error-level toast SHALL be shown with the error message
- **AND** `aiState` SHALL return to `aiIdle`

#### Scenario: AI tag error shows toast

- **WHEN** an AI tag operation returns an error
- **THEN** an Error-level toast SHALL be shown with the error message
- **AND** `aiState` SHALL return to `aiIdle`

#### Scenario: AI tag apply error shows toast

- **WHEN** applying selected AI tags fails (link or resolve error)
- **THEN** an Error-level toast SHALL be shown with the error message
- **AND** `aiState` SHALL return to `aiIdle`

### Requirement: Sync warnings display as toast

When `Projector.Sync()` produces unresolved references, the TUI SHALL display them as a Warning-level toast with group-based aggregation.

#### Scenario: Sync produces unresolved references

- **WHEN** a sync operation completes with 3 unresolved references
- **THEN** a Warning-level toast SHALL be shown
- **AND** the items SHALL use group key "unresolved refs"
- **AND** the toast SHALL display "⚠ 3 unresolved refs"
