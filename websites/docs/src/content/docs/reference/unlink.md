---
title: tmd relation unlink
description: Remove a relation between two objects.
sidebar:
  order: 7
---

Removes a relation. Use `--both` to remove the inverse side as well.

Object IDs support prefix matching — you can omit the ULID suffix if the prefix uniquely identifies an object.

```bash
tmd relation unlink book/golang-in-action author person/alan-donovan --both
```

## Options

| Option | Description |
|--------|-------------|
| `--both` | Remove the inverse side of the relation as well |
