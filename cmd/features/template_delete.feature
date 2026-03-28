Feature: Template delete
  As a user
  I want to delete templates via CLI
  So that I can remove templates I no longer need

  Background:
    Given a vault is ready

  Scenario: Delete template with force flag
    Given a template "book/review" with body "## Review"
    When I run template delete "book/review" with force
    Then the command should succeed
    And the output should contain "Deleted template book/review"
    And the template file "book/review" should not exist

  Scenario: Delete nonexistent template
    When I run template delete "book/nonexistent" with force
    Then the command should fail

  Scenario: Delete template with invalid argument format
    When I run template delete "review" with force
    Then the command should fail with "type/name"
