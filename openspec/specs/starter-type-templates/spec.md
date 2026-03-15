## ADDED Requirements

### Requirement: Embedded starter type definitions

The system SHALL embed starter type schemas as YAML files within the binary. Each starter type SHALL have a name, emoji, description, and a valid type schema YAML definition.

The following starter types SHALL be available:
- `idea` (💡) — Capture and develop ideas
- `note` (📝) — Quick notes and thoughts
- `book` (📚) — Track your reading

#### Scenario: List available starter types

- **WHEN** the system queries available starter types
- **THEN** it SHALL return exactly 3 starter types: idea, note, book
- **AND** each SHALL include name, emoji, description, and valid YAML bytes

#### Scenario: Starter type YAML is valid

- **WHEN** a starter type's YAML is parsed as a TypeSchema
- **THEN** it SHALL pass type schema validation without errors

### Requirement: Interactive starter type selection during init

During `tmd init`, the system SHALL present a Bubble Tea interactive checkbox selector for choosing which starter types to install. All starter types SHALL be selected by default.

#### Scenario: Default init shows starter picker

- **WHEN** user runs `tmd init` in an interactive terminal
- **THEN** the system SHALL display a checkbox selector with all starter types listed
- **AND** all items SHALL be selected (checked) by default

#### Scenario: User confirms default selection

- **WHEN** user presses Enter without changing any selection
- **THEN** all 3 starter types SHALL be written to `.typemd/types/`

#### Scenario: User deselects some types

- **WHEN** user deselects `note` and confirms
- **THEN** only `idea` and `book` SHALL be written to `.typemd/types/`
- **AND** `note.yaml` SHALL NOT exist in `.typemd/types/`

#### Scenario: User selects none

- **WHEN** user deselects all items and confirms (or presses q/Esc)
- **THEN** no starter type files SHALL be written
- **AND** the vault SHALL be initialized with an empty `.typemd/types/` directory

#### Scenario: Select all shortcut

- **WHEN** user presses `a` during selection
- **THEN** all starter types SHALL become selected

#### Scenario: Deselect all shortcut

- **WHEN** user presses `n` during selection
- **THEN** all starter types SHALL become deselected

### Requirement: Non-interactive mode with --no-starters

The `tmd init` command SHALL accept a `--no-starters` flag that skips the starter type selection entirely.

#### Scenario: Init with --no-starters

- **WHEN** user runs `tmd init --no-starters`
- **THEN** the vault SHALL be initialized without any starter type files
- **AND** no interactive UI SHALL be displayed

#### Scenario: Default behavior without flag

- **WHEN** user runs `tmd init` without `--no-starters`
- **THEN** the interactive starter type selector SHALL be displayed

### Requirement: Starter types written as regular type schema files

Selected starter types SHALL be written as regular `.typemd/types/<name>.yaml` files, identical in format to user-created type schemas. They are fully owned and editable by the user after creation.

#### Scenario: Written files are standard type schemas

- **WHEN** starter type `book` is selected and written
- **THEN** `.typemd/types/book.yaml` SHALL exist
- **AND** its content SHALL be a valid type schema YAML parseable by `LoadType("book")`

#### Scenario: Init output lists created types

- **WHEN** starter types are written during init
- **THEN** the system SHALL print each created type with its emoji and name
