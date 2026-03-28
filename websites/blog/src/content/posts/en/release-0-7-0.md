---
title: "TypeMD 0.7.0 — Edit Everything Inline"
description: "v0.7.0 brings inline editing for every property type — no more leaving the TUI to change your data."
date: 2026-03-28
tags: [release, devlog]
---

0.6.0 brought AI into your vault. But even after AI generated descriptions and tags for you, you still had to open an editor to change property values. This time we tore down that last wall — **edit everything directly in the TUI**.

## Properties? Just Press Enter

Select a property, press Enter, edit it. No opening files, no remembering YAML syntax.

Every property type is supported:

- **string / text** — inline text input
- **number** — numeric input
- **select** — option list
- **checkbox** — toggle switch
- **url** — URL input
- **date / datetime** — segmented input with calendar picker

Press Enter to save, Esc to cancel. That's it.

## Dates Get Their Own Calendar Picker

Date properties aren't just text fields anymore — we built a full date picker:

- Segmented input (year/month/day), Tab moves between segments
- Inline calendar popup, navigate with arrow keys
- Supports `date_format` and `datetime_format` settings for custom display:

```yaml
# .typemd/config.yaml
date_format: "2006/01/02"
datetime_format: "2006/01/02 15:04"
```

## Relation Picking: Fuzzy Search, Done

When editing a relation property, a fuzzy-search picker pops up. Type a few characters, find your target object, press Enter. Multi-value relations let you keep adding.

Link a book's author to a person? Two seconds.

## Table View Is Now Editable

Last release added table view. This release makes it editable. Select a cell in the table, press Enter, and use the exact same editing widgets as the properties panel.

Browse and modify data in bulk — table view is now a real workspace.

## Lock Objects to Prevent Accidents

Don't want important objects to be accidentally changed? Add one line to frontmatter:

```yaml
locked: true
```

The TUI shows a lock icon and blocks all editing. Unlock by setting `locked` back to `false`.

## Local Properties No Longer Disappear

If you added properties to frontmatter that weren't defined in the schema, they used to be silently ignored during sync. Now these "local properties" show up in a separate section and are fully preserved during sync.

Add fields freely, no schema changes needed.

## Bring Your Own LLM

0.6.0's AI features were tied to Claude CLI. Now you can connect any OpenAI-compatible provider — Ollama, LM Studio, vLLM, you name it:

```yaml
# .typemd/config.yaml
ai:
  enabled: true
  default: ollama
  providers:
    ollama:
      type: openai-compatible
      base_url: http://localhost:11434
      model: llama3.2
```

Want to switch providers? Change one line in `default`. Existing `ai.enabled` / `ai.model` configs auto-migrate.

## What Else

- **`tmd serve`** — start an HTTP server with `tmd serve` to browse your vault in the browser, with a built-in three-theme system (warm/dark/light)
- **Structured Logging** — all packages now use `slog`; TUI logs go to `.typemd/logs/`, CLI uses `--debug` for output
- **Focus Mode** — press `.` to toggle single full-width body panel for distraction-free writing

## What's Next

The editing experience is in place. Next up in 0.8.0, we're polishing the developer experience — template CLI, MCP extensions, and markdown syntax highlighting.

To upgrade:

```bash
brew upgrade typemd/tap/typemd-cli
```

Or grab pre-built binaries from [GitHub Releases](https://github.com/typemd/typemd/releases/tag/v0.7.0).
