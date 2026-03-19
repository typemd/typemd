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

Defines what kind of Object you're creating. Every Object belongs to exactly one Type (e.g. `book`, `person`, `note`). Types are defined by YAML schema files in `.typemd/types/`.

[Learn more →](/concepts/types)

## Property

A named field on an Object, defined by its Type schema. Properties have types (`string`, `number`, `date`, `select`, `relation`, etc.) and are stored in YAML frontmatter. Properties can have optional attributes like `emoji`, `pin`, and `default`.

[Learn more →](/basics/properties)

## Relation

A named, typed link between two Objects. Unlike a simple hyperlink, Relations are defined in the Type schema and can be bidirectional — updating one side automatically updates the other.

[Learn more →](/concepts/relations)

## Link

An inline reference to another Object using `[[type/slug-ulid]]` syntax in the Markdown body. TypeMD tracks these links and their reverse connections (backlinks).

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

A saved configuration that controls how Objects of a Type are displayed — including sort order, filter rules, grouping, and layout. Views are stored as YAML files in `.typemd/types/<name>/views/`. Every Type has an implicit default View (list layout, sorted by name).

[Learn more →](/advanced/file-structure#views)

## System Property

A property managed by TypeMD that exists on every Object: `name`, `description`, `created_at`, `updated_at`, and `tags`. System properties appear first in frontmatter in a fixed order. User-authored system properties (`name`, `description`, `tags`) can be overridden by templates; auto-managed ones (`created_at`, `updated_at`) cannot.

[Learn more →](/basics/properties#system-properties)

## Shared Property

A reusable property definition in `.typemd/properties.yaml` that can be referenced from type schemas via the `use` keyword. Shared properties ensure consistent definitions across types. When referencing, you can override `pin`, `emoji`, and `description` per type.

[Learn more →](/basics/properties#shared-properties)

## Query

A structured filter-and-sort request against the index. Queries use filter rules (property + operator + value) and sort rules (property + direction) to retrieve matching Objects of a given Type. Exposed via the CLI (`tmd query`) and used internally by Views.

[Learn more →](/basics/queries)

## Search

A full-text search across Object names, descriptions, and body content. Unlike Queries which filter on structured properties, Search matches free-text input against the search index. Exposed via the CLI (`tmd search`) and the TUI search bar.

## Slug

The human-readable part of an Object's filename. For example, `golang-in-action` in `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md`. Derived from the name you provide when creating an Object.

## ULID

A 26-character Universally Unique Lexicographically Sortable Identifier, automatically appended to the slug by the CLI. Ensures uniqueness even when two Objects share the same slug. Objects created manually (without the CLI) do not require a ULID.

## Object ID

The full identifier for an Object in `type/slug-ulid` format. For example: `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz`. Used in Link targets and relation references.

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
