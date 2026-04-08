---
title: Glossary
description: Definitions of key terms used throughout TypeMD.
sidebar:
  order: 10
---

A quick reference for the terms used across TypeMD documentation, CLI, and configuration files.

## Vault

A directory managed by TypeMD. Contains `.typemd/` (configuration), `objects/` (content), and optionally `templates/` (object templates). Each Vault is self-contained and can be version-controlled.

## Object

The basic unit of knowledge in TypeMD. Instead of thinking in "files" or "notes", you think in Objects — a book, a person, an idea, a meeting. Each Object is a Markdown file with YAML frontmatter (properties) and a free-form body.

[Learn more →](/concepts/objects)

## Type

Defines what kind of Object you're creating. Every Object belongs to exactly one Type (e.g. `book`, `person`, `note`). Types are defined by YAML schema files in `types/`. Optional schema fields include `plural` (grammatically correct display name for collections), `emoji`, `color`, `description`, `version`, and `unique`.

[Learn more →](/concepts/types)

## Property

A named field on an Object, defined by its Type schema. Properties have types (`string`, `number`, `date`, `select`, `relation`, etc.) and are stored in YAML frontmatter. Properties can have optional attributes like `emoji`, `pin`, and `default`.

[Learn more →](/basics/properties)

## Relation

A named, typed link between two Objects. Unlike a simple hyperlink, Relations are defined in the Type schema and can be bidirectional — updating one side automatically updates the other.

[Learn more →](/concepts/relations)

## Link

An inline reference to another Object written using wiki-link syntax (`[[...]]`) in the Markdown body. TypeMD tracks these links and their reverse connections (backlinks). Also called a **wiki-link**.

[Learn more →](/concepts/links)

## Backlink

A reverse connection automatically computed from Links. If Object A links to Object B, then B has a backlink from A. Backlinks are read-only and maintained by the index.

[Learn more →](/concepts/links#backlinks)

## Tag

A first-class Object of the built-in `tag` type. Tags have enforced name uniqueness and can be linked to any Object via the `tags` system property. Tags support `color` and `icon` properties.

[Learn more →](/basics/tags)

## Template

Provides default frontmatter and body content when creating new Objects. Object templates are stored as Markdown files in `templates/<type>/`. Name templates use `{{ date:FORMAT }}` syntax to auto-generate Object names.

[Learn more →](/basics/templates)

## View

A saved configuration that controls how Objects of a Type are displayed — including sort order, filter rules, grouping, and layout. Views are stored as YAML files in `types/<name>/views/`. Every Type has an implicit default View (list layout, sorted by name).

[Learn more →](/advanced/file-structure#views)

## Slug

The human-readable part of an Object's filename. For example, `golang-in-action` in `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md`. Derived from the name you provide when creating an Object.

When you create an Object with a natural-language name (e.g., "Clean Code"), TypeMD converts it to a slug using these rules:

1. Convert to lowercase
2. Replace spaces and underscores with hyphens
3. Remove characters that are not letters, digits, or hyphens (Unicode letters like CJK and accented characters are preserved)
4. Collapse consecutive hyphens into a single hyphen
5. Trim leading and trailing hyphens

| Input | Slug |
|-------|------|
| `Clean Code` | `clean-code` |
| `What's the plan?` | `whats-the-plan` |
| `my_great_idea` | `my-great-idea` |
| `Chapter 3 Notes` | `chapter-3-notes` |
| `clean-code` | `clean-code` (already valid) |

The original input is always preserved as the `name` property — slug conversion only affects the filename.

## ULID

A 26-character Universally Unique Lexicographically Sortable Identifier, automatically appended to the slug by the CLI. Ensures uniqueness even when two Objects share the same slug. Objects created manually (without the CLI) do not require a ULID.

## Object ID

The full identifier for an Object in `type/slug-ulid` format. For example: `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`. Used in Link targets and relation references.

## Display ID

The human-readable form of an Object ID with the ULID suffix stripped. For example: `book/golang-in-action` from `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`. Used when rendering wiki-links without display text.

## Frontmatter

The YAML metadata block at the top of an Object file, delimited by `---`. Contains properties in a fixed order: system properties first, then schema-defined properties.

[Learn more →](/advanced/frontmatter)

## Index

A SQLite database (`.typemd/index.db`) that caches Object metadata for fast querying and full-text search. Automatically built and updated — never needs manual editing.

## Sync

The process of reading Object files and updating the Index. Happens automatically when opening a Vault. During sync, missing `name` properties are added, tag references are resolved, and the search index is refreshed.

## Validation

Checks that schemas, Objects, Relations, Links, and name uniqueness constraints are all correct. Run with `tmd type validate`.

[Learn more →](/basics/validation)
