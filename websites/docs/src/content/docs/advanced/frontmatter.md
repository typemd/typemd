---
title: Frontmatter
description: YAML frontmatter format for manually editing Object files.
sidebar:
  order: 2
---

TypeMD stores Object properties as YAML frontmatter — the block between `---` delimiters at the top of a Markdown file. This page covers the format rules for manually creating or editing Object files.

## What is frontmatter?

Frontmatter is a YAML block enclosed by triple-dash delimiters (`---`) at the very beginning of a Markdown file:

```markdown
---
name: Clean Code
description: A handbook of agile software craftsmanship
created_at: 2025-03-10T14:30:00+08:00
updated_at: 2025-03-12T09:15:00+08:00
tags:
  - tag/programming-01jqr3k5mpbvn8e0f2g7h9txyz
locked: true
status: reading
author: person/robert-martin-01jqr3k5mpbvn8e0f2g7h9txyz
---

Your notes and content go here...
```

Everything above the second `---` is YAML metadata (properties). Everything below is the Markdown body.

## Property ordering

Properties must follow a specific order in the frontmatter. System properties always come first in a fixed sequence, followed by schema-defined properties:

1. `name`
2. `description`
3. `created_at`
4. `updated_at`
5. `tags`
6. `locked`
7. *(schema-defined properties in schema order)*

TypeMD preserves this ordering when saving files. If you add properties out of order, they will be reordered on the next save via the CLI or TUI.

## Property value formats

### Strings

Plain text values. Quotes are optional unless the value contains special YAML characters:

```yaml
name: Clean Code
description: "A book about writing better code"
```

### Dates

Timestamps in [RFC 3339](https://datatracker.ietf.org/doc/html/rfc3339) format:

```yaml
created_at: 2025-03-10T14:30:00+08:00
updated_at: 2025-03-12T09:15:00+08:00
```

### Arrays

Lists of values, written as YAML arrays:

```yaml
tags:
  - tag/programming-01jqr3k5mpbvn8e0f2g7h9txyz
  - tag/software-01jqr3k5mpbvn8e0f2g7h9txyz
```

### Relation references

Relations reference other Objects by their full Object ID (`type/slug-ulid`):

```yaml
author: person/robert-martin-01jqr3k5mpbvn8e0f2g7h9txyz
related_books:
  - book/refactoring-01jqr3k5mpbvn8e0f2g7h9txyz
```

### Select values

Select properties use plain string values matching one of the defined options:

```yaml
status: reading
```

### Numbers

Numeric values without quotes:

```yaml
rating: 5
pages: 464
```

## Manual file creation

You can create Object files manually without using the CLI. Place a `.md` file in the appropriate `objects/<type>/` directory:

```
objects/book/my-new-book.md
```

The ULID suffix is optional for manually created files. TypeMD will recognize the file with or without it.

A minimal manually created file looks like:

```markdown
---
name: My New Book
---

Notes about this book...
```

## What happens on sync

When TypeMD processes your files (on vault open), it applies these rules:

- **Missing `name`** — if omitted, TypeMD generates it from the filename slug
- **Timestamps** — `created_at` and `updated_at` are preserved if present; if missing, they are set based on file metadata
- **Unknown properties** — properties not defined in the type schema are kept in the file but filtered out of the index (they won't appear in search or query results)
- **Ordering** — properties are reordered to match the canonical order on the next save

Files are always the source of truth. The index is rebuilt from files, never the other way around. See [Data Model](/developers/data-model) for details on the indexing mechanism.
