## ADDED Requirements

### Requirement: Type creation form renders in title panel

When type creation is active, the title panel SHALL transform into an inline creation form showing an emoji input (optional), a name text input (required), and a plural text input (optional). The sidebar SHALL NOT render any creation-specific UI.

#### Scenario: Title panel shows creation form

- **WHEN** the user selects `+ New Type` and presses Enter
- **THEN** the title panel SHALL display `[emoji___] new type · [name___]  plural: [plural___]`
- **AND** the name field SHALL be focused by default
- **AND** the sidebar SHALL show the normal type/object list without any creation UI appended

#### Scenario: Title panel with emoji filled

- **WHEN** the user enters emoji "📝" in the emoji field
- **THEN** the title panel SHALL display `📝 new type · [name___]  plural: [plural___]`

#### Scenario: Title panel visible during creation even without prior selection

- **WHEN** no object or type was previously selected (title panel was hidden)
- **AND** the user initiates type creation
- **THEN** the title panel SHALL appear with the creation form

### Requirement: Tab navigates between creation fields

During type creation, Tab SHALL cycle focus between the emoji, name, and plural input fields. The focused field SHALL be visually indicated.

#### Scenario: Tab from name to plural

- **WHEN** the name field is focused
- **AND** the user presses Tab
- **THEN** the plural field SHALL become focused

#### Scenario: Tab from plural to emoji

- **WHEN** the plural field is focused
- **AND** the user presses Tab
- **THEN** the emoji field SHALL become focused

#### Scenario: Tab from emoji to name

- **WHEN** the emoji field is focused
- **AND** the user presses Tab
- **THEN** the name field SHALL become focused

### Requirement: Live preview of type schema

During type creation, the right panel SHALL display a read-only preview of the type schema being created. The preview SHALL update in real time as the user types in any field.

#### Scenario: Preview shows name

- **WHEN** the user types "meeting" in the name field
- **THEN** the right panel SHALL display a type preview with Name: "meeting"

#### Scenario: Preview shows emoji and plural

- **WHEN** the user enters emoji "📋" and plural "meetings"
- **THEN** the right panel SHALL display a type preview with Emoji: "📋" and Plural: "meetings"

#### Scenario: Preview shows empty fields as placeholders

- **WHEN** the emoji and plural fields are empty
- **THEN** the right panel SHALL display the preview with Emoji and Plural as empty or placeholder values

### Requirement: Name validation on submit

When the user presses Enter to create the type, the system SHALL validate the name. If validation fails, the error SHALL be displayed inline in the title panel.

#### Scenario: Empty name rejected

- **WHEN** the user presses Enter with an empty name field
- **THEN** no type SHALL be created
- **AND** the creation form SHALL remain active

#### Scenario: Duplicate name rejected

- **WHEN** the user presses Enter with name "book"
- **AND** a type named "book" already exists
- **THEN** the title panel SHALL display error `type "book" already exists`
- **AND** no type SHALL be created

#### Scenario: Valid name accepted

- **WHEN** the user presses Enter with name "meeting"
- **AND** no type named "meeting" exists
- **THEN** the type SHALL be created with the provided name, emoji, and plural
- **AND** the type editor SHALL open for the new type

### Requirement: Escape cancels type creation

When the user presses Escape during type creation, the creation form SHALL be dismissed without creating a type.

#### Scenario: Escape cancels creation

- **WHEN** the creation form is displayed in the title panel
- **AND** the user presses Escape
- **THEN** the creation form SHALL be dismissed
- **AND** no type SHALL be created
- **AND** the title panel SHALL return to its previous state

### Requirement: Created type opens in type editor

After successful type creation, the type editor SHALL open with the newly created type, pre-populated with the emoji and plural values entered during creation.

#### Scenario: Type editor opens with pre-populated fields

- **WHEN** the user creates a type with name "meeting", emoji "📋", and plural "meetings"
- **THEN** the type editor SHALL open showing the "meeting" type
- **AND** the type editor SHALL display Emoji: "📋" and Plural: "meetings"
- **AND** the focus SHALL move to the right panel (type editor)

#### Scenario: Type editor opens with empty optional fields

- **WHEN** the user creates a type with name "meeting" and no emoji or plural
- **THEN** the type editor SHALL open showing the "meeting" type
- **AND** the Emoji and Plural fields SHALL be empty

### Requirement: Help bar reflects type creation mode

The help bar SHALL display context-appropriate keybinding hints during type creation.

#### Scenario: Help bar during type creation

- **WHEN** the type creation form is active
- **THEN** the help bar SHALL display hints including "tab: next field  enter: create  esc: cancel"
