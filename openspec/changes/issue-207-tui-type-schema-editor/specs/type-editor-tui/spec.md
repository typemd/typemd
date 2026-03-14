## ADDED Requirements

### Requirement: Type editor opens when Enter is pressed on a type header

The TUI SHALL open a type editor panel in the right side when the user presses `Enter` on a type group header in the sidebar. The type editor SHALL display the type schema for the selected type group.

#### Scenario: Enter on type header opens editor
- **WHEN** the cursor is on a type group header "book" in the sidebar
- **AND** the user presses `Enter`
- **THEN** the right panel SHALL display the type editor for the "book" type schema

#### Scenario: Enter on type header while object is selected
- **WHEN** an object is currently displayed in the right panel
- **AND** the user moves the cursor to a type header and presses `Enter`
- **THEN** the right panel SHALL switch from object detail view to type editor view

#### Scenario: Esc in type editor returns to empty panel
- **WHEN** the type editor is open
- **AND** the user presses `Esc`
- **THEN** the type editor SHALL close and the right panel SHALL return to the empty state

### Requirement: Type editor displays meta fields

The type editor SHALL display the following meta fields at the top: Name, Plural, Emoji, and Unique. Each field SHALL show its current value from the TypeSchema.

#### Scenario: All meta fields displayed
- **WHEN** the type editor opens for type "book" with plural "books", emoji "📖", and unique false
- **THEN** the editor SHALL display Name: "book", Plural: "books", Emoji: "📖", Unique: "no"

#### Scenario: Missing optional fields
- **WHEN** the type editor opens for a type with no plural and no emoji
- **THEN** the editor SHALL display Name with the type name, Plural as empty, Emoji as empty, Unique as "no"

### Requirement: Unified cursor navigates meta fields and properties

The type editor SHALL use a single cursor that moves through meta fields (Name at index 0, Plural at index 1, Emoji at index 2, Unique at index 3) and then the property list (starting at index 4). Section separator lines SHALL be skipped by the cursor.

#### Scenario: Cursor moves through meta and properties
- **WHEN** the type has 2 properties and the cursor is on Unique (index 3)
- **AND** the user presses `↓`
- **THEN** the cursor SHALL move to the first property (index 4), skipping any section separator

#### Scenario: Cursor wraps at boundaries
- **WHEN** the cursor is on Name (index 0)
- **AND** the user presses `↑`
- **THEN** the cursor SHALL remain on Name (no wrap)

### Requirement: Name field is not editable

The Name meta field SHALL NOT be editable. When the cursor is on the Name field and the user presses `e`, the editor SHALL not enter edit mode for that field.

#### Scenario: Edit key on Name field
- **WHEN** the cursor is on the Name field
- **AND** the user presses `e`
- **THEN** nothing SHALL happen (no edit mode entered)

### Requirement: Plural field is editable via inline text input

When the cursor is on the Plural field and the user presses `e`, the field SHALL become an editable text input. Pressing `Enter` SHALL confirm the edit and save. Pressing `Esc` SHALL cancel the edit.

#### Scenario: Edit Plural field
- **WHEN** the cursor is on the Plural field showing "books"
- **AND** the user presses `e`, clears the field, types "Books", and presses `Enter`
- **THEN** the Plural field SHALL update to "Books" and the type schema SHALL be saved

#### Scenario: Cancel Plural edit
- **WHEN** the user is editing the Plural field
- **AND** the user presses `Esc`
- **THEN** the Plural field SHALL revert to its previous value

### Requirement: Emoji field is editable via inline text input

When the cursor is on the Emoji field and the user presses `e`, the field SHALL become an editable text input. Pressing `Enter` SHALL confirm and save. Pressing `Esc` SHALL cancel.

#### Scenario: Edit Emoji field
- **WHEN** the cursor is on the Emoji field
- **AND** the user presses `e`, types "📚", and presses `Enter`
- **THEN** the Emoji field SHALL update to "📚" and the type schema SHALL be saved

### Requirement: Unique field is toggled via edit key

When the cursor is on the Unique field and the user presses `e`, the value SHALL toggle between true and false and save immediately.

