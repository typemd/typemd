## 1. BDD: CLI disambiguation scenarios

- [x] 1.1 Write BDD feature file `cmd/features/disambiguate.feature` covering: interactive picker display, user selection, user cancel, non-interactive fallback, and all commands (show, link, unlink)
- [x] 1.2 Implement BDD step definitions in `cmd/bdd_steps_disambiguate_test.go` (initially failing)

## 2. Disambiguation picker UI

- [x] 2.1 Create `cmd/disambiguate_picker.go` with Bubble Tea model: items (name + ID), cursor navigation (↑/↓/j/k), Enter to select, Esc/q to cancel, styled View
- [x] 2.2 Add unit tests for picker model in `cmd/disambiguate_picker_test.go` (Init, Update key handling, View rendering, selected/cancelled state)

## 3. Shared resolve helpers

- [x] 3.1 Add `resolveIDInteractive(vault, prefix) (string, error)` to `cmd/helper.go` — wraps `vault.ResolveID()`, catches `AmbiguousMatchError`, detects TTY, launches picker or falls back to error
- [x] 3.2 Add `resolveObjectInteractive(vault, prefix) (*core.Object, error)` to `cmd/helper.go` — wraps `vault.ResolveObject()` with same disambiguation logic
- [x] 3.3 Add unit tests for resolve helpers in `cmd/helper_test.go` (mock non-interactive path)

## 4. Integrate into CLI commands

- [x] 4.1 Update `cmd/show.go` to use `resolveObjectInteractive` instead of `vault.ResolveObject`
- [x] 4.2 Update `cmd/link.go` to use `resolveIDInteractive` for both from-id and to-id
- [x] 4.3 Update `cmd/unlink.go` to use `resolveIDInteractive` for both from-id and to-id

## 5. Verify and polish

- [x] 5.1 Run full test suite (`make test`) and ensure all BDD scenarios pass
- [x] 5.2 Manual testing: create duplicate-prefix objects in a test vault, verify picker appears and works for show/link/unlink
