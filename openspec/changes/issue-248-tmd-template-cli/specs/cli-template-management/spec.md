## ADDED Requirements

### Requirement: CLI provides template command group

The CLI SHALL provide a `tmd template` parent command that groups template management subcommands (`list`, `show`, `create`, `delete`).

#### Scenario: Template command with no subcommand shows help

- **WHEN** `tmd template` is executed with no subcommand
- **THEN** it SHALL display usage help listing available subcommands

### Requirement: CLI lists templates filtered by type

The `tmd template list` command SHALL list all available templates. When an optional `[type]` argument is provided, it SHALL list only templates for that type. Output SHALL show one template per line in `type/name` format. When `--json` flag is provided, output SHALL be a JSON array of objects with `type` and `name` fields.

#### Scenario: List all templates across types

- **WHEN** `tmd template list` is executed
- **AND** `templates/book/` contains `review.md` and `summary.md`
- **AND** `templates/note/` contains `meeting.md`
- **THEN** it SHALL output:
  ```
  book/review
  book/summary
  note/meeting
  ```

#### Scenario: List templates filtered by type

- **WHEN** `tmd template list book` is executed
- **AND** `templates/book/` contains `review.md` and `summary.md`
- **THEN** it SHALL output:
  ```
  book/review
  book/summary
  ```

#### Scenario: List templates for type with no templates

- **WHEN** `tmd template list book` is executed
- **AND** `templates/book/` does not exist
- **THEN** it SHALL output nothing (empty output, exit 0)

#### Scenario: List all templates when none exist

- **WHEN** `tmd template list` is executed
- **AND** no templates directory exists
- **THEN** it SHALL output nothing (empty output, exit 0)

#### Scenario: List templates as JSON

- **WHEN** `tmd template list --json` is executed
- **AND** `templates/book/` contains `review.md`
- **THEN** it SHALL output a JSON array:
  ```json
  [{"type": "book", "name": "review"}]
  ```

### Requirement: CLI shows template content

The `tmd template show` command SHALL accept a `<type/name>` argument and display the template's frontmatter properties and body content. Properties SHALL be displayed under a "Properties" heading, and body under a "Body" heading.

#### Scenario: Show template with properties and body

- **WHEN** `tmd template show book/review` is executed
- **AND** `templates/book/review.md` contains frontmatter `status: draft` and body `## Notes`
- **THEN** it SHALL display the template identifier, properties section showing `status: draft`, and body section showing `## Notes`

#### Scenario: Show template with body only

- **WHEN** `tmd template show book/simple` is executed
- **AND** `templates/book/simple.md` contains only body `## My Template`
- **THEN** it SHALL display the template identifier, properties section showing `(none)`, and body section showing `## My Template`

#### Scenario: Show nonexistent template

- **WHEN** `tmd template show book/nonexistent` is executed
- **AND** `templates/book/nonexistent.md` does not exist
- **THEN** it SHALL return an error

#### Scenario: Show template with invalid argument format

- **WHEN** `tmd template show review` is executed (missing type prefix)
- **THEN** it SHALL return an error indicating the expected `type/name` format

### Requirement: CLI creates template files

The `tmd template create` command SHALL accept a `<type/name>` argument, create a new template file at `templates/<type>/<name>.md` with frontmatter delimiters, and open it in the user's editor. The editor SHALL be resolved from `$EDITOR`, falling back to `$VISUAL`, then `vi`.

#### Scenario: Create new template and open in editor

- **WHEN** `tmd template create book/review` is executed
- **AND** `templates/book/review.md` does not exist
- **THEN** it SHALL create `templates/book/review.md` with empty frontmatter delimiters
- **AND** open the file in the user's editor

#### Scenario: Create template when file already exists

- **WHEN** `tmd template create book/review` is executed
- **AND** `templates/book/review.md` already exists
- **THEN** it SHALL return an error indicating the template already exists

#### Scenario: Create template creates type directory if missing

- **WHEN** `tmd template create newtype/first` is executed
- **AND** `templates/newtype/` does not exist
- **THEN** it SHALL create `templates/newtype/` directory
- **AND** create `templates/newtype/first.md`

#### Scenario: Create template with invalid argument format

- **WHEN** `tmd template create review` is executed (missing type prefix)
- **THEN** it SHALL return an error indicating the expected `type/name` format

### Requirement: CLI deletes template files with confirmation

The `tmd template delete` command SHALL accept a `<type/name>` argument and delete the template file at `templates/<type>/<name>.md`. In interactive terminals, it SHALL prompt for confirmation before deleting. The `--force` / `-f` flag SHALL skip the confirmation prompt.

#### Scenario: Delete template with confirmation

- **WHEN** `tmd template delete book/review` is executed in an interactive terminal
- **AND** `templates/book/review.md` exists
- **AND** the user confirms deletion
- **THEN** it SHALL delete `templates/book/review.md`
- **AND** print a confirmation message

#### Scenario: Delete template with force flag

- **WHEN** `tmd template delete book/review --force` is executed
- **AND** `templates/book/review.md` exists
- **THEN** it SHALL delete `templates/book/review.md` without prompting
- **AND** print a confirmation message

#### Scenario: Delete nonexistent template

- **WHEN** `tmd template delete book/nonexistent` is executed
- **AND** `templates/book/nonexistent.md` does not exist
- **THEN** it SHALL return an error

#### Scenario: Delete template with invalid argument format

- **WHEN** `tmd template delete review` is executed (missing type prefix)
- **THEN** it SHALL return an error indicating the expected `type/name` format

### Requirement: CLI provides shell completion for template commands

The `tmd template` subcommands SHALL provide shell completion. `list` SHALL complete type names. `show`, `create`, and `delete` SHALL complete `type/name` pairs from existing templates.

#### Scenario: Completion for template list suggests type names

- **WHEN** the user types `tmd template list ` and requests completion
- **THEN** it SHALL suggest available type names that have templates

#### Scenario: Completion for template show suggests type/name pairs

- **WHEN** the user types `tmd template show ` and requests completion
- **THEN** it SHALL suggest existing `type/name` template identifiers
