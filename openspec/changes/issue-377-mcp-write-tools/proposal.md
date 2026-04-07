## Why

The MCP server currently only exposes read-only tools (`search`, `get_object`). To enable LLM-maintained wiki workflows — where an LLM agent can create, update, and link objects — the MCP server needs write tools. This is the foundational piece that unlocks AI-driven knowledge base maintenance.

## What Changes

- Add `create_object` tool — create a new object with type, name, optional template, properties, and body
- Add `update_object` tool — update an existing object's properties and/or body (merge semantics)
- Add `link_objects` tool — create a relation between two objects
- Add `unlink_objects` tool — remove a relation between two objects
- Add `list_types` read tool — list available type schemas (needed for object creation workflows)

## Capabilities

### New Capabilities

- `mcp-write-tools`: MCP server write tools for creating, updating, linking, and unlinking objects
- `mcp-list-types`: MCP tool for listing available type schemas

### Modified Capabilities

_(none — existing read tools are unchanged)_

## Impact

- `mcp/tools.go` — add five new tool registrations and handlers
- `mcp/tools_test.go` — add unit tests for all new tools
- No changes to `core/` — all write operations delegate to existing Vault facade methods (`NewObject`, `SaveObject`, `LinkObjects`, `UnlinkObjects`, `ListTypes`)
