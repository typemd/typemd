---
title: Wiki-links
description: Inline references between Objects using [[target]] syntax.
sidebar:
  order: 5
---

A Wiki-link is an inline reference to another Object, written directly in your Markdown body. It's the simplest way to connect ideas as you write.

## What is a Wiki-link?

When you mention another Object in your notes, you can wrap its ID in double brackets:

```markdown
---
title: Go in Action
---

# Notes

Great introduction to Go concurrency patterns.
See also [[note/go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc]].
```

You can also add display text:

```markdown
See [[note/go-concurrency-patterns-01jqr5n8oqdxp0g4h1i2kuvabc|Go Concurrency Patterns]] for details.
```

The target must be a full Object ID (including the ULID suffix).

## How it works

1. When the index is synced (`tmd --reindex` or auto-sync), the indexer parses `[[...]]` patterns from each Object's body
2. Targets are resolved against existing Objects in the database
3. Wiki-link records are stored in the SQLite index for fast backlink lookups
4. On re-sync, wiki-links that have been removed from the body are automatically cleaned up — their backlinks are also removed
5. Unresolvable targets are stored as broken links (detectable via `tmd type validate`)

## Backlinks

Wiki-links are tracked automatically. If Object A contains `[[B]]`, then B knows that A references it — this is called a **backlink**.

Backlinks appear as a built-in property when viewing an Object:

```
backlinks: ⟵ A
```

You don't need to maintain backlinks manually. TypeMD computes them from the wiki-link syntax in your Markdown body.

## Wiki-links vs. Relations

TypeMD has two ways to connect Objects. They serve different purposes:

| | Wiki-link | Relation |
|---|-----------|----------|
| Defined in | Markdown body | Type schema (frontmatter) |
| Structured | No — freeform inline reference | Yes — named, typed, queryable |
| Bidirectional | Via backlinks (read-only, automatic) | Configurable per schema |
| Schema required | No | Yes (`type: relation`) |
| Use case | Informal references (see also, mentioned in) | Formal connections (author, project members) |

**Rule of thumb**: use Relations for connections that are part of your data model. Use wiki-links for casual references in your notes.

## Validation

`tmd type validate` includes a wiki-link validation phase that detects broken links — references to Objects that do not exist.
