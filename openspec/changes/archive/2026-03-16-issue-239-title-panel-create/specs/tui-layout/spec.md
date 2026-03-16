## MODIFIED Requirements

### Requirement: Title panel hidden when no object selected

The title panel SHALL NOT be displayed when no object is selected, UNLESS object creation mode is active. When creation mode is active, the title panel SHALL be visible to display the creation form.

#### Scenario: No selection

- **WHEN** no object is selected in the list
- **AND** no creation mode is active
- **THEN** the title panel SHALL be hidden and the body panel SHALL display the default placeholder message

#### Scenario: No selection but creation active

- **WHEN** no object is selected in the list
- **AND** object creation mode is active (user pressed `n` or `N`)
- **THEN** the title panel SHALL be visible displaying the creation form
