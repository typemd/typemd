## ADDED Requirements

### Requirement: Vault provides template path helpers

The `Vault` struct SHALL provide `TemplatesDir()`, `TypeTemplatesDir(typeName)`, and `TemplatePath(typeName, templateName)` methods that return paths under `templates/` at the vault root.

#### Scenario: TemplatesDir returns vault root templates directory

- **WHEN** `TemplatesDir()` is called on a vault rooted at `/my/vault`
- **THEN** it SHALL return `/my/vault/templates`

#### Scenario: TypeTemplatesDir returns type-specific templates directory

- **WHEN** `TypeTemplatesDir("book")` is called on a vault rooted at `/my/vault`
- **THEN** it SHALL return `/my/vault/templates/book`

#### Scenario: TemplatePath returns full template file path

- **WHEN** `TemplatePath("book", "review")` is called on a vault rooted at `/my/vault`
- **THEN** it SHALL return `/my/vault/templates/book/review.md`

### Requirement: ListTemplates discovers available templates for a type

The `Vault` SHALL provide a `ListTemplates(typeName)` method that returns a list of template names available for the given type. Template names SHALL be derived from filenames (without the `.md` extension) in `templates/<type>/`. If the directory does not exist or is empty, it SHALL return an empty list.

#### Scenario: Type has multiple templates

- **WHEN** `templates/book/` contains `review.md` and `summary.md`
- **THEN** `ListTemplates("book")` SHALL return `["review", "summary"]`

#### Scenario: Type has one template

- **WHEN** `templates/book/` contains `default.md`
- **THEN** `ListTemplates("book")` SHALL return `["default"]`

#### Scenario: Type has no templates directory

- **WHEN** `templates/book/` does not exist
- **THEN** `ListTemplates("book")` SHALL return an empty list

#### Scenario: Type templates directory is empty

- **WHEN** `templates/book/` exists but contains no `.md` files
- **THEN** `ListTemplates("book")` SHALL return an empty list

### Requirement: LoadTemplate reads and parses a template file

The `Vault` SHALL provide a `LoadTemplate(typeName, templateName)` method that reads the template file and returns its parsed frontmatter properties and body content. The template file SHALL be parsed using the same frontmatter parser used for object files.

#### Scenario: Template with frontmatter and body

- **WHEN** `LoadTemplate("book", "review")` is called
- **AND** `templates/book/review.md` contains frontmatter `status: draft` and body `## Notes`
- **THEN** it SHALL return properties `{"status": "draft"}` and body `"## Notes\n"`

#### Scenario: Template with body only

- **WHEN** `LoadTemplate("book", "simple")` is called
- **AND** `templates/book/simple.md` contains only body content `## My Book`
- **THEN** it SHALL return empty properties and body `"## My Book\n"`

#### Scenario: Template with frontmatter only

- **WHEN** `LoadTemplate("book", "preset")` is called
- **AND** `templates/book/preset.md` contains only frontmatter `status: reading`
- **THEN** it SHALL return properties `{"status": "reading"}` and empty body

#### Scenario: Template file not found

- **WHEN** `LoadTemplate("book", "nonexistent")` is called
- **AND** `templates/book/nonexistent.md` does not exist
- **THEN** it SHALL return an error

### Requirement: NewObject applies template when templateName is provided

The `NewObject(typeName, name, templateName)` method SHALL load and apply the specified template when `templateName` is non-empty. Template frontmatter properties SHALL override schema defaults. Template body SHALL become the initial body of the new object. When `templateName` is empty, behavior SHALL be identical to the current implementation (empty body, schema defaults only).

#### Scenario: Object created with template providing body and properties

- **WHEN** `NewObject("book", "my-book", "review")` is called
- **AND** the `review` template has frontmatter `status: draft` and body `## Review Notes`
- **THEN** the created object SHALL have `status: draft` and body `## Review Notes`

#### Scenario: Object created with template overriding schema default

- **WHEN** `NewObject("book", "my-book", "review")` is called
- **AND** the schema defines `status` with default `to-read`
- **AND** the `review` template has frontmatter `status: draft`
- **THEN** the created object SHALL have `status: draft` (template wins over schema default)

#### Scenario: Object created without template

- **WHEN** `NewObject("book", "my-book", "")` is called
- **THEN** the created object SHALL have an empty body and schema default values (current behavior)

#### Scenario: Template provides mutable system property name

- **WHEN** `NewObject("book", "", "daily")` is called
- **AND** the `daily` template has frontmatter `name: daily-reading`
- **THEN** the created object SHALL have `name: daily-reading`

#### Scenario: Template provides mutable system property description

- **WHEN** `NewObject("book", "my-book", "review")` is called
- **AND** the `review` template has frontmatter `description: A book review`
- **THEN** the created object SHALL have `description: A book review`

#### Scenario: Template provides immutable system property created_at

- **WHEN** `NewObject("book", "my-book", "review")` is called
- **AND** the `review` template has frontmatter `created_at: 2020-01-01T00:00:00Z`
- **THEN** the created object SHALL have `created_at` set to the actual creation time, NOT the template value

#### Scenario: Template provides immutable system property updated_at

