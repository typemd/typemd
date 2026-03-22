## ADDED Requirements

### Requirement: Embedded skill registry
The system SHALL embed marketplace SKILL.md files into the Go binary and provide a registry of available skills. Each skill SHALL have a name, description, and instructions body parsed from SKILL.md frontmatter.

#### Scenario: List available skills
- **WHEN** user runs `tmd instructions` with no arguments
- **THEN** the system SHALL display all available skill names with their descriptions

#### Scenario: List skills as JSON
- **WHEN** user runs `tmd instructions --json` with no arguments
- **THEN** the system SHALL output a JSON array of objects with `name` and `description` fields

### Requirement: Enriched skill output
The system SHALL output skill instructions enriched with vault context as JSON when a skill name is provided.

#### Scenario: Output skill with vault context
- **WHEN** user runs `tmd instructions explore` in an initialized vault with types defined
- **THEN** the system SHALL output JSON with `name`, `description`, `instructions` (body without frontmatter), and `context` containing type summaries

#### Scenario: Context includes type summaries
- **WHEN** user runs `tmd instructions explore` in a vault with a "book" type having properties
- **THEN** the `context.types` array SHALL contain an entry with `name`, `emoji`, `description`, and `properties` (each with `name`, `type`, and `description`)

#### Scenario: Empty vault context
- **WHEN** user runs `tmd instructions explore` in a vault with no custom types
- **THEN** the `context.types` array SHALL include only built-in types (tag, page)

### Requirement: Raw skill output
The system SHALL support a `--skill` flag that outputs the raw SKILL.md content including frontmatter, suitable for saving into a skills directory.

#### Scenario: Raw output with --skill flag
- **WHEN** user runs `tmd instructions explore --skill`
- **THEN** the system SHALL output the complete SKILL.md content including YAML frontmatter delimiters and body

#### Scenario: Raw output works without vault
- **WHEN** user runs `tmd instructions explore --skill` outside any vault
- **THEN** the system SHALL output the SKILL.md content (no vault context needed)

### Requirement: Vault override
The system SHALL check `.typemd/instructions/<skill>.md` for a vault-level override before using the embedded skill. The override file follows the same SKILL.md frontmatter format.

#### Scenario: Override replaces embedded skill
- **WHEN** a file exists at `.typemd/instructions/explore.md` and user runs `tmd instructions explore`
- **THEN** the system SHALL use the override file's instructions body instead of the embedded version

#### Scenario: Override without frontmatter
- **WHEN** a file exists at `.typemd/instructions/explore.md` with no YAML frontmatter
- **THEN** the system SHALL treat the entire file as instructions body and use the embedded skill's name and description

#### Scenario: Override with --skill flag
- **WHEN** a file exists at `.typemd/instructions/explore.md` and user runs `tmd instructions explore --skill`
- **THEN** the system SHALL output the override file's raw content

### Requirement: Unknown skill error
The system SHALL return an error with a list of available skills when an unknown skill name is provided.

#### Scenario: Unknown skill name
- **WHEN** user runs `tmd instructions nonexistent`
- **THEN** the system SHALL return an error message listing all available skill names
