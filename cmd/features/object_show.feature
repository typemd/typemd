Feature: Object display
  As a user
  I want to show object details via CLI
  So that I can inspect properties and body content

  Scenario: Show object detail
    Given a vault with a book "Clean Code"
    When I run object show for the created book
    Then the command should succeed
    And the output should contain "Properties"
    And the output should contain "Body"

  Scenario: Show object with empty body
    Given a vault with a book "Empty Book"
    When I run object show for the created book
    Then the command should succeed
    And the output should contain "(empty)"

  Scenario: Show non-existent object
    Given a vault is ready
    When I run object show "book/nonexistent"
    Then the command should fail
