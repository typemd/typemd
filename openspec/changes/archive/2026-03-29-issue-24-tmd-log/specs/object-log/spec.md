## ADDED Requirements

### Requirement: View object commit history

The CLI SHALL provide a `tmd log <object-id>` command that displays the git commit history for a specific object file. The command SHALL resolve the object ID to its file path and execute `git log --follow` on that path.

#### Scenario: Show log for an object with commits

- **WHEN** the user runs `tmd log book/my-book` and the object exists with git history
- **THEN** the command displays the git log output for the object's file

#### Scenario: Show log with prefix matching

- **WHEN** the user runs `tmd log book/my` and the prefix uniquely matches one object
- **THEN** the command resolves the prefix and displays the git log for the matched object

#### Scenario: Object not found

- **WHEN** the user runs `tmd log book/nonexistent` and no object matches
- **THEN** the command exits with an error message indicating the object was not found

### Requirement: Compact output mode

The CLI SHALL support a `--oneline` flag that displays each commit on a single line (abbreviated hash and subject).

#### Scenario: Oneline output

- **WHEN** the user runs `tmd log --oneline book/my-book`
- **THEN** the command displays the git log in oneline format

### Requirement: Git repository detection

The CLI SHALL detect when the vault is not inside a git repository and display a clear error message.

#### Scenario: Vault not in a git repository

- **WHEN** the user runs `tmd log book/my-book` in a vault that is not inside a git repository
- **THEN** the command exits with an error: "vault is not inside a git repository"

### Requirement: Uncomitted object handling

The CLI SHALL handle objects that exist but have no git commits gracefully.

#### Scenario: Object with no commits

- **WHEN** the user runs `tmd log book/my-book` and the object file has not been committed to git
- **THEN** the command displays a message: "no commits found for this object"
