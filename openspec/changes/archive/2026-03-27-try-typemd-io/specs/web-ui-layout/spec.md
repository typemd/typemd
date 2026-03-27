## ADDED Requirements

### Requirement: Three-panel layout mirroring TUI
The web UI SHALL use a three-panel layout matching the TUI: left sidebar (object list), center body (markdown content), and optional right panel (properties).

#### Scenario: Default layout with object selected
- **WHEN** an object is selected
- **THEN** the left panel shows the grouped object list, the title panel shows `{emoji} {type} · {name}`, the center panel renders the markdown body, and the properties panel is hidden by default

#### Scenario: Toggle properties panel
- **WHEN** the user clicks the properties toggle button
- **THEN** the right properties panel SHALL appear showing the object's schema properties, relations, and backlinks

#### Scenario: Responsive narrow screen
- **WHEN** the viewport width is below the mobile breakpoint
- **THEN** the layout SHALL collapse to a single column with navigation between list view and detail view

### Requirement: Sidebar with grouped type list
The left sidebar SHALL display objects grouped by type with expandable/collapsible headers, matching the TUI list behavior.

#### Scenario: Display type groups
- **WHEN** the vault is loaded
- **THEN** each type is shown as a group header with format `{emoji} {typeName} ({count})`, and groups are sorted alphabetically

#### Scenario: Expand and collapse groups
- **WHEN** a group header is clicked
- **THEN** the group toggles between expanded (showing objects) and collapsed (hiding objects)

#### Scenario: Select object from list
- **WHEN** an object name is clicked in the sidebar
- **THEN** that object is selected, its content loads in the body panel, and the title panel updates

### Requirement: Markdown body rendering
The body panel SHALL render the object's markdown content with proper formatting.

#### Scenario: Render markdown content
- **WHEN** an object is selected
- **THEN** its markdown body is rendered with headings, lists, code blocks, links, and other standard markdown elements

#### Scenario: Render wiki-links as clickable links
- **WHEN** the body contains `[[type/name-ulid]]` syntax
- **THEN** it SHALL render as a clickable link displaying the name (ULID stripped), and clicking it SHALL navigate to that object

#### Scenario: Render wiki-links with display text
- **WHEN** the body contains `[[type/name-ulid|Display Text]]` syntax
- **THEN** it SHALL render as a clickable link displaying "Display Text"

#### Scenario: Empty body
- **WHEN** the selected object has no body content
- **THEN** the body panel SHALL show an "(empty)" placeholder

### Requirement: Title panel
The title panel SHALL display the selected object's type emoji and name above the body panel.

#### Scenario: Show title with emoji
- **WHEN** an object of a type with an emoji is selected
- **THEN** the title panel shows `{emoji} {typeName} · {displayName}`

#### Scenario: Show title without emoji
- **WHEN** an object of a type without an emoji is selected
- **THEN** the title panel shows `{typeName} · {displayName}`

### Requirement: Properties panel
The properties panel SHALL display the selected object's schema properties, formatted values, and backlinks.

#### Scenario: Display schema properties
- **WHEN** the properties panel is visible and an object is selected
- **THEN** it SHALL list each property defined in the type schema with its formatted value

#### Scenario: Display relation properties as links
- **WHEN** a property value is a relation reference (e.g., `person/robert-martin-01abc...`)
- **THEN** it SHALL render as a clickable link that navigates to that object

#### Scenario: Display backlinks
- **WHEN** other objects link to the selected object via wiki-links
- **THEN** the properties panel SHALL show a "backlinks" section listing those objects as clickable links

### Requirement: Search
The sidebar SHALL include a search input that filters the object list.

#### Scenario: Instant search filtering
- **WHEN** the user types into the search input
- **THEN** the sidebar SHALL show matching objects in a flat list (no group headers), filtering as the user types

#### Scenario: Clear search
- **WHEN** the user clears the search input
- **THEN** the sidebar SHALL return to the grouped type list view
