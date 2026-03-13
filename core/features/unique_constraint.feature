Feature: Unique constraint on type schemas
  Types with unique: true enforce that no two objects share the same name.

  Background:
    Given a vault is ready

  # --- Schema parsing ---

  Scenario: Type schema with unique true
    Given a type schema "person" with unique constraint
    When I load the type schema "person"
    Then the loaded schema should have unique true

  Scenario: Type schema without unique field defaults to false
    When I load the type schema "book"
    Then the loaded schema should have unique false

  # --- Creation-time enforcement ---

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
