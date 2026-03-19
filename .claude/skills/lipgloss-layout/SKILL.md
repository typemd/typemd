---
name: lipgloss-layout
description: Use when calculating panel widths, heights, or split layouts in the TUI using lipgloss v2. Use when panels don't fill the screen, content overflows borders, or split panels have gaps. Triggered by "Width() semantics", "border width calculation", "panel doesn't fill", "content overflows panel".
---

# Lipgloss v2 Layout Calculation

## Overview

In lipgloss v2 (charm.land/lipgloss/v2), `Width()` and `Height()` set the **total rendered size including border**. This is the single most important fact for all layout math.

## The Key Rule

```
Width(N) → rendered output is exactly N columns wide (border included)
         → content area is N - border_size wide
```

For `RoundedBorder()`: border = 2 (1 left + 1 right). So:

```
Width(50) → rendered = 50, content = 48
```

## Quick Reference

| What you want | Formula |
|---------------|---------|
| Single panel fills `m.width` | `Width(m.width - bdr)` where `bdr = 2` |
| Two panels fill `m.width` | `Width(panelA)` + `Width(panelB)` where `panelA + panelB = m.width - bdr` |
| Content width for `vm.SetSize` | `panelWidth - bdr` (subtract border from the Width value) |
| `lipgloss.Width(rendered)` | Returns the Width value (total including border) |

## Single Panel (Full Width)

The title bar and full-width body both use this pattern:

```go
bdr := 2

// Title: fills entire terminal width
titleStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    Width(m.width - bdr)
// rendered = m.width - bdr + bdr... NO!
// Width IS the total. rendered = m.width - bdr.
```

**Wait — why `m.width - bdr` and not `m.width`?**

Because `JoinVertical` and terminal rendering add implicit constraints. In practice, `Width(m.width - bdr)` fills the screen. This was determined empirically:
- `Width(m.width)` → overflows (wraps to next line)
- `Width(m.width - bdr)` → fills correctly

This suggests there's a 2-char overhead from the terminal or Bubble Tea framework. Treat `m.width - bdr` as the usable "full width".

## Two-Panel Split

Both panels must sum to the same usable width:

```go
bdr := 2
totalContent := m.width - bdr  // usable full width (same as title)

previewW := totalContent / 2
tableW := totalContent - previewW
// previewW + tableW = totalContent ✓

tableStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    Width(tableW)

previewStyle := lipgloss.NewStyle().
    Border(lipgloss.RoundedBorder()).
    Width(previewW)
```

**Content size for sub-models:**

```go
// vm.SetSize wants content width (what to render inside the border)
vm.SetSize(tableW - bdr, bodyH - bdr)
vm.previewWidth = previewW - bdr
```

## Verification with logf

When debugging layout, log the key values:

```go
logf("split: m.width=%d totalContent=%d tableW=%d previewW=%d sum=%d",
    m.width, totalContent, tableW, previewW, tableW+previewW)

rendered := lipgloss.JoinHorizontal(lipgloss.Top, tableRendered, previewRendered)
logf("rendered: table=%d preview=%d joined=%d",
    lipgloss.Width(tableRendered), lipgloss.Width(previewRendered),
    lipgloss.Width(rendered))
```

Expected: `joined = totalContent = m.width - bdr`.

## Common Mistakes

| Mistake | Symptom | Fix |
|---------|---------|-----|
| `Width(m.width)` for full-width panel | Content overflows terminal | Use `Width(m.width - bdr)` |
| Two panels: `panelA + panelB = m.width` | Panels overflow | Sum must equal `m.width - bdr` |
| Two panels: `panelA + panelB = m.width - 2*bdr` | Gap on right side | Don't double-count border; sum = `m.width - bdr` |
| `vm.SetSize(panelWidth, ...)` | Table rows overflow panel | Pass content width: `vm.SetSize(panelWidth - bdr, ...)` |
| Mixing content/total width in calculation | Off-by-2 or off-by-4 errors | Be explicit: `totalContent` for Width(), `contentWidth` for SetSize() |

## The Three-Number Pattern

For any panel, track three numbers:

```
totalContent = m.width - bdr          // what Width() gets for full-width
panelWidth   = share of totalContent  // what Width() gets for split panel
contentWidth = panelWidth - bdr       // what vm.SetSize() gets
```

If you keep these three straight, layout math works.
