## 1. BDD: Write feature file and step definitions

- [ ] 1.1 Write BDD scenarios in `cmd/features/log.feature` covering: log with commits, prefix matching, object not found, --oneline flag, vault not in git repo, object with no commits
- [ ] 1.2 Implement step definitions in `cmd/bdd_steps_log_test.go`

## 2. CLI: Implement tmd log command

- [ ] 2.1 Create `cmd/log.go` with `logCmd` registered on `rootCmd`: resolve object ID, detect git repo, execute `git log --follow`, support `--oneline` flag, handle no-commits case
- [ ] 2.2 Add shell completion for object ID argument

## 3. Unit tests

- [ ] 3.1 Add unit tests in `cmd/log_test.go` for edge cases: git not installed, empty object ID argument validation
