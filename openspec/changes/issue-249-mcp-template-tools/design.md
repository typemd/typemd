## Context

The MCP server (`mcp/tools.go`) registers tools on an `MCPServer` instance with handler functions that delegate to `core.Vault`. Currently two tools exist: `search` and `get_object`. The pattern is straightforward — define a `mcplib.NewTool`, register it with `s.AddTool`, and implement a handler that calls Vault methods and returns JSON.

Template CRUD already exists in core: `Vault.ListTemplates(typeName)` returns `[]string` of template names for a type, and `Vault.LoadTemplate(typeName, templateName)` returns `*Template` (Name, Properties, Body).

## Goals / Non-Goals

**Goals:**
- Expose `list_templates` and `get_template` as MCP tools
- Follow the existing tool registration pattern in `mcp/tools.go`
- Enable MCP clients to discover templates and read their content

**Non-Goals:**
- Template creation/deletion via MCP (read-only access)
- Template application during object creation via MCP
- Changes to existing `search` or `get_object` tools

## Decisions

### 1. `list_templates` accepts optional `type` parameter

When `type` is provided, list templates for that specific type. When omitted, list templates across all types by iterating `Vault.ListTypes()` and calling `ListTemplates` for each. This gives MCP clients a single call to discover all available templates.

**Alternative considered:** Require `type` parameter — rejected because clients would need to know type names upfront, adding a round-trip.

### 2. Return format includes type name in listing

Each entry in `list_templates` response includes `type` and `name` fields, so clients can unambiguously identify templates without additional lookups.

### 3. `get_template` requires both `type` and `name` parameters

Both are required to uniquely identify a template. This matches the core API signature (`LoadTemplate(typeName, templateName)`).

## Risks / Trade-offs

- **[Risk] Type with many templates slows `list_templates` without filter** → Acceptable for typical vault sizes; templates are lightweight metadata. No mitigation needed now.
- **[Trade-off] Read-only access** → Simplifies implementation and avoids write-path complexity. Template mutation can be added in a future issue if needed.
