## MODIFIED Requirements

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
