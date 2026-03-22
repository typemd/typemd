## 1. Core: Property validation helper

- [x] 1.1 Add `ValidatePropertyValue(propType string, options []Option, input string) error` to core package
- [x] 1.2 Add unit tests for ValidatePropertyValue (number, date, datetime, url, select, multi_select edge cases)

## 2. Core: Checkbox display format

- [x] 2.1 Update `DisplayProperty.FormatValue()` to use ☐/☑ instead of ✓/empty for checkbox type
- [x] 2.2 Add unit tests for checkbox FormatValue (true→☑, false→☐, nil→☐)

## 3. TUI: Property cursor navigation

- [x] 3.1 Write unit tests for property cursor navigation (cursor appears, moves up/down, skips read-only, boundaries)
- [x] 3.2 Create `propEditor` sub-model with cursor state, editable property filtering, and navigation logic
- [x] 3.3 Integrate `propEditor` into main model — activate on focusProps, render cursor indicator in renderProperties
- [x] 3.4 Verify unit tests pass

## 4. TUI: String/number/date/datetime/url inline editing (textinput types)

- [x] 4.1 Write unit tests for textinput-based editing (activate, confirm, cancel, validation error)
- [x] 4.2 Add textinput editing state to propEditor — Enter activates, pre-fills value, Enter confirms, Esc cancels
- [x] 4.3 Add type-specific validation on confirm using core.ValidatePropertyValue, show toast on error
- [x] 4.4 Add auto-save on successful edit (update obj.Properties, set dirty, save)
- [x] 4.5 Verify unit tests pass

## 5. TUI: Checkbox toggle

- [x] 5.1 Write unit tests for checkbox toggle (Enter toggles, Space toggles, display format)
- [x] 5.2 Implement checkbox toggle in propEditor — Enter/Space toggles value, immediate save
- [x] 5.3 Verify unit tests pass

## 6. TUI: Select property picker

- [x] 6.1 Write unit tests for select picker (open, navigate, select, cancel, current value highlighted)
- [x] 6.2 Implement select option picker in propEditor — Enter opens, j/k navigates, Enter selects, Esc cancels
- [x] 6.3 Verify unit tests pass

## 7. TUI: Multi-select property picker

- [x] 7.1 Write unit tests for multi-select picker (open, toggle with Space, confirm with Enter, cancel)
- [x] 7.2 Implement multi-select option picker in propEditor — Enter opens, Space toggles, Enter confirms, Esc cancels
- [x] 7.3 Verify unit tests pass

## 8. TUI: Visual indicators and help bar

- [x] 8.1 Update help bar to show property editing hints when focusProps is active
- [x] 8.2 Add edit border color (orange) when actively editing a property
- [x] 8.3 Update renderProperties to use propEditor cursor state for highlighting
