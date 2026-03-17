## 1. Core: createTypeState struct and helpers

- [x] 1.1 Write BDD scenarios for type creation title panel — skipped (TUI BDD infra not yet available; covered by unit tests in create_type_test.go)
- [x] 1.2 Create `tui/create_type.go` with `createTypeState` struct (emojiInput, nameInput, pluralInput textinput.Models; errMsg; field enum for tab cycling)
- [x] 1.3 Add `createTypeField` enum (createTypeFieldEmoji, createTypeFieldName, createTypeFieldPlural) and Tab cycling logic
- [x] 1.4 Add `startCreateType()` method on model — initialize `createTypeState`, focus name field, set title panel visible
- [x] 1.5 Add `renderCreateTypeTitleContent()` — render `[emoji] new type · [name___]  plural: [plural___]` with error display
- [x] 1.6 Add `renderCreateTypePreview()` — render read-only type schema preview for right panel

## 2. Integration: Wire into TUI model

- [x] 2.1 Replace `model.newTypeMode bool` and `model.newTypeName textinput.Model` with `model.createType *createTypeState`
- [x] 2.2 Update `startNewType()` in `app.go` to call `startCreateType()` instead
- [x] 2.3 Add `updateCreateType()` handler — dispatch Enter (validate + create), Esc (cancel), Tab (cycle fields), and field-specific key routing
- [x] 2.4 Update mode priority in `Update()` (`app.go`) — add `m.createType != nil` check at the same priority level as `m.create != nil`
- [x] 2.5 Update `hasTitlePanel()` to return true when `createType != nil`
- [x] 2.6 Update title panel rendering in `View()` to call `renderCreateTypeTitleContent()` when `createType != nil`
- [x] 2.7 Update right panel rendering to show type preview when `createType != nil`

## 3. Validation and creation logic

- [x] 3.1 Implement name validation in `updateCreateType()` Enter handler — empty name, duplicate name check via `vault.ListTypes()`
- [x] 3.2 Create `TypeSchema` with name, emoji, and plural from form fields, call `vault.SaveType()`
- [x] 3.3 After successful creation: refresh data, load type, open type editor, set focus to right panel

## 4. Help bar and cleanup

- [x] 4.1 Add `renderCreateTypeHelpBar()` — display "NEW TYPE" mode with "tab: next field  enter: create  esc: cancel"
- [x] 4.2 Update help bar rendering to call `renderCreateTypeHelpBar()` when `createType != nil`
- [x] 4.3 Remove old `newTypeMode`/`newTypeName` references from sidebar rendering in `list.go`
- [x] 4.4 Add unit tests for `renderCreateTypeTitleContent()`, `renderCreateTypeHelpBar()`, and field cycling logic

## 5. Verification

- [x] 5.1 Run `go test ./tui/...` — all tests pass
- [x] 5.2 Run `go test ./...` — full test suite passes
- [x] 5.3 Run `go build ./...` — clean build
