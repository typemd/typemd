## Purpose

Name templates allow type schemas to define auto-generated object names using placeholders like `{{ date:YYYY-MM-DD }}`. This enables consistent naming patterns for objects like journals, meeting notes, or time-based entries.

## Requirements

### Requirement: Type schema supports name template for auto-generated object names

A type schema MAY include a `name` entry in its `properties` array with a `template` field. The template defines a pattern for auto-generating the object's `name` property at creation time. The `name` entry SHALL only allow the `template` field — no `type`, `options`, `pin`, `emoji`, or other property fields.

#### Scenario: Type schema with name template

- **WHEN** a type schema YAML contains `- name: name` with `template: "日記 {{ date:YYYY-MM-DD }}"`
- **THEN** the loaded TypeSchema SHALL have its NameTemplate field set to "日記 {{ date:YYYY-MM-DD }}"

#### Scenario: Type schema without name template

- **WHEN** a type schema YAML does not contain a `name` entry in its properties
- **THEN** the loaded TypeSchema SHALL have its NameTemplate field set to an empty string

#### Scenario: Name entry with disallowed fields rejected

- **WHEN** a type schema YAML contains `- name: name` with `type: string`
- **THEN** schema validation SHALL return an error indicating only `template` is allowed on the `name` system property entry

#### Scenario: Name entry with template and other fields rejected

- **WHEN** a type schema YAML contains `- name: name` with `template: "{{ date:YYYY-MM-DD }}"` and `emoji: 📝`
- **THEN** schema validation SHALL return an error indicating only `template` is allowed on the `name` system property entry

### Requirement: Template evaluation at object creation time

When `NewObject()` is called without an explicit name (empty string) and the type schema has a NameTemplate, the template SHALL be evaluated and the result used as the object's `name` property value and slug.

#### Scenario: Object created with template-generated name

- **WHEN** a type schema has NameTemplate "日記 {{ date:YYYY-MM-DD }}"
- **AND** `NewObject("journal", "")` is called on 2026-03-14
- **THEN** the created object SHALL have `name: "日記 2026-03-14"` in its properties
- **AND** the filename SHALL contain the slugified template result

#### Scenario: Object created with explicit name overrides template

- **WHEN** a type schema has NameTemplate "日記 {{ date:YYYY-MM-DD }}"
- **AND** `NewObject("journal", "我的日記")` is called
- **THEN** the created object SHALL have `name: "我的日記"` in its properties

#### Scenario: Object creation fails when no name and no template

- **WHEN** a type schema has no NameTemplate
- **AND** `NewObject("book", "")` is called
- **THEN** it SHALL return an error indicating a name is required

### Requirement: Date placeholder supports user-friendly format syntax

The `{{ date:FORMAT }}` placeholder SHALL accept format tokens `YYYY`, `MM`, `DD`, `HH`, `mm`, `ss` and convert them to Go reference time equivalents for formatting. The date used SHALL be the current time at object creation.

#### Scenario: Date placeholder with YYYY-MM-DD format

- **WHEN** a template contains `{{ date:YYYY-MM-DD }}`
- **AND** the current date is 2026-03-14
- **THEN** the placeholder SHALL be replaced with "2026-03-14"

#### Scenario: Date placeholder with YYYY-MM format

- **WHEN** a template contains `{{ date:YYYY-MM }}`
- **AND** the current date is 2026-03-14
- **THEN** the placeholder SHALL be replaced with "2026-03"

#### Scenario: Date placeholder with datetime format

- **WHEN** a template contains `{{ date:YYYY-MM-DD HH:mm }}`
- **AND** the current time is 2026-03-14 09:30
- **THEN** the placeholder SHALL be replaced with "2026-03-14 09:30"

#### Scenario: Template with static text and date placeholder

- **WHEN** a template is "日記 {{ date:YYYY-MM-DD }}"
- **AND** the current date is 2026-03-14
- **THEN** the result SHALL be "日記 2026-03-14"

#### Scenario: Template with no placeholders used as literal

- **WHEN** a template is "Weekly Review"
- **THEN** the result SHALL be "Weekly Review"

### Requirement: CLI name argument is optional when template exists

The `tmd object create` command SHALL accept 1 or 2 arguments. When only the type argument is provided and the type has a NameTemplate, the template SHALL be used to generate the name. When no template exists and no name is provided, the command SHALL return an error.

#### Scenario: Create with template, no name argument

- **WHEN** user runs `tmd object create journal`
- **AND** the journal type has a name template
- **THEN** the object SHALL be created with the template-generated name

#### Scenario: Create with template, name argument provided

- **WHEN** user runs `tmd object create journal "my-journal"`
- **AND** the journal type has a name template
- **THEN** the object SHALL be created with name "my-journal" (template ignored)

#### Scenario: Create without template, no name argument

- **WHEN** user runs `tmd object create book`
- **AND** the book type has no name template
- **THEN** the command SHALL return an error indicating a name is required
