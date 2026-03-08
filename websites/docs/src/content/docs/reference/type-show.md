---
title: tmd type show
description: Show details of a type schema.
sidebar:
  order: 9
---

Displays the details of a type schema, including its properties and their configurations.

```bash
tmd type show book
```

Example output:

```
Type: 📚 book

Properties
──────────
  title: string
  status: enum (reading, read, to-read)
  rating: number
  author: relation → person
```
