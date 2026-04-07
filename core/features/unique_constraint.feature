Feature: Unique constraint on type schemas
  Types with unique: true enforce that no two objects share the same name.

  Background:
    Given a vault is ready

  Scenario: First object with a name succeeds on unique type
    Given a type schema "person" with unique constraint
    When I create a "person" object named "john-doe"
    Then no error should occur

  Scenario: Duplicate name rejected on unique type
    Given a type schema "person" with unique constraint
    And a "person" object named "john-doe" exists
    When I create a "person" object named "john-doe"
    Then an error should occur

  Scenario: Same name allowed for different unique types
    Given a type schema "person" with unique constraint
    Given a type schema "character" with unique constraint
    And a "person" object named "john-doe" exists
    When I create a "character" object named "john-doe"
    Then no error should occur

  Scenario: Duplicate name allowed on non-unique type
    Given a "book" object named "clean-code" exists
    When I create another "book" object named "clean-code"
    Then no error should occur

  Scenario: Different case names are not duplicates
    Given a type schema "person" with unique constraint
    And a "person" object named "john-doe" exists
    When I create a "person" object named "John-Doe"
    Then no error should occur

  # --- Validation ---

  Scenario: Validation passes with no duplicates on unique type
    Given a type schema "person" with unique constraint
    And a "person" object named "alice" exists
    And a "person" object named "bob" exists
    When I validate name uniqueness
    Then there should be no name uniqueness errors

  Scenario: Validation reports duplicates on unique type
    Given a type schema "person" with unique constraint
    And a "person" object named "john-doe" exists
    And a raw duplicate object of type "person" named "john-doe" exists
    When I validate name uniqueness
    Then there should be name uniqueness errors

  Scenario: Validation skips non-unique types
    Given a "book" object named "clean-code" exists
    And a raw duplicate object of type "book" named "clean-code" exists
    When I validate name uniqueness
    Then there should be no name uniqueness errors
