---
name: importer
description: Convert existing markdown files into typemd objects.
disable-model-invocation: true
---

# Importer

Convert existing markdown files into typemd objects.

## Steps

1. **Load vault-guide** — run `/typemd:vault-guide` to get the full typemd reference
2. **Get importer instructions** — run `tmd instructions importer` to get JSON with instructions and vault context
3. **Follow the instructions** — parse the returned JSON and follow the `instructions` field exactly
