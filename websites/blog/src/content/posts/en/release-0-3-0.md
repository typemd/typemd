---
title: "TypeMD 0.3.0 — Your Objects, Your Rules"
description: "v0.3.0 brings object templates, name templates, uniqueness constraints, system properties, and a visual type schema editor right in the TUI."
date: 2026-03-14
tags: [release, devlog]
---

If 0.2.0 was about "defining types for your properties," 0.3.0 is about "defining rules for your vault."

Last release, you could describe what fields a book has and what types they are. This release, you can tell TypeMD: "book names can't repeat," "journals should be named by date automatically," "every new book review starts with this template." You write the rules, TypeMD enforces them.

## Object templates: start every object in the right place

Tired of manually setting `status: draft` and pasting the same Markdown structure every time you create a book review? That's over.

Drop a Markdown file in `templates/<type>/` with your desired frontmatter defaults and body skeleton:

```yaml
# templates/book/review.md
---
status: to-read
---

## Why I Want to Read This

## Notes
```

It gets applied automatically when creating an object:

```bash
tmd object create book "Clean Code" --template review
```

If a type has only one template, you don't even need `--template` — it applies automatically. Multiple templates? You'll be prompted to choose.

Templates can override user-editable system properties (`name`, `description`), but `created_at` and `updated_at` are system-managed — template values for those are ignored.

## Name templates: let TypeMD name things for you

Some objects have predictable names. A daily journal is just "Journal 2026-03-14." A meeting note is "Meeting 2026-03-14." Typing that out every time is tedious and error-prone.

Now you can define a name template in your type schema:

```yaml
name: journal
properties:
  - name: name
    template: "Journal {{ date:YYYY-MM-DD }}"
```

Create an object without providing a name, and TypeMD fills it in. You can always specify one manually to override.

## Unique constraint: prevent duplicates

Add `unique: true` to a type schema, and no two objects of that type can share the same name:

```yaml
name: person
unique: true
```

The first `person/john-doe` is fine. The second gets rejected. Different types are independent — `person/john-doe` and `character/john-doe` can coexist.

`tmd type validate` also gained tag name uniqueness checking. Tags are inherently unique, and the validator catches duplicates for you.

## Plural display names

Type schemas now support a `plural` field:

```yaml
name: book
plural: books
```

Wherever a collection label is needed (like TUI group headers), TypeMD uses `books` instead of `book`. A small touch, but it reads better every time.

## System properties: things TypeMD manages for you

Starting with 0.3.0, every object automatically has five system properties:

| Property | Description | Editable? |
|----------|-------------|-----------|
| `name` | Display name | Yes |
| `description` | Description text | Yes |
| `created_at` | Creation time | No (set automatically) |
| `updated_at` | Last update time | No (updated automatically) |
| `tags` | Tags (linked to built-in `tag` type) | Yes |

These properties don't need to be — and can't be — defined in type schemas. They're reserved names. If your schema manually defines `description`, `created_at`, `updated_at`, or `tags`, validation will fail after upgrading. Remove them first.

Tags are a brand new built-in type. Reference a non-existent tag in an object's frontmatter, and the tag object gets created automatically during sync.

## Goodbye, built-in types

In 0.2.0, new vaults came with `book`, `person`, and `note` pre-installed. We found that these defaults confused new users more than they helped — "Are these required? Can I delete them?"

The answer is: **your vault, your types.** In 0.3.0, all default types are removed. `tmd init` now creates a clean, empty vault. The only built-in type that remains is `tag`, since it's the foundation of the tagging system.

If you're already using `book`, `person`, or `note` — don't worry. Your objects and type schemas are unaffected. New vaults just won't come with them anymore.

## TUI type editor

Previously, editing type schemas meant opening YAML files by hand. Now, navigate to any type header in the TUI sidebar and you'll enter a visual type editor:

- **Browse** the property list with type, emoji, pin, and other metadata at a glance
- **Add** properties with full configuration — name, type, emoji, pin, and more
- **Edit** any field of existing properties
- **Delete** properties (with confirmation)
- **Reorder** properties by moving them up and down

No leaving the TUI. No remembering YAML syntax.

## Upgrade

```bash
brew upgrade typemd/tap/typemd-cli
```

Or grab a pre-built binary from [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.3.0).

Before upgrading:
1. Remove any manually defined `description`, `created_at`, `updated_at`, or `tags` properties from your type schemas
2. Run `tmd type validate` to confirm everything is clean

## What's next

0.4.0 is about "moving in" — vault health checks (`tmd doctor`), continuous validation (`tmd type validate --watch`), and a better TUI creation experience. Making your vault comfortable for daily use.

If you have any questions or ideas, [open an issue](https://github.com/typemd/typemd/issues) — we read them all.
