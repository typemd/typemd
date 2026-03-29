Feature: Watch validate
  The --watch flag on tmd type validate enters continuous validation mode.

  Scenario: Watch flag is accepted
    Given a vault is ready
    When I run watch validate briefly
    Then the watch command should have started successfully

  Scenario: Without watch flag runs once and exits
    Given a vault is ready
    When I run command "type validate"
    Then the command should succeed
    And the output should contain "Validation passed"

  Scenario: Validation errors are shown in watch mode
    Given a vault is ready
    And a broken schema exists
    When I run watch validate briefly
    Then the watch output should contain "Schema errors"
