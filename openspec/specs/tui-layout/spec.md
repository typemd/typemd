## Requirements

### Requirement: Title panel displays object identity
The TUI detail view SHALL display a dedicated title panel above the body and properties panels showing the type emoji, type name, and object name (from `GetName()`).

#### Scenario: Title panel with emoji
- **WHEN** an object of type "book" with emoji "📖" and name "Clean Code" is selected
- **THEN** the title panel SHALL display "📖 book · Clean Code"

#### Scenario: Title panel without emoji
- **WHEN** an object of type "note" with no emoji and name "My Note" is selected
- **THEN** the title panel SHALL display "note · My Note"

### Requirement: Title panel spans body and properties width
The title panel SHALL span the full width of both the body panel and properties panel combined.

#### Scenario: Properties visible
- **WHEN** the properties panel is visible
- **THEN** the title panel width SHALL equal the combined width of body and properties panels (including borders)

#### Scenario: Properties hidden
- **WHEN** the properties panel is hidden
- **THEN** the title panel width SHALL equal the body panel width only

### Requirement: Title panel hidden when no object selected
The title panel SHALL NOT be displayed when no object is selected.

#### Scenario: No selection
- **WHEN** no object is selected in the list
- **THEN** the title panel SHALL be hidden and the body panel SHALL display the default placeholder message

### Requirement: Body panel no longer contains title header
The body panel SHALL NOT render the object title or separator line. The body panel SHALL display only the markdown body content.

#### Scenario: Body content without title
- **WHEN** an object is selected and the body panel is rendered
- **THEN** the body panel SHALL start directly with the markdown body content, without a title line or separator

### Requirement: Title panel height is fixed
The title panel SHALL occupy exactly 3 lines of vertical space (1 content line + 2 border lines).

#### Scenario: Vertical space allocation
- **WHEN** the TUI detail view is rendered with an object selected
- **THEN** the body and properties panels SHALL have their content height reduced by 3 lines compared to the no-title-panel state

### Requirement: Property supports optional pin field

The Property struct SHALL support an optional `pin` field that stores a positive integer value. When a property definition in a type schema YAML includes a `pin` field, it SHALL be parsed and stored. When the field is omitted, the pin SHALL default to zero (not pinned).

#### Scenario: Property with pin defined
- **WHEN** a type schema property definition contains `pin: 1`
- **THEN** the loaded Property SHALL have its Pin field set to 1

#### Scenario: Property without pin defined
- **WHEN** a type schema property definition does not contain a `pin` field
- **THEN** the loaded Property SHALL have its Pin field set to 0

### Requirement: Pin values must be positive integers

When a property has a pin value set, it SHALL be a positive integer (greater than zero). Schema validation SHALL reject negative pin values.

#### Scenario: Positive pin value accepted
- **WHEN** a type schema property has `pin: 3`
- **THEN** schema validation SHALL accept it without error

#### Scenario: Negative pin value rejected
- **WHEN** a type schema property has `pin: -1`
- **THEN** schema validation SHALL return an error indicating invalid pin value

### Requirement: Pin values unique within type scope

Within a single type schema, no two properties SHALL have the same non-zero pin value. Schema validation SHALL reject duplicate pin values.

#### Scenario: Unique pin values accepted
- **WHEN** a type schema has properties with pin values 1 and 2
- **THEN** schema validation SHALL accept it without error

#### Scenario: Duplicate pin values rejected
- **WHEN** a type schema has two properties both with `pin: 1`
- **THEN** schema validation SHALL return an error indicating duplicate pin value 1

#### Scenario: Unpinned properties do not conflict
- **WHEN** a type schema has three properties where two have no pin and one has `pin: 1`
- **THEN** schema validation SHALL accept it without error

### Requirement: Pinned properties displayed at top of body panel

Properties with a non-zero `pin` value SHALL be rendered at the top of the body panel, above the markdown body content. Pinned properties SHALL be sorted by pin value ascending (lower number first). When a property has an emoji defined, it SHALL be displayed alongside the pinned value.

#### Scenario: Pinned property rendered in body panel
- **WHEN** a type schema has property "status" with `pin: 1` and `emoji: 📋`
- **AND** the object has `status: reading`
- **THEN** the body panel SHALL display `📋 status: reading` at the top, before the markdown body

#### Scenario: Separator between pinned properties and body
- **WHEN** an object has pinned properties and non-empty body content
- **THEN** a horizontal separator SHALL appear between the pinned properties and the body content

