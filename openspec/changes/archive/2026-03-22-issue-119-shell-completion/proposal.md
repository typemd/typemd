## Why

Users must type or copy-paste full object IDs including ULID suffixes (e.g. `book/clean-code-01kk39c30x27ck7ahyc7ct4nyn`). Without tab completion, the CLI is tedious for any command that accepts object IDs. Shell completion would make the CLI practical for daily use.

## What Changes

- Add `ValidArgsFunction` to all commands that accept object IDs (`show`, `link`, `unlink`) with progressive two-stage completion: type prefix first, then object name
- Add `ValidArgsFunction` to commands that accept type names (`migrate`, `type show`) and `--type` flag completion (`stats`, `format`)
- Add context-aware relation name completion for `link`/`unlink` second argument
- Add `tmd completion` subcommand to generate shell-specific completion scripts for bash, zsh, and fish

## Capabilities

### New Capabilities

- `shell-completion`: Tab completion for object IDs, type names, relation names, and shell script generation via `tmd completion`

### Modified Capabilities

_(none — no existing spec requirements change)_

## Impact

- **Code**: `cmd/` package — new `completion.go` command, modifications to `show.go`, `link.go`, `unlink.go`, `migrate.go`, `type_show.go`, `stats.go`, `format.go`
- **Dependencies**: No new dependencies; uses Cobra's built-in completion framework
- **APIs**: No breaking changes; purely additive CLI feature
