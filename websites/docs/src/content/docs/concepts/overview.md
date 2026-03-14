---
title: Concepts Overview
description: Key terms and concepts in TypeMD.
sidebar:
  order: 1
---

TypeMD uses a small set of core concepts to organize your knowledge. This page gives a quick definition of each term. For deeper explanations, see the dedicated pages linked below.

## Vault

A Vault is a directory managed by TypeMD. It contains a `.typemd/` configuration folder, an `objects/` folder for your content, and an optional `templates/` folder for object templates. Each Vault is self-contained — you can move, copy, or version-control it like any folder.

## Object

An Object is the basic unit of knowledge in TypeMD. Instead of thinking in "files" or "notes", you think in Objects — a book, a person, an idea, a meeting.

Each Object is stored as a Markdown file with YAML frontmatter (properties) and a body for free-form content.

[Learn more about Objects →](/concepts/objects)

## Type

A Type defines what kind of Object you're creating. Every Object belongs to exactly one Type (e.g. `book`, `person`, `note`). Types are defined by schema files that specify which properties an Object can have.

[Learn more about Types →](/concepts/types)

## Relation

A Relation is a named, typed link between two Objects. Unlike a simple hyperlink, Relations are defined in the Type schema and can be bidirectional — updating one side automatically updates the other.

[Learn more about Relations →](/concepts/relations)

## Property

A Property is a named field on an Object, defined by its Type schema. Properties have types like `string`, `number`, `date`, `select`, or `relation`. They are stored in the YAML frontmatter of the Markdown file. Properties can also have optional attributes like `emoji` (visual icon), `pin` (prominent display), and `default` (default value).

## Slug

A Slug is the human-readable part of an Object's filename, e.g. `golang-in-action` in `book/golang-in-action-01jqr3k5mpbvn8e0f2g7h9txyz.md`. Slugs are derived from the name you give when creating an Object.

## ULID

A ULID (Universally Unique Lexicographically Sortable Identifier) is a 26-character ID automatically appended to the slug when you create an Object via the CLI. It ensures uniqueness even if two Objects share the same slug. Objects created manually (without the CLI) do not require a ULID.

## Wiki-link

A Wiki-link is an inline reference to another Object using `[[type/slug]]` syntax in the Markdown body. TypeMD tracks these links and their backlinks, so you can discover connections between Objects.

[Learn more about Wiki-links →](/concepts/wiki-links)

## Index

The Index is a SQLite database (`.typemd/index.db`) that caches Object metadata for fast querying and searching. It is automatically built and updated — you never need to edit it by hand.
