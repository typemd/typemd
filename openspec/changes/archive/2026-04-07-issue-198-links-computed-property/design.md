## Context

Wiki-links (`[[target]]`) are already parsed, stored in SQLite (`wikilinks` table), and used for backlink queries. The `backlinks` computed property is exposed in `BuildDisplayProperties` as individual `DisplayProperty` entries (one per backlink, with `IsBacklink: true`). The `links` property is registered in the system property registry (`Computed: true`) but has no resolver or display integration — it's a placeholder.

The infrastructure is ready:
- `ObjectIndex.ListWikiLinks(objectID)` returns outgoing `StoredWikiLink` records
- `QueryService.ListWikiLinks(objectID)` wraps the index call
- `DisplayProperty` has `IsBacklink` and `FromID` fields used for backlinks display

## Goals / Non-Goals

**Goals:**

- Expose outgoing wiki-links as the `links` computed system property in display properties
- Follow the same pattern as `backlinks` (one `DisplayProperty` per link, with a type flag)
- Ensure `SetProperty("links", ...)` returns an error (already enforced by computed property infrastructure)

**Non-Goals:**

- Changing wiki-link parsing, storage, or sync logic
- Adding link validation or broken-link indicators in display
- TUI-specific rendering changes (TUI already renders `DisplayProperty` generically)

## Decisions

### 1. Display as individual entries (same as backlinks)

Each outgoing link becomes a separate `DisplayProperty` with `IsLink: true`. This mirrors the backlinks pattern and allows TUI/MCP consumers to render each link individually.

**Alternative considered:** A single `links` entry with a `[]string` value. Rejected because it breaks the established one-entry-per-reference pattern used by backlinks and reverse relations.

### 2. Insert links before backlinks in display order

Display order: schema properties → reverse relations → **links** → backlinks → local properties. Links appear before backlinks because outgoing references are more actionable than incoming ones.

### 3. Add `IsLink` field to `DisplayProperty`

Add a boolean `IsLink` field parallel to `IsBacklink`. This keeps the type system explicit and allows consumers to distinguish links from backlinks without string-matching the key.

### 4. Show resolved target IDs only (skip broken links)

Only include links with a non-empty `ToID` (successfully resolved). Broken links (empty `ToID`) are excluded from display properties since they represent parse artifacts, not meaningful references.

**Alternative considered:** Include broken links with a visual indicator. Rejected — broken link detection belongs in validation (`tmd type validate`), not in display properties.

## Risks / Trade-offs

- **[Performance]** Objects with many outgoing links produce many `DisplayProperty` entries. → Acceptable: backlinks already have the same characteristic, and typical objects have <20 links.
- **[Display noise]** Showing both links and backlinks could clutter the property panel. → Mitigated by placing links in a dedicated section; TUI already handles multiple backlink entries gracefully.
