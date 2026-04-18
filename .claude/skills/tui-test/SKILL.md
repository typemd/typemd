---
name: tui-test
description: Use when a typemd code change needs end-to-end TUI verification — keybinding behaviour, layout rendering, config overrides, error recovery — and "open TTY failed" or "non-interactive sandbox" is about to be treated as a blocker. Trigger on symptoms like "can't test TUI without a terminal", "need to verify what the user sees", "press ctrl+X and check X appears", "help popup should show Y", "sidebar should look like Z after config change".
---

# TUI Manual Testing via tmux PTY

## Overview

Automated `go test` covers logic. It does NOT confirm what the user sees. The TUI is Bubble Tea on top of a terminal — layout math, dispatch routing, toast timing, colour output can all go wrong in ways unit tests miss.

`tmux` gives the agent a real PTY. With it you can launch `tmd`, press keys, `capture-pane` the visible screen, and assert on what the user would actually see. **This is the default; there is no "we can't test this in a sandbox" excuse.**

## The Rule

```
If you changed TUI dispatch, layout, a popup, a toast, or a config-driven visual,
you verify it with a tmux session BEFORE claiming the work is done.
```

Unit tests proving `helpEntries(km, false)` returns the right slice ≠ proof that the help popup actually renders on screen with the right key.

## When to Use

Use when any of these symptoms appear in the change:

- Keybinding added, removed, or remapped
- New popup, toast, help overlay, or status bar content
- Layout change (new panel mode, new width math, new title row)
- Config option added that affects TUI rendering (theme, debounce, keybindings, stats layout, etc.)
- A bug report of the form "pressing X does Y when it should do Z"
- Error state handling (what the user sees when something is mis-configured)

