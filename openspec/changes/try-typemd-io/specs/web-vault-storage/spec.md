## ADDED Requirements

### Requirement: VaultStorage interface defines read operations
The system SHALL provide a VaultStorage interface with methods: `listTypes()`, `getTypeSchema(typeName)`, `listObjects(filter?)`, `getObject(id)`, and `searchObjects(keyword)`. All methods SHALL return Promises.

#### Scenario: List available types
- **WHEN** `listTypes()` is called
- **THEN** the system returns an array of TypeSchema objects parsed from `.typemd/types/*.yaml`

#### Scenario: Get a specific type schema
- **WHEN** `getTypeSchema("book")` is called
- **THEN** the system returns the parsed TypeSchema for `book` including its properties and relations

#### Scenario: List all objects
- **WHEN** `listObjects()` is called with no filter
- **THEN** the system returns an array of ObjectSummary with id, type, displayName, and emoji for every object in the vault

#### Scenario: Filter objects by type
- **WHEN** `listObjects({ type: "book" })` is called
- **THEN** the system returns only ObjectSummary entries where type is "book"

#### Scenario: Get full object
- **WHEN** `getObject("book/clean-code-01abc...")` is called
- **THEN** the system returns the full Object including parsed frontmatter properties, raw body, and display properties

#### Scenario: Search objects by keyword
- **WHEN** `searchObjects("clean")` is called
- **THEN** the system returns ObjectSummary entries whose name or body contains "clean" (case-insensitive)

### Requirement: GitHubBackend implements VaultStorage
The GitHubBackend SHALL implement the VaultStorage interface by reading data from a GitHub repository via the REST API.

#### Scenario: Fetch repo tree on initialization
- **WHEN** GitHubBackend is initialized with a repo and optional PAT
- **THEN** it SHALL fetch the repo tree via `GET /repos/{owner}/{repo}/git/trees/{sha}?recursive=1` to discover all files

#### Scenario: Parse type schemas from GitHub
- **WHEN** `listTypes()` is called on GitHubBackend
- **THEN** it SHALL fetch and parse YAML files from `.typemd/types/` directory in the repo

#### Scenario: Lazy-load object content
- **WHEN** `getObject(id)` is called
- **THEN** the GitHubBackend SHALL fetch the object's markdown file content on demand (not preloaded)

#### Scenario: Public repo without token
- **WHEN** GitHubBackend is initialized with a public repo and no PAT
- **THEN** all read operations SHALL work using unauthenticated GitHub API calls

#### Scenario: Private repo with token
- **WHEN** GitHubBackend is initialized with a private repo and a valid PAT
- **THEN** all API calls SHALL include the `Authorization: token <PAT>` header

### Requirement: In-memory object index
The system SHALL build an in-memory index from the repo tree and fetched content, enabling filtering and search without a database.

#### Scenario: Build index from file tree
- **WHEN** the repo tree is fetched
- **THEN** the system SHALL parse object paths (`objects/<type>/<name>.md`) to populate the index with id, type, and displayName

#### Scenario: Search uses in-memory index
- **WHEN** `searchObjects(keyword)` is called
- **THEN** the system SHALL match against the in-memory index using case-insensitive substring matching on name and body content

### Requirement: Frontmatter parsing in browser
The system SHALL parse YAML frontmatter from markdown files in the browser to extract object properties.

#### Scenario: Parse standard frontmatter
- **WHEN** a markdown file with `---` delimited YAML frontmatter is loaded
- **THEN** the system SHALL extract properties as key-value pairs and the remaining content as the body

#### Scenario: Parse wiki-links from body
- **WHEN** a markdown body contains `[[type/name-ulid]]` or `[[type/name-ulid|Display Text]]` syntax
- **THEN** the system SHALL extract wiki-link targets and display text
