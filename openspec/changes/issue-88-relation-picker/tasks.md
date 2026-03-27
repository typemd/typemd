## 1. Core: Enable relation properties as editable

- [ ] 1.1 Write BDD scenarios for relation picker behavior (activation, selection, clearing, cancel) in `tui/features/`
- [ ] 1.2 Remove relation/tags from `isPropertyEditable()` skip list — make forward relations and tags navigable by the property cursor
- [ ] 1.3 Add `propModeRelationPick` and `propModeRelationMultiPick` to `propEditMode` enum

## 2. TUI: Relation picker activation and candidate loading

- [ ] 2.1 Implement `activateRelationPicker` — query candidate objects of target type via `QueryService.Query`, populate picker state (candidates, search input, cursor, checked state for multi)
- [ ] 2.2 Wire Enter on relation property to call `activateRelationPicker` (single or multi based on `Multiple` flag)
- [ ] 2.3 Wire Enter on `tags` property to call `activateRelationPicker` with target type `tag` and multi mode

## 3. TUI: Relation picker update (keyboard handling)

- [ ] 3.1 Implement `updateRelationPick` — handle j/k navigation, search text input, backspace, Esc to cancel
- [ ] 3.2 Implement single-select confirm (Enter selects candidate, triggers link command)
- [ ] 3.3 Implement "(none)" option for clearing single-value relations (triggers unlink command)
- [ ] 3.4 Implement multi-select toggle (Space toggles candidate check state)
- [ ] 3.5 Implement multi-select confirm (Enter triggers link/unlink diff, produces commands)

## 4. TUI: Relation picker rendering

- [ ] 4.1 Implement `renderRelationPicker` — search input field + filtered candidate list with cursor, scrollable viewport
- [ ] 4.2 Display candidate names without ULID; add type prefix for untyped relations
- [ ] 4.3 Show checkmarks (☑/☐) for multi-value picker
- [ ] 4.4 Show "(none)" option at top of single-value picker
- [ ] 4.5 Update help bar to show `[PICK]` during relation picker mode

## 5. TUI: Relation mutation messages

- [ ] 5.1 Define `relationLinkedMsg` and `relationUnlinkedMsg` tea.Msg types for async link/unlink results
- [ ] 5.2 Implement tea.Cmd wrappers that call `Vault.LinkObjects` / `Vault.UnlinkObjects` and return messages
- [ ] 5.3 Handle `relationLinkedMsg` / `relationUnlinkedMsg` in propEditor Update — refresh display properties, show error toasts on failure

## 6. Integration and edge cases

- [ ] 6.1 Verify locked object guard blocks relation picker activation with toast
- [ ] 6.2 Verify reverse relations and backlinks remain skipped by property cursor
- [ ] 6.3 Add unit tests for candidate filtering logic (substring match, case-insensitive, empty search)
- [ ] 6.4 Add unit tests for link/unlink diff calculation in multi-value confirm
- [ ] 6.5 Run full test suite and verify no regressions
