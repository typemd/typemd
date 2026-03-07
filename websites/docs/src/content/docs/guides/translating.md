---
title: Translating
description: How to contribute translations to TypeMD documentation.
sidebar:
  order: 5
---

TypeMD documentation supports multiple languages using [Starlight's i18n system](https://starlight.astro.build/guides/i18n/).

## Supported Languages

| Language | Code | Directory |
|----------|------|-----------|
| English (default) | `en` | `docs/src/content/docs/` |
| Traditional Chinese | `zh-TW` | `docs/src/content/docs/zh-tw/` |

## Adding a Translation

1. Find the English source file under `docs/src/content/docs/`
2. Create the corresponding file under `docs/src/content/docs/zh-tw/` with the same path
3. Translate the content, keeping code blocks, commands, and technical terms in English
4. Update the `title` and `description` in the YAML frontmatter

## Adding a New Language

1. Add the locale to `docs/astro.config.mjs` under `locales`
2. Add sidebar translations under each sidebar group's `translations` field
3. Create the directory under `docs/src/content/docs/<locale>/`
4. Mirror the English file structure with translated content
5. Create a `README.<LOCALE>.md` at the project root

## Translation Guidelines

- **Keep in English:** TypeMD, Object, Type, Relation, CLI commands, property names, code blocks, YAML keys, file paths
- **Translate:** Prose descriptions, section headings, table descriptions, frontmatter `title` and `description`
- **Spacing:** Add a space between full-width CJK characters and half-width alphanumeric characters