**Do NOT use when:**
- The change is purely internal (no user-visible rendering change)
- CI / hooks are the goal (use `go test` — don't script tmux for things unit tests can assert)

## Quick Reference

| Task | Command |
|------|---------|
| Check tmux available | `command -v tmux` |
| Build binary | `go build -o /tmp/tmd-test ./cmd/tmd` |
| Init vault non-interactively | `<binary> --vault <dir> init --no-starters` |
| Launch TUI in tmux | `tmux new-session -d -s <name> -x 140 -y 40 "<binary> --vault <dir>; sleep 10"` |
| Send a key | `tmux send-keys -t <name> C-d` (control-D), `"?"` (literal), `"Tab"`, `"Enter"`, `"Escape"` |
| Capture visible screen | `tmux capture-pane -t <name> -p` |
| Kill the session | `tmux kill-session -t <name>` |

## Standard Flow

1. **Build once** — `go build -o /tmp/tmd-test ./cmd/tmd`. A fresh binary avoids cache surprises.
2. **Create a disposable vault** — `mktemp -d` or `/tmp/tmd-<issue>-vault`; init with `--no-starters` to avoid the interactive starter picker blocking.
3. **Write the config** — cat-heredoc into `<vault>/.typemd/config.yaml`.
4. **Launch in tmux** — `-x 140 -y 40` gives the TUI a realistic terminal size; `sleep 10` keeps the session alive after `tmd` exits so you can still capture.
5. **Wait for first paint** — `sleep 2` after launch. Bubble Tea's initial render needs a frame.
6. **Send keys → capture → assert** in a loop (see patterns below).
7. **Kill the session** — every test must end with `tmux kill-session -t <name> 2>/dev/null` even on failure.
8. **Clean up** — delete the vault dir and the binary when finished.

## The Four Test Categories

### 1. Keybinding / dispatch

Goal: prove a key press triggers the expected handler.

**Always test TWO directions for rebind / remap changes:**

1. **Positive**: the NEW key triggers the action.
2. **Negative**: the OLD default key NO LONGER triggers it.

Skipping the negative case is how I shipped a Stats-rebind bug to PR #394: the new key worked, but the old key *also* still worked — a rebind acted as an add, not a replace. Unit tests missed it; so did my positive-only manual test.

```bash
BIN=/tmp/tmd-test
V=$(mktemp -d)
$BIN --vault "$V" init --no-starters >/dev/null
cat > "$V/.typemd/config.yaml" <<EOF
tui:
  keybindings:
    stats: "ctrl+d"
EOF

# Positive: new key triggers stats
tmux new-session -d -s pos -x 140 -y 40 "$BIN --vault $V; sleep 10"
sleep 2
tmux send-keys -t pos C-d
sleep 1
tmux capture-pane -t pos -p | grep -q "Vault Statistics" && echo "PASS: C-d opens stats" || echo "FAIL"
tmux kill-session -t pos

# Negative: old default key no longer triggers stats
tmux new-session -d -s neg -x 140 -y 40 "$BIN --vault $V; sleep 10"
sleep 2
tmux send-keys -t neg C-s
sleep 1
if tmux capture-pane -t neg -p | grep -q "Vault Statistics"; then
  echo "FAIL: C-s still opens stats — rebind not working as replace"
else
  echo "PASS: C-s no longer opens stats"
fi
tmux kill-session -t neg
```

Assert on a *stable string* that only appears in the target state ("Vault Statistics" is a panel title, not accidentally present elsewhere). **Use a fresh tmux session per assertion** — the previous one has the stats panel open, which taints the next capture.

### 2. Layout / render

Goal: prove the rendered screen matches expectation after a resize or panel toggle.

```bash
tmux new-session -d -s layout -x 140 -y 40 "$BIN --vault $V; sleep 10"
sleep 2
tmux send-keys -t layout "p"            # toggle props panel
sleep 1
# Count vertical rule characters — three panels should show two dividers
tmux capture-pane -t layout -p | head -1 | grep -o "│" | wc -l
tmux kill-session -t layout
```

When asserting on layout, prefer **structural signals** (border chars, emoji from a type row) over pixel counts. Terminals wrap; widths vary.

### 3. Config-override behaviour

Goal: prove a config value actually changes what the user sees.

```bash
# Run the same scenario twice: once with default config, once with override.
# Capture help popup in both cases. Diff.

for CONFIG in "" "stats: ctrl+d"; do
  cat > "$V/.typemd/config.yaml" <<EOF
tui:
  keybindings:
    $CONFIG
EOF
  tmux new-session -d -s cfgtest -x 140 -y 40 "$BIN --vault $V; sleep 10"
  sleep 2
  tmux send-keys -t cfgtest "?"
  sleep 1
  tmux capture-pane -t cfgtest -p | grep -E "stats|search" > /tmp/help-${CONFIG:-default}.txt
  tmux kill-session -t cfgtest
done
diff /tmp/help-default.txt /tmp/help-stats:*.txt
```

### 4. Error-state recovery

Goal: prove the TUI does NOT crash on bad input, and does surface a warning.

```bash
# Inject an invalid key into config. Startup must succeed, toast must fire.
cat > "$V/.typemd/config.yaml" <<EOF
tui:
  keybindings:
    stats: "crtl+s"              # typo
EOF

tmux new-session -d -s errtest -x 140 -y 40 "$BIN --vault $V; sleep 10"
sleep 2
SCREEN=$(tmux capture-pane -t errtest -p)
echo "$SCREEN" | grep -q "Invalid key" && echo "PASS: warning toast" || echo "FAIL: no toast"
echo "$SCREEN" | grep -q "pages\|sources\|tags" && echo "PASS: TUI still rendered" || echo "FAIL: TUI crashed"
tmux kill-session -t errtest
```

## Key-string Reference (tmux send-keys)

| Key | Literal for send-keys |
|-----|-----------------------|
| `a`-`z`, `0`-`9`, punctuation | `"a"`, `"/"`, `","` (quoted) |
| Capital letter (shifted) | `"A"` |
| Control | `C-d`, `C-c`, `C-e`, `C-s` |
| Tab | `Tab` |
| Enter | `Enter` (prefer this over `"\r"`) |
| Escape | `Escape` |
| Space | `Space` |
| Arrow keys | `Up`, `Down`, `Left`, `Right` |
| Function key | `F1`-`F12` |

Multiple keys in one call: `tmux send-keys -t s "?" Space "q"`. **Don't chain a string of letters** with `send-keys -t s "help"` expecting it as one action — tmux will send h-e-l-p as individual key events, which can trigger unintended shortcuts.

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Running `go run ./cmd/tmd` directly, seeing "could not open TTY", giving up | Always wrap in `tmux new-session` — this IS the PTY |
| Forgetting `--no-starters` on `init`, hanging forever | Init is interactive by default; always pass the flag in tests |
| `sleep 1` after launch — popup not painted yet | Use `sleep 2` after launch, `sleep 1` after a key press |
| Tiny session `-x 80 -y 24` makes layout collapse | Use `-x 140 -y 40` to reflect a realistic terminal |
| Session dies before `capture-pane` runs | Always append `; sleep 10` to the launched command |
| Leaving tmux sessions around after failure | Always `tmux kill-session -t <name> 2>/dev/null` in cleanup, even on error |
| Asserting on decorative characters that differ across terminals (╭, │, ↓) | Assert on stable strings (panel titles, type names, config-visible keys) |
| Comparing full `capture-pane` output — noisy | `grep -E` for the one line that matters, or filter then diff |

## Red Flags — STOP

When you catch yourself writing or thinking any of these, you're about to skip the verification:

- "The environment doesn't have a TTY, so this can't be tested" — **it does; use tmux**
- "Unit tests cover this" — unit tests don't prove rendering
- "The user can verify this locally" — that's the exact task they delegated to you
- "I'll just trust the code review" — rendering bugs pass code review all the time
- "tmux isn't installed" — check first: `command -v tmux`. If genuinely absent, escalate; otherwise proceed.
- "teatest / expect would be better" — those are fine; tmux works today. Don't punt on "better".
- "The positive case works, that's enough" — a rebind / remap / toggle change has TWO sides. Assert the OLD behaviour is gone, not just that the NEW one works.

## Reporting

After a tmux test, always report:

1. **What you sent** (the keys)
2. **What you captured** (relevant lines from `capture-pane`, not the full screen)
3. **Pass or fail** for each assertion
4. **Cleanup confirmation** (session killed, temp dir deleted)

## When It Genuinely Can't Work

Rare cases where tmux still won't suffice:

- Rendering tests that depend on true colour / font choices (use visual regression tools the user sets up)
- Timing-sensitive animations shorter than ~500ms (Bubble Tea re-renders faster than `capture-pane` can read)

In those cases, state the specific limitation AND still run the tmux session to verify everything that CAN be checked.
