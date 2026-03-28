## ADDED Requirements

### Requirement: Inline property editing
The properties panel SHALL support inline editing of property values by double-clicking.

#### Scenario: Edit string property
- **WHEN** the user double-clicks a string property value
- **THEN** a text input SHALL appear with the current value, and pressing Enter SHALL save the change

#### Scenario: Edit number property
- **WHEN** the user double-clicks a number property value
- **THEN** a number input SHALL appear, and pressing Enter SHALL save the change

#### Scenario: Edit date property
- **WHEN** the user double-clicks a date property value
- **THEN** a date picker input SHALL appear, and selecting a date or pressing Enter SHALL save the change

#### Scenario: Toggle checkbox property
- **WHEN** the user clicks a checkbox property
- **THEN** the value SHALL toggle and save immediately without entering edit mode

#### Scenario: Cancel editing
- **WHEN** the user presses Escape while editing a property
- **THEN** the edit SHALL be cancelled and the original value SHALL be restored

#### Scenario: Non-editable properties
- **WHEN** the user double-clicks `created_at`, `updated_at`, a reverse relation, or a backlink property
- **THEN** no editor SHALL appear because these properties are not editable

#### Scenario: Locked object
- **WHEN** the object has `locked: true`
- **THEN** no property SHALL be editable

### Requirement: Body editing
The body panel SHALL support editing the object's markdown body.

#### Scenario: Enter edit mode
- **WHEN** the user clicks the "Edit" button in the body header
- **THEN** the body SHALL switch to a textarea editor with the current content

#### Scenario: Save body
- **WHEN** the user clicks "Save" or presses ⌘+Enter / Ctrl+Enter while editing
- **THEN** the body SHALL be saved and the panel SHALL return to view mode

#### Scenario: Cancel body edit
- **WHEN** the user clicks "Cancel" or presses Escape while editing
- **THEN** the body SHALL revert to the original content and return to view mode

#### Scenario: Locked object hides edit button
- **WHEN** the object has `locked: true`
- **THEN** the "Edit" button SHALL NOT be displayed

### Requirement: Object creation
The web UI SHALL support creating new objects through a dialog.

#### Scenario: Open create dialog
- **WHEN** the user presses `n` outside of an input field, or clicks "+ New" in the sidebar
- **THEN** a create object dialog SHALL appear

#### Scenario: Create with type and name
- **WHEN** the user selects a type and enters a name in the create dialog and submits
- **THEN** a new object SHALL be created and selected in the sidebar

#### Scenario: Create with template
- **WHEN** the selected type has templates available
- **THEN** the dialog SHALL show a template selector, and selecting a template SHALL apply it to the new object

#### Scenario: Cancel creation
- **WHEN** the user presses Escape or clicks Cancel in the create dialog
- **THEN** the dialog SHALL close without creating an object

### Requirement: Sidebar refresh after mutation
The sidebar SHALL reflect changes after object creation.

#### Scenario: New object appears
- **WHEN** a new object is created via the create dialog
- **THEN** the sidebar SHALL refresh to include the new object and update the type count
