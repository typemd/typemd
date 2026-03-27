## ADDED Requirements

### Requirement: Property cursor navigation
When the properties panel is focused, the system SHALL display a cursor highlighting the current property. The cursor SHALL be navigable with j/k or arrow keys.

#### Scenario: Cursor appears on focus
- **WHEN** the user presses Tab to focus the properties panel
- **THEN** the first editable property SHALL be highlighted with a cursor indicator

#### Scenario: Cursor moves down
- **WHEN** the properties panel is focused and the user presses j or down arrow
- **THEN** the cursor SHALL move to the next property in the list

#### Scenario: Cursor moves up
- **WHEN** the properties panel is focused and the user presses k or up arrow
- **THEN** the cursor SHALL move to the previous property in the list

#### Scenario: Cursor wraps at boundaries
- **WHEN** the cursor is on the last property and the user presses j
- **THEN** the cursor SHALL remain on the last property (no wrap)

#### Scenario: Cursor wraps at top
- **WHEN** the cursor is on the first property and the user presses k
- **THEN** the cursor SHALL remain on the first property (no wrap)

### Requirement: Read-only property skipping
The cursor SHALL skip read-only properties during navigation. Read-only properties include: `created_at`, `updated_at`, reverse relations, and backlinks.

#### Scenario: Skip immutable system properties
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** `created_at` and `updated_at` properties SHALL be skipped

#### Scenario: Skip reverse relations
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** reverse relation properties (IsReverse=true) SHALL be skipped

#### Scenario: Skip backlinks
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** backlink properties (IsBacklink=true) SHALL be skipped

#### Scenario: Forward relations are editable
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** forward relation properties (IsRelation=true, not reverse or backlink) SHALL receive the cursor

#### Scenario: Tags property is editable
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** the `tags` property SHALL receive the cursor (editable via relation picker)

### Requirement: Pinned properties excluded from cursor
Pinned properties (Pin > 0) are displayed in the body panel, not the properties panel. They SHALL NOT appear in cursor navigation.

#### Scenario: Pinned properties not in cursor list
- **WHEN** a property has Pin > 0
- **THEN** it SHALL NOT appear in the properties panel cursor navigation

### Requirement: Name property excluded from cursor
The `name` property is displayed in the title panel and edited via rename. It SHALL NOT appear in cursor navigation.

#### Scenario: Name property not in cursor list
- **WHEN** the properties panel is focused
- **THEN** the `name` property SHALL NOT appear in cursor navigation
