## Context

The TUI's "New Type" creation currently uses a minimal inline text input at the bottom of the sidebar (`newTypeMode` + `newTypeName` fields on the model). This only collects the type name, creating an empty `TypeSchema` that the user must then populate via the type editor.

The "New Object" creation flow (PR #242) already follows a title panel pattern with multi-field input and live preview. This design brings the "New Type" flow to the same level of polish by reusing the same architectural pattern.

### Current implementation

- `model.newTypeMode bool` — flag to indicate type creation is active
- `model.newTypeName textinput.Model` — single text input for name
- `updateNewType()` in `update.go` — handles Enter (create) and Esc (cancel)
- `startNewType()` in `app.go` — initializes the text input
- Sidebar renders the input inline at the bottom of the left panel

### Reference implementation (New Object)

- `model.create *createState` — struct holding all creation state
- `createState` has `field createField` for Tab-switching between name/template
- Title panel renders via `renderCreateTitleContent()` in `create.go`
- Live preview updates body and properties viewports
- Help bar adapts to current creation context

## Goals / Non-Goals

**Goals:**
- Move type creation input to the title panel, matching the object creation UX pattern
- Support three fields: emoji (optional), name (required), plural (optional)
- Tab navigation between fields
- Live preview of the resulting type schema in the right panel
- Inline validation with error messages
- After creation, open the type editor with emoji/plural pre-populated

**Non-Goals:**
- Adding properties during creation (that stays in the type editor)
- Changing the core `SaveType()` API or `TypeSchema` validation
- Batch/quick create mode for types (unlike objects, types are created infrequently)
- Starter type template selection during creation (that belongs to `tmd init`)

## Decisions

### 1. Introduce `createTypeState` struct (parallel to `createState`)

**Decision:** Create a new `createTypeState` struct in a new file `tui/create_type.go`, mirroring the pattern of `tui/create.go` for objects.

**Rationale:** The object creation and type creation flows share the title panel pattern but have different fields and behavior. A separate struct avoids polluting `createState` with type-specific concerns. Keeping them in separate files follows the existing code organization.

**Alternatives considered:**
- Extending `createState` with a discriminator field — rejected because the two flows share almost no state (templates vs emoji/plural)
- Using the existing `newTypeMode`/`newTypeName` fields — rejected because it cannot support multi-field input or title panel rendering

### 2. Three-field layout: emoji · name · plural

**Decision:** The title panel renders `[emoji] new type · [name___]  plural: [plural___]` with Tab cycling through emoji → name → plural → emoji.

**Rationale:** This matches the visual density of the object creation title panel. The emoji field comes first (visually leftmost, matching how type headers display emoji), followed by name (required, primary), then plural (optional, secondary).

**Field order for Tab navigation:** emoji → name → plural (wraps around). Name is focused by default since it's the required field.

### 3. Live preview in right panel shows type editor-style view

**Decision:** During type creation, the right panel shows a read-only preview of what the type schema will look like, rendered similarly to the type editor's view mode but without interactive elements.

**Rationale:** This gives immediate visual feedback as the user fills in fields, mirroring how object creation previews template body/properties.

### 4. Replace `newTypeMode`/`newTypeName` with `createType *createTypeState`

**Decision:** Remove the `newTypeMode bool` and `newTypeName textinput.Model` fields from the model. Replace with `createType *createTypeState` (nil when not creating a type).

**Rationale:** This follows the same pattern as `create *createState` for objects. A nil check is cleaner than a bool flag, and the struct encapsulates all creation state.

### 5. Validation approach

**Decision:** Validate on Enter (not on every keystroke). Show error inline in the title panel after the plural field.

**Validations:**
- Empty name → "name is required"
- Duplicate name → `type "X" already exists`
- Name validation delegated to `ValidateSchema()` (same as current flow)

**Rationale:** Keystroke validation would be noisy for short type names. Validating on Enter is consistent with the current behavior and the object creation flow.

## Risks / Trade-offs

- **Additional code for a rarely used flow** → Type creation is infrequent compared to object creation. However, the polish improvement is worth it for first-time UX and consistency. The code mirrors existing patterns, so maintenance cost is low.
- **Emoji input limitations** → Terminal emoji input varies by platform. Users may need to paste emoji from clipboard. This is acceptable since emoji is optional and can be added later in the type editor.
- **Title panel height** → The title panel is fixed at 3 lines (1 content + 2 borders). Three fields must fit in one line. With emoji (2-3 chars) + name input + plural input, this should fit comfortably in standard terminal widths (80+ columns).
