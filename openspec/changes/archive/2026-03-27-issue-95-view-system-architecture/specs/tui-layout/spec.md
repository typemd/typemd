## ADDED Requirements

### Requirement: TUI supports View panel mode

The TUI SHALL support a new right panel mode `panelView` that replaces the standard three-panel layout with a full-width View rendering. The existing panel modes (`panelEmpty`, `panelObject`, `panelTypeEditor`, `panelTemplate`) SHALL remain unchanged.

#### Scenario: View mode active

- **WHEN** the TUI is in `panelView` mode
- **THEN** the sidebar, body panel, and properties panel SHALL be hidden, and the full terminal width SHALL be used for the View rendering

#### Scenario: Exit view mode restores layout

- **WHEN** the user exits View mode by pressing Esc from the View list
- **THEN** the TUI SHALL restore the standard three-panel layout with the previous sidebar state

### Requirement: TUI supports nested navigation in View mode

When in View mode, pressing Enter on an object SHALL display the object detail in the standard three-panel layout, but within the View context. Pressing Esc from the object detail SHALL return to the View list, not to the sidebar.

#### Scenario: Object detail within View context

- **WHEN** the user presses Enter on an object in View mode
- **THEN** the TUI SHALL display the object detail (title panel + body + properties) as if in normal `panelObject` mode

#### Scenario: Esc from object detail returns to View

- **WHEN** the user presses Esc from the object detail displayed within View mode
- **THEN** the TUI SHALL return to the full-width View list, not to the sidebar

### Requirement: View mode supports split layout with preview

When the preview panel is active in View mode, the layout SHALL split horizontally: table on the left (~60% width), preview panel on the right (~40% width). Both panels SHALL have their own borders.

#### Scenario: Preview split layout

- **WHEN** the user activates the preview panel in View mode
- **THEN** the full-width area SHALL split into a table panel (left) and a preview panel (right) with rounded borders

#### Scenario: Preview closed restores full-width

- **WHEN** the user deactivates the preview panel
- **THEN** the table SHALL expand to use the full terminal width
