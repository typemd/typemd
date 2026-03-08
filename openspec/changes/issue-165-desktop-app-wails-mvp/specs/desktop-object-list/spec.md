## ADDED Requirements

### Requirement: Frontend displays object list grouped by type

The web frontend SHALL display objects from the vault grouped by their type name.

#### Scenario: Objects rendered in groups
- **WHEN** the desktop app window opens with a vault containing objects
- **THEN** the frontend displays objects grouped under type headings (e.g., "book", "note")

#### Scenario: Empty vault
- **WHEN** the desktop app opens with an empty vault (no objects)
- **THEN** the frontend displays a message indicating no objects found

### Requirement: Frontend communicates with Go backend via Wails bindings

The frontend SHALL call Go backend methods through Wails-generated bindings, not via HTTP or other protocols.

#### Scenario: Object data loaded via bindings
- **WHEN** the frontend page loads
- **THEN** it calls the `ListObjects` binding to fetch object data from the Go backend
