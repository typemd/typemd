## 1. Core: Add FormatValue and refactor Format

- [x] 1.1 Write unit tests for FormatValue() covering all property types (string, checkbox, date, datetime, multi_select, relation, backlink, reverse, nil)
- [x] 1.2 Add FormatValue() method to DisplayProperty in core/display.go
- [x] 1.3 Refactor Format() to delegate to FormatValue() (key + ": " + FormatValue())
- [x] 1.4 Update existing Format() test expectations for checkbox (from [x]/[ ] to ✓/empty)

## 2. TUI: Replace formatPropValue with FormatValue

- [x] 2.1 Replace formatPropValue() calls in view_mode.go table rows with DisplayProperty.FormatValue()
- [x] 2.2 Replace formatPropValue() calls in view_mode.go preview panel with DisplayProperty.FormatValue()
- [x] 2.3 Remove the formatPropValue() function from view_mode.go

## 3. Verify

- [x] 3.1 Run go test ./... and confirm all tests pass
- [x] 3.2 Run go build ./... and confirm clean build
