Feature: tmd stats command
  The stats command shows aggregate statistics about the vault.

  Background:
    Given a vault is ready

  Scenario: Vault-wide stats shows type summary
    Given a book object "Clean Code" exists
    When I run stats
    Then the command should succeed
    And the output should contain "books"
    And the output should contain "Total"

  Scenario: Vault-wide stats on empty vault
    When I run stats
    Then the command should succeed
    And the output should contain "No objects in vault."

  Scenario: Per-type stats shows property aggregations
    Given a book object "Clean Code" exists
    When I run stats with type "book"
    Then the command should succeed
    And the output should contain "book"
    And the output should contain "objects"

  Scenario: Per-type stats with invalid type
    When I run stats with type "nonexistent"
    Then the command should fail

  Scenario: JSON output
    Given a book object "Clean Code" exists
    When I run stats with json
    Then the command should succeed
    And the output should start with "{"