#### Scenario: Toggle Unique from false to true
- **WHEN** the cursor is on the Unique field showing "no"
- **AND** the user presses `e`
- **THEN** the Unique field SHALL change to "yes" and the type schema SHALL be saved

#### Scenario: Toggle Unique from true to false
- **WHEN** the cursor is on the Unique field showing "yes"
- **AND** the user presses `e`
- **THEN** the Unique field SHALL change to "no" and the type schema SHALL be saved

### Requirement: Property list displays in two sections

The type editor SHALL display properties in two sections: "Pinned (Header)" for properties with pin > 0 (sorted by pin value ascending), and "Properties" for properties with pin = 0 (in schema definition order). Section headers SHALL be displayed as non-selectable separator lines.

#### Scenario: Properties split into sections
- **WHEN** a type has properties: author (pin:1), genre (pin:2), rating (pin:0), isbn (pin:0)
- **THEN** the Pinned section SHALL show author, genre and the Properties section SHALL show rating, isbn

#### Scenario: Empty Pinned section
- **WHEN** a type has no pinned properties
- **THEN** the Pinned section header SHALL still be displayed with no items below it

### Requirement: Property emoji editable via inline input

When the cursor is on a property and the user presses `e`, an inline text input SHALL appear for editing the property's emoji. Pressing `Enter` confirms and saves. Pressing `Esc` cancels.

#### Scenario: Edit property emoji
- **WHEN** the cursor is on property "author" with emoji "👤"
- **AND** the user presses `e`, types "👨‍💻", and presses `Enter`
- **THEN** the property emoji SHALL update to "👨‍💻" and the type schema SHALL be saved

#### Scenario: Clear property emoji
- **WHEN** the cursor is on a property with emoji "👤"
- **AND** the user presses `e`, clears the field, and presses `Enter`
- **THEN** the property emoji SHALL be cleared and the type schema SHALL be saved

### Requirement: Pin toggle via p key

When the cursor is on a property in the Properties section and the user presses `p`, the property SHALL move to the Pinned section with pin value set to `max(existing pin values) + 1`. When the cursor is on a property in the Pinned section and the user presses `p`, the property SHALL move to the Properties section with pin value set to 0.

#### Scenario: Pin a property
- **WHEN** existing pinned properties have pin values 1 and 2
- **AND** the cursor is on an unpinned property "rating"
- **AND** the user presses `p`
- **THEN** rating SHALL move to the Pinned section with pin value 3 and the type schema SHALL be saved

#### Scenario: Unpin a property
- **WHEN** the cursor is on a pinned property "author" with pin value 1
- **AND** the user presses `p`
- **THEN** author SHALL move to the Properties section with pin value 0 and the type schema SHALL be saved

#### Scenario: Pin first property when no pins exist
- **WHEN** no properties are pinned
- **AND** the user presses `p` on a property
- **THEN** the property SHALL be pinned with pin value 1

### Requirement: Move mode for reordering properties

When the cursor is on a property and the user presses `m`, the editor SHALL enter move mode. In move mode, `↑`/`↓` SHALL swap the property with its neighbor. `Enter` or `Esc` SHALL exit move mode and save.

#### Scenario: Reorder within same section
- **WHEN** the Properties section has rating (index 0) and isbn (index 1)
- **AND** the user moves cursor to rating, presses `m`, then presses `↓`
- **THEN** isbn SHALL appear before rating in the Properties section

#### Scenario: Move across sections
- **WHEN** the user is in move mode on a property in the Properties section
- **AND** the user presses `↑` past the section boundary into Pinned
- **THEN** the property SHALL gain a pin value and appear in the Pinned section

#### Scenario: Exit move mode saves
- **WHEN** the user has reordered properties in move mode
- **AND** the user presses `Enter`
- **THEN** the type schema SHALL be saved with the new property order

### Requirement: Add Property wizard

When the user presses `a`, the type editor SHALL start a multi-step Add Property wizard rendered inline in the right panel.

#### Scenario: Wizard Step 1 — property name
- **WHEN** the user presses `a`
- **THEN** the panel SHALL display a text input for "Property name"
- **AND** pressing `Enter` with a non-empty name SHALL advance to Step 2
- **AND** pressing `Esc` SHALL cancel and return to the property list

