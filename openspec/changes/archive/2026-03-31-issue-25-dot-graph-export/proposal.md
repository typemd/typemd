## Why

Objects in typemd are connected by relations and wiki-links, but there is no way to visualize the full relationship graph. Users must inspect objects one at a time to understand connections. Exporting to DOT format enables rendering with external tools (Graphviz, d3-graphviz) and serves as a foundation for future Web UI graph visualization.

## What Changes

- Add `tmd graph` CLI command that exports the object relation graph in DOT format to stdout
- Nodes represent objects, labeled with display names and type emoji
- Edges represent relations (labeled with property name) and wiki-links (labeled "wikilink")
- Support `--type` flag to filter the graph to specific object types
- Support `--no-wikilinks` and `--no-relations` flags to control edge types
- Output is valid DOT that can be piped to `dot -Tpng` or similar tools

## Capabilities

### New Capabilities

- `graph-export`: DOT format graph export of object relations and wiki-links via `tmd graph` command

### Modified Capabilities

(none)

## Impact

- **New files**: `cmd/graph.go` (CLI command), `core/graph.go` (DOT generation logic)
- **Modified files**: `cmd/root.go` (register command)
- **New dependencies**: none (DOT format is plain text)
