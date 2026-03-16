## Context

The TUI title panel is a 3-line bordered panel (`titlePanelHeight = 3`) above the body and properties panels. It currently renders `📚 book · Clean Code` via `renderTitleContent()`. During creation (when `m.create != nil`), the sidebar currently appends `renderCreateUI()` at the bottom showing template selection or name input.

The `createState` struct already tracks mode, step, templates, cursor, template selection, nameInput, flash, and error state. The `createStep` enum has `createStepTemplate` and `createStepName` as sequential steps.

## Goals / Non-Goals

**Goals:**
- Move name input and template selector into the title panel as concurrent fields (not sequential steps)
- Live preview: switching templates updates body viewport and properties viewport in real time
- Remove all creation UI from the sidebar
- Preserve all existing creation logic (executeCreate, tryNameTemplateSkip, batch mode, flash, errors)

**Non-Goals:**
- Changing the title panel height (it stays at 3 lines — 1 content line is sufficient for `[name input] 📝 template`)
- Adding template content editing during creation
- Changing core/ APIs

## Decisions

### 1. Replace sequential steps with concurrent fields

Remove `createStepTemplate` and `createStepName`. Instead, add a `createField` enum to track which field has focus within the title panel:

```go
type createField int

const (
    createFieldName     createField = iota // name text input (default)
    createFieldTemplate                     // template cycling selector
)
```

The `createState` struct gets:
- `field createField` — which field is focused
- `previewBody string` — cached template body for viewport
- `previewProps []core.DisplayProperty` — cached template properties for viewport

On `startCreate`, both name input and template are initialized simultaneously. The title panel always shows both (when templates exist).

### 2. Title panel creation form layout

```
╭──────────────────────────────────────────╮
│ 📚 book · [name█              ] 📝 review│
╰──────────────────────────────────────────╯
```

When `createFieldName` is focused, the text input is active. When `createFieldTemplate` is focused, `←`/`→` cycles templates. `Tab` switches between fields.

For types with 0 templates: no template selector shown.
For types with 1 template: template shown but not interactive (auto-selected, grayed out or static).
For types with 2+ templates: template is interactive.

The template display format: `📝 <name>` where 📝 is a fixed icon. For `(none)`: show `📝 (none)`.

### 3. Live template preview

When the template selection changes (via `←`/`→`), load the template and update the body and properties viewports:

```go
func (m *model) updateCreatePreview() {
    cs := m.create
    tmplName := cs.selectedTemplateName()
    if tmplName != "" {
        tmpl, _ := m.vault.LoadTemplate(cs.typeName, tmplName)
        if tmpl != nil {
            cs.previewBody = tmpl.Body
            // Build display props from template frontmatter + schema defaults
            cs.previewProps = m.buildTemplatePreviewProps(cs.typeName, tmpl)
        }
    } else {
        cs.previewBody = ""
        cs.previewProps = m.buildSchemaDefaultProps(cs.typeName)
    }
    // Update viewports
    m.bodyViewport.SetContent(renderPreviewBody(cs.previewBody))
    m.propsViewport.SetContent(renderPreviewProps(cs.previewProps))
}
```

This is called on `startCreate` (initial preview) and on every template switch.

### 4. Key handling changes

| Key | createFieldName focused | createFieldTemplate focused |
|-----|------------------------|---------------------------|
| Text keys | Input to nameInput | Ignored |
| `Tab` | Switch to createFieldTemplate | Switch to createFieldName |
| `←`/`→` | Passed to nameInput (cursor move) | Cycle template |
| `↑`/`↓` | Ignored | Cycle template (alternative) |
| `Enter` | Create object | Create object |
| `Esc` | Cancel (or exit batch) | Cancel (or exit batch) |

### 5. Sidebar cleanup

Remove the `renderCreateUI()` call from the sidebar View rendering. The sidebar shows the normal object list during creation — no appended UI at the bottom.

### 6. hasTitlePanel during creation

`hasTitlePanel()` must return `true` when `m.create != nil` (even if no object is selected), so the title panel space is allocated for the creation form.

### 7. Name template pre-fill instead of auto-skip

For types with name templates, instead of skipping the name input entirely, pre-fill the name input with the evaluated template value. The user can accept it (Enter) or edit it. This replaces `nameTemplateSkip` with `nameTemplatePrefill`.

## Risks / Trade-offs

- **Risk: Title panel width** — On narrow terminals, `[name input] 📝 template` might not fit. → Mitigation: Template selector truncates or hides on very narrow widths.
- **Risk: Template loading on every switch** — `LoadTemplate()` reads from disk on each `←`/`→`. → Mitigation: Template files are small; cache loaded templates in `createState.templateCache map[string]*core.Template`.
- **Trade-off: No visual cursor on template selector** — Unlike the old list with `>` cursor, the cycling selector just shows the current value. The user must know `←`/`→` cycles it. → Mitigation: Help bar shows `◀▶: template` hint.
