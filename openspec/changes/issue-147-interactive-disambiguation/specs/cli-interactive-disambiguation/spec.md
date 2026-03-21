## ADDED Requirements

### Requirement: Interactive disambiguation picker for ambiguous object ID prefixes

When a CLI command resolves an object ID prefix that matches multiple objects, and the terminal is interactive, the system SHALL display a Bubble Tea selection picker listing all matching candidates. The user SHALL select one candidate using keyboard navigation (↑/↓ to move, Enter to confirm). The selected object's full ID SHALL be used to continue the command's operation.

#### Scenario: Prefix matches two objects in interactive terminal

- **WHEN** `tmd object show book/clean-code` is run in an interactive terminal
- **AND** `book/clean-code` matches `book/clean-code-01aaa` and `book/clean-code-01bbb`
- **THEN** the system SHALL display a selection picker with both candidates
- **AND** each candidate SHALL show the object's `name` property and full ID
- **AND** the user SHALL be able to select one with Enter to proceed

#### Scenario: User selects an object from the picker

- **WHEN** the disambiguation picker is displayed with multiple candidates
- **AND** the user navigates to a candidate and presses Enter
- **THEN** the command SHALL proceed as if the user had provided the selected object's full ID

#### Scenario: User cancels the picker

- **WHEN** the disambiguation picker is displayed
- **AND** the user presses Esc or q
- **THEN** the command SHALL exit with the original `AmbiguousMatchError` message
- **AND** the exit code SHALL be non-zero

### Requirement: Non-interactive fallback preserves existing error behavior

When the terminal is not interactive (stdin is not a TTY), the system SHALL NOT display the picker and SHALL return the existing `AmbiguousMatchError` error message unchanged.

#### Scenario: Ambiguous prefix in piped input

- **WHEN** `echo "book/clean-code" | tmd object show book/clean-code` is run (stdin is piped)
- **AND** the prefix matches multiple objects
- **THEN** the system SHALL print the existing error message listing all candidates
- **AND** the system SHALL exit with a non-zero exit code
- **AND** no interactive picker SHALL be displayed

### Requirement: All ID-accepting commands support interactive disambiguation

The disambiguation picker SHALL be available in all CLI commands that resolve object IDs: `object show`, `relation link` (both from-id and to-id arguments), and `relation unlink` (both from-id and to-id arguments).

#### Scenario: Disambiguation in show command

- **WHEN** `tmd object show <ambiguous-prefix>` is run interactively
- **THEN** the disambiguation picker SHALL appear for the ambiguous prefix

#### Scenario: Disambiguation in link command with ambiguous from-id

- **WHEN** `tmd relation link <ambiguous-prefix> author person/someone` is run interactively
- **AND** the from-id prefix is ambiguous
- **THEN** the disambiguation picker SHALL appear for the from-id

#### Scenario: Disambiguation in link command with ambiguous to-id

- **WHEN** `tmd relation link book/known author <ambiguous-prefix>` is run interactively
- **AND** the to-id prefix is ambiguous
- **THEN** the disambiguation picker SHALL appear for the to-id

#### Scenario: Disambiguation in unlink command

- **WHEN** `tmd relation unlink <ambiguous-prefix> author person/someone` is run interactively
- **AND** a prefix is ambiguous
- **THEN** the disambiguation picker SHALL appear for the ambiguous prefix

### Requirement: Picker displays human-readable object names

Each candidate in the disambiguation picker SHALL display the object's `name` property (from the index) as the primary label, with the full object ID shown as a secondary line beneath it.

#### Scenario: Candidate with name property

- **WHEN** the picker displays a candidate with ID `book/clean-code-01aaa`
- **AND** the object has `name: "Clean Code"`
- **THEN** the picker SHALL show "Clean Code" as the primary label
- **AND** the full ID `book/clean-code-01aaa` SHALL be shown below

#### Scenario: Candidate without name property

- **WHEN** the picker displays a candidate whose name cannot be resolved
- **THEN** the picker SHALL show the full object ID as the primary label
