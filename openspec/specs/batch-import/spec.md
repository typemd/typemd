# batch-import Specification

## Purpose
TBD - created by archiving change issue-381-onboarding-skill. Update Purpose after archive.
## Requirements
### Requirement: Execute a conversion plan
The system SHALL execute a confirmed plan file: create type schemas, then create objects in dependency order, then resolve wiki-links.

#### Scenario: Execute a plan with new types and objects
- **WHEN** user runs `tmd import execute plan.json` with a plan containing 2 new types and 10 objects
- **THEN** system creates the 2 type schemas first, then creates the 10 objects in dependency order, outputting progress for each

#### Scenario: Execute a plan with no new types
- **WHEN** executing a plan where all types already exist
- **THEN** system skips type creation and proceeds directly to object creation

### Requirement: Create type schemas from plan
The system SHALL create type schemas listed in the plan's `types` array using the specified schema definition.

#### Scenario: Create a new type schema
- **WHEN** the plan includes a type `recipe` with emoji `🍳`, plural `recipes`, and properties `ingredients` (string) and `prep_time` (number)
- **THEN** system creates `types/recipe/schema.yaml` with the specified schema

#### Scenario: Type already exists during execution
- **WHEN** the plan lists a type that was created between planning and execution
- **THEN** system skips that type and continues

### Requirement: Create objects in dependency order
The system SHALL create objects following the plan's `order` array, using `Vault.NewObject()` for each, then populating properties from the plan's property mapping.

#### Scenario: Object created with mapped properties
- **WHEN** executing a plan entry that maps source file `clean-code.md` to type `book` with properties `{name: "Clean Code", author: "Robert C. Martin"}`
- **THEN** system creates a `book` object with those properties populated in frontmatter

#### Scenario: Object body from source content
- **WHEN** executing a plan entry for a source file with markdown body content
- **THEN** the created object's body contains the source file's markdown content (excluding frontmatter)

### Requirement: Handle conflicts during execution
The system SHALL respect the `conflict` field on each plan entry: `skip` skips the file, `overwrite` replaces the existing object, `none` creates normally.

#### Scenario: Skip conflicting file
- **WHEN** executing a plan entry with `conflict: "skip"`
- **THEN** system skips that entry and reports it as skipped in the output

#### Scenario: Overwrite conflicting file
- **WHEN** executing a plan entry with `conflict: "overwrite"`
- **THEN** system updates the existing object with the new content and properties

### Requirement: Resolve wiki-links after all objects are created
The system SHALL trigger a reconciliation pass after all objects are created to resolve wiki-links in object bodies.

#### Scenario: Wiki-links resolved after batch creation
- **WHEN** all objects in the plan have been created and some contain `[[name]]` style wiki-links
- **THEN** system runs reconciliation to resolve shorthand wiki-links to full object IDs

### Requirement: Report execution progress
The system SHALL output JSON progress for each operation: type created, object created, object skipped, object failed.

#### Scenario: Progress output during execution
- **WHEN** executing a plan with 5 objects
- **THEN** system outputs a JSON line for each operation with status (`created`, `skipped`, `failed`) and the object ID or source path

