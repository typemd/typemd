---
title: "TypeMD 0.4.0 — Laying the Foundation"
description: "v0.4.0 focuses on core formats, configuration, and creation UX: a built-in page type, vault health checks, config command, starter templates, and a completely redesigned TUI creation flow."
date: 2026-03-18
tags: [release, devlog]
---

If 0.3.0 was about "defining rules for your vault," 0.4.0 is about laying the foundation.

We didn't rush to add flashy new features. Instead, we went back and settled the fundamentals — stable data formats, persistent configuration, faster object creation. None of this is glamorous, but skip any one of them and daily use starts to feel rough.

## Built-in page type: just write

Not everything fits neatly into a `book` or `person`. Sometimes you just want to jot down a paragraph, an idea, a link.

`page` is a new built-in type (📄) that, like `tag`, exists without a YAML file. No custom properties — just the simplest possible free-form container:

```bash
tmd object create page "Weekend Trip Checklist"
```

Pair it with vault config, and you can make `page` the default type. After that, you don't even need to specify the type:

```bash
tmd config set cli.default_type page
tmd object create "Quick Note"
```

## tmd doctor: health check your vault

Vaults accumulate cruft over time — a type gets deleted but its directory lingers, objects reference types that no longer exist. Before, you'd have to dig through directories yourself. Now it's one command:

```bash
tmd doctor
```

It scans for orphan directories, missing types, and other integrity issues, presenting results in categorized reports.

## tmd config: settings that stick

Previously, many behaviors could only be controlled through flags — use them once, forget them next time. Now your vault has its own config file (`.typemd/config.yaml`), managed through `tmd config`:

```bash
tmd config set cli.default_type idea
tmd config get cli.default_type
tmd config list
```

`tmd init` also creates this file automatically, defaulting to `default_type: page`.

## Starter type templates: don't start from scratch

In 0.3.0, we removed built-in types and gave you a clean, empty vault. But for new users, a blank slate can be more overwhelming than a set of defaults.

Now `tmd init` asks if you'd like to add starter types:

- 💡 **idea** — with a status property (seed → growing → ready)
- 📝 **note** — the simplest note type
- 📚 **book** — full book tracking (status, rating, author relation)

Pick what you want, skip what you don't. Your vault, your call.

## New type schema fields

This release adds three practical fields to type schemas:

**version** — track schema versions in semver format (`"1.0"`), laying groundwork for automatic migrations:

```yaml
name: book
version: "1.0"
```

**color** — preset names or hex codes for visual distinction:

```yaml
name: book
color: blue
```

**description** — document types and properties inline, no comments needed:

```yaml
name: book
description: "Track reading lists and book reviews"
properties:
  - name: rating
    type: number
    description: "1-5 star rating"
```

## TUI creation experience, redesigned

The flows for creating objects and types have been completely rethought.

**Creating objects**: instead of a popup input box, you type the name directly in the title panel. As you type, the right panel shows a live template preview so you know exactly what you'll get before confirming. Names support natural language and auto-convert to valid slugs.

**Creating types**: a new multi-field wizard lets you Tab between emoji, name, and plural name fields. The right panel shows a live schema preview, and after creation the type editor opens automatically.

**Template management**: the type editor now includes a Templates section where you can browse, edit, create, and delete object templates. No need to leave the TUI for the file system.

## Upgrade

```bash
brew upgrade typemd/tap/typemd-cli
```

Or grab a pre-built binary from [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.4.0).

No breaking changes in this release — upgrade and go.

## What's next

The foundation is laid. Now we start building on it. The direction for 0.5.0 is still taking shape, but aliases, a richer search experience, and the Web UI are all on our radar.

If you have any questions or ideas, [open an issue](https://github.com/typemd/typemd/issues) — we read them all.
