Feature: Object listing
  As a user
  I want to list objects via CLI
  So that I can see what is in my vault

  Scenario: List objects in plain format
    Given a vault with 2 books
    When I run object list
    Then the command should succeed
    And the output should have 2 lines containing "book/"

  Scenario: List objects in JSON format
    Given a vault with 2 books
    When I run object list with json
    Then the command should succeed
    And the output should start with "["

  Scenario: List objects in empty vault
    Given a vault is ready
    When I run object list
    Then the command should succeed
    And the output should be empty
