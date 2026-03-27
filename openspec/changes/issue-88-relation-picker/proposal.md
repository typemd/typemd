## Why

Relation properties are currently read-only in the TUI — users can see links (e.g., `author: → person/alan`) but cannot create, change, or remove them without editing the Markdown file directly. This is the last major editing gap in the inline property editor (#87 completed), blocking the v0.7.0 "Edit Everything Inline" release.

## What Changes

- Relation properties become editable in the properties panel via a fuzzy-search picker overlay
- Single-value relations (`Multiple: false`) use a single-select picker with search
- Multi-value relations (`Multiple: true`, including `tags`) use a multi-select picker with search
- Reverse relations and backlinks remain read-only (display only)
- The picker lists all objects of the relation's target type, with real-time substring filtering
- Object display uses human-readable names (ULID stripped) in the picker
- Clearing a relation value is supported (select empty / deselect all)
- Lock guards prevent editing relations on locked objects

## Capabilities

### New Capabilities
- `relation-picker`: Fuzzy-search picker overlay for selecting target objects when editing relation properties, supporting both single and multi-value relations

### Modified Capabilities
- `inline-property-editing`: Relation properties (including tags) transition from read-only to editable via the new relation picker

## Impact

- **tui/prop_editor.go** — Remove relation/tags from `isPropertyEditable()` skip list; add relation picker activation logic
- **tui/prop_editor_update.go** — Add picker navigation, search filtering, selection handling for relations
- **tui/prop_editor.go (render)** — Add picker overlay rendering using `widget.OverlayPopup`
- **tui/model.go** — Wire relation picker messages (link/unlink results) through main Update loop
- **core/vault.go or query_service.go** — May need a convenience method to list objects by type for picker population
- No breaking changes to core APIs — uses existing `ObjectService.Link/Unlink` and `QueryService.Query`
