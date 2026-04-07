## ADDED Requirements

### Requirement: List types via MCP

The MCP server SHALL provide a `list_types` tool that returns all available type schemas in the vault. The tool SHALL accept no required parameters. The response SHALL include each type's name, plural form, emoji, and property definitions. This enables LLM agents to discover available types before creating objects.

#### Scenario: List types in vault with custom types

- **WHEN** the vault has custom type schemas defined
- **AND** client calls `list_types`
- **THEN** the tool returns all types including built-in types (tag, page) and custom types with their schema metadata

#### Scenario: List types in vault with only built-in types

- **WHEN** the vault has no custom type schemas
- **AND** client calls `list_types`
- **THEN** the tool returns at least the built-in types (tag, page)
