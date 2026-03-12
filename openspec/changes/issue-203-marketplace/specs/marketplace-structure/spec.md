## ADDED Requirements

### Requirement: Marketplace directory layout

The `marketplace/` directory SHALL follow Claude Code's native plugin marketplace format with `.claude-plugin/marketplace.json` at the root and plugins organized under `plugins/`.

#### Scenario: Valid marketplace directory structure
- **WHEN** the `marketplace/` directory is inspected
- **THEN** it SHALL contain `.claude-plugin/marketplace.json`, `plugins/` directory, `README.md`, and `CONTRIBUTING.md`

### Requirement: marketplace.json catalog format

The `marketplace.json` SHALL define the marketplace name as `typemd-marketplace`, include owner information, and list all plugins with relative path sources.

#### Scenario: marketplace.json contains required fields
- **WHEN** `marketplace/.claude-plugin/marketplace.json` is parsed
- **THEN** it SHALL contain `name` set to `"typemd-marketplace"`, `owner` with `name` field, `metadata` with `description` and `pluginRoot`, and a `plugins` array

#### Scenario: Plugin entry uses relative path source
- **WHEN** a plugin entry is listed in marketplace.json
- **THEN** its `source` field SHALL be a relative path starting with `./plugins/`

### Requirement: Plugin directory structure

Each plugin under `plugins/` SHALL follow Claude Code's plugin format with `.claude-plugin/plugin.json` and at least one component directory.

#### Scenario: Plugin has required files
- **WHEN** a plugin directory under `plugins/<name>/` is inspected
- **THEN** it SHALL contain `.claude-plugin/plugin.json`, at least one of `skills/` or `commands/`, and `README.md`

#### Scenario: plugin.json contains required metadata
- **WHEN** a plugin's `.claude-plugin/plugin.json` is parsed
- **THEN** it SHALL contain `name`, `description`, and `version` fields

### Requirement: Plugin naming convention

Plugin names SHALL use kebab-case. Names prefixed with `tmd-` or `typemd-` are reserved for official plugins.

#### Scenario: Community plugin with valid name
- **WHEN** a community contributor submits a plugin named `reading-notes`
- **THEN** the name SHALL be accepted

#### Scenario: Community plugin with reserved prefix
- **WHEN** a community contributor submits a plugin named `tmd-reading-notes`
- **THEN** the name SHALL be rejected by CI validation

### Requirement: Contributing guidelines

`CONTRIBUTING.md` SHALL document the submission process, naming rules, plugin structure requirements, and quality expectations.

#### Scenario: CONTRIBUTING.md exists with required sections
- **WHEN** `marketplace/CONTRIBUTING.md` is inspected
- **THEN** it SHALL contain sections for submission process, naming rules, plugin structure, and quality requirements
