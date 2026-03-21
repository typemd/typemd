## ADDED Requirements

### Requirement: Format all objects
The system SHALL provide a `tmd format` command that rewrites all object Markdown files with canonical frontmatter formatting.

#### Scenario: Format objects with out-of-order properties
- **WHEN** an object file has properties in non-canonical order (e.g., `created_at` before `name`)
- **THEN** the formatter SHALL rewrite the file with properties in canonical order: system properties first (name, description, created_at, updated_at, tags), then schema-defined properties in schema order, then extra properties alphabetically

#### Scenario: Format objects with non-canonical YAML style
- **WHEN** an object file has valid but non-canonical YAML formatting (e.g., different quoting style or indentation)
- **THEN** the formatter SHALL rewrite the file with `yaml.v3` default serialization style

#### Scenario: Already-formatted objects are not rewritten
- **WHEN** an object file already has canonical property order and YAML formatting
- **THEN** the formatter SHALL NOT modify the file

#### Scenario: Body content is preserved
- **WHEN** an object file has body content after the frontmatter
- **THEN** the formatter SHALL preserve the body content unchanged

### Requirement: Format does not update timestamps
The system SHALL NOT modify the `updated_at` property during formatting. Formatting is a pure layout change with no semantic modification.

#### Scenario: updated_at unchanged after format
- **WHEN** an object file is reformatted
- **THEN** the `updated_at` value in the output SHALL be identical to the input value

### Requirement: Filter formatting by type
The system SHALL support a `--type <name>` flag to format only objects of the specified type.

#### Scenario: Format only objects of a specific type
- **WHEN** `tmd format --type book` is run
- **THEN** only object files under `objects/book/` SHALL be considered for formatting

#### Scenario: Invalid type name
- **WHEN** `tmd format --type nonexistent` is run and no schema exists for `nonexistent`
- **THEN** the command SHALL return an error indicating the type does not exist

### Requirement: Dry-run mode
The system SHALL support a `--dry-run` flag that lists files needing formatting without writing changes.

#### Scenario: Dry-run lists unformatted files
- **WHEN** `tmd format --dry-run` is run and some files need formatting
- **THEN** the command SHALL print the paths of files that would be changed and exit with code 1

#### Scenario: Dry-run with all files formatted
- **WHEN** `tmd format --dry-run` is run and all files are already formatted
- **THEN** the command SHALL print nothing (or a success message) and exit with code 0

### Requirement: Format type schemas
The system SHALL also format type schema YAML files (`.typemd/types/<name>/schema.yaml`) by round-tripping through `MarshalTypeSchema`.

#### Scenario: Format schema with non-canonical YAML
- **WHEN** a schema file has valid but non-canonical YAML formatting
- **THEN** the formatter SHALL rewrite the schema file with canonical output from `MarshalTypeSchema`

#### Scenario: Schema filtered by type
- **WHEN** `tmd format --type book` is run
- **THEN** only the `book` schema file SHALL be considered for formatting

#### Scenario: Built-in types without schema files are skipped
- **WHEN** formatting schemas and a type is built-in without a custom schema file
- **THEN** the formatter SHALL skip that type without error

### Requirement: Format output summary
The system SHALL display a summary of formatting results after execution.

#### Scenario: Files were formatted
- **WHEN** formatting completes and 3 files were changed
- **THEN** the command SHALL print a summary like "Formatted 3 file(s)."

#### Scenario: No files needed formatting
- **WHEN** formatting completes and no files needed changes
- **THEN** the command SHALL print "All files are already formatted. No changes needed."
