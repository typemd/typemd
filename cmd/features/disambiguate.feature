Feature: Interactive disambiguation
  As a user
  I want to select the right object from ambiguous prefix matches
  So that I can continue my command without re-typing

  Background:
    Given a vault with two books sharing prefix "clean-code"

  Scenario: Ambiguous prefix fails in non-interactive mode
    When I run object show "book/clean-code" in non-interactive mode
    Then the command should fail with "ambiguous object ID"
    And the command should fail with "matches 2 objects"

  Scenario: Disambiguation applies to show command
    When I run object show "book/clean-code" selecting candidate 1
    Then the command should succeed
    And the output should contain "Properties"

  Scenario: Disambiguation applies to link command
    Given a vault with a person "Robert Martin"
    When I run relation link with ambiguous from-id "book/clean-code" selecting candidate 1
    Then the command should succeed
    And the output should contain "Linked"

  Scenario: Disambiguation applies to unlink command
    Given a vault with a person "Robert Martin"
    And a link from the first book to the person via "author"
    When I run relation unlink with ambiguous from-id "book/clean-code" selecting candidate 1
    Then the command should succeed
    And the output should contain "Unlinked"

  Scenario: User cancels disambiguation picker
    When I run object show "book/clean-code" and cancel the picker
    Then the command should fail with "ambiguous object ID"
