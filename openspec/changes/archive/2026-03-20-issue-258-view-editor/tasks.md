## 1. Core: GroupRule struct and ViewConfig migration

- [x] 1.1 Write BDD scenarios for GroupRule (single, multiple, empty, YAML serialization, backward-compat deserialization)
- [x] 1.2 Implement BDD step definitions for GroupRule scenarios
- [x] 1.3 Add GroupRule struct, change ViewConfig.GroupBy to []GroupRule, implement custom UnmarshalYAML for backward compatibility
- [x] 1.4 Add unit tests for GroupRule edge cases (legacy string migration, empty string, nil handling, round-trip YAML)
- [x] 1.5 Update DefaultView() to return nil GroupBy and fix all existing references to the old string GroupBy field

## 2. TUI: Multi-level grouping in view mode

- [x] 2.1 Write BDD scenarios for multi-level grouping (compound labels, missing values, single level, no grouping)
- [x] 2.2 Implement BDD step definitions for multi-level grouping scenarios
- [x] 2.3 Rewrite buildGroups() to support []GroupRule with compound labels ("·" separator)
- [x] 2.4 Update visibleRows() and View() to use len(view.GroupBy) > 0 instead of GroupBy == ""
- [x] 2.5 Add unit tests for compound label generation edge cases (empty values, special characters, 3+ levels)

## 3. TUI: View editor sub-model scaffold

- [x] 3.1 Create viewEditor struct with sections (filter, sort, group), cursor, mode, and layout fields
- [x] 3.2 Implement View() rendering: three sections with current rules displayed and "+ Add" action rows
- [x] 3.3 Implement Update() routing: section navigation (Tab), rule navigation (↑/↓), and mode dispatch
- [x] 3.4 Implement HelpBar() with context-sensitive help text
- [x] 3.5 Implement SetSize() for layout calculations

## 4. TUI: View editor pickers

- [x] 4.1 Implement property picker (text input + scrollable filtered list of schema properties including system properties)
- [x] 4.2 Implement operator picker (scrollable list of type-aware operators from validOperators registry, no text input)

## 5. TUI: View editor rule management

- [x] 5.1 Implement add filter rule flow (property picker → operator picker → value input → confirm, skip value for is_empty/is_not_empty)
- [x] 5.2 Implement add sort rule flow (property picker → direction toggle asc/desc → confirm)
- [x] 5.3 Implement add group rule flow (property picker → confirm)
- [x] 5.4 Implement edit existing rule (Enter on rule, reuse add flow with pre-populated values)
- [x] 5.5 Implement delete rule (x/d key, immediate removal)
- [x] 5.6 Implement move rule order (Shift+K up, Shift+J down, boundary check)
- [x] 5.7 Implement auto-save: persist ViewConfig to YAML after each add/edit/delete/move

## 6. TUI: View editor integration with view mode

- [x] 6.1 Add viewEditor field to viewMode struct, wire `e` key to open editor (close preview if active)
- [x] 6.2 Integrate split panel rendering: table (60%) + editor (40%) using JoinHorizontal, with narrow terminal fallback
- [x] 6.3 Implement live reload: after each rule change, re-query objects and refresh left table
- [x] 6.4 Wire Esc in editor to close editor and return to full-width table
- [x] 6.5 Implement delete view action (D key with y/n confirmation, exit view mode on delete)
- [x] 6.6 Update view mode HelpBar() to show editor-specific help when editor is open

## 7. Testing and cleanup

- [x] 7.1 Run full test suite (go test ./...) and fix any regressions from GroupBy type change
- [x] 7.2 Manual testing: create views with filter/sort/group rules via editor, verify YAML output and live table updates
