## Why

The TUI supports selecting templates during object creation with live preview, but provides no way to browse, edit, or delete templates. Users must switch to a file manager or text editor to manage template files, breaking the TUI workflow.

## What Changes

- Add template section to the Type Editor showing available templates for each type
- Add `panelTemplate` right panel mode for viewing and editing template detail (body + properties)
- Add inline template editing: body via textarea, properties via property editor
- Add template creation flow (name input → empty `.md` file)
- Add template deletion with confirmation
- Add `Vault.SaveTemplate()` and `Vault.DeleteTemplate()` methods to core

## Capabilities

### New Capabilities
- `tui-template-management`: TUI template CRUD — listing templates in type editor, template detail panel with body/props view, inline editing (body textarea + property editor), creation wizard, and delete confirmation flow

### Modified Capabilities
- `object-templates`: Add `SaveTemplate` and `DeleteTemplate` write operations to complement existing read-only `ListTemplates`/`LoadTemplate` API

## Impact

- **core/**: New `SaveTemplate` and `DeleteTemplate` methods on `Vault` (delegating to `ObjectRepository`/`LocalObjectRepository`)
- **tui/**: New `panelTemplate` right panel mode, new `templateEditor` sub-model, template section added to type editor view
- **Files**: `templates/<type>/<name>.md` — create and delete operations
