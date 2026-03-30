## Why

The MCP server currently only exposes `search` and `get_object` tools. Template functionality is not available to MCP clients, preventing AI assistants from discovering and using templates when helping users create objects. Adding template tools closes this gap and enables MCP clients to suggest appropriate templates during object creation workflows.

## What Changes

- Add `list_templates` MCP tool that lists available templates, optionally filtered by type name
- Add `get_template` MCP tool that retrieves a specific template's content (properties + body)
- Both tools delegate to existing `Vault.ListTemplates()` and `Vault.LoadTemplate()` core methods

## Capabilities

### New Capabilities
- `mcp-template-tools`: MCP tools for listing and retrieving object templates (`list_templates` with optional type filter, `get_template` with type and template name)

### Modified Capabilities
_(none — this adds new MCP tools without changing existing tool behavior or specs)_

## Impact

- **Code**: `mcp/tools.go` — new tool registrations and handlers
- **Tests**: `mcp/tools_test.go` — new test cases for template tools
- **APIs**: Two new MCP tools exposed to clients (`list_templates`, `get_template`)
- **Dependencies**: No new dependencies — uses existing `core.Vault` template methods
