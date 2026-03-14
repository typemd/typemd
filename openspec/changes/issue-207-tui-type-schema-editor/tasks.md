## 1. Core: TypeSchema YAML Serialization

- [x] 1.1 Write BDD scenarios for TypeSchema serialization (complete schema, NameTemplate, omit zero-value fields, round-trip fidelity)
- [x] 1.2 Implement BDD step definitions for serialization scenarios
- [x] 1.3 Implement `MarshalTypeSchema(schema *TypeSchema) ([]byte, error)` function in `core/type_schema.go`
- [x] 1.4 Add unit tests for serialization edge cases (empty properties, relation with all fields, select with options, NameTemplate round-trip)

## 2. Core: ObjectRepository.DeleteSchema

- [x] 2.1 Write BDD scenarios for DeleteSchema (delete existing, delete non-existent)
- [x] 2.2 Implement BDD step definitions for DeleteSchema scenarios
- [x] 2.3 Add `DeleteSchema(name string) error` to `ObjectRepository` interface in `core/object_repository.go`
- [x] 2.4 Implement `DeleteSchema` in `LocalObjectRepository` in `core/local_object_repository.go`
- [x] 2.5 Add unit tests for DeleteSchema edge cases

## 3. Core: Vault.SaveType, DeleteType, CountObjectsByType

- [x] 3.1 Write BDD scenarios for SaveType (valid save, validation failure, overwrite existing)
- [x] 3.2 Write BDD scenarios for DeleteType (user-defined type, built-in type rejected, non-existent type)
- [x] 3.3 Write BDD scenarios for CountObjectsByType (type with objects, type with no objects)
- [x] 3.4 Implement BDD step definitions for Vault type CRUD scenarios
- [x] 3.5 Implement `Vault.SaveType(schema *TypeSchema) error` — validate, marshal, write
- [x] 3.6 Implement `Vault.DeleteType(name string) error` — check built-in, delegate to repo
- [x] 3.7 Implement `Vault.CountObjectsByType(typeName string) (int, error)` — query index
- [x] 3.8 Add unit tests for edge cases (save empty schema, delete tag, count with empty index)

## 4. TUI: Right Panel View Mode Enum and Sidebar Key Split

- [x] 4.1 Add `rightPanelMode` enum (`panelEmpty`, `panelObject`, `panelTypeEditor`) to `tui/app.go`
- [x] 4.2 Add `rightPanel` field to `model` struct and update `View()` to route rendering based on panel mode
- [x] 4.3 Split `Enter`/`Space` behavior in `updateNormal()`: `Enter` on header sets `panelTypeEditor`, `Space` toggles expand/collapse
- [x] 4.4 Add `rowNewType` kind to `listRow` and render "+ New Type" at bottom of sidebar in `list.go`
- [x] 4.5 Add unit tests for panel mode transitions and Enter/Space key split behavior

## 5. TUI: Type Editor Sub-model — View Mode

- [x] 5.1 Create `tui/type_editor.go` with `typeEditor` struct, `newTypeEditor()`, `Update()`, and `View()` methods
- [x] 5.2 Implement meta fields display (Name, Plural, Emoji, Unique) with unified cursor navigation
- [x] 5.3 Implement property list display split into Pinned (Header) and Properties sections with section separators
- [x] 5.4 Wire type editor into parent model: create on `Enter` header, destroy on `Esc`, delegate Update/View
- [x] 5.5 Implement help bar context switching (show type editor keybindings when editor is active)
- [x] 5.6 Add unit tests for cursor navigation (meta fields, properties, section separator skipping, boundary clamping)

## 6. TUI: Type Editor — Meta Field Editing

- [x] 6.1 Implement Plural field inline text input (e to edit, enter to confirm, esc to cancel) with save-on-confirm
- [x] 6.2 Implement Emoji field inline text input with save-on-confirm
- [x] 6.3 Implement Unique field toggle (e to toggle, immediate save)
- [x] 6.4 Implement Name field no-op on `e` key
- [x] 6.5 Add unit tests for each meta field edit operation (confirm, cancel, toggle)

## 7. TUI: Type Editor — Property Emoji Edit and Pin Toggle

- [x] 7.1 Implement property emoji inline input (e on property → emoji input → enter saves, esc cancels)
- [x] 7.2 Implement `p` key pin toggle: unpin (clear pin to 0) / pin (assign max+1), save on toggle
- [x] 7.3 Ensure property list re-renders correctly after pin toggle (property moves between sections)
- [x] 7.4 Add unit tests for emoji edit and pin toggle (pin first, unpin, pin value assignment)

## 8. TUI: Type Editor — Move Mode (Reorder)

- [x] 8.1 Implement move mode: `m` key enters, `↑↓` swaps property position, `enter`/`esc` exits and saves
- [x] 8.2 Implement cross-section movement (moving property from Properties to Pinned auto-sets pin, and vice versa)
- [x] 8.3 Update help bar to show move mode keybindings
- [x] 8.4 Add unit tests for reorder within section, across sections, and boundary conditions

## 9. TUI: Type Editor — Add Property Wizard

- [x] 9.1 Implement wizard Step 1: property name text input with duplicate name validation
- [x] 9.2 Implement wizard Step 2: property type selection list (string, number, date, datetime, url, checkbox, select, multi_select, relation)
- [x] 9.3 Implement wizard Step 2b: options input for select/multi_select types
- [x] 9.4 Implement wizard Step 3: relation config (target type selector, multiple toggle, bidirectional toggle, inverse name input)
- [x] 9.5 Wire wizard completion to save (append new property to schema, call SaveType)
- [x] 9.6 Add unit tests for each wizard step, cancellation, and validation

## 10. TUI: Type Editor — Delete Property and Delete Type

- [x] 10.1 Implement delete property: `d` on property → confirmation prompt → remove property → save
- [x] 10.2 Implement delete type: keybinding → confirmation with object count → call DeleteType → close editor
- [x] 10.3 Implement built-in type protection (tag cannot be deleted, show error message)
- [x] 10.4 Add unit tests for delete property confirm/cancel, delete type with objects, delete built-in rejection

## 11. TUI: New Type Creation

- [x] 11.1 Implement "+ New Type" selection: inline name input in sidebar → validate unique name → call SaveType with empty schema → open type editor
- [x] 11.2 Rebuild sidebar groups after type creation (refresh from vault)
- [x] 11.3 Add unit tests for new type creation, duplicate name rejection, and cancel
