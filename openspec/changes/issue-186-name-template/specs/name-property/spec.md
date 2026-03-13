## MODIFIED Requirements

### Requirement: All objects have a required name property

Every object SHALL have a `name` property that serves as its primary display title. The `name` property is implicit — it does not need to be declared in type schemas. The `name` property SHALL be stored in the object's YAML frontmatter. When a type schema defines a name template, the `name` property SHALL be auto-populated from the evaluated template at creation time if no explicit name is provided.

#### Scenario: New object has name populated from slug

- **WHEN** a new object is created with slug "golang-in-action"
- **THEN** the object's frontmatter SHALL contain `name: golang-in-action`

#### Scenario: Object with explicit name

- **WHEN** an object's frontmatter contains `name: "Go 語言實戰"`
- **THEN** `GetName()` SHALL return "Go 語言實戰"

#### Scenario: New object has name populated from template

- **WHEN** a new object is created with type "journal" that has name template "日記 {{ date:YYYY-MM-DD }}"
- **AND** no explicit name is provided
- **AND** the current date is 2026-03-14
- **THEN** the object's frontmatter SHALL contain `name: 日記 2026-03-14`
