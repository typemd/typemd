## ADDED Requirements

### Requirement: Create & Edit mode creates one object and enters edit mode

Pressing `n` on a type header or object in the sidebar SHALL initiate Create & Edit mode. After the object is created, the TUI SHALL select the new object, focus the body panel, and enter edit mode automatically.

#### Scenario: Create & Edit with name input

- **WHEN** the user presses `n` on a type header for type "book" (no name template, no templates)
- **AND** types "clean-code" and presses Enter
- **THEN** a new book object named "clean-code" SHALL be created
- **AND** the object SHALL be selected in the sidebar
- **AND** the body panel SHALL be focused in edit mode

#### Scenario: Create & Edit cancelled

- **WHEN** the user presses `n` on a type header for type "book"
- **AND** presses Escape before entering a name
- **THEN** no object SHALL be created
- **AND** the sidebar SHALL return to normal mode

#### Scenario: Create & Edit in read-only mode

- **WHEN** the TUI is in read-only mode
- **AND** the user presses `n`
- **THEN** no creation flow SHALL be initiated

### Requirement: Quick Create mode supports batch object creation

Pressing `N` (Shift+n) on a type header or object in the sidebar SHALL initiate Quick Create mode. After each object is created, the name input SHALL be cleared and remain focused for the next entry. Pressing Escape SHALL exit Quick Create mode.

#### Scenario: Create multiple objects in batch

- **WHEN** the user presses `N` on a type header for type "book"
- **AND** types "book-one" and presses Enter
- **AND** types "book-two" and presses Enter
- **AND** presses Escape
- **THEN** two book objects SHALL be created ("book-one" and "book-two")
- **AND** the last created object ("book-two") SHALL be selected in the sidebar

#### Scenario: Quick Create shows success flash

- **WHEN** the user creates an object in Quick Create mode
- **THEN** a success flash message SHALL be displayed (e.g., "Created: book-one")
- **AND** the flash SHALL auto-dismiss after approximately 2 seconds

#### Scenario: Quick Create in read-only mode

- **WHEN** the TUI is in read-only mode
- **AND** the user presses `N`
- **THEN** no creation flow SHALL be initiated

### Requirement: Template selection when multiple templates exist

When a type has two or more object templates, the creation flow SHALL present a template selection list before the name input step. The list SHALL include all available templates plus a "(none)" option. The user SHALL navigate with arrow keys and confirm with Enter.

#### Scenario: Type with multiple templates

- **WHEN** the user initiates object creation for type "book"
- **AND** type "book" has templates "review" and "summary"
- **THEN** a template selection list SHALL be displayed with options: "review", "summary", "(none)"
- **AND** the user SHALL select one before proceeding to name input

#### Scenario: Template selection cancelled

- **WHEN** the template selection list is displayed
- **AND** the user presses Escape
- **THEN** the creation flow SHALL be cancelled entirely

#### Scenario: Template none selected

- **WHEN** the user selects "(none)" from the template list
- **THEN** the creation flow SHALL proceed to name input without a template

### Requirement: Single template auto-applied

When a type has exactly one object template, the creation flow SHALL automatically select that template without prompting the user.

#### Scenario: Auto-apply single template

- **WHEN** the user initiates object creation for type "book"
- **AND** type "book" has exactly one template "default"
- **THEN** the "default" template SHALL be automatically selected
- **AND** the creation flow SHALL proceed to name input (or object creation if name template exists)

### Requirement: No templates skips template step

When a type has no object templates, the creation flow SHALL skip the template selection step entirely.

#### Scenario: No templates available

- **WHEN** the user initiates object creation for type "note"
- **AND** type "note" has no templates
- **THEN** the creation flow SHALL proceed directly to name input

### Requirement: Name template auto-skip in Create & Edit mode

In Create & Edit mode (`n`), when the type defines a name template, the name input step SHALL be skipped. The object SHALL be created with an auto-generated name from the name template.

#### Scenario: Name template skips name input

- **WHEN** the user presses `n` on type "journal" which has name template `{{ date:YYYY-MM-DD }}`
- **AND** no template selection is needed (0 or 1 templates)
- **THEN** the object SHALL be created with name derived from the name template (e.g., "2026-03-16")
- **AND** the TUI SHALL enter edit mode on the new object's body

#### Scenario: Name template with template selection

- **WHEN** the user presses `n` on type "journal" which has name template `{{ date:YYYY-MM-DD }}`
- **AND** type "journal" has multiple templates
- **THEN** the template selection step SHALL be shown
- **AND** after template selection, the object SHALL be created without name input

### Requirement: Quick Create always requires name input

In Quick Create mode (`N`), the name input step SHALL always be shown regardless of whether a name template is defined. Name templates SHALL NOT be auto-applied in Quick Create mode.

#### Scenario: Quick Create ignores name template

- **WHEN** the user presses `N` on type "journal" which has name template `{{ date:YYYY-MM-DD }}`
- **THEN** the name input SHALL be displayed
- **AND** the user SHALL type a name manually

### Requirement: Batch mode reuses selected template

In Quick Create mode, the template selection (if applicable) SHALL occur only once at the beginning. All subsequent objects in the batch SHALL use the same template.

#### Scenario: Template persists across batch

- **WHEN** the user presses `N` on type "book" with templates "review" and "summary"
- **AND** selects "review" from the template list
- **AND** creates "book-one" and "book-two"
- **THEN** both objects SHALL be created using the "review" template

### Requirement: Unique constraint validation with inline feedback

When a type has `unique: true` and the user enters a name that already exists, the creation flow SHALL display an inline error message below the input and SHALL NOT create the object. The error SHALL clear when the user modifies the input.

#### Scenario: Duplicate name rejected

- **WHEN** the user enters name "existing-book" for type "book" with `unique: true`
- **AND** a book named "existing-book" already exists
- **THEN** an error message SHALL be displayed below the input (e.g., 'book name "existing-book" already exists')
- **AND** the object SHALL NOT be created
- **AND** the input SHALL remain focused for correction

#### Scenario: Error clears on input change

- **WHEN** a duplicate name error is displayed
- **AND** the user modifies the input text
- **THEN** the error message SHALL be cleared

### Requirement: Empty name rejected

When the user presses Enter with an empty or whitespace-only name, the creation flow SHALL not create an object and SHALL remain in the current input state.

#### Scenario: Empty name ignored

- **WHEN** the user presses Enter with an empty name input
- **THEN** no object SHALL be created
- **AND** the input SHALL remain focused

### Requirement: Help bar reflects creation mode

The help bar SHALL display context-appropriate keybinding hints during object creation.

#### Scenario: Help bar during template selection

- **WHEN** the template selection list is displayed
- **THEN** the help bar SHALL display navigation hints (e.g., "↑↓: select  enter: confirm  esc: cancel")

#### Scenario: Help bar during name input in Create & Edit mode

- **WHEN** the name input is displayed in Create & Edit mode
- **THEN** the help bar SHALL display "enter: create & edit  esc: cancel"

#### Scenario: Help bar during name input in Quick Create mode

- **WHEN** the name input is displayed in Quick Create mode
- **THEN** the help bar SHALL display "enter: create  esc: done"
