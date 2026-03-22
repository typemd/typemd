---
name: explore
description: Analyze existing markdown files and suggest typemd type schemas.
disable-model-invocation: true
---

# Explore

Analyze existing markdown files and suggest typemd type schemas and properties.

## Steps

1. **Load vault-guide** — run `/typemd:vault-guide` to get the full typemd reference
2. **Get explore instructions** — run `tmd instructions explore` to get JSON with instructions and vault context
3. **Follow the instructions** — parse the returned JSON and follow the `instructions` field exactly
