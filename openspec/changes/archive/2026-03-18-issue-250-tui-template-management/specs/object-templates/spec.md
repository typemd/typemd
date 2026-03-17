## ADDED Requirements

### Requirement: SaveTemplate writes a template file to disk

The `Vault` SHALL provide a `SaveTemplate(typeName, templateName string, tmpl *Template) error` method that writes the template to `templates/<type>/<name>.md`. The file SHALL contain YAML frontmatter (from `tmpl.Properties`) followed by the body content. If the `templates/<type>/` directory does not exist, it SHALL be created. If the template file already exists, it SHALL be overwritten.

#### Scenario: Save template with properties and body

- **WHEN** `SaveTemplate("book", "review", &Template{Properties: {"status": "draft"}, Body: "## Notes\n"})` is called
- **THEN** `templates/book/review.md` SHALL contain frontmatter `status: draft` and body `## Notes\n`

#### Scenario: Save template with body only

- **WHEN** `SaveTemplate("book", "simple", &Template{Properties: {}, Body: "## My Template\n"})` is called
- **THEN** `templates/book/simple.md` SHALL contain only body `## My Template\n` without frontmatter delimiters

#### Scenario: Save template with properties only

- **WHEN** `SaveTemplate("book", "preset", &Template{Properties: {"status": "reading"}, Body: ""})` is called
- **THEN** `templates/book/preset.md` SHALL contain frontmatter `status: reading` and empty body

#### Scenario: Save template creates type directory if missing

- **WHEN** `SaveTemplate("newtype", "first", &Template{Body: "hello\n"})` is called
- **AND** `templates/newtype/` does not exist
- **THEN** `templates/newtype/` SHALL be created
- **AND** `templates/newtype/first.md` SHALL be written

#### Scenario: Save template overwrites existing file

- **WHEN** `templates/book/review.md` already exists with content "old"
- **AND** `SaveTemplate("book", "review", &Template{Body: "new\n"})` is called
- **THEN** `templates/book/review.md` SHALL contain "new\n"

### Requirement: DeleteTemplate removes a template file from disk

The `Vault` SHALL provide a `DeleteTemplate(typeName, templateName string) error` method that removes the template file at `templates/<type>/<name>.md`. If the file does not exist, it SHALL return an error. After deletion, if the type's template directory is empty, it SHALL be removed.

#### Scenario: Delete existing template

- **WHEN** `templates/book/review.md` exists
- **AND** `DeleteTemplate("book", "review")` is called
- **THEN** `templates/book/review.md` SHALL be removed

#### Scenario: Delete last template in type directory

- **WHEN** `templates/book/review.md` is the only file in `templates/book/`
- **AND** `DeleteTemplate("book", "review")` is called
- **THEN** `templates/book/review.md` SHALL be removed
- **AND** `templates/book/` directory SHALL be removed

#### Scenario: Delete nonexistent template

- **WHEN** `templates/book/nonexistent.md` does not exist
- **AND** `DeleteTemplate("book", "nonexistent")` is called
- **THEN** it SHALL return an error

### Requirement: ObjectRepository interface includes template write operations

The `ObjectRepository` interface SHALL include `SaveTemplate(typeName, name string, tmpl *Template) error` and `DeleteTemplate(typeName, name string) error` methods alongside the existing `GetTemplate` and `ListTemplates` methods.

#### Scenario: Interface includes SaveTemplate

- **WHEN** a type implements `ObjectRepository`
- **THEN** it SHALL implement `SaveTemplate(typeName, name string, tmpl *Template) error`

#### Scenario: Interface includes DeleteTemplate

- **WHEN** a type implements `ObjectRepository`
- **THEN** it SHALL implement `DeleteTemplate(typeName, name string) error`
