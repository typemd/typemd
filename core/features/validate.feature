Feature: Validation
  Schemas, objects, relations, and wiki-links can be validated for correctness.

  Scenario: Valid schema passes validation
    Given a vault is ready
    And a type schema "book" with a "title" string property
    When I validate all schemas
    Then schema "book" should have no errors

  Scenario: Invalid schema is detected
    Given a vault is ready
    And a type schema "bad" with an enum property missing values
    When I validate all schemas
    Then schema "bad" should have errors

  Scenario: Orphaned relation target is detected
    Given a vault is ready
    And an orphaned relation from "book/test-book" to "person/ghost" exists
    When I validate relations
    Then there should be 1 relation error

  Scenario: Valid wiki-links pass validation
    Given a vault is ready
    And two linked notes exist
    When I validate wiki-links
    Then there should be no wiki-link errors

  Scenario: Broken wiki-link is detected
    Given a vault is ready
    And a note with a broken wiki-link exists
    When I validate wiki-links
    Then there should be 1 wiki-link error
    And the error should mention "broken wiki-link"
