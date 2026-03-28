## 1. TUI: Date edit sub-model (core widget)

- [x] 1.1 Create `tui/date_edit.go` with `dateEdit` struct, dateEditMode constants, segment mode state, calendar mode state, and mode toggle
- [x] 1.2 Implement segmented input: segment navigation, increment/decrement with carry, digit input, day-of-week display
- [x] 1.3 Implement calendar mode: month grid rendering, day navigation, month switching (H/L), today jump (t), today marker
- [x] 1.4 Implement confirm (Enter) and cancel (Esc) for both modes, outputting validated time.Time
- [x] 1.5 Implement View() rendering for both modes (segment display and calendar grid)
- [x] 1.6 Add unit tests for dateEdit: segment navigation, increment/decrement with carry, digit input auto-advance, leap year clamping, calendar navigation, mode toggle, confirm/cancel

## 2. TUI: Integrate date picker into property editor

- [x] 2.1 Add `propModeDateSegment` and `propModeDateCalendar` mode constants to `propEditMode`
- [x] 2.2 Add `dateEdit *dateEdit` field to `propEditor` struct
- [x] 2.3 Route `date` type in `activateEdit()` to create and activate `dateEdit` instead of textinput
- [x] 2.4 Handle date picker key events in `prop_editor_update.go` (delegate to `dateEdit.Update()`)
- [x] 2.5 Render date picker in properties panel (delegate to `dateEdit.View()`)
- [x] 2.6 Update help bar to show `[DATE]` for `propModeDateSegment` and `[CAL]` for `propModeDateCalendar`
- [x] 2.7 Add unit tests for prop editor date integration (activate opens date picker, confirm saves, cancel restores)

## 3. TUI: Integrate date picker into table view cell editing

- [x] 3.1 Add `cellModeDateSegment` and `cellModeDateCalendar` mode constants to `cellEditMode`
- [x] 3.2 Add `dateEdit *dateEdit` field to `cellEdit` struct
- [x] 3.3 Route `date` type in `activateCellEdit()` to create and activate `dateEdit` instead of textinput
- [x] 3.4 Handle date picker key events in `view_mode_cell_edit.go` (delegate to `dateEdit.Update()`)
- [x] 3.5 Render date picker as overlay in table view (delegate to `dateEdit.View()`)
- [x] 3.6 Update help bar to show `[DATE]` for `cellModeDateSegment` and `[CAL]` for `cellModeDateCalendar`

## 4. Integration testing and polish

- [x] 4.1 Run full test suite and fix any regressions
- [x] 4.2 Verify date picker works with existing date validation (`core.ValidatePropertyValue`)
- [x] 4.3 Verify auto-save behavior on confirm in both property editor and table view
