## 1. BDD: Watch Validate Scenarios

- [x] 1.1 Write BDD feature file for watch-validate (`cmd/features/watch_validate.feature`) covering: flag acceptance, initial validation run, re-validation on file change, graceful exit
- [x] 1.2 Implement BDD step definitions for watch-validate scenarios

## 2. Core: Watch Loop Implementation

- [x] 2.1 Add `--watch` / `-w` bool flag to `validateCmd` in `cmd/validate.go`
- [x] 2.2 Implement the watch loop: fsnotify watcher on types dir, properties file, and objects dir with 200ms debounce
- [x] 2.3 Implement terminal clear + timestamp display + validation re-run on each cycle
- [x] 2.4 Implement vault re-sync (`vault.Sync()`) before each validation run
- [x] 2.5 Implement graceful shutdown via `signal.NotifyContext` (SIGINT/SIGTERM)
- [x] 2.6 Handle missing watched directories gracefully (skip without crashing)

## 3. Unit Tests

- [x] 3.1 Add unit tests for edge cases: missing directories, rapid changes, signal handling
