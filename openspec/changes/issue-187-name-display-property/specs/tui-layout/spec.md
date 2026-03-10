## MODIFIED Requirements

### Requirement: Title panel displays object identity
The TUI detail view SHALL display a dedicated title panel above the body and properties panels showing the type emoji, type name, and object name (from `GetName()`).

#### Scenario: Title panel with emoji
- **WHEN** an object of type "book" with emoji "📖" and name "Clean Code" is selected
- **THEN** the title panel SHALL display "📖 book · Clean Code"

#### Scenario: Title panel without emoji
- **WHEN** an object of type "note" with no emoji and name "My Note" is selected
- **THEN** the title panel SHALL display "note · My Note"
