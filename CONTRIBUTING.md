# Contributing to typemd

Thank you for your interest in contributing to typemd!

## Reporting Bugs

Please [open an issue](https://github.com/typemd/typemd/issues/new) and include:

- What you did (steps to reproduce)
- What you expected to happen
- What actually happened
- Go version (`go version`) and OS

For questions or ideas, use [GitHub Discussions](https://github.com/typemd/typemd/discussions) instead of the issue tracker.

## Submitting a Pull Request

1. Fork the repository and create a branch from `main`
2. Make your changes and add tests
3. Ensure all tests pass: `go test ./...`
4. Ensure the build is clean: `go build ./...`
5. Push to your fork and open a pull request

Please link the related issue in your PR description (e.g. `Closes #123`).

## Commit Message Convention

This project follows [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <description>
```

Examples:

```
feat(core): add wiki-links and backlinks support
fix(tui): use DisplayID in detail panel title
test(core): add BDD tests for core business logic
refactor(core,cmd): unify naming, extract helpers, improve error handling
```

Common types: `feat`, `fix`, `refactor`, `test`, `docs`, `chore`

## Prerequisites

- Go 1.25+
- SQLite (system built-in is fine)

## Build

```bash
go build ./...
```

## Testing

This project uses two layers of testing:

- **BDD** ([Godog](https://github.com/cucumber/godog), Cucumber for Go) — Define behaviors and establish shared vocabulary. Gherkin scenarios describe **what** a feature does from the user's perspective, not implementation details.
- **Unit tests** — Verify precise logic: edge cases, output formats, exact values, error conditions.

When deciding where a test belongs: if it defines a behavior or names a concept, write a BDD scenario. If it validates an implementation detail (e.g. JSON output format, lowercase ULID, flag edge cases), write a unit test.

### BDD scope by package

| Package | Testing approach |
|---------|-----------------|
| `core/` | BDD (`core/features/`) + unit tests |
| `tui/`  | BDD (`tui/features/`, planned) + unit tests |
| `web/`  | BDD (`web/features/`, future) |
| `cmd/`  | BDD (`cmd/features/`) + unit tests |
| `mcp/`  | Unit tests — BDD TBD |

### Run all tests

```bash
go test ./...
```

### Run tests by package

```bash
go test ./core/...   # Core logic (BDD + unit)
go test ./tui/...    # Terminal UI
go test ./mcp/...    # MCP server
```

### Verbose output

```bash
go test -v ./...
```

### Writing BDD Tests

All new behavior should be described as Gherkin scenarios in `core/features/*.feature`:

```gherkin
Feature: Vault initialization and lifecycle
  A vault stores objects as Markdown files and manages a SQLite index.

  Scenario: Initialize a new vault
    When I initialize a new vault
    Then the vault directory structure should exist
    And the SQLite database should exist

  Scenario: Double initialization fails
    Given a vault is initialized
    When I initialize the vault again
    Then an error should occur
```

**Steps to add a new test:**

1. Write or edit a `.feature` file in `core/features/` — describe the behavior in plain language first
2. Run `go test ./core/... -run TestFeatures` — Godog will report undefined steps with suggested snippets
3. Implement the step definitions in `core/bdd_steps_test.go`
4. Run again to verify all scenarios pass

**Guidelines for writing scenarios:**

- Use the `Given` / `When` / `Then` structure to separate preconditions, actions, and assertions
- Write scenarios from the user's perspective, not the implementation's
- Each scenario should define one behavior — keep them focused
- Focus on **what** the system does, not **how** — avoid asserting implementation details like exact output formats, specific error messages, or internal data structures
- Reuse existing step definitions when possible (check `core/bdd_steps_test.go`)

### Project structure

```
core/
├── features/           # Gherkin feature files
│   ├── vault.feature
│   ├── object.feature
│   ├── relation.feature
│   ├── query.feature
│   ├── validate.feature
│   ├── wikilink.feature
│   └── frontmatter.feature
├── bdd_test.go         # Test entry point (TestFeatures)
└── bdd_steps_test.go   # Step definitions
```

### Running BDD tests only

```bash
go test ./core/... -run TestFeatures
```

### Unit tests

Unit tests use traditional Go `testing` style and handle precise validation:

- Exact output formats (e.g. `--json` flag output)
- Edge cases and boundary conditions (e.g. empty input, invalid values)
- Implementation-specific logic (e.g. ULID generation, frontmatter parsing details)

Use `t.TempDir()` for test isolation. Test files are named `*_test.go` alongside the source files.
