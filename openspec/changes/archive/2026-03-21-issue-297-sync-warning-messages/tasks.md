## 1. TUI: Reason-based toast grouping

- [x] 1.1 Update unit tests in `tui/sync_toast_test.go` to expect reason-based group keys ("not found", "ambiguous") instead of the single "unresolved refs" group, and add test cases for mixed reasons and unknown reason fallback
- [x] 1.2 Update `refreshData()` in `tui/app.go` to map `UnresolvedRelation.Reason` to group keys: "not_found" → "not found", "ambiguous" → "ambiguous", fallback → "unresolved"
- [x] 1.3 Run `make test` to verify all tests pass
