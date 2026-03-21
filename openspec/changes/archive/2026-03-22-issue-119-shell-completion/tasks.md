## 1. Core: Completion Helper Functions

- [x] 1.1 Write BDD scenarios for object ID completion (two-stage: type prefix, object name)
- [x] 1.2 Write BDD scenarios for type name completion
- [x] 1.3 Write BDD scenarios for relation name completion
- [x] 1.4 Create `cmd/completion.go` with helper functions: `completeObjectID`, `completeTypeName`, `completeRelationName`
- [x] 1.5 Implement BDD step definitions and make scenarios pass
- [x] 1.6 Add unit tests for edge cases (empty vault, invalid prefix, vault path flag)

## 2. Wire Completion to Commands

- [x] 2.1 Add `ValidArgsFunction` to `show` command (object ID completion)
- [x] 2.2 Add `ValidArgsFunction` to `link` command (positional: object ID / relation name / object ID)
- [x] 2.3 Add `ValidArgsFunction` to `unlink` command (same as link)
- [x] 2.4 Add `ValidArgsFunction` to `migrate` command (type name completion)
- [x] 2.5 Add `ValidArgsFunction` to `type show` command (type name completion)
- [x] 2.6 Register `--type` flag completion for `stats` and `format` commands

## 3. Completion Subcommand

- [x] 3.1 Enable Cobra's built-in completion command on rootCmd
- [x] 3.2 Add unit test verifying `tmd completion bash/zsh/fish` produces output

## 4. Verification

- [x] 4.1 Run full test suite (`make test`) and verify all pass
- [ ] 4.2 Manual smoke test: build binary and test completion in bash/zsh/fish
