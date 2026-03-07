---
title: tmd unlink
description: Remove a relation between two objects.
sidebar:
  order: 7
---

Removes a relation. Use `--both` to remove the inverse side as well.

```bash
tmd unlink book/golang-in-action-01jqr3k5mp... author person/alan-donovan-01jqr3k5mp... --both
```

## Options

| Option | Description |
|--------|-------------|
| `--both` | Remove the inverse side of the relation as well |
