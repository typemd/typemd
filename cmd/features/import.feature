Feature: Import command
  As a user
  I want to import external markdown files into my vault
  So that I can migrate existing notes into typemd

  Scenario: Scan a directory with markdown files
    Given a vault is ready
    And a source directory with markdown files
    When I run import scan on the source directory
    Then the command should succeed
    And the output should contain "file_count"
    And the output should contain "sources"

  Scenario: Scan reports an error for missing path
    Given a vault is ready
    When I run import scan on a nonexistent path
    Then the command should fail

  Scenario: Execute a valid plan file
    Given a vault is ready
    And a plan file with a page object
    When I run import execute with the plan file
    Then the command should succeed
    And the output should contain "objects_created"

  Scenario: Execute reports an error for missing plan file
    Given a vault is ready
    When I run import execute "nonexistent.json"
    Then the command should fail

  Scenario: Plan from classifications file
    Given a vault is ready
    And a classifications file with a page
    When I run import plan with the classifications file
    Then the command should succeed
    And the output should contain "objects"
    And the output should contain "order"
