Feature: Object relations
  Objects can be linked through typed relations defined in type schemas.
  Relations support bidirectional linking with inverse properties.

  Background:
    Given a vault is ready with relation schemas

  Scenario: Link objects with bidirectional relation
    Given a "book" object named "golang-in-action" exists
    And a "person" object named "alan-donovan" exists
    When I link "golang-in-action" to "alan-donovan" via "author"
    Then the "author" property of "golang-in-action" should reference "alan-donovan"
    And the "books" property of "alan-donovan" should contain "golang-in-action"

  Scenario: Link with type mismatch fails
    Given a "book" object named "book-a" exists
    And a "book" object named "book-b" exists
    When I link "book-a" to "book-b" via "author"
    Then an error should occur

  Scenario: Overwrite single-value relation
    Given a "book" object named "test" exists
    And a "person" object named "alan" exists
    And a "person" object named "brian" exists
    When I link "test" to "alan" via "author"
    And I link "test" to "brian" via "author"
    Then the "author" property of "test" should reference "brian"

  Scenario: Unlink both directions
    Given a "book" object named "test" exists
    And a "person" object named "alan" exists
    And I link "test" to "alan" via "author"
    When I unlink "test" from "alan" via "author" with both flag
    Then the "author" property of "test" should be empty
    And the "books" property of "alan" should be empty

  Scenario: List relations for a linked object
    Given a "book" object named "test" exists
    And a "person" object named "alan" exists
    When I link "test" to "alan" via "author"
    Then listing relations for "test" should return 2 entries

  Scenario: List relations for an unlinked object
    Given a "book" object named "test" exists
    Then listing relations for "test" should return 0 entries
