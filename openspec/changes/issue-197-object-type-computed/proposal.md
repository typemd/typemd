## Why

The computed system property infrastructure was added in #371, but no computed properties are actually resolvable yet. `object_type` is the simplest computed property — derived from the file path convention `objects/<type>/<name>.md` — and formalizing it as a resolvable system property enables querying, filtering, and displaying the type through the property API. This is the first concrete use of the computed property infrastructure and establishes the pattern for `links`, `backlinks`, `created_by`, and `updated_by`.

## What Changes

- Add a `GetProperty(name)` method on `Object` that resolves both stored and computed properties
- For `object_type`, return `Object.Type` (already available, zero computation cost)
- `SetProperty("object_type", ...)` must reject writes (already handled by immutability check)
- Expose `object_type` in query results and object detail views

## Capabilities

### New Capabilities

- `computed-object-type`: Resolve `object_type` as a computed system property via `GetProperty()`, with read-only enforcement and consistent access across CLI, TUI, MCP, and Web

### Modified Capabilities

_(none)_

## Impact

- **core/object.go** — new `GetProperty()` method on `Object`
- **core/query_service.go** / **core/object_service.go** — use `GetProperty()` where `object_type` is needed
- **cmd/**, **tui/**, **mcp/**, **web/** — any consumer that displays object properties may surface `object_type`
- No breaking changes — `Object.Type` field remains, `GetProperty` is additive
