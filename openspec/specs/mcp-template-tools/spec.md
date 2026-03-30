# mcp-template-tools Specification

## Purpose
TBD - created by archiving change issue-249-mcp-template-tools. Update Purpose after archive.
## Requirements
### Requirement: List templates via MCP
The MCP server SHALL expose a `list_templates` tool that returns available templates. When the optional `type` parameter is provided, the tool SHALL return only templates for that type. When `type` is omitted, the tool SHALL return templates across all types.

#### Scenario: List templates for a specific type
- **WHEN** MCP client calls `list_templates` with `type` set to a type that has templates
- **THEN** the tool returns a JSON array of objects, each with `type` and `name` fields, containing only templates for that type

#### Scenario: List all templates across types
- **WHEN** MCP client calls `list_templates` without the `type` parameter
- **THEN** the tool returns a JSON array of objects, each with `type` and `name` fields, containing templates from all types

#### Scenario: List templates for type with no templates
- **WHEN** MCP client calls `list_templates` with `type` set to a type that has no templates
- **THEN** the tool returns an empty JSON array

#### Scenario: List templates for invalid type name
- **WHEN** MCP client calls `list_templates` with `type` set to a non-existent type
- **THEN** the tool returns an error message

### Requirement: Get template via MCP
The MCP server SHALL expose a `get_template` tool that retrieves a specific template's content. Both `type` and `name` parameters are required.

#### Scenario: Get existing template
- **WHEN** MCP client calls `get_template` with valid `type` and `name`
- **THEN** the tool returns a JSON object with `type`, `name`, `properties` (map), and `body` (string) fields

#### Scenario: Get template with empty body
- **WHEN** MCP client calls `get_template` for a template that has properties but no body content
- **THEN** the tool returns the template with `properties` populated and `body` as an empty string

#### Scenario: Get non-existent template
- **WHEN** MCP client calls `get_template` with a `name` that does not exist for the given type
- **THEN** the tool returns an error message

