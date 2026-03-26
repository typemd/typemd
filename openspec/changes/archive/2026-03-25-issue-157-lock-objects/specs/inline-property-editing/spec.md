## ADDED Requirements

### Requirement: Locked objects disable property editing
The TUI property editor SHALL NOT activate for locked objects. When a locked object is displayed, the properties panel SHALL show all properties in read-only mode without the cursor indicator.

#### Scenario: Property editor not available for locked objects
- **WHEN** the user views a locked object and presses Tab to focus properties
- **THEN** the properties panel SHALL remain in read-only display mode
- **AND** no property cursor (▸) SHALL be shown

#### Scenario: Toast notification on edit attempt of locked object
- **WHEN** the user attempts to activate property editing on a locked object
- **THEN** a toast notification SHALL display "Object is locked. Unlock to edit."
