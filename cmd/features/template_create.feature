Feature: Template create
  As a user
  I want to create templates via CLI
  So that I can scaffold new template files

  Background:
    Given a vault is ready

  Scenario: Create new template
    When I run template create "book/review"
    Then the command should succeed
    And the output should contain "Created template book/review"
    And the template file "book/review" should exist

  Scenario: Create template when file already exists
    Given a template "book/review" with body "## Existing"
    When I run template create "book/review"
    Then the command should fail with "already exists"

  Scenario: Create template creates type directory
    When I run template create "newtype/first"
    Then the command should succeed
    And the template file "newtype/first" should exist

  Scenario: Create template with invalid argument format
    When I run template create "review"
    Then the command should fail with "type/name"
