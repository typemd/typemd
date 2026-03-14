## Context

The TUI currently uses a three-panel layout (sidebar, body, properties) designed exclusively for viewing and editing objects. Type schemas are only modifiable via manual YAML editing. Issue #207 requires adding a type schema editor to the TUI, which means the right panel must support a second content mode and the core layer needs write/delete APIs for type schemas.

The existing TUI architecture uses `focusPanel` (left/body/props) for focus management and a flat `model` struct. The core layer has `ObjectRepository.GetSchema()` and `WriteSchema(name, []byte)` but no high-level Vault facade methods for type CRUD.

## Goals / Non-Goals

**Goals:**
- Enable full type schema CRUD from the TUI (create, view, edit properties, delete)
- Maintain clean separation: type editor as an independent sub-model, core API additions on Vault facade
- Keep the existing object detail view and all current keybindings working unchanged
- Save-on-operation-complete semantics (no explicit save/cancel workflow)

**Non-Goals:**
- Type rename (requires moving object files and updating references)
- Shared properties (`use` references) visual editing — displayed as resolved, not editable
- System property editing (name, description, created_at, updated_at, tags)
- Property name or type modification (only emoji and pin are editable)
- Undo/redo support

## Decisions

### 1. Right panel view mode enum

**Decision**: Introduce `rightPanelMode` enum (`panelEmpty`, `panelObject`, `panelTypeEditor`) to control what the right panel renders.

**Alternatives considered**:
- (A) Use `selected` / `selectedType` mutual exclusion — simpler but implicit, hard to extend for future panel types (settings, welcome screen)
- (B) View mode enum — explicit state machine, trivially extensible

**Rationale**: The enum makes the panel routing explicit in both `Update()` and `View()`. Adding future panel types requires only a new enum value, not new nil-checking logic.

### 2. Enter/Space key split on sidebar headers

**Decision**: `Enter` on a type header opens the type editor in the right panel. `Space` toggles expand/collapse. Currently both keys have identical behavior (toggle).

**Rationale**: Follows tree-view conventions (space=toggle, enter=open). The split gives headers a "primary action" (enter=inspect/edit) distinct from the structural action (space=expand). Object rows retain `Enter` = select object (unchanged).

### 3. Independent `typeEditor` sub-model

**Decision**: The type editor is a separate struct with its own `Update(msg) (typeEditor, tea.Cmd)` and `View() string` methods. The parent `model` holds a `*typeEditor` field; when non-nil, key events and rendering are delegated to it.

**Alternatives considered**:
- Flat fields on `model` (consistent with existing `searchMode`, `editMode`) — rejected because the type editor has substantial internal state (cursor, wizard, mode) that would bloat `model` and make testing harder

**Rationale**: Encapsulation. The type editor can be tested independently. The parent model only needs to know when to create/destroy it and how to route messages.

### 4. Unified cursor across meta fields and properties

**Decision**: A single cursor integer moves through meta fields (Name, Plural, Emoji, Unique) and then the property list. Name is at index 0 but not editable. Properties start at index 4 (after 4 meta fields).

**Rationale**: Consistent navigation (↑↓ always works the same way). The `e` key behavior adapts to cursor context: text input for Plural/Emoji, toggle for Unique, emoji input for properties, no-op for Name.

### 5. Property list split into Pinned/Properties sections

**Decision**: The property list displays two visual sections: "Pinned (Header)" for pin > 0 and "Properties" for pin = 0. Section headers are non-selectable visual separators.

**Rationale**: Mirrors the actual semantic distinction (pinned properties appear in the body panel). Inspired by Anytype's Header/Properties panel split. Makes pin status immediately visible without entering edit mode.

### 6. Pin toggle via `p` key

**Decision**: `p` key toggles a property between Pinned and Properties sections. Moving to Pinned assigns `max(existing pins) + 1`. Moving from Pinned clears pin to 0.

**Alternatives considered**:
- Edit pin value in a form field — rejected as overly complex for a simple toggle operation

**Rationale**: Most users just want to pin/unpin. Auto-assigned pin values are reordered via `m` mode. The `p` key is mnemonic and fast.

### 7. Move mode for reordering

**Decision**: `m` key enters move mode; `↑↓` moves the current property; `enter`/`esc` exits move mode. Movement works within and across sections (crossing section boundary auto-toggles pin).

**Rationale**: Avoids keybinding conflicts (↑↓ normally navigates). Move mode provides clear visual feedback (e.g., highlighted row moves with cursor). Cross-section movement unifies reorder and pin management.

### 8. Add Property wizard inline in panel

**Decision**: Multi-step wizard (Step 1: name input → Step 2: type selection → Step 3: relation config if applicable) renders inline in the right panel, replacing the property list temporarily.

**Rationale**: No modal overlay needed. The panel is already dedicated to type editing. Steps are simple enough that inline rendering works well in a terminal. `esc` cancels and returns to the property list.

### 9. Core API: TypeSchema YAML serialization

**Decision**: Add a `MarshalTypeSchema(schema *TypeSchema) ([]byte, error)` function that produces YAML matching the expected file format. Special handling for `NameTemplate`: if set, emit a `name` property entry with only `template` field.

**Rationale**: `TypeSchema` has `NameTemplate` as `yaml:"-"` (not auto-marshaled). The custom marshaler ensures round-trip fidelity: load YAML → modify in memory → write YAML produces equivalent output.

### 10. Core API: Vault.SaveType and Vault.DeleteType

**Decision**: Add `Vault.SaveType(schema *TypeSchema) error` (validate → marshal → write) and `Vault.DeleteType(name string) error` (check not built-in → delete file). Add `ObjectRepository.DeleteSchema(name string) error` to the repository interface.

**Rationale**: Follows existing patterns: Vault as facade, ObjectRepository as infrastructure. Validation before write prevents invalid schemas from being persisted.

### 11. Save-on-operation-complete

**Decision**: Each discrete operation (edit a field, add/delete property, reorder, toggle pin) triggers an immediate write via `Vault.SaveType()`. No dirty state tracking needed.

**Rationale**: Simplest mental model for users — "I did something, it's saved." No risk of losing work by forgetting to save. Each operation is atomic: if write fails, error shown in status bar, in-memory state rolls back.

### 12. Sidebar "+ New Type" row

**Decision**: Add a new `rowKind` (`rowNewType`) to `listRow`. The row always appears at the bottom of the sidebar. Selecting it prompts for a type name (inline text input in the sidebar or a minimal wizard), creates the type, and opens the editor.

**Rationale**: Discoverable entry point for creating types. Positioned at the bottom to avoid interfering with the existing type group layout.

## Risks / Trade-offs

- **[TypeSchema YAML round-trip fidelity]** → Custom marshaling logic must handle edge cases (NameTemplate, use entries, comment preservation). Mitigation: comprehensive BDD scenarios comparing input/output YAML.
- **[Cursor indexing complexity]** → Unified cursor across meta fields (fixed count) and properties (variable count) with section separators requires careful index math. Mitigation: unit tests for boundary conditions; separators excluded from cursor range.
- **[Terminal emoji width]** → Emoji rendering varies across terminals (some display as 1 char width, others as 2). Mitigation: use lipgloss width calculation; accept minor alignment issues in rare terminals.
- **[No undo]** → Save-on-operation-complete means accidental changes are immediately persisted. Mitigation: type schema files are in a git repo; users can revert. Delete type requires explicit confirmation.
