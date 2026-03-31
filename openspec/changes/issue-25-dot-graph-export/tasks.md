## 1. Core: DOT Graph Generation

- [ ] 1.1 Write BDD scenarios for graph export (empty vault, single object, relations, wiki-links, type filter, edge flags)
- [ ] 1.2 Implement BDD step definitions for graph export scenarios
- [ ] 1.3 Create `core/graph.go` with `GraphOptions` struct and `ExportDOT` function (make scenarios pass)
- [ ] 1.4 Add unit tests for DOT output edge cases (special characters in names, deduplication)

## 2. CLI: `tmd graph` Command

- [ ] 2.1 Write BDD scenarios for `tmd graph` CLI (output format, --type flag, --no-wikilinks, --no-relations)
- [ ] 2.2 Implement BDD step definitions for graph CLI scenarios
- [ ] 2.3 Create `cmd/graph.go` with Cobra command, flags, and flag completion
