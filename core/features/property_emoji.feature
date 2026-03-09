Feature: Property emoji
  Properties in type schemas can have an optional emoji for compact display.

  Scenario: Property with emoji defined
    Given a vault is ready
    And a type schema "item" with property "title" having emoji "📖"
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Property without emoji defined
    Given a vault is ready
    And a type schema "item" with a "title" string property
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Unique property emojis accepted
    Given a vault is ready
    And a type schema "item" with properties having unique emojis
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Duplicate property emojis rejected
    Given a vault is ready
    And a type schema "item" with properties having duplicate emojis
    When I validate all schemas
    Then schema "item" should have errors

  Scenario: Empty emojis do not conflict
    Given a vault is ready
    And a type schema "item" with some properties missing emojis
    When I validate all schemas
    Then schema "item" should have no errors
