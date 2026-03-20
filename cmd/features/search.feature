Feature: Full-text search
  As a user
  I want to search objects via CLI
  So that I can find objects by keyword

  Scenario: Search finds matching objects
    Given a vault with a book "Clean Code"
    When I run search "clean"
    Then the command should succeed
    And the output should contain "Clean Code"

  Scenario: Search with no results
    Given a vault with a book "Clean Code"
    When I run search "zzzznotfound"
    Then the command should succeed
    And the output should be empty

  Scenario: Search with JSON output
    Given a vault with a book "Clean Code"
    When I run search "clean" with json
    Then the command should succeed
    And the output should start with "["

  Scenario: Search requires a keyword argument
    Given a vault is ready
    When I run search with no arguments
    Then the command should fail
