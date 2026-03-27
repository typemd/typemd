## ADDED Requirements

### Requirement: Single-value relation picker
The system SHALL display a searchable picker when the user activates editing on a single-value relation property (`Multiple: false`). The picker SHALL list all objects of the relation's target type. Selecting an object SHALL set it as the relation value via `LinkObjects`.

#### Scenario: Activate single-value relation picker
- **WHEN** the cursor is on a single-value relation property and the user presses Enter
- **THEN** a relation picker SHALL appear showing a search input and a list of candidate objects of the target type

#### Scenario: Filter candidates by search text
- **WHEN** the user types in the search input
- **THEN** the candidate list SHALL filter to show only objects whose display name contains the search text (case-insensitive substring match)

#### Scenario: Select a relation target
- **WHEN** the user navigates to a candidate and presses Enter
- **THEN** the relation SHALL be set to the selected object via `LinkObjects` and the picker SHALL close

#### Scenario: Clear single-value relation
- **WHEN** the user selects the "(none)" option at the top of the picker
- **THEN** the existing relation SHALL be removed via `UnlinkObjects` and the picker SHALL close

#### Scenario: Cancel single-value relation picker
- **WHEN** the user presses Esc in the relation picker
- **THEN** the picker SHALL close and the relation value SHALL remain unchanged

#### Scenario: Current value highlighted on open
- **WHEN** the single-value relation picker opens and a relation value already exists
- **THEN** the picker cursor SHALL be positioned on the currently linked object

### Requirement: Multi-value relation picker
The system SHALL display a searchable multi-select picker when the user activates editing on a multi-value relation property (`Multiple: true`). The picker SHALL list all objects of the relation's target type with checkmarks for currently linked objects.

#### Scenario: Activate multi-value relation picker
- **WHEN** the cursor is on a multi-value relation property and the user presses Enter
- **THEN** a relation picker SHALL appear showing a search input and a list of candidate objects with checkmarks

#### Scenario: Toggle relation in multi-picker
- **WHEN** the user presses Space on a candidate in the multi-value relation picker
- **THEN** the candidate SHALL toggle between selected and unselected

#### Scenario: Confirm multi-value relation changes
- **WHEN** the user presses Enter in the multi-value relation picker
- **THEN** newly selected objects SHALL be linked via `LinkObjects`, newly deselected objects SHALL be unlinked via `UnlinkObjects`, and the picker SHALL close

#### Scenario: Cancel multi-value relation picker
- **WHEN** the user presses Esc in the multi-value relation picker
- **THEN** the picker SHALL close and all relation values SHALL remain unchanged

#### Scenario: Currently linked objects shown as checked
- **WHEN** the multi-value relation picker opens
- **THEN** objects that are already linked SHALL show a checkmark (☑) and unlinked objects SHALL show an empty checkbox (☐)

### Requirement: Tags editing via relation picker
The `tags` system property SHALL be editable using the multi-value relation picker. The picker SHALL list all objects of type `tag`.

#### Scenario: Edit tags property
- **WHEN** the cursor is on the `tags` property and the user presses Enter
- **THEN** a multi-value relation picker SHALL appear listing all tag objects

#### Scenario: Add a tag
- **WHEN** the user selects an unchecked tag and presses Enter to confirm
- **THEN** the tag SHALL be linked to the object

#### Scenario: Remove a tag
- **WHEN** the user deselects a checked tag and presses Enter to confirm
- **THEN** the tag SHALL be unlinked from the object

### Requirement: Relation picker navigation
The relation picker SHALL support keyboard navigation consistent with existing picker patterns.

#### Scenario: Navigate candidates with j/k
- **WHEN** the relation picker is open
- **THEN** the user SHALL navigate the candidate list using j/k or arrow keys

#### Scenario: Search input receives keystrokes
- **WHEN** the relation picker is open and the user types alphanumeric characters
- **THEN** the search input SHALL update and the candidate list SHALL re-filter in real time

#### Scenario: Backspace in search input
- **WHEN** the user presses Backspace in the search input
- **THEN** the last character SHALL be removed and the candidate list SHALL re-filter

#### Scenario: Empty search shows all candidates
- **WHEN** the search input is empty
- **THEN** all candidate objects of the target type SHALL be shown

### Requirement: Relation picker display
The relation picker SHALL display candidate objects with human-readable names.

#### Scenario: Display name without ULID
- **WHEN** a candidate object is displayed in the picker
- **THEN** its name SHALL be shown without the ULID suffix

#### Scenario: Type prefix for untyped relations
- **WHEN** a relation has no target type constraint
- **THEN** candidate names SHALL be prefixed with their type (e.g., `book/golang-in-action`)

#### Scenario: Help bar shows picker mode
- **WHEN** the relation picker is active
- **THEN** the help bar SHALL show `[PICK]` to indicate picker mode

### Requirement: Locked objects block relation editing
Relation editing SHALL be blocked for locked objects, consistent with other property editing.

#### Scenario: Relation picker blocked on locked object
- **WHEN** the user attempts to edit a relation property on a locked object
- **THEN** the picker SHALL NOT open and a toast notification SHALL display "Object is locked. Unlock to edit."

### Requirement: Reverse relations and backlinks remain read-only
Reverse relations (`IsReverse: true`) and backlinks (`IsBacklink: true`) SHALL NOT be editable. The property cursor SHALL skip them during navigation.

#### Scenario: Reverse relation not editable
- **WHEN** the cursor navigation reaches a reverse relation property
- **THEN** the property SHALL be skipped (not focusable)

#### Scenario: Backlink not editable
- **WHEN** the cursor navigation reaches a backlink property
- **THEN** the property SHALL be skipped (not focusable)
