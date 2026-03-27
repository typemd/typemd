## ADDED Requirements

### Requirement: TUI supports View mode as full-width table

The TUI SHALL support a View mode that replaces the standard three-panel layout with a full-width table display. When View mode is active, the sidebar, body, and properties panels SHALL be hidden. The table SHALL display object names followed by property columns (pinned properties first, then unpinned). The number of columns SHALL adjust to the terminal width.

#### Scenario: Enter View mode

- **WHEN** the user activates a View for type "book"
- **THEN** the TUI SHALL display a full-width table of book objects with property columns, hiding the sidebar, body, and properties panels

#### Scenario: View mode header

- **WHEN** View mode is active with view name "by-rating" for type "book" with emoji "📚"
- **THEN** the top of the screen SHALL display the view identity (e.g., "📚 book · by-rating")

#### Scenario: Table column header

- **WHEN** View mode is active for type "book" with properties "status", "rating", "author"
- **THEN** the table SHALL display a column header row with "NAME" followed by property names in uppercase

### Requirement: View mode entered from type editor

The type editor's Views section SHALL list all saved views for the current type. Pressing Enter on a view SHALL activate View mode for that view. A "+ Add View" action row SHALL allow creating new views inline.

#### Scenario: Type editor shows views

- **WHEN** the type editor is displayed for type "book" which has views "default" and "by-rating"
- **THEN** the Views section SHALL list both views and a "+ Add View" row

#### Scenario: Enter view from type editor

- **WHEN** the user presses Enter on "by-rating" in the type editor Views section
- **THEN** the TUI SHALL enter View mode with the "by-rating" view configuration applied

#### Scenario: Create view from type editor

- **WHEN** the user presses Enter on "+ Add View" and enters "reading-now"
- **THEN** a new view SHALL be created with default config (list layout, sort by name asc) and appear in the Views section

### Requirement: View mode entered via keyboard shortcut

The TUI SHALL support a keyboard shortcut (configurable, default `v`) to enter View mode for the currently selected type. If the type has multiple views, a selection popup SHALL appear. If the type has only the default view, it SHALL enter View mode directly.

#### Scenario: Shortcut with single view

- **WHEN** the user presses `v` while a "note" type header or object is selected, and "note" has no saved views
- **THEN** the TUI SHALL enter View mode with the default view for "note"

#### Scenario: Shortcut with multiple views

- **WHEN** the user presses `v` while a "book" object is selected, and "book" has views "default" and "by-rating"
- **THEN** the TUI SHALL show a view selection popup listing both views

### Requirement: View mode applies filter and sort

When View mode is active, the displayed objects SHALL be filtered and sorted according to the ViewConfig. The TUI SHALL call `QueryService.Query()` with the filter and sort from the ViewConfig.

#### Scenario: View with filter

- **WHEN** the active view has filter [{property: "status", operator: "is", value: "reading"}]
- **THEN** only objects matching the filter SHALL be displayed

#### Scenario: View with sort

- **WHEN** the active view has sort [{property: "rating", direction: "desc"}]
- **THEN** objects SHALL be displayed in descending rating order

#### Scenario: View with filter and sort

- **WHEN** the active view has filter on status=reading and sort by rating desc
- **THEN** only reading objects SHALL be displayed, sorted by rating descending

### Requirement: View mode applies group_by

When a ViewConfig has a non-empty `group_by` field, the View list SHALL group objects by the specified property value. Each group SHALL have a collapsible header showing the property value.

#### Scenario: Group by select property

- **WHEN** the active view has group_by "status" and objects have status values "reading", "finished", "want-to-read"
- **THEN** the View list SHALL display three groups with headers "reading", "finished", "want-to-read", each containing their respective objects

#### Scenario: Group by with empty values

- **WHEN** some objects have no value for the group_by property
- **THEN** those objects SHALL be grouped under a "(none)" or empty group header

#### Scenario: No group_by

- **WHEN** the active view has no group_by
- **THEN** objects SHALL be displayed as a flat list without group headers

### Requirement: View mode navigation stack

View mode SHALL maintain a navigation stack: sidebar → view list → object detail. Esc SHALL navigate back one level.

#### Scenario: Enter object detail from view

- **WHEN** View mode is active and the user presses Enter on an object
- **THEN** the TUI SHALL show the object detail (three-panel layout) within the View context

#### Scenario: Return to view list from object detail

- **WHEN** the user is viewing an object detail within View mode and presses Esc
- **THEN** the TUI SHALL return to the View list with the cursor position preserved

#### Scenario: Exit view mode

- **WHEN** the user is on the View list and presses Esc
- **THEN** the TUI SHALL return to the standard sidebar browsing mode

### Requirement: Default view always available

Every type SHALL have an accessible default view. When no saved views exist, the default view (layout: list, sort by name ascending) SHALL be available in the type editor and via keyboard shortcut.

#### Scenario: Type with no saved views

- **WHEN** type "note" has no saved views
- **THEN** the type editor SHALL show "default" in the Views section, and the keyboard shortcut SHALL enter View mode with the implicit default

#### Scenario: Type with customized default

- **WHEN** type "book" has a saved `views/default.yaml` with sort by rating desc
- **THEN** the default view SHALL use the saved configuration instead of the implicit default

### Requirement: View mode supports preview panel

View mode SHALL support a toggleable preview panel on the right side. When active, the layout splits into a table on the left and a preview panel on the right showing the cursor object's properties and body.

#### Scenario: Toggle preview on

- **WHEN** the user presses `p` in View mode with no preview active
- **THEN** the right side SHALL display a preview panel showing the current cursor object's properties and body content

#### Scenario: Preview follows cursor

- **WHEN** preview is active and the user moves the cursor to a different object
- **THEN** the preview panel SHALL update to show the newly selected object's content

#### Scenario: Toggle preview off

- **WHEN** the user presses `p` in View mode with preview active
- **THEN** the preview panel SHALL be hidden and the table SHALL expand to full width

### Requirement: View selection popup for multiple views

When the `v` shortcut is pressed and the type has more than one saved view, a selection popup SHALL appear. The popup SHALL use the huh Select component for keyboard-navigable selection.

#### Scenario: Popup with multiple views

- **WHEN** the user presses `v` on type "book" which has views "default", "by-rating", and "reading-now"
- **THEN** a centered popup SHALL appear listing all three views for selection

#### Scenario: Popup confirm selection

- **WHEN** the user selects "by-rating" in the popup and presses Enter
- **THEN** the popup SHALL close and View mode SHALL open with the "by-rating" view

#### Scenario: Popup cancel

- **WHEN** the user presses Esc in the popup
- **THEN** the popup SHALL close and the TUI SHALL remain in sidebar mode
