Feature: Object creation
  As a user
  I want to create objects via CLI
  So that I can add entries to my vault

  Scenario: Create object with type and name
    Given a vault is ready
    When I run object create "book" "Clean Code"
    Then the command should succeed
    And the output should contain "Created book/"
    And the output should contain "clean-code"

  Scenario: Create object with default type
    Given a vault is ready
    When I run object create "My Page"
    Then the command should succeed
    And the output should contain "Created page/"

  Scenario: Create object without type or default
    Given a vault with no default type
    When I run object create with no arguments
    Then the command should fail with "type is required"

  Scenario: Create object with template
    Given a vault is ready
    And a vault with a "review" template for "book"
    When I run object create "book" "My Book" with template "review"
    Then the command should succeed
    And the output should contain "Created book/"
    When I run object show for the created book
    Then the output should contain "Template Author"
    And the output should contain "review template body"
