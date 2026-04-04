---
title: "TypeMD 0.2.0 — Property Types, Shared Properties, and a Better TUI"
description: "v0.2.0 brings 9 formally defined property types, a shared properties mechanism, and a new TUI experience with emoji, pinned fields, and session persistence."
date: 2026-03-11
tags: [release, devlog]
---

After the first release of TypeMD, the next thing we wanted to tackle was: "what property types are actually supported?" In v0.1.0, the `type` field existed on properties — but was never formally defined. You could write `type: number` and things would mostly work, but there was no real contract behind it.

0.2.0 changes that. The core of this release is making TypeMD's type system something you can actually rely on.

## Nine property types, formally defined

Type schemas now support 9 explicit property types:

| Type | Description |
|------|-------------|
| `string` | Short text |
| `number` | Numeric values |
| `date` | Calendar dates |
| `datetime` | Date and time |
| `url` | Web links |
| `checkbox` | Boolean (yes/no) |
| `select` | Single choice |
| `multi_select` | Multiple choices |
| `relation` | Link to objects of another type |

Here's what a book type looks like:

```yaml
name: book
emoji: 📚
properties:
  - name: status
    type: select
    options:
      - value: to-read
      - value: reading
      - value: done
  - name: rating
    type: number
  - name: published
    type: date
  - name: author
    type: relation
    target: person
    bidirectional: true
    inverse: books
```

With explicit types, validation, migration, and display all get smarter.

## Stop repeating yourself: shared properties

If multiple types share the same properties — say, every resource needs `tags` and `favorite` — you used to have to define them in each type schema separately. Not anymore.

Define them once in `properties/properties.yaml`:

```yaml
properties:
  - name: tags
    type: multi_select
    emoji: 🔖
    options:
      - value: programming
      - value: design
  - name: favorite
    type: checkbox
    emoji: ❤️
```

Then reference them in any type via `use`:

```yaml
properties:
  - name: rating
    type: number
  - use: tags
  - use: favorite
```

Change it in one place, and it applies everywhere.

## A TUI that looks like your vault

### Emoji (if you want it)

Both types and properties now support an `emoji` field:

```yaml
name: book
emoji: 📚
properties:
  - name: rating
    type: number
    emoji: ⭐
```

Type emoji show up in object list group headers; property emoji appear in the detail panel. Your vault starts to feel like *your* vault.

### A dedicated title bar

When viewing an object, the TUI now shows a title bar at the top with the type emoji and object name. A small change, but it makes a big difference in orientation.

### Pin your most important properties

Mark a property with a `pin` value and it floats to the top of the detail panel in that order. No more scrolling past a wall of fields to find the one you care about:

```yaml
  - name: status
    type: select
    pin: 1
  - name: rating
    type: number
    pin: 2
```

### Pick up where you left off

Close the TUI and reopen it — cursor position, selected object, and panel state are all restored. The TUI starts feeling less like a viewer and more like a workspace.

## Other improvements

- **`name` is now a system property** — automatically populated from the object slug. You don't need to (and can't) define it manually in type schemas anymore.
- **`--readonly` flag** — launch the TUI without editing enabled. Handy when you just want to browse without accidentally changing things.
- **`--reindex` flag** — replaces the old `tmd reindex` subcommand. Pass it at startup to force a full index rebuild.
- **Prefix matching** — when referencing objects by ID, a short prefix of the ULID is enough. No need for the full ID.
- **Homebrew** — `brew install typemd/tap/typemd-cli`. The easiest way to install on macOS.

## Upgrade

```bash
brew upgrade typemd/tap/typemd-cli
```

Or grab a pre-built binary from [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.2.0).

## What's next

0.3.0 will focus on the Web UI — a browser-based way to browse your vault without the terminal. Same local-first, file-based foundation, with a graphical interface on top.

If you have any questions or ideas, [open an issue](https://github.com/typemd/typemd/issues) — we read them all.
