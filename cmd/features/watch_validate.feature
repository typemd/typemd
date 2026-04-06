Feature: Validate type schemas
  Users validate type schemas to catch configuration errors early.

  Scenario: Validate schemas with no errors
    Given a vault is ready
    When I run command "type validate"
    Then the command should succeed
    And the output should contain "Validation passed"

  Scenario: Validation errors are reported in watch mode
    Given a vault is ready
    And a broken schema exists
    When I run watch validate briefly
    Then the watch output should contain "Schema errors"