#### Scenario: No separator when no body content
- **WHEN** an object has pinned properties but empty body content
- **THEN** the pinned properties SHALL be displayed without a trailing separator

### Requirement: Pinned properties excluded from Properties panel

Properties with a non-zero `pin` value SHALL NOT appear in the Properties panel. Only unpinned properties (pin = 0) SHALL be displayed in the Properties panel.

#### Scenario: Pinned property absent from Properties panel
- **WHEN** a type schema has property "status" with `pin: 1`
- **THEN** "status" SHALL NOT appear in the Properties panel

### Requirement: Properties panel displays property values
The properties panel SHALL display unpinned, non-name properties. When the properties panel is focused, properties SHALL display with a cursor indicator on the currently selected property. Editable properties SHALL be visually distinguished from read-only properties.

#### Scenario: Cursor indicator on focused panel
- **WHEN** the properties panel gains focus via Tab
- **THEN** the first editable property SHALL be highlighted with a cursor indicator (e.g., `▸` prefix)

#### Scenario: Read-only properties shown without cursor
- **WHEN** the properties panel is focused
- **THEN** read-only properties (created_at, updated_at, reverse relations, backlinks, relations) SHALL be displayed but SHALL NOT receive cursor highlight during navigation

#### Scenario: Edit mode border color
- **WHEN** a property is being actively edited (textinput visible or picker open)
- **THEN** the properties panel border SHALL use the edit border color (orange)

### Requirement: TUI startup initializes from restored state
The TUI `Start()` function SHALL attempt to load session state from `.typemd/tui-state.yaml` before applying default values. Restored state values take precedence over hardcoded defaults. If no state file exists or loading fails, the TUI SHALL use the current default behavior (first group expanded, first object selected).

#### Scenario: Startup with saved state
- **WHEN** the TUI starts and `.typemd/tui-state.yaml` contains valid state
- **THEN** the TUI SHALL initialize with the restored state instead of hardcoded defaults

#### Scenario: Startup without saved state
- **WHEN** the TUI starts and no `.typemd/tui-state.yaml` exists
- **THEN** the TUI SHALL initialize with current defaults (first group expanded, cursor at top, focus on left panel)

### Requirement: ViewConfig supports table layout constant
The core layer SHALL define a `ViewLayoutTable` constant with value `"table"` alongside the existing `ViewLayoutList`.

#### Scenario: Table layout in YAML
- **WHEN** a view YAML contains `layout: table`
- **THEN** ViewConfig.Layout SHALL be `ViewLayoutTable`

### Requirement: View editor supports layout selection
The view editor SHALL provide a Layout section allowing users to switch between `list` and `table` layouts.

#### Scenario: Switch from list to table
- **WHEN** the user selects `table` in the Layout section of the view editor
- **THEN** the view SHALL save with `layout: table` and the display SHALL switch to columnar table format immediately

#### Scenario: Switch from table to list
- **WHEN** the user selects `list` in the Layout section of the view editor
- **THEN** the view SHALL save with `layout: list` and the display SHALL switch to inline name format immediately

### Requirement: Group header displays type emoji

The TUI object list panel SHALL display the type's emoji prefix in group headers when the type schema defines an emoji field. Object list loading SHALL use `[]FilterRule` parameters when calling `Vault.QueryObjects()`.

#### Scenario: Type with emoji defined
- **WHEN** a type schema has an emoji field (e.g., book with 📚)
- **THEN** the group header displays as `▼ 📚 book (N)` where N is the object count

#### Scenario: Type without emoji defined
- **WHEN** a type schema does not have an emoji field
- **THEN** the group header displays as `▼ book (N)` with no extra spacing or placeholder

#### Scenario: Object list loading uses structured filter
- **WHEN** the TUI loads the object list (via `app.go` or `view_mode.go`)
- **THEN** it SHALL call `Vault.QueryObjects([]FilterRule{...})` instead of passing a filter string

### Requirement: Normal mode help bar shows both creation keybindings

When the sidebar is focused in normal mode, the help bar SHALL include both `n` (new) and `N` (quick create) keybinding hints.

#### Scenario: Sidebar focused help bar
- **WHEN** the sidebar is focused in normal mode with a type header or object selected
- **THEN** the help bar SHALL include hints for both `n: new` and `N: quick create`
