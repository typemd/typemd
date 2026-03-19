## 1. Core: ViewConfig changes

- [x] 1.1 Add `ViewLayoutTable` constant to core/view.go
- [x] 1.2 Add `Columns []string` field to ViewConfig struct (with yaml tag `columns,omitempty`)
- [x] 1.3 Add unit tests for ViewConfig YAML serialization with table layout and columns field
- [x] 1.4 Update DefaultView() to confirm it returns ViewLayoutList (already does, verify)

## 2. TUI: Layout dispatch and list rendering

- [x] 2.1 Extract current `View()` body into `viewTable()` method (no behavior change)
- [x] 2.2 Add `viewList()` method: emoji + name + inline column values with ` · ` separator
- [x] 2.3 Update `View()` to dispatch based on `vm.view.Layout` (list → viewList, table → viewTable, default → viewList)
- [x] 2.4 Update `viewColumns()` to respect `Columns` config when non-empty (skip pinned/unpinned auto-detection)
- [x] 2.5 Add sort indicators (`↑`/`↓`) to table column headers
- [x] 2.6 Add unit tests for viewList rendering (name only, with columns, empty values omitted)
- [x] 2.7 Add unit tests for viewColumns with explicit Columns config

## 3. TUI: View editor Layout and Columns sections

- [x] 3.1 Add Layout section to view editor (single-choice: list / table)
- [x] 3.2 Add Columns section to view editor (property picker to add/remove/reorder)
- [x] 3.3 Wire layout and columns changes to auto-save via vault.SaveView()

## 4. Documentation and verification

- [x] 4.1 Update CLAUDE.md data model section to document `columns` field and `table` layout
- [x] 4.2 Run full test suite and verify all tests pass
- [ ] 4.3 Manual verification: create view → switch to table → verify columns → switch back to list
