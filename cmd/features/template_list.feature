Feature: Template listing
  As a user
  I want to list templates via CLI
  So that I can see what templates are available

  Background:
    Given a vault is ready

  Scenario: List all templates across types
    Given a template "book/review" with body "## Review"
    And a template "book/summary" with body "## Summary"
    And a template "note/meeting" with body "## Meeting"
    When I run template list
    Then the command should succeed
    And the output should contain "book/review"
    And the output should contain "book/summary"
    And the output should contain "note/meeting"

  Scenario: List templates filtered by type
    Given a template "book/review" with body "## Review"
    And a template "book/summary" with body "## Summary"
    And a template "note/meeting" with body "## Meeting"
    When I run template list with type "book"
    Then the command should succeed
    And the output should contain "book/review"
    And the output should contain "book/summary"
    And the output should not contain "note/meeting"

  Scenario: List templates for type with no templates
    When I run template list with type "book"
    Then the command should succeed
    And the output should be empty

  Scenario: List all templates when none exist
    When I run template list
    Then the command should succeed
    And the output should be empty

  Scenario: List templates as JSON
    Given a template "book/review" with body "## Review"
    When I run template list with json
    Then the command should succeed
    And the output should start with "["
    And the output should contain "book"
    And the output should contain "review"
