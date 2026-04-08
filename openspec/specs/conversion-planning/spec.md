# conversion-planning Specification

## Purpose
TBD - created by archiving change issue-381-onboarding-skill. Update Purpose after archive.
## Requirements
### Requirement: Generate a conversion plan from scan data
The system SHALL produce a JSON plan file that maps each source file to a target type, property mapping, and import order. The plan SHALL be written to a file for user review before execution.

#### Scenario: Generate a plan for scanned files
- **WHEN** user runs `tmd import plan ./notes/`
- **THEN** system outputs a JSON plan file containing `types` (schemas to create), `objects` (source-to-object mapping), and `order` (import sequence)

#### Scenario: Plan references existing vault types
- **WHEN** scanning files that match an existing `book` type schema
- **THEN** the plan maps those files to the existing `book` type without including it in the `types` array (no schema creation needed)

### Requirement: Determine import order by dependency
The plan SHALL order object imports so that tags and relation targets are created before objects that reference them.

#### Scenario: Tags imported before objects that use them
- **WHEN** a plan includes tag objects and objects with `tags` property referencing those tags
- **THEN** tag objects appear before referencing objects in the `order` array

#### Scenario: Relation targets imported before sources
- **WHEN** object A has a relation property pointing to object B
- **THEN** object B appears before object A in the `order` array

#### Scenario: Circular dependency between objects
- **WHEN** object A references object B and object B references object A
- **THEN** the plan breaks the cycle by ordering one arbitrarily first, and the relation is resolved in the reconciliation pass

### Requirement: Detect existing objects for incremental import
The plan SHALL check existing vault objects by name + type match and mark duplicates with conflict status.

#### Scenario: Source file matches an existing object
- **WHEN** planning import for `clean-code.md` and a `book/clean-code-*` object already exists
- **THEN** the plan entry for that file has `conflict: "skip"` by default

#### Scenario: No matching existing object
- **WHEN** planning import for a file with no matching vault object
- **THEN** the plan entry has `conflict: "none"`

### Requirement: Plan includes type schemas to create
The plan SHALL list new type schemas that need to be created before object import can proceed.

#### Scenario: Plan with new type schemas
- **WHEN** the AI classifies files into a `recipe` type that does not exist in the vault
- **THEN** the plan's `types` array includes a `recipe` entry with the suggested schema (emoji, plural, properties)

#### Scenario: Plan with only existing types
- **WHEN** all classified files map to existing vault types
- **THEN** the plan's `types` array is empty

