## Why

The current TUI object creation flow (#229) renders name input and template selection at the bottom of the sidebar, which is cramped and disconnected from where the object will actually appear. Users cannot preview template content before selecting, making template choice a blind guess. Moving the creation UI to the title panel — where the object's identity will ultimately display — creates a "what you see is where it goes" experience and frees the sidebar from creation-specific UI.

## What Changes

- When `n`/`N` is pressed, the title panel transforms into an inline creation form with a name text input and a template cycling selector (e.g., `📚 book · [name█] 📝 review`)
- `Tab` switches focus between the name field and the template selector; `←`/`→` cycles templates
- Body and properties panels update in real time to preview the selected template's content and frontmatter properties
- Sidebar no longer renders any creation UI (no `renderCreateUI` at the bottom)
- The `createState` struct and creation logic (`executeCreate`, `updateCreate`) are preserved; only the rendering and input location change
- Template selection is no longer a separate step — name input and template selector coexist in the title panel

## Capabilities

### Modified Capabilities
- `tui-object-creation`: Replace sidebar-bottom creation UI with title-panel inline form; add live template preview in body/properties panels; replace sequential template→name steps with concurrent title-panel fields
- `tui-layout`: Title panel renders differently during creation mode (inline form instead of static title)

## Impact

- `tui/create.go` — Refactor `createState` to remove step-based flow; add `createField` (name vs template) for Tab switching; add preview state. Refactor `renderCreateUI` → `renderCreateTitlePanel`. Remove sidebar rendering.
- `tui/app.go` — Title panel View rendering branches on `m.create != nil` to show creation form. Body/props panels show template preview during creation.
- `tui/update.go` — No structural change; `updateCreate` dispatches to the same handlers but with different field focus logic.
- No core/ changes needed — all template loading and object creation APIs remain the same.
