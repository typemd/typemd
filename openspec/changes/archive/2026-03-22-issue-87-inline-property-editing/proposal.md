## Why

The TUI properties panel is currently read-only. Users can see property values but cannot edit them without leaving the TUI and manually editing Markdown files. This forces a context-switch that breaks the flow of knowledge management. Inline editing is the core feature of v0.7.0 and the foundation for making the TUI a complete editing environment.

## What Changes

- Add cursor navigation between individual property fields in the properties panel (j/k or arrow keys)
- Add inline editing for all non-relation property types with type-appropriate input widgets:
  - **string/number/date/datetime/url**: textinput component with type-specific validation
  - **checkbox**: direct toggle between ☐ and ☑ via Enter/Space
  - **select**: options list picker with j/k navigation
  - **multi_select**: options list multi-picker with Space toggle and Enter confirm
- Skip read-only properties during navigation (reverse relations, backlinks, immutable system properties like `created_at`, `updated_at`)
- Auto-save edited values on confirm (consistent with existing body edit behavior)
- Validate input per property type before accepting (number format, date format, URL format)

## Capabilities

### New Capabilities

- `property-cursor-navigation`: Cursor-based navigation between individual property fields in the properties panel, with read-only field skipping
- `inline-property-editing`: Type-aware inline editing widgets for all non-relation property types (string, number, date, datetime, url, checkbox, select, multi_select)

### Modified Capabilities

- `tui-layout`: Properties panel gains interactive editing state with visual indicators (edit border color, active field highlight)

## Impact

- **tui/**: Major changes — new property cursor model, edit state management, input widget rendering, keyboard handling
- **core/**: Minor changes — may need validation helpers exposed for TUI consumption (number parsing, date format validation, URL validation, select option validation)
- **No breaking changes**: Existing read-only behavior is preserved; editing is additive
