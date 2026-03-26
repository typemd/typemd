Feature: tmd format command
  The format command normalizes object frontmatter and schema YAML files.

  Background:
    Given a vault is ready

  Scenario: Already formatted vault shows no-op message
    Given the vault formatting is stable
    When I run format
    Then the command should succeed
    And the output should contain "All files are already formatted."

  Scenario: Formatting files shows count
    Given an object with out-of-order frontmatter
    When I run format
    Then the command should succeed
    And the output should contain "Formatted"
    And the output should contain "file(s)."

  Scenario: Format with type filter
    When I run format with type "book"
    Then the command should succeed

  Scenario: Format with invalid type
    When I run format with type "nonexistent"
    Then the command should fail

  Scenario: Dry-run lists files needing format
    When I run format with dry-run
    Then the command should fail
