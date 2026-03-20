---
title: "TypeMD 0.5.0 — Your Data, Your Perspective"
description: "v0.5.0 introduces the View system: list and table layouts, filtering, sorting, grouping, plus the new tmd stats command."
date: 2026-03-20
tags: [release, devlog]
---

At the end of 0.4.0, we mentioned aliases and search improvements were on our radar. But after using it for a while, we realized that "how to look at your data" was a more pressing pain point than "how to find it" — as objects pile up, the sidebar just doesn't cut it anymore.

So 0.5.0 tackles that head-on with the View system: **see your data the way you want**.

## View System: One Type, Many Perspectives

Every type can now have multiple Views. A View is a saved configuration for "how I want to see these objects" — which layout, what filters, how to sort, how to group.

Views live at `.typemd/types/<name>/views/<view>.yaml`, plain text files just like type schemas — fully git-trackable.

Every type has an implicit default view (list layout, sorted by name). Want to customize it? Press `v` in the TUI to enter view mode.

## List vs. Table

Two layouts, pick the one that fits:

- **List** (`layout: list`) — a clean name list, great for quick scanning.
- **Table** (`layout: table`) — a columnar view that spreads properties out. Use `columns` to pick which properties to show:

```yaml
name: reading-list
layout: table
columns: [status, rating, author]
sort:
  - property: rating
    order: desc
```

Table mode lets you see every book's status, rating, and author at a glance — no clicking required.

## Filtering: See Only What Matters

Views support structured filter rules. Each property type has its own operators — strings get `contains`, `starts_with`; numbers get `gt`, `lte`; dates get `before`, `after`:

```yaml
filter:
  - property: status
    operator: is
    value: reading
  - property: rating
    operator: gte
    value: "4"
```

Filters are applied at the query pipeline level, not in the display layer — so performance scales.

## Grouping: Nested Structure at a Glance

`group_by` supports multi-level grouping to organize objects by property values:

```yaml
group_by:
  - property: status
  - property: author
```

Group by status first, then by author — clear structure, instantly.

## View Editor: Edit Right in the TUI

Don't want to write YAML by hand? Press `e` in view mode and a side-panel editor opens up. It has property pickers and operator pickers, and every change auto-saves.

## tmd stats: The Big Picture in One Command

Want to know how many objects are in your vault? Which type has the most? Which properties are used most often?

```bash
tmd stats
```

It prints aggregate statistics so you can quickly grasp the shape of your vault.

## ⚠️ Removed tmd query

The `tmd query` command has been removed in this version. Filtering is now part of the View system's filter rules; for full-text search, use `tmd search`. If you have scripts using `tmd query`, you'll need to migrate them.

## What Else

- **Incremental index sync** — TUI data refresh switched from full reindex to incremental sync (`Projector.SyncFiles()`), much faster (#231)
- **View mode session persistence** — close the TUI and reopen it, you're right back in the view you left (#265)
- **TUI overlay refactor** — help overlay now uses lipgloss Layer/Compositor for more stable rendering (#276)
- **Fixes** — CJK text alignment, sidebar name display, duplicate columns in view editor, and more

## What's Next

With the View system in place, we'll be heading toward more advanced territory — additional layout options, cross-type queries, and the Web UI.

To upgrade:

```bash
brew upgrade typemd/tap/tmd
```

Or download pre-built binaries from [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.5.0).
