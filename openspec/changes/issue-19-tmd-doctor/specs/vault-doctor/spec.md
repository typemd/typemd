## ADDED Requirements

### Requirement: Doctor checks schema validity
The system SHALL validate all type schemas and shared properties as part of the doctor check, reusing existing `ValidateAllSchemas()`.

#### Scenario: Schema with invalid property type
- **WHEN** a type schema defines a property with an unsupported type
- **THEN** doctor reports it as an error under the "Schemas" category

#### Scenario: All schemas are valid
- **WHEN** all type schemas are structurally correct
- **THEN** doctor shows the "Schemas" category as passing with a count of schemas checked

### Requirement: Doctor checks object conformance
The system SHALL validate all objects against their type schemas as part of the doctor check, reusing existing `ValidateAllObjects()`.

#### Scenario: Object has property with wrong type
- **WHEN** an object has a property value that does not match its schema-defined type
- **THEN** doctor reports it as an error under the "Objects" category

#### Scenario: Object references unknown type
- **WHEN** an object's type directory does not correspond to any known type schema
- **THEN** doctor reports it as an error under the "Objects" category

### Requirement: Doctor checks relation integrity
The system SHALL validate that all relation endpoints reference existing objects, reusing existing `ValidateRelations()`.

#### Scenario: Relation references non-existent target
- **WHEN** a relation's target object ID does not exist in the index
- **THEN** doctor reports it as an error under the "Relations" category

#### Scenario: All relations are valid
- **WHEN** all relation endpoints reference existing objects
- **THEN** doctor shows the "Relations" category as passing

### Requirement: Doctor checks wiki-link resolution
The system SHALL check that all wiki-links resolve to existing objects, reusing existing `ValidateWikiLinks()`.

#### Scenario: Wiki-link target does not exist
- **WHEN** a wiki-link references a target that cannot be resolved
- **THEN** doctor reports it as an error under the "Wiki-links" category

### Requirement: Doctor checks name uniqueness
The system SHALL check that types with `unique: true` have no duplicate names, reusing existing `ValidateNameUniqueness()`.

#### Scenario: Two objects of a unique type share the same name
- **WHEN** two objects of a type with `unique: true` have the same name value
- **THEN** doctor reports it as an error under the "Uniqueness" category

### Requirement: Doctor detects corrupted files
The system SHALL scan all `.md` files under `objects/` and report files with unparseable YAML frontmatter. This check MUST detect files that `Walk()` silently skips.

#### Scenario: File with malformed YAML frontmatter
- **WHEN** an object file exists with invalid YAML in the frontmatter section
- **THEN** doctor reports it as an error under the "Files" category with the file path and parse error

#### Scenario: File with no frontmatter delimiter
- **WHEN** an object file exists without the `---` frontmatter delimiters
- **THEN** doctor reports it as an error under the "Files" category

#### Scenario: All files parse successfully
- **WHEN** all object files have valid YAML frontmatter
- **THEN** doctor shows the "Files" category as passing

### Requirement: Doctor checks index-disk synchronization
The system SHALL check whether the SQLite index is in sync with files on disk, and auto-rebuild if out of sync.

#### Scenario: Index is out of sync
- **WHEN** the SQLite index does not match the current state of files on disk
- **THEN** doctor automatically rebuilds the index and reports it as auto-fixed under the "Index" category

#### Scenario: Index is in sync
- **WHEN** the SQLite index matches the current state of files on disk
- **THEN** doctor shows the "Index" category as passing

### Requirement: Doctor detects orphan object directories
The system SHALL scan `objects/` for subdirectories that do not correspond to any known type schema (custom types from `.typemd/types/*.yaml` or built-in types).

#### Scenario: Directory for deleted type
- **WHEN** `objects/note/` exists but there is no `note` type schema
- **THEN** doctor reports it as a warning under the "Orphans" category with the directory path

#### Scenario: All object directories match known types
- **WHEN** every subdirectory under `objects/` corresponds to a known type schema
- **THEN** the "Orphans" category has no directory-related warnings

### Requirement: Doctor detects orphan template directories
The system SHALL scan `templates/` for subdirectories that do not correspond to any known type schema.

#### Scenario: Template directory for non-existent type
- **WHEN** `templates/note/` exists but there is no `note` type schema
- **THEN** doctor reports it as a warning under the "Orphans" category with the directory path

#### Scenario: No templates directory
- **WHEN** the `templates/` directory does not exist
- **THEN** the orphan template check is skipped without error

### Requirement: Doctor produces grouped summary output
The system SHALL output results grouped by category, with a status indicator per category and a final summary line.

#### Scenario: Healthy vault
- **WHEN** all checks pass with no issues
- **THEN** every category shows `✓` and the summary reads "No issues found"

#### Scenario: Vault with mixed results
- **WHEN** some checks find issues and index was auto-fixed
- **THEN** each category with issues shows `✗` with issue details indented below, the index category shows auto-fix status, and the summary line shows total counts (e.g., "5 issues found, 1 auto-fixed")

### Requirement: Doctor exit code reflects health status
The system SHALL exit with code 0 when no issues are found, and exit with code 1 when any issues (errors or warnings) are found. Auto-fixed issues that are fully resolved SHALL NOT count toward the exit code.

#### Scenario: Exit code on healthy vault
- **WHEN** doctor finds no issues
- **THEN** the process exits with code 0

#### Scenario: Exit code on unhealthy vault
- **WHEN** doctor finds at least one error or warning
- **THEN** the process exits with code 1

#### Scenario: Exit code when only auto-fixed
- **WHEN** the only issue was index out of sync and it was auto-fixed
- **THEN** the process exits with code 0
