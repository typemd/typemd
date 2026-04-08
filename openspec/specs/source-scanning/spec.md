# source-scanning Specification

## Purpose

Defines how `tmd import scan` collects source markdown files and extracts frontmatter patterns to feed into the AI-assisted import workflow. The scan output provides the raw material for conversion planning.

## Requirements
### Requirement: Scan source directories for markdown files
The system SHALL accept one or more source paths (directories or files) and scan them for markdown files. The scan SHALL collect file count, directory structure, and size distribution.

#### Scenario: Scan a directory with markdown files
- **WHEN** user runs `tmd import scan ./notes/`
- **THEN** system outputs a JSON result containing all `.md` files found, their paths, sizes, and the directory structure

#### Scenario: Scan multiple sources
- **WHEN** user runs `tmd import scan ./notes/ ./docs/ ./ideas.md`
- **THEN** system scans all three sources and produces a single combined result

#### Scenario: Scan a non-existent path
- **WHEN** user runs `tmd import scan ./nonexistent/`
- **THEN** system returns an error indicating the path does not exist

#### Scenario: Scan a directory with no markdown files
- **WHEN** user runs `tmd import scan ./images/` (contains only `.png` files)
- **THEN** system outputs a result with `file_count: 0` and an empty sources list

### Requirement: Extract frontmatter patterns from scanned files
The system SHALL parse YAML frontmatter from each markdown file and collect aggregate statistics: which keys appear, their frequency, and sample values.

#### Scenario: Files with YAML frontmatter
- **WHEN** scanning a directory where 8 of 10 files have a `title` key and 5 have an `author` key
- **THEN** the patterns section shows `title` with count 8 and `author` with count 5, with sample values for each

#### Scenario: Files without frontmatter
- **WHEN** scanning files that have no YAML frontmatter
- **THEN** the patterns section shows an empty key list and reports the count of files without frontmatter

### Requirement: Analyze against existing vault types
The system SHALL, when run inside an initialized vault, include existing type schemas in the scan output so the AI skill can compare source files against existing types.

#### Scenario: Scan in a vault with existing types
- **WHEN** user runs `tmd import scan ./notes/` inside a vault that has `book` and `person` types
- **THEN** the scan result includes an `existing_types` field listing the vault's type schemas with their properties

#### Scenario: Scan outside a vault
- **WHEN** user runs `tmd import scan ./notes/` in a directory without `.typemd/`
- **THEN** the scan succeeds with `existing_types` as an empty list

