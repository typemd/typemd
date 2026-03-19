## 1. TUI: Extend SessionState struct

- [x] 1.1 Write BDD scenarios for view mode state persistence (save on quit, restore on launch, fallbacks)
- [x] 1.2 Implement BDD step definitions for view mode persistence scenarios
- [x] 1.3 Add view mode fields to SessionState struct (`ViewTypeName`, `ViewName`, `ViewCursor`, `ViewScroll`, `ViewExpandedGroups`)
- [x] 1.4 Add unit tests for SessionState YAML serialization (view fields present/absent with omitempty)

## 2. TUI: Capture view mode state

- [x] 2.1 Add `expandedGroupLabels()` method to viewMode to expose expanded group labels
- [x] 2.2 Update `captureState()` to populate view fields when `rightPanel == panelView`
- [x] 2.3 Add unit tests for captureState in view mode vs sidebar mode

## 3. TUI: Restore view mode on startup

- [x] 3.1 Add view mode restoration logic in `applySessionState()` / `Start()` flow
- [x] 3.2 Implement fallback chain: type deleted → sidebar, view deleted → default view → sidebar
- [x] 3.3 Implement cursor/scroll clamping to valid range
- [x] 3.4 Implement view expanded groups matching with silent ignore for stale labels
- [x] 3.5 Add unit tests for fallback scenarios (type deleted, view deleted, cursor out of range, stale groups)

## 4. Verification

- [x] 4.1 Run full test suite and verify all BDD scenarios pass
- [ ] 4.2 Manual verification: enter view mode → quit → relaunch → confirm view mode restored
