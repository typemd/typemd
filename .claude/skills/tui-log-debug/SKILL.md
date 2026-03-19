---
name: tui-log-debug
description: Use when debugging TUI layout, rendering, or interaction issues that cannot be diagnosed from code inspection alone. Use when you need to observe runtime values (widths, heights, state) in the terminal UI. Triggered by visual bugs in screenshots, panel sizing issues, or "it looks wrong but the code seems right" situations.
---

# TUI Log Debug

## Overview

The TUI writes debug logs to `.typemd/logs/{YYYY-MM-DD}.log`. Use `logf()` to instrument code, run the TUI, then read the log to observe runtime values.

**Core principle:** Don't guess layout math — log it and read the actual numbers.

## Infrastructure

The log system is in `tui/log.go`:

- `initLog(vaultRoot)` — called in `Start()`, creates `.typemd/logs/` and opens daily log file
- `logf(format, args...)` — write a formatted log line (no-op if logging not initialized)
- Log file: `.typemd/logs/2026-03-20.log` (date-based)

## When to Use

- Panel width/height calculations produce unexpected results
- Content overflows or doesn't fill a panel
- State transitions aren't working (e.g., rightPanel mode, focus)
- Resize behavior is broken
- Any TUI bug where "the code looks correct" but the visual output is wrong

## Workflow

1. **Add `logf()` calls** at the point you suspect is wrong:

```go
logf("preview split: m.width=%d previewW=%d tableW=%d vmWidth=%d",
    m.width, previewW, tableW, vm.width)
```

2. **Build and run** the TUI:

```bash
go build ./... && go run ./cmd/tmd
```

3. **Trigger the behavior** (resize, open preview, enter view mode, etc.)

4. **Read the log** (the vault might be in a subdirectory like `examples/book-vault/`):

```bash
cat .typemd/logs/2026-03-20.log
# or find it:
find . -path "*/.typemd/logs/*.log" -newer /tmp/marker
```

5. **Analyze values** — compare expected vs actual, look for off-by-one, double-subtraction, etc.

6. **Remove `logf()` calls** after fixing the bug. Don't commit debug logs.

## Tips

- Log **before and after** a calculation to see what changed
- Log **rendered widths** with `lipgloss.Width(rendered)` to verify actual output size
- For split layouts, log both panel widths and their sum vs target
- The log file location depends on which vault is opened — check `examples/` subdirectories
- Multiple renders may fire per frame — look at timestamps to distinguish sessions

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Forgetting to remove `logf` before commit | Search for `logf(` before committing |
| Looking for log in wrong vault directory | Use `find . -path "*/.typemd/logs/*.log"` |
| Not triggering the specific state | Make sure to reproduce the exact scenario (e.g., press `p` for preview) |
| Logging in `View()` without context | Include the method name or context in the log message |
