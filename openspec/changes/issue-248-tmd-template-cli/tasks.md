## 1. BDD: Template CLI Scenarios

- [x] 1.1 Write BDD feature file for `tmd template list` (all types, filtered by type, empty, JSON output)
- [x] 1.2 Write BDD feature file for `tmd template show` (with properties+body, body only, nonexistent, invalid format)
- [x] 1.3 Write BDD feature file for `tmd template create` (new template, already exists, creates directory, invalid format)
- [x] 1.4 Write BDD feature file for `tmd template delete` (with force, nonexistent, invalid format)

## 2. BDD: Step Definitions

- [x] 2.1 Implement step definitions for template list scenarios
- [x] 2.2 Implement step definitions for template show scenarios
- [x] 2.3 Implement step definitions for template create scenarios
- [x] 2.4 Implement step definitions for template delete scenarios

## 3. Implementation: Template Command Group

- [x] 3.1 Create `cmd/template.go` with parent `tmd template` command and `parseTemplateArg` helper
- [x] 3.2 Implement `tmd template list` subcommand with optional type filter and `--json` flag
- [x] 3.3 Implement `tmd template show` subcommand with properties and body display
- [x] 3.4 Implement `tmd template create` subcommand with editor integration
- [x] 3.5 Implement `tmd template delete` subcommand with confirmation prompt and `--force` flag
- [x] 3.6 Add shell completion for template commands (type names, type/name pairs)

## 4. Unit Tests: Edge Cases

- [x] 4.1 Unit test `parseTemplateArg` for valid/invalid formats
- [x] 4.2 Unit test editor resolution (`$EDITOR` → `$VISUAL` → `vi` fallback)
