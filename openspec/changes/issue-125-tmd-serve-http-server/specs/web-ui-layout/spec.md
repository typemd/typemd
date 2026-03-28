## ADDED Requirements

### Requirement: Three-panel layout
The web UI SHALL display a three-panel layout: sidebar (left), body (center), and properties (right).

#### Scenario: Default view
- **WHEN** the web UI is loaded and an object is selected
- **THEN** all three panels SHALL be visible: sidebar with type groups, body with object content, and properties with object attributes

#### Scenario: No object selected
- **WHEN** no object is selected
- **THEN** the sidebar SHALL be visible, the body SHALL show a placeholder message, and the properties panel SHALL be hidden

### Requirement: Sidebar displays type groups
The sidebar SHALL display all vault types as expandable groups, with objects listed under each expanded group.

#### Scenario: Type group header
- **WHEN** the sidebar loads
- **THEN** each type SHALL be displayed with its emoji (if present), name (plural form if available), and a disclosure arrow

#### Scenario: Expand type group
- **WHEN** a type group header is clicked
- **THEN** the group SHALL expand to show all objects of that type, sorted by name ascending

#### Scenario: Collapse type group
- **WHEN** an expanded type group header is clicked again
- **THEN** the group SHALL collapse and hide its object list

#### Scenario: Object count on hover
- **WHEN** the user hovers over a type group header
- **THEN** the object count SHALL be displayed

### Requirement: Sidebar object selection
Clicking an object in the sidebar SHALL select it and display its content in the body and properties panels.

#### Scenario: Select object
- **WHEN** an object name is clicked in the sidebar
- **THEN** the body panel SHALL display the object's markdown body and the properties panel SHALL display the object's properties

#### Scenario: Selected object highlight
- **WHEN** an object is selected
- **THEN** it SHALL be visually highlighted in the sidebar

### Requirement: Focus mode
The web UI SHALL support a focus mode that hides the sidebar and properties panels, showing only the body at full width.

#### Scenario: Toggle focus mode
- **WHEN** the user presses `.` (period key) outside of an input field
- **THEN** the sidebar and properties panels SHALL toggle between hidden and visible

### Requirement: Properties panel toggle
The web UI SHALL support toggling the properties panel visibility.

#### Scenario: Toggle properties
- **WHEN** the user presses `p` outside of an input field
- **THEN** the properties panel SHALL toggle between hidden and visible

### Requirement: Body panel displays object content
The body panel SHALL display the object's name as a header and its markdown body as plain text.

#### Scenario: Object with body
- **WHEN** an object with body content is selected
- **THEN** the body panel SHALL display the object name as a header and the body text in a monospace font with preserved whitespace

#### Scenario: Object without body
- **WHEN** an object with no body content is selected
- **THEN** the body panel SHALL display the object name and an italic "No content" placeholder

### Requirement: Properties panel displays formatted attributes
The properties panel SHALL display object properties with key labels and formatted values.

#### Scenario: Pinned properties
- **WHEN** an object has pinned properties (pin > 0)
- **THEN** pinned properties SHALL appear in a separate section at the top, sorted by pin order

#### Scenario: Unpinned properties
- **WHEN** an object has unpinned properties
- **THEN** unpinned properties SHALL appear below the pinned section, labeled "Details" if pinned properties exist

#### Scenario: Property display format
- **WHEN** a property is displayed
- **THEN** the key SHALL appear above the value, with emoji prefix if the property has one
