## 1. Create State Machine Foundation

- [x] 1.1 Write unit tests for createState struct and step transitions (createStepTemplate → createStepName, mode determination, reset behavior)
- [x] 1.2 Add `createState` struct, `createMode`, and `createStep` types to `tui/app.go`; replace `newObjectMode`, `newObjectName`, `newObjectType` fields with a single `create *createState` field
- [x] 1.3 Add `startCreate(groupIndex int, mode createMode)` method that initializes `createState` with mode, type name, fetches templates, and determines initial step
- [x] 1.4 Write unit tests for `startCreate` step determination: 0 templates → createStepName, 1 template → auto-select + createStepName, 2+ templates → createStepTemplate

## 2. Template Selection UI

- [x] 2.1 Write unit tests for template selection key handling (up/down cursor movement, enter to confirm, esc to cancel, cursor wrapping)
- [x] 2.2 Implement `updateCreateTemplateStep()` in `tui/update.go` — handle arrow keys for cursor, Enter to select template (transition to createStepName or create if name template), Esc to cancel
- [x] 2.3 Implement template list rendering in sidebar View — show "Select template:" header, template names with cursor indicator, "(none)" option at end
- [x] 2.4 Write unit test verifying "(none)" selection sets empty template and proceeds to name step

## 3. Name Input Step

- [x] 3.1 Write unit tests for name input in Create & Edit mode (Enter creates object and enters edit mode, Esc cancels)
- [x] 3.2 Write unit tests for name input in Quick Create mode (Enter creates object and clears input, Esc exits and selects last object)
- [x] 3.3 Implement `updateCreateNameStep()` in `tui/update.go` — handle Enter (create object with selected template), Esc (cancel or exit batch), text input delegation
- [x] 3.4 Implement name input rendering in sidebar View — show "New {type}: {input}" with mode-appropriate help bar

## 4. Name Template Auto-Skip

- [x] 4.1 Write unit tests for name template auto-skip in Create & Edit mode (type with name template skips name input, object created with auto-generated name)
- [x] 4.2 In `startCreate()`, detect name template on the type schema; in Create & Edit mode, if name template exists and template step is resolved, skip to immediate creation
- [x] 4.3 Write unit test verifying Quick Create mode always shows name input even when name template exists

## 5. Post-Creation Behavior

- [x] 5.1 Write unit tests for Create & Edit post-creation (object selected, body focused, edit mode active)
- [x] 5.2 Implement Create & Edit post-creation: select object, set `focus = focusBody`, set `editMode = true`, position body textarea cursor
- [x] 5.3 Write unit tests for Quick Create post-creation (flash message displayed, input cleared and refocused, last object selected on Esc)
- [x] 5.4 Implement Quick Create post-creation: set flash message on `createState`, clear input, track last created object; on Esc select it
- [x] 5.5 Implement flash message rendering and auto-dismiss via `tea.Tick` (2 second timeout)

## 6. Error Handling and Validation

- [x] 6.1 Write unit tests for unique constraint error display (error shown inline, clears on input change)
- [x] 6.2 Write unit tests for empty name rejection (no object created, input remains focused)
- [x] 6.3 Implement inline error display in `createState.errMsg` — render below input, clear on key press that modifies input
- [x] 6.4 Capture error from `vault.Objects.Create()` and set `createState.errMsg` instead of `model.saveErr`

## 7. Help Bar and Keybindings

- [x] 7.1 Update `updateNormal()` to handle `N` key — call `startCreate(groupIndex, createModeBatch)`
- [x] 7.2 Route key messages through `updateCreate()` dispatcher when `create != nil` (replaces `updateNewObject` in the mode priority chain)
- [x] 7.3 Update help bar rendering for all creation states: template selection hints, name input hints per mode, normal mode `n`/`N` hints
- [x] 7.4 Write unit test verifying `n` and `N` keybindings are ignored in read-only mode

## 8. Cleanup and Integration

- [x] 8.1 Remove old `newObjectMode`, `newObjectName`, `newObjectType` fields and `updateNewObject()` function
- [x] 8.2 Update `Update()` mode priority chain to use `create != nil` check instead of `newObjectMode`
- [x] 8.3 Run full test suite (`go test ./...`) and fix any regressions
- [x] 8.4 Manual TUI testing: verify both modes, template selection, name template skip, batch creation, error display, flash messages
