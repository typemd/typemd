## 1. TUI: State Data Model & Serialization

- [x] 1.1 Write BDD scenarios for session state save/load (state file created on exit, restored on launch, missing file fallback, corrupt file fallback, partial state fallback)
- [x] 1.2 Implement BDD step definitions for session state scenarios
- [x] 1.3 Define `SessionState` struct and YAML marshal/unmarshal in `tui/state.go` (make BDD scenarios pass)
- [x] 1.4 Add unit tests for `SessionState` edge cases (empty fields, unknown YAML keys, zero values vs missing)

## 2. TUI: Save State on Exit

- [x] 2.1 Write BDD scenarios for state save on quit (state includes selected object ID, expanded groups, panel widths, focus, scroll offset; excludes search state)
- [x] 2.2 Implement BDD step definitions for save scenarios
- [x] 2.3 Add `saveState()` method to model, call from quit handler in `Update()` (make BDD scenarios pass)
- [x] 2.4 Add unit tests for save edge cases (no object selected, no groups expanded, write permission error)

## 3. TUI: Restore State on Launch

- [x] 3.1 Write BDD scenarios for state restore (restore selected object, restore expanded groups, restore panel dimensions with clamp, restore focus panel)
- [x] 3.2 Implement BDD step definitions for restore scenarios
- [x] 3.3 Add `loadState()` function called from `Start()`, apply restored values before building initial model (make BDD scenarios pass)
- [x] 3.4 Add unit tests for restore edge cases (terminal size changed, all groups collapsed in state)

## 4. TUI: Fallback Logic

- [x] 4.1 Write BDD scenarios for fallback (object deleted with same type existing, object deleted with type removed, stale type groups in expandedGroups ignored)
- [x] 4.2 Implement BDD step definitions for fallback scenarios
- [x] 4.3 Implement fallback logic in `loadState()` / `Start()` — locate object by ID, fall back to same type group first object, then overall first object (make BDD scenarios pass)
- [x] 4.4 Add unit tests for fallback edge cases (empty vault, single object vault, object ID with no type match)
