## ADDED Requirements

### Requirement: Object ID completion uses progressive two-stage resolution
The CLI SHALL provide tab completion for object IDs in two stages: first completing the type prefix (with trailing `/`), then completing the object name within that type.

#### Scenario: Complete type prefix when no slash present
- **WHEN** user types `tmd show b<TAB>` and type `book` exists
- **THEN** completion offers `book/` (with trailing slash, no space appended)

#### Scenario: Complete object name after type prefix
- **WHEN** user types `tmd show book/clean<TAB>` and object `book/clean-code-01abc` exists
- **THEN** completion offers `book/clean-code-01abc`

#### Scenario: Multiple type matches
- **WHEN** user types `tmd show <TAB>` and types `book`, `blog` exist
- **THEN** completion offers both `book/` and `blog/`

#### Scenario: No matches
- **WHEN** user types `tmd show xyz<TAB>` and no type starts with `xyz`
- **THEN** completion offers no candidates

### Requirement: Object ID completion applies to all object-accepting commands
The `show`, `link` (args 0 and 2), and `unlink` (args 0 and 2) commands SHALL support object ID completion via `ValidArgsFunction`.

#### Scenario: show command completes its argument
- **WHEN** user types `tmd object show book/<TAB>`
- **THEN** object ID completion is triggered for argument 0

#### Scenario: link command completes from-id and to-id
- **WHEN** user types `tmd relation link <TAB>`
- **THEN** object ID completion is triggered for argument 0
- **WHEN** user types `tmd relation link book/x rel <TAB>`
- **THEN** object ID completion is triggered for argument 2

#### Scenario: unlink command completes from-id and to-id
- **WHEN** user types `tmd relation unlink <TAB>`
- **THEN** object ID completion is triggered for argument 0
- **WHEN** user types `tmd relation unlink book/x rel <TAB>`
- **THEN** object ID completion is triggered for argument 2

### Requirement: Type name completion for type-accepting commands
Commands that accept type names SHALL provide tab completion from the list of available types.

#### Scenario: migrate command completes type name
- **WHEN** user types `tmd migrate b<TAB>` and type `book` exists
- **THEN** completion offers `book`

#### Scenario: type show command completes type name
- **WHEN** user types `tmd type show p<TAB>` and types `page`, `person` exist
- **THEN** completion offers `page` and `person`

#### Scenario: stats --type flag completes type name
- **WHEN** user types `tmd stats --type b<TAB>` and type `book` exists
- **THEN** completion offers `book`

#### Scenario: format --type flag completes type name
- **WHEN** user types `tmd format --type b<TAB>` and type `book` exists
- **THEN** completion offers `book`

### Requirement: Relation name completion for link and unlink
The `link` and `unlink` commands SHALL provide tab completion for the relation name (argument 1) based on the source object's type schema.

#### Scenario: Complete relation name from source object schema
- **WHEN** user types `tmd relation link book/clean-code <TAB>` and `book` type has relation properties `author` and `series`
- **THEN** completion offers `author` and `series`

#### Scenario: No relation completion when source object is invalid
- **WHEN** user types `tmd relation link invalid/ <TAB>` and the source object cannot be resolved
- **THEN** completion offers no candidates (fails gracefully)

### Requirement: Completion script generation
The CLI SHALL provide a `tmd completion` subcommand that generates shell-specific completion scripts for bash, zsh, and fish.

#### Scenario: Generate bash completion script
- **WHEN** user runs `tmd completion bash`
- **THEN** a valid bash completion script is printed to stdout

#### Scenario: Generate zsh completion script
- **WHEN** user runs `tmd completion zsh`
- **THEN** a valid zsh completion script is printed to stdout

#### Scenario: Generate fish completion script
- **WHEN** user runs `tmd completion fish`
- **THEN** a valid fish completion script is printed to stdout

### Requirement: Completion does not require SQLite index
All completion functions SHALL work without opening the SQLite index. Completion reads from the filesystem only.

#### Scenario: Completion works without index.db
- **WHEN** `.typemd/index.db` does not exist
- **THEN** object ID and type name completion still function correctly

### Requirement: File completion is suppressed
All completion functions SHALL return `ShellCompDirectiveNoFileComp` to prevent the shell from falling back to file path completion.

#### Scenario: No file completion fallback
- **WHEN** completion returns no matches for a prefix
- **THEN** the shell does not fall back to completing file paths