#### Scenario: Wizard Step 2 — property type
- **WHEN** the user has entered a property name
- **THEN** the panel SHALL display a selectable list of property types (string, number, date, datetime, url, checkbox, select, multi_select, relation)
- **AND** pressing `Enter` on a non-relation type SHALL create the property and save
- **AND** pressing `Enter` on "relation" SHALL advance to Step 3

#### Scenario: Wizard Step 3 — relation config
- **WHEN** the user selected "relation" as property type
- **THEN** the panel SHALL display: target type selector, multiple toggle (y/n), bidirectional toggle (y/n)
- **AND** when bidirectional is yes, an inverse name text input SHALL appear
- **AND** pressing `Enter` SHALL create the relation property and save

#### Scenario: Wizard Step 2 — select/multi_select options
- **WHEN** the user selected "select" or "multi_select" as property type
- **THEN** the panel SHALL display an input for adding options (comma-separated or one per line)
- **AND** pressing `Enter` SHALL create the property with the given options and save

#### Scenario: Duplicate property name rejected
- **WHEN** the user enters a property name that already exists in the type schema
- **THEN** the wizard SHALL display an error message and remain on Step 1

### Requirement: Delete property with confirmation

When the cursor is on a property and the user presses `d`, the editor SHALL prompt for confirmation before removing the property.

#### Scenario: Confirm delete property
- **WHEN** the cursor is on property "isbn"
- **AND** the user presses `d`
- **THEN** the editor SHALL display "Delete property 'isbn'? [y/n]"
- **AND** pressing `y` SHALL remove the property and save

#### Scenario: Cancel delete property
- **WHEN** the delete confirmation is shown
- **AND** the user presses `n` or `Esc`
- **THEN** the property SHALL NOT be removed

### Requirement: New Type creation from sidebar

A "+ New Type" item SHALL appear at the bottom of the sidebar. Selecting it SHALL prompt for a type name and create a new empty type schema.

#### Scenario: Create new type
- **WHEN** the user selects "+ New Type" and enters "project"
- **THEN** a new type schema file SHALL be created at `.typemd/types/project.yaml`
- **AND** the type editor SHALL open for the new type

#### Scenario: Cancel new type creation
- **WHEN** the user selects "+ New Type" and presses `Esc`
- **THEN** no type SHALL be created and the sidebar SHALL return to normal

#### Scenario: Duplicate type name rejected
- **WHEN** the user enters a type name that already exists
- **THEN** an error message SHALL be displayed and the prompt SHALL remain open

### Requirement: Delete type with confirmation

When the type editor is open and the user triggers delete type (via a keybinding), the editor SHALL show a confirmation dialog with the count of existing objects of that type.

#### Scenario: Delete type with existing objects
- **WHEN** the user deletes type "note" which has 5 existing objects
- **THEN** the confirmation SHALL show "Delete type 'note'? 5 existing objects will become orphaned. [y/n]"
- **AND** pressing `y` SHALL delete `.typemd/types/note.yaml` and close the editor

#### Scenario: Delete built-in type rejected
- **WHEN** the user attempts to delete the "tag" type
- **THEN** the editor SHALL display "Cannot delete built-in type 'tag'" and not proceed

#### Scenario: Delete type with no objects
- **WHEN** the user deletes type "project" which has 0 objects
- **THEN** the confirmation SHALL show "Delete type 'project'? [y/n]" without orphan warning

### Requirement: Type editor keybindings shown in help bar

When the type editor is active, the help bar at the bottom SHALL display the available keybindings for the current context.

#### Scenario: View mode help bar
- **WHEN** the type editor is in view mode
- **THEN** the help bar SHALL display "[e]dit [a]dd [d]elete [m]ove [p]in esc: back"

#### Scenario: Move mode help bar
- **WHEN** the type editor is in move mode
- **THEN** the help bar SHALL display "[MOVE] ↑↓: reorder enter/esc: done"

#### Scenario: Edit mode help bar
- **WHEN** the type editor is in inline edit mode
- **THEN** the help bar SHALL display "[EDIT] enter: save esc: cancel"
