## ADDED Requirements

### Requirement: Plugin identity and metadata

The `markdown-import` plugin SHALL be a Claude Code plugin with skill that helps users convert existing markdown files into typemd objects.

#### Scenario: Plugin is installable from marketplace
- **WHEN** a user runs `/plugin install markdown-import@typemd-marketplace`
- **THEN** the plugin SHALL install successfully and the skill SHALL be available

### Requirement: Markdown file analysis

The skill SHALL analyze existing markdown files to determine the appropriate typemd type, extract metadata for frontmatter, and identify potential relations.

#### Scenario: Import a markdown file with clear topic
- **WHEN** the user invokes the skill with a markdown file about a book (e.g., containing title, author, notes)
- **THEN** the skill SHALL suggest creating a `book` type object, extract name and relevant properties for frontmatter, preserve the markdown body content, and place the file at `objects/book/<slug>.md`

#### Scenario: Import a markdown file with ambiguous content
- **WHEN** the user invokes the skill with a markdown file that doesn't clearly map to an existing type
- **THEN** the skill SHALL ask the user which type to use or suggest creating a new type schema

### Requirement: Frontmatter generation

The skill SHALL generate YAML frontmatter with system properties (`name`, `description`, `created_at`, `updated_at`) and type-specific properties based on the vault's type schemas.

#### Scenario: Generated frontmatter includes system properties
- **WHEN** a markdown file is converted to a typemd object
- **THEN** the frontmatter SHALL include `name` (derived from the file content), `description` (optional, extracted if available), `created_at`, and `updated_at` in that order

#### Scenario: Generated frontmatter includes type-specific properties
- **WHEN** the target type has a schema with defined properties (e.g., `book` type with `author`, `isbn`)
- **THEN** the skill SHALL attempt to extract matching values from the markdown content and populate them in the frontmatter

### Requirement: Vault context awareness

The skill SHALL read the vault's type schemas from `.typemd/types/*.yaml` to understand available types and their properties.

#### Scenario: Skill uses existing type schemas
- **WHEN** the vault has type schemas defined (e.g., `book.yaml`, `person.yaml`)
- **THEN** the skill SHALL reference these schemas when suggesting types and generating frontmatter

#### Scenario: Vault with no type schemas
- **WHEN** the vault has no type schemas
- **THEN** the skill SHALL inform the user and suggest creating basic type schemas first, or proceed with system properties only

### Requirement: Relation discovery

The skill SHALL identify potential relations between the imported object and existing objects in the vault.

#### Scenario: Detect potential wiki-links
- **WHEN** the markdown content mentions names or concepts that match existing objects in the vault
- **THEN** the skill SHALL suggest converting those mentions to wiki-links (`[[type/name-ulid]]`)

### Requirement: Batch import guidance

The skill SHALL support guiding the user through importing multiple markdown files in sequence.

#### Scenario: User has multiple files to import
- **WHEN** the user indicates they have a directory of markdown files to import
- **THEN** the skill SHALL guide the user through importing them one at a time, maintaining consistency in type assignment and relation building across files

### Requirement: Non-destructive operation

The skill SHALL NOT modify or delete original markdown files. It SHALL only create new typemd object files.

#### Scenario: Original file preserved
- **WHEN** a markdown file is imported
- **THEN** the original file SHALL remain unchanged and a new file SHALL be created under `objects/<type>/`

### Requirement: Usage documentation

The plugin SHALL include a README.md that explains the skill's purpose, usage examples, and prerequisites.

#### Scenario: README covers essential information
- **WHEN** a user reads the plugin's README.md
- **THEN** it SHALL explain what the skill does, how to invoke it, what inputs it expects, example workflows, and any prerequisites (e.g., type schemas)
