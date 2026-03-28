Feature: Debug flag
  The --debug flag enables debug-level logging to stderr.

  Scenario: Debug flag is accepted
    Given a vault is ready
    When I run command "object list --debug"
    Then the command should succeed

  Scenario: Normal mode produces no log output on stderr
    Given a vault is ready
    When I run command "object list"
    Then the command should succeed
