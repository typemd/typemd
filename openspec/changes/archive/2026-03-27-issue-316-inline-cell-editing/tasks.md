## 1. TUI: Cell Navigation

- [x] 1.1 ~~Write BDD scenarios for table cell navigation~~ — deferred to #114 (TUI BDD infra not yet established; unit tests cover behavior)
- [x] 1.2 Add `colCursor` field to `viewMode` struct and implement left/right/h/l key handling
- [x] 1.3 Implement Tab/Shift+Tab navigation between editable cells (skip read-only)
- [x] 1.4 Remap Enter from open-object to edit-cell; add `o` key for opening object detail

## 2. TUI: Crosshair Rendering

- [x] 2.1 ~~Write BDD scenarios for crosshair highlighting~~ — deferred to #114
- [x] 2.2 Modify `viewTable()` to render per-cell styles: active cell highlight, dim row highlight, dim column header tint
- [x] 2.3 Add unit tests for crosshair style selection logic

## 3. TUI: Cell Edit State

- [x] 3.1 Define `cellEdit` struct (editing bool, row/col index, property name/type, textinput model, picker state)
- [x] 3.2 Implement `viewMode.activateCellEdit()` — determine property type and initialize appropriate widget
- [x] 3.3 Implement `viewMode.cancelCellEdit()` — discard edit and return to navigation
- [x] 3.4 Implement read-only cell detection (relation, created_at, updated_at) — Enter is no-op on these

## 4. TUI: Cell Editing — Text-based Types

- [x] 4.1 ~~Write BDD scenarios for text-based cell editing~~ — deferred to #114
- [x] 4.2 Implement textinput cell editing for string, number, date, datetime, url property types
- [x] 4.3 Implement validation via `core.ValidatePropertyValue()` with toast error on failure
- [x] 4.4 Implement `viewMode.applyCellValue()` save pipeline (update object, call vault.SaveObject)

## 5. TUI: Cell Editing — Picker Types

- [x] 5.1 ~~Write BDD scenarios for select/multi-select cell editing~~ — deferred to #114
- [x] 5.2 Implement select option picker for select property cells
- [x] 5.3 Implement multi-select option picker with Space toggle for multi_select property cells

## 6. TUI: Cell Editing — Checkbox and NAME

- [x] 6.1 ~~Write BDD scenarios for checkbox/NAME editing~~ — deferred to #114
- [x] 6.2 Implement checkbox direct toggle on Enter/Space in table cells
- [x] 6.3 Implement NAME column editing via textinput (update object name and save)

## 7. TUI: Edge Cases and Integration

- [x] 7.1 Handle cancel-on-external-file-change: cancel active cell edit when file watcher triggers object refresh
- [x] 7.2 Disable cell editing when preview panel or view editor is open
- [x] 7.3 Update help bar text to reflect new key bindings (h/l, Enter=edit, o=open, Tab)
- [x] 7.4 Add unit tests for edge cases (empty schema, single NAME column, column cursor reset on view reload)
