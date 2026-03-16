## 1. Refactor createState for concurrent fields

- [x] 1.1 Write unit tests for createField enum and Tab switching between createFieldName and createFieldTemplate
- [x] 1.2 Replace `createStep` with `createField` in createState; remove sequential template→name step flow; add `field`, `previewBody`, `previewProps`, and `templateCache` fields
- [x] 1.3 Refactor `startCreate()` to initialize both name input and template simultaneously (no more step-based branching); pre-fill name from name template when applicable
- [x] 1.4 Write unit tests for startCreate: 0 templates (no template field), 1 template (auto-select, static), 2+ templates (interactive), name template pre-fill

## 2. Title panel creation form rendering

- [x] 2.1 Write unit tests for `renderCreateTitleContent()` output: emoji + type + name input + template selector format
- [x] 2.2 Implement `renderCreateTitleContent()` that renders `📚 book · [nameInput] 📝 template` with focused field highlighting
- [x] 2.3 Update `View()` in app.go: when `m.create != nil`, render creation form in title panel instead of static title; ensure `hasTitlePanel()` returns true during creation
- [x] 2.4 Remove `renderCreateUI()` call from sidebar View rendering (left panel no longer shows creation UI)

## 3. Live template preview

- [x] 3.1 Write unit tests for `updateCreatePreview()`: template body → body viewport, template props → props viewport, (none) → empty/defaults
- [x] 3.2 Implement `updateCreatePreview()` that loads template via `Vault.LoadTemplate()`, builds preview body and display properties, updates viewports
- [x] 3.3 Implement `buildTemplatePreviewProps()` that merges template frontmatter with schema defaults for display
- [x] 3.4 Call `updateCreatePreview()` from `startCreate()` (initial preview) and on every template switch

## 4. Key handling for title panel fields

- [x] 4.1 Write unit tests for Tab switching: name→template, template→name, Tab with 0 templates (no-op)
- [x] 4.2 Write unit tests for ←/→ template cycling when template field focused (wrapping, preview update)
- [x] 4.3 Refactor `updateCreate()` to dispatch based on `createField` instead of `createStep`; handle Tab, ←/→ for template cycling, text input for name
- [x] 4.4 Write unit tests for Enter from both fields (creates object), Esc from both fields (cancels)

## 5. Help bar updates

- [x] 5.1 Update `renderCreateHelpBar()` to show field-specific hints: name focused shows "tab: template", template focused shows "◀▶: switch  tab: name"
- [x] 5.2 Write unit tests verifying help bar content for each field focus state and mode combination

## 6. Batch mode and flash in title panel

- [x] 6.1 Write unit tests for batch mode: flash displays in title panel area, name clears after creation, template persists
- [x] 6.2 Update batch mode post-creation to show flash in title panel (integrate with `renderCreateTitleContent`)
- [x] 6.3 Write unit test for name template pre-fill in batch mode (pre-fill shown but editable, not auto-skipped)

## 7. E2E integration tests

- [x] 7.1 Update E2E tests: `TestE2E_CreateAndEdit_WithTemplateSelection` to verify title panel rendering and live preview
- [x] 7.2 Update E2E tests: `TestE2E_QuickCreate_BatchFlow` to verify flash in title panel
- [x] 7.3 Add E2E test: Tab between name and template, ←/→ cycles template, body viewport updates
- [x] 7.4 Run full test suite (`go test ./...`) and fix any regressions
