## 1. Core: Template Write Operations

- [x] 1.1 Write BDD scenarios for SaveTemplate and DeleteTemplate
- [x] 1.2 Implement step definitions for template write scenarios
- [x] 1.3 Add SaveTemplate and DeleteTemplate to ObjectRepository interface
- [x] 1.4 Implement SaveTemplate on LocalObjectRepository (file write + mkdir)
- [x] 1.5 Implement DeleteTemplate on LocalObjectRepository (file delete + empty dir cleanup)
- [x] 1.6 Add SaveTemplate and DeleteTemplate facade methods on Vault
- [x] 1.7 Add unit tests for edge cases (overwrite, empty dir cleanup, nonexistent delete)

## 2. TUI: Type Editor Template Section

- [x] 2.1 Write BDD scenarios for template listing in type editor
- [x] 2.2 Add template list state to typeEditor (template names, cursor tracking)
- [x] 2.3 Render Templates section in type editor View() (📝 prefix, "(none)" placeholder, "+ Add Template")
- [x] 2.4 Add teModeTemplateList mode with cursor navigation (j/k within templates section)
- [x] 2.5 Add Enter handler to transition from template list to panelTemplate mode

## 3. TUI: Template Detail View (panelTemplate)

- [x] 3.1 Add panelTemplate to rightPanelMode enum
- [x] 3.2 Create templateEditor sub-model struct with mode enum (tmplModeView, tmplModeEditBody, tmplModeEditProp, tmplModeDelete)
- [x] 3.3 Implement templateEditor View() for read-only mode: body viewport + properties viewport
- [x] 3.4 Implement schema-aware property list (union of template props and schema props, exclude immutable system props)
- [x] 3.5 Implement templateEditor HelpBar() with context-sensitive text per mode
- [x] 3.6 Implement title panel rendering for panelTemplate ("📝 type · template-name")
- [x] 3.7 Integrate panelTemplate into main model Update() and View() (routing, layout, Esc→back to type editor)

## 4. TUI: Template Body Editing

- [x] 4.1 Implement tmplModeEditBody: activate textarea with template body on `e` key
- [x] 4.2 Implement save-on-Esc (call Vault.SaveTemplate) and cancel-on-Ctrl+C
- [x] 4.3 Track dirty state and display save errors

## 5. TUI: Template Property Editing

- [x] 5.1 Implement props panel cursor navigation (j/k to move between properties)
- [x] 5.2 Implement tmplModeEditProp: Enter on property opens inline text input with current value
- [x] 5.3 Implement property save on Enter (update template, call Vault.SaveTemplate)
- [x] 5.4 Implement property clear (empty input removes property from frontmatter)
- [x] 5.5 Implement cancel on Esc

## 6. TUI: Template Creation

- [x] 6.1 Add teModeAddTemplate mode to type editor with name text input
- [x] 6.2 Implement creation flow: Enter confirms name, creates empty template via Vault.SaveTemplate
- [x] 6.3 Add duplicate name validation (error message if template already exists)
- [x] 6.4 Refresh template list after creation

## 7. TUI: Template Deletion

- [x] 7.1 Implement tmplModeDelete: `d` key shows confirmation prompt
- [x] 7.2 Implement `y` to confirm (call Vault.DeleteTemplate, return to type editor)
- [x] 7.3 Implement `n`/Esc to cancel (return to template view)
- [x] 7.4 Refresh template list after deletion
