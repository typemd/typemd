## Why

Type schemas can only be modified by manually editing YAML files in `.typemd/types/`. The TUI displays types as group headers in the sidebar but provides no way to inspect or modify them. Users need a visual way to create, view, edit, and delete types directly from the terminal.

## What Changes

- Type group headers in the sidebar become selectable via `Enter` (opens type editor in right panel); `Space` retains toggle expand/collapse behavior
- Right panel gains a new view mode (`panelTypeEditor`) alongside the existing `panelObject` mode, controlled by a `rightPanelMode` enum
- New `typeEditor` sub-model with its own `Update()`/`View()` methods, featuring:
  - Unified cursor across meta fields (Name, Plural, Emoji, Unique) and property list
  - Property list split into **Pinned (Header)** and **Properties** sections (mirroring pin semantics)
  - Inline editing for safe meta fields (Plural, Emoji as text input; Unique as toggle); Name is read-only
  - `p` key to toggle pin (move property between Pinned/Properties sections, auto-assign pin values)
  - `e` key for inline emoji editing on properties
  - `m` mode + `↑↓` for reordering properties within and across sections
  - `a` key for multi-step Add Property wizard (name → type → relation config) rendered inline in the panel
  - `d` key for delete property with confirmation
- "**+ New Type**" item at the bottom of the sidebar; creates a new type (prompts for name) and opens editor
- Delete type with confirmation dialog showing existing object count; only `tag` (the sole built-in type) is protected
- Save-on-operation-complete: each discrete edit (field change, add/delete property, reorder) writes YAML immediately
- Core API additions: `Vault.SaveType()`, `Vault.DeleteType()`, `Vault.CountObjectsByType()`, TypeSchema YAML serialization

## Capabilities

### New Capabilities
- `type-editor-tui`: TUI type schema editor panel with full CRUD for types and their properties
- `type-crud-api`: Core API for saving, deleting, and counting types (Vault facade methods + TypeSchema serialization)

### Modified Capabilities
- `tui-layout`: Right panel now supports multiple view modes (object detail vs type editor); sidebar Enter/Space key semantics change
- `type-schema`: TypeSchema gains YAML serialization (struct → file) and delete capability

## Impact

- **tui/**: New `type_editor.go` sub-model; changes to `app.go` (rightPanelMode enum, Update/View routing), `update.go` (Enter/Space key split), `list.go` (new rowKind for "+ New Type")
- **core/**: New methods on Vault/ObjectService; new `ObjectRepository.DeleteSchema()` interface method; `LocalObjectRepository` gains delete and improved write support; TypeSchema YAML marshaling logic
- **core/type_schema.go**: Serialization must handle NameTemplate → name property entry round-trip
- **No breaking changes**: Existing CLI commands and data model are unaffected
