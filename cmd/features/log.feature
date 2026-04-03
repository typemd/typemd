Feature: tmd log command
  The log command shows git commit history for a specific object.

  Background:
    Given a vault is ready

  Scenario: Show log for an object with commits
    Given a book object "Clean Code" exists
    And the vault is a git repository with committed objects
    When I run log "book/clean-code"
    Then the command should succeed
    And the output should contain "add objects"

  Scenario: Show log with oneline flag
    Given a book object "Clean Code" exists
    And the vault is a git repository with committed objects
    When I run log "book/clean-code" with oneline
    Then the command should succeed

  Scenario: Object not found
    When I run log "book/nonexistent"
    Then the command should fail

  Scenario: Vault not in a git repository
    Given a book object "Clean Code" exists
    When I run log with the first created object
    Then the command should fail with "not inside a git repository"

  Scenario: Object with no commits
    Given a book object "Clean Code" exists
    And the vault is a git repository
    When I run log with the first created object
    Then the command should succeed
    And the output should contain "no commits found"
