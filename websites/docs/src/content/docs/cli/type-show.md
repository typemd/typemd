---
title: tmd type show
description: Show details of a type schema.
sidebar:
  order: 10
---

Displays the details of a type schema, including its name, plural form (if defined), and all properties and their configurations.

```bash
tmd type show book
```

Example output:

```
Type: 📚 book
Plural: books

Properties
──────────
  title (string)
  status (select) [to-read, reading, done]
  rating (number)
  author (relation) -> person (bidirectional) inverse=books
```

When a type schema defines a `plural` field, it will be displayed after the type name. If the field is not defined, this line is omitted.
