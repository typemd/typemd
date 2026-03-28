---
title: Links
description: Inline references between Objects using [[target]] syntax.
sidebar:
  order: 5
---

A link is an inline reference to another Object, written directly in your Markdown body using wiki-link syntax (`[[...]]`). It's the simplest way to connect ideas as you write.

## Syntax

When you mention another Object in your notes, you can wrap its ID in double brackets:

```markdown
---
name: Go in Action
---

# Notes

Great introduction to Go concurrency patterns.
See also [[note/go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc]].
```

You can also add display text:

```markdown
See [[note/go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc|Go Concurrency Patterns]] for details.
```

### Shorthand syntax

You don't have to type the full ID every time. TypeMD supports shorter forms that are resolved automatically during sync:

| Format | Example | Resolution |
|--------|---------|------------|
| `[[type/name-ulid]]` | `[[book/clean-code-01abc...]]` | Exact match (full ID) |
| `[[type/name]]` | `[[book/clean-code]]` | Resolved by name within the specified type |
| `[[name]]` | `[[clean-code]]` | Resolved by name within the **source object's type** |

```markdown
See [[book/clean-code]] for more.
Also check [[my-other-note]].
```

When a shorthand link resolves to exactly one object, it is automatically expanded to the full ID in your file during sync. If a shorthand matches multiple objects, it is treated as a broken link — use `tmd type validate` to see the ambiguous matches.

## How it works

1. When the index is synced (on vault open), the indexer parses `[[...]]` patterns from each Object's body
2. Targets are resolved against existing Objects — full IDs are matched exactly, shorthand targets are resolved by name within the appropriate type
3. Resolved shorthand targets are written back to the source file as full IDs (like relation name expansion)
4. Link records are stored in the SQLite index for fast backlink lookups
5. On re-sync, links that have been removed from the body are automatically cleaned up — their backlinks are also removed
6. Unresolvable or ambiguous targets are stored as broken links (detectable via `tmd type validate`)

## Backlinks

Every link creates a reverse connection automatically. If Object A contains `[[B]]`, then B knows that A references it — this is called a **backlink**.

Backlinks are a powerful way to discover related content without any manual effort. You write a link in one direction, and TypeMD tracks the reverse for you.

Backlinks appear as a built-in property when viewing an Object:

```
backlinks: ⟵ A
```

You don't need to maintain backlinks manually. TypeMD computes them from the wiki-link syntax in your Markdown body. When a link is removed from the source Object's body, the corresponding backlink is automatically cleaned up during the next sync.

## Links vs. Relations

TypeMD has two ways to connect Objects. They serve different purposes:

| | Link | Relation |
|---|------|----------|
| Defined in | Markdown body | Type schema (frontmatter) |
| Structured | No — freeform inline reference | Yes — named, typed, queryable |
| Bidirectional | Via backlinks (read-only, automatic) | Configurable per schema |
| Schema required | No | Yes (`type: relation`) |
| Use case | Informal references (see also, mentioned in) | Formal connections (author, project members) |

**Rule of thumb**: use Relations for connections that are part of your data model. Use links for casual references in your notes.

## Fixing shorthand links

Use `tmd fix wikilinks` to expand all shorthand wiki-links in your vault to full IDs in a single pass. This is useful after importing content or bulk-editing files.

## Validation

`tmd type validate` includes a link validation phase that detects broken links — references to Objects that do not exist. For shorthand links that match multiple objects (ambiguous), validation reports the ambiguous target along with the list of matching full IDs.