- **WHEN** `NewObject("book", "my-book", "review")` is called
- **AND** the `review` template has frontmatter `updated_at: 2020-01-01T00:00:00Z`
- **THEN** the created object SHALL have `updated_at` set to the actual creation time, NOT the template value

#### Scenario: Template specifies nonexistent template name

- **WHEN** `NewObject("book", "my-book", "nonexistent")` is called
- **AND** no template named `nonexistent` exists for type `book`
- **THEN** it SHALL return an error

#### Scenario: Template property not in schema is ignored

- **WHEN** `NewObject("book", "my-book", "review")` is called
- **AND** the `review` template has frontmatter `unknown_prop: value`
- **AND** the schema does not define `unknown_prop`
- **THEN** the created object SHALL NOT include `unknown_prop` in its properties

### Requirement: CLI create command supports template flag

The `tmd object create` command SHALL accept an optional `-t` / `--template` flag to specify a template name. When the flag is provided, it SHALL pass the template name to `NewObject`.

#### Scenario: Create with explicit template flag

- **WHEN** `tmd object create book my-book -t review` is executed
- **THEN** the object SHALL be created using the `review` template

#### Scenario: Create with long template flag

- **WHEN** `tmd object create book my-book --template review` is executed
- **THEN** the object SHALL be created using the `review` template

### Requirement: CLI auto-applies single template

When a type has exactly one template and no `-t` flag is specified, `tmd object create` SHALL automatically apply that template.

#### Scenario: Single template auto-applied

- **WHEN** `tmd object create book my-book` is executed
- **AND** `templates/book/` contains only `default.md`
- **THEN** the object SHALL be created using the `default` template

#### Scenario: No templates available

- **WHEN** `tmd object create book my-book` is executed
- **AND** `templates/book/` does not exist
- **THEN** the object SHALL be created with empty body and schema defaults (current behavior)

### Requirement: CLI prompts for template selection when multiple exist

When a type has multiple templates and no `-t` flag is specified, `tmd object create` SHALL present an interactive selection prompt listing available template names. The user selects one, and that template is applied.

#### Scenario: Multiple templates trigger interactive selection

- **WHEN** `tmd object create book my-book` is executed
- **AND** `templates/book/` contains `review.md` and `summary.md`
- **THEN** the command SHALL prompt the user to select between `review` and `summary`
- **AND** apply the selected template

#### Scenario: Multiple templates with -t flag skips prompt

- **WHEN** `tmd object create book my-book -t review` is executed
- **AND** `templates/book/` contains `review.md` and `summary.md`
- **THEN** the object SHALL be created using `review` without prompting

### Requirement: CLI create command supports type flag

The `tmd object create` command SHALL accept an optional `--type` flag to specify the object type. When `--type` is provided, all positional arguments are treated as the name. The `--type` flag SHALL NOT have a `-t` short form (reserved by `--template`).

#### Scenario: Create with --type flag and name

- **WHEN** `tmd object create --type note "Meeting Notes"` is executed
- **THEN** the object SHALL be created with type `note` and name `"Meeting Notes"`

#### Scenario: Create with --type flag and no name

- **WHEN** `tmd object create --type idea` is executed
- **AND** the `idea` type has a name template
- **THEN** the object SHALL be created with the auto-generated name

#### Scenario: Create with --type flag overriding config default

- **WHEN** config has `cli.default_type: idea`
- **AND** `tmd object create --type note "Meeting Notes"` is executed
- **THEN** the object SHALL be created with type `note` (flag overrides config)

### Requirement: CLI create command type argument is optional

The `tmd object create` command SHALL accept 0 to 2 positional arguments. When the type is not provided as a positional argument, it SHALL be resolved from the `--type` flag or `cli.default_type` config.

#### Scenario: Zero args with config default type and name template

- **WHEN** `tmd object create` is executed with no arguments
- **AND** config has `cli.default_type: idea`
- **AND** the `idea` type has a name template
- **THEN** the object SHALL be created with type `idea` and auto-generated name

#### Scenario: Zero args without config or flag

- **WHEN** `tmd object create` is executed with no arguments
- **AND** no `--type` flag is provided
- **AND** no `cli.default_type` is configured
- **THEN** the command SHALL return an error indicating type is required

#### Scenario: One arg resolved as type (backward compatible)

- **WHEN** `tmd object create book` is executed
- **AND** `book` is a valid type in the vault
- **THEN** the object SHALL be created with type `book` (backward compatible behavior)

#### Scenario: One arg resolved as name with config default

- **WHEN** `tmd object create "Some Thought"` is executed
- **AND** `"Some Thought"` is NOT a valid type in the vault
- **AND** config has `cli.default_type: idea`
- **THEN** the object SHALL be created with type `idea` and name `"Some Thought"`

#### Scenario: One arg not a valid type and no default

- **WHEN** `tmd object create "Some Thought"` is executed
- **AND** `"Some Thought"` is NOT a valid type in the vault
- **AND** no `cli.default_type` is configured and no `--type` flag provided
- **THEN** the command SHALL return an error indicating unknown type

#### Scenario: Two args (backward compatible)

- **WHEN** `tmd object create book "Clean Code"` is executed
- **THEN** the object SHALL be created with type `book` and name `"Clean Code"` (unchanged behavior)
