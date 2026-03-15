## 1. Core: Embedded Starter Types

- [x] 1.1 Write BDD scenarios for starter type listing and YAML validity (`core/features/`)
- [x] 1.2 Create starter YAML files: `core/starters/idea.yaml`, `note.yaml`, `book.yaml`
- [x] 1.3 Implement `core/starters.go` with `//go:embed`, `StarterType` struct, and `StarterTypes()` function
- [x] 1.4 Implement BDD step definitions to make scenarios pass
- [x] 1.5 Add unit tests for edge cases (embed integrity, schema validation of each starter)

## 2. Core: WriteStarterTypes Method

- [x] 2.1 Write BDD scenarios for writing starter types to vault (`core/features/`)
- [x] 2.2 Implement `Vault.WriteStarterTypes(names []string) error` in `core/vault.go`
- [x] 2.3 Implement BDD step definitions to make scenarios pass
- [x] 2.4 Add unit tests for edge cases (empty names, duplicate names, already-existing file)

## 3. Cmd: Bubble Tea Starter Picker

- [x] 3.1 Implement `cmd/starter_picker.go` with Bubble Tea checkbox model (Init/Update/View)
- [x] 3.2 Add unit tests for picker model (key handling, toggle, select all, deselect all, confirm, quit)

## 4. Cmd: Init Command Integration

- [x] 4.1 Add `--no-starters` flag to `cmd/init.go`
- [x] 4.2 Integrate starter picker into init flow: show picker → call `WriteStarterTypes` → print results
- [x] 4.3 Add unit tests for init flag parsing and non-interactive mode
