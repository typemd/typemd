## Context

typemd objects are connected by two types of edges: **relations** (typed, schema-defined property links between objects) and **wiki-links** (inline `[[target]]` references in markdown body). Both are indexed in SQLite. Currently there is no way to visualize the full graph.

The `tmd graph` command exports this graph as DOT format text, suitable for piping to Graphviz (`dot -Tpng`) or other visualization tools.

## Goals / Non-Goals

**Goals:**
- `tmd graph` outputs valid DOT digraph to stdout
- Nodes = objects, labeled with display name and type emoji
- Edges = relations (labeled with property name) and wiki-links (labeled "wikilink")
- `--type` flag filters to objects of specific types (multiple allowed)
- `--no-wikilinks` / `--no-relations` flags to exclude edge types
- Empty vault produces valid empty digraph
- Isolated objects (no edges) appear as nodes

**Non-Goals:**
- Rendering or visualization (external tools handle this)
- Interactive graph exploration
- Web UI graph view (future)
- Subgraph clustering by type (can be added later)
- Performance optimization for very large vaults (document limitation)

## Decisions

### Core logic in `core/graph.go`, CLI wrapper in `cmd/graph.go`

Graph generation logic lives in the core package so it can be reused by MCP, Web UI, or TUI in the future. The CLI command is a thin wrapper that calls core and writes to stdout.

### Use `QueryObjects` + `ListRelations` / `ListWikiLinks` per object

Query all objects via `Vault.QueryObjects()` (with optional type filter), then iterate to collect relations and wiki-links. This reuses existing query infrastructure. For relations, we query per object and deduplicate bidirectional pairs by keeping only edges where `FromID` matches the current object (or both directions for directed representation).

### Directed graph with deduplicated relation edges

Relations are stored as `(FromID, ToID)` pairs. Bidirectional relations produce two rows in the index. To avoid duplicate edges in the DOT output, we track seen `(from, to, name)` triples and skip duplicates. Wiki-links are inherently directional (source → target).

### Node IDs use full object ID, labels use display name

DOT node identifiers use the full object ID (e.g., `"book/clean-code-01abc"`) for uniqueness. Node labels show the display name with type emoji (e.g., `📖 Clean Code`).

### GraphOptions struct for configuration

A `GraphOptions` struct encapsulates filter flags. This keeps the core function signature clean and extensible.

## Risks / Trade-offs

- **N+1 queries** — `ListRelations` and `ListWikiLinks` are called per object. For vaults with hundreds of objects this is acceptable. A batch query could optimize this later.
- **Broken wiki-links** — Wiki-links with empty `ToID` (unresolved targets) are skipped since we cannot produce a valid edge.
- **Large output** — For very large vaults the DOT output may be unwieldy. This is documented as a known limitation; users can use `--type` to filter.
