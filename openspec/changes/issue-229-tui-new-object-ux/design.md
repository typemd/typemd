## Context

The TUI currently has a single creation mode triggered by `n` — it opens an inline text input for the object name and creates the object with no template. The `model` struct tracks this with `newObjectMode bool`, `newObjectName textinput.Model`, and `newObjectType string`. All creation logic lives in `updateNewObject()` in `tui/update.go`.

Core already provides all necessary APIs: `ObjectService.Create(type, name, template)`, `Vault.ListTemplates(type)`, `Vault.LoadTemplate(type, name)`, and name template evaluation via `EvaluateNameTemplate()`.

## Goals / Non-Goals

**Goals:**
- Two distinct creation modes: `n` (Create & Edit) and `N` (Quick Create / batch)
- Template selection UI when type has 2+ templates; auto-apply for 1 template; skip for 0
- Auto-skip name input when type defines a name template (Create & Edit mode only)
- Inline error feedback for unique constraint violations
- Visual confirmation on creation success (flash message in batch mode)

**Non-Goals:**
- Changing core `ObjectService.Create()` API — TUI passes template name, core handles the rest
- Adding new Bubble Tea dependencies — use existing textinput + manual list rendering
- Supporting template preview (showing template contents before selection)
- Inline editing of template properties during creation

## Decisions

### 1. Multi-step creation state machine

Replace the flat `newObjectMode bool` with a struct-based state machine that tracks the current step in the creation flow.

```
type createStep int

const (
    createStepTemplate createStep = iota  // selecting template (if applicable)
    createStepName                         // entering name (if no name template)
)

type createMode int

const (
    createModeNone   createMode = iota
    createModeSingle                    // n: Create & Edit
    createModeBatch                     // N: Quick Create
)

type createState struct {
    mode      createMode
    step      createStep
    typeName  string
    templates []string       // available templates for this type
    cursor    int            // template selection cursor
    template  string         // selected template name
    nameInput textinput.Model
    flash     string         // success message (batch mode)
    flashTime time.Time      // when flash was set
    errMsg    string         // validation error
}
```

**Rationale**: A state machine is cleaner than adding more booleans. The `createState` struct encapsulates all creation flow state, making it easy to reset and reason about. This follows the pattern used by `typeEditor` — an independent sub-model with its own state.

**Alternative considered**: Separate `newObjectMode` booleans for each step. Rejected — becomes a mess with 2 modes × 2 steps × error states.

### 2. Template selection as a simple cursor list

Template selection uses arrow keys to move a cursor through a list rendered inline in the sidebar (below the group list). No new Bubble Tea component needed — just render the list and handle up/down/enter/esc.

```
  ┌─────────────────────────┐
  │ ▼ 📚 book (3)           │
  │   clean-code            │
  │   design-patterns       │
  │   go-programming        │
  │                         │
  │ Select template:        │
  │   > review              │  ← cursor here
  │     summary             │
  │     (none)              │  ← always last option
  │                         │
  │ enter: select  esc: cancel│
  └─────────────────────────┘
```

The list always includes a `(none)` option at the end so users can skip template selection.

**Rationale**: Matches existing TUI patterns (sidebar is text-rendered, no complex widgets). Template lists are typically small (1-5 items).

**Alternative considered**: Popup modal dialog. Rejected — the TUI doesn't have a modal system and adding one is out of scope.

### 3. Batch mode keeps input focused after each creation

In Quick Create mode (`N`), after pressing Enter to create an object:
1. Show a flash message (e.g., `✓ Created: my-book`) above the input
2. Clear the text input
3. Keep the text input focused for the next object
4. Template selection (if done) persists — all objects in the batch use the same template
5. Esc exits batch mode and selects the last created object

**Rationale**: The whole point of batch mode is speed. Re-selecting templates per object defeats the purpose.

### 4. Name template auto-skip (Create & Edit mode only)

When a type has a name template and the mode is Create & Edit (`n`):
1. Skip the name input step entirely
2. Pass empty string as filename to `ObjectService.Create()` — core handles name template evaluation
3. Proceed directly to object creation → enter edit mode

In Quick Create mode (`N`), name templates are ignored — the user always types a name. This is because batch creating auto-named objects (e.g., journals) doesn't make sense.

**Rationale**: Name templates exist for types like "journal" where the name is always a date. Forcing the user to type a name defeats the purpose. But in batch mode, auto-naming would create identical names.

### 5. Error display as inline text below input

Validation errors (unique constraint, empty name) display as a styled error line below the text input, replacing the help bar temporarily.

```
  New book: duplicate-name█
  ✗ book name "duplicate-name" already exists
```

The error clears when the user modifies the input.

**Rationale**: Simpler than a toast/notification system. Error is contextually relevant (right next to the input). Follows the principle of least surprise.

### 6. Flash message with auto-dismiss

The success flash in batch mode uses a `time.After` tick to auto-dismiss after 2 seconds. This uses Bubble Tea's `tea.Tick` command pattern.

**Rationale**: Flash must be temporary to not clutter the UI. 2 seconds is enough to confirm without slowing down.

## Risks / Trade-offs

- **Risk: State complexity** — Two modes × multi-step flow adds state to track. → Mitigation: Encapsulate all state in `createState` struct with clear transitions.
- **Risk: Flash timing in batch mode** — If user creates objects faster than 2s, flashes stack. → Mitigation: Replace (don't stack) flash messages. Only show the latest.
- **Trade-off: No template preview** — Users can't see template contents before selecting. Acceptable for v1 — template names should be descriptive enough.
- **Trade-off: Batch mode ignores name template** — Could be surprising for types that define one. → Mitigation: Show a hint like "name template available, use `n` instead" if applicable.
