Feature: Template show
  As a user
  I want to view a template's content via CLI
  So that I can inspect its properties and body

  Background:
    Given a vault is ready

  Scenario: Show template with properties and body
    Given a template "book/review" with property "status" value "draft" and body "## Notes"
    When I run template show "book/review"
    Then the command should succeed
    And the output should contain "book/review"
    And the output should contain "status: draft"
    And the output should contain "## Notes"

  Scenario: Show template with body only
    Given a template "book/simple" with body "## My Template"
    When I run template show "book/simple"
    Then the command should succeed
    And the output should contain "(none)"
    And the output should contain "## My Template"

  Scenario: Show nonexistent template
    When I run template show "book/nonexistent"
    Then the command should fail

  Scenario: Show template with invalid argument format
    When I run template show "review"
    Then the command should fail with "type/name"
