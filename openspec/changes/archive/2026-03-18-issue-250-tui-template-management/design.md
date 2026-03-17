## Context

Templates are Markdown files (`templates/<type>/<name>.md`) with YAML frontmatter property overrides and body content. The core layer currently supports read-only operations (`ListTemplates`, `LoadTemplate`) used during object creation. The TUI has a live preview feature for templates during object creation, but no way to browse, manage, or edit templates.

The type editor sub-model (`typeEditor` in `tui/type_editor.go`) provides an established pattern for independent panel-level sub-models with their own mode state machines.

## Goals / Non-Goals

**Goals:**
- Users can view all templates for a type from within the type editor
- Users can view template detail (body + properties) in a dedicated panel
- Users can create new templates (name input → empty `.md` file)
- Users can edit template body (textarea) and properties (inline) within the TUI
- Users can delete templates with confirmation
- Core provides `SaveTemplate` and `DeleteTemplate` operations

**Non-Goals:**
- Template renaming (requires updating references; future work)
- Template duplication / copy-from-object
- Cross-type template sharing
- Template versioning or migration
- File-watch for external template edits (can be added later like object file watch)

## Decisions

### 1. New `panelTemplate` right panel mode

**Decision:** Add `panelTemplate` as a fourth `rightPanelMode` value, alongside `panelEmpty`, `panelObject`, `panelTypeEditor`.

**Rationale:** Template detail view needs body + props layout similar to `panelObject`, but with different data source (template vs object) and different edit semantics (no ULID, no system properties like `created_at`). A dedicated panel mode keeps the concerns separated.

**Alternative considered:** Embedding template view as a nested mode inside `typeEditor`. Rejected because the template detail view is structurally a body+props panel (like object detail), not a type schema editor. Mixing these would bloat `typeEditor` which already has 7 modes.

### 2. `templateEditor` sub-model pattern

**Decision:** Create a `templateEditor` sub-model following the `typeEditor` pattern — independent struct with `Update()`, `View()`, `HelpBar()`, `CanQuit()` methods.

**State:**
```
templateEditor {
    typeName   string
    template   *Template
    schema     *TypeSchema
    vault      *Vault

    mode       templateEditorMode
    // View mode: read-only body + props display
    // Edit body mode: textarea editing
    // Edit props mode: property value editing
    // Delete mode: confirmation

    bodyViewport   viewport.Model
    bodyTextarea   textarea.Model
    propsViewport  viewport.Model
    propsCursor    int
    propEditInput  textinput.Model

    dirty    bool
    saveErr  string
}
```

**Modes:**
- `tmplModeView` — read-only display of body and properties
- `tmplModeEditBody` — body editing via textarea (mirrors object edit mode)
- `tmplModeEditProp` — single property value editing via text input
- `tmplModeDelete` — delete confirmation

**Rationale:** Following the established `typeEditor` pattern ensures consistency and keeps template editing logic isolated from the main model.

### 3. Template section in type editor

**Decision:** Add a "Templates" section to the type editor view, between the properties section and the delete type action. Display template names with 📝 prefix. Enter on a template transitions to `panelTemplate` mode.

**Navigation flow:**
```
Sidebar (type header) → Type Editor → Template list → Template Detail
         ←Esc──────────── ←Esc──────── ←Esc──────────
```

The type editor gains a new mode `teModeTemplateList` for navigating the templates section, plus `teModeAddTemplate` for the creation flow (name input).

### 4. Template body + props layout reuse

**Decision:** The `panelTemplate` view uses the same body+props split layout as `panelObject`. The title panel shows `📝 type · template-name`. Body panel shows template markdown. Props panel shows template frontmatter property values.

**Rationale:** Consistent visual language — users already understand the body+props layout from object viewing.

### 5. Core write operations on `ObjectRepository`

**Decision:** Add `SaveTemplate(typeName, name string, tmpl *Template) error` and `DeleteTemplate(typeName, name string) error` to `ObjectRepository` interface. `LocalObjectRepository` implements these as file write/delete. `Vault` exposes them as facade methods.

**SaveTemplate serialization:** Uses the same frontmatter + body format as object files. Properties are serialized as YAML frontmatter, body as markdown content after the `---` delimiter.

**Rationale:** Template write operations mirror the existing read operations (`GetTemplate`, `ListTemplates`) and follow the repository pattern.

### 6. Property editing approach

**Decision:** Template properties are edited one-at-a-time using a text input overlay (similar to type editor's meta field editing). The props panel shows property names and values; moving cursor to a property and pressing Enter opens an inline text input for that value.

Properties shown are the union of: (a) properties already in the template frontmatter, and (b) all properties defined in the type schema. Schema-defined properties not in the template show as empty/placeholder. Immutable system properties (`created_at`, `updated_at`) are excluded.

**Rationale:** Keeps the editing experience simple and consistent. Full property type awareness (select options, relation targets) is valuable but can be added incrementally.

## Risks / Trade-offs

**[Risk] Type editor complexity grows** → Mitigated by keeping template detail in a separate `templateEditor` sub-model. Type editor only gains template list rendering and navigation to `panelTemplate`.

**[Risk] Template save format inconsistency** → Mitigated by reusing the existing frontmatter serialization (`yaml.Marshal` + body concatenation) used for objects.

**[Risk] No file watch for external template edits** → Accepted trade-off. Templates are edited less frequently than objects. Users can re-enter the template to see updated content. File watch can be added as a follow-up.

**[Risk] Property editing is text-only (no type-aware widgets)** → Accepted for v1. Select options, relation pickers, and date pickers are future enhancements. Text input works for all property types as a baseline.
