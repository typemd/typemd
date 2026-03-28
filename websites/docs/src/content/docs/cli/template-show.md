---
title: tmd template show
description: Display a template's content.
sidebar:
  order: 9.7
---

Displays a template's frontmatter properties and body content.

```bash
tmd template show book/review
```

Example output:

```
book/review

Properties
──────────
  status: draft

Body
────
  ## Review Notes
  Write your review here.
```

The argument must be in `type/name` format (e.g., `book/review`). Properties are shown sorted alphabetically. If the template has no properties, `(none)` is displayed. If the body is empty, `(empty)` is displayed.
