Feature: Instructions command
  As a user
  I want to output skill instructions from the CLI
  So that I can use them with LLM integrations

  Scenario: List available skills
    When I run instructions with no arguments
    Then the command should succeed
    And the output should contain "explore"
    And the output should contain "importer"
    And the output should contain "importer"

  Scenario: List skills as JSON
    When I run instructions with json flag
    Then the command should succeed
    And the output should contain "name"
    And the output should contain "description"

  Scenario: Get skill with vault context
    Given a vault is ready
    When I run instructions "explore"
    Then the command should succeed
    And the output should contain "instructions"
    And the output should contain "context"

  Scenario: Get raw skill content
    When I run instructions "explore" with skill flag
    Then the command should succeed
    And the output should contain "---"
    And the output should contain "name: explore"

  Scenario: Unknown skill returns error
    When I run instructions "nonexistent"
    Then the command should fail with "unknown skill"
