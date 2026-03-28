## MODIFIED Requirements

### Requirement: Read-only property skipping
The cursor SHALL skip read-only properties during navigation. Read-only properties include: `created_at`, `updated_at`, reverse relations, backlinks, relation properties, and **local properties** (`IsLocal = true`).

#### Scenario: Skip immutable system properties
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** `created_at` and `updated_at` properties SHALL be skipped

#### Scenario: Skip reverse relations
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** reverse relation properties (IsReverse=true) SHALL be skipped

#### Scenario: Skip backlinks
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** backlink properties (IsBacklink=true) SHALL be skipped

#### Scenario: Skip relation properties
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** relation properties (IsRelation=true) SHALL be skipped (handled by #88)

#### Scenario: Skip tags property
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** the `tags` property SHALL be skipped (relation to tag type, handled by #88)

#### Scenario: Skip local properties
- **WHEN** the user navigates with j/k in the properties panel
- **THEN** local properties (IsLocal=true) SHALL be skipped
