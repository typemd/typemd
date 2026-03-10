Feature: Pinned properties
  Properties in type schemas can have an optional pin value for prominent display.

  Scenario: Property with pin defined
    Given a vault is ready
    And a type schema "item" with property "status" having pin 1
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Property without pin defined
    Given a vault is ready
    And a type schema "item" with a "title" string property
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Unique pin values accepted
    Given a vault is ready
    And a type schema "item" with properties having unique pins
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Duplicate pin values rejected
    Given a vault is ready
    And a type schema "item" with properties having duplicate pins
    When I validate all schemas
    Then schema "item" should have errors

  Scenario: Unpinned properties do not conflict
    Given a vault is ready
    And a type schema "item" with some properties unpinned
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Negative pin value rejected
    Given a vault is ready
    And a type schema "item" with property "status" having pin -1
    When I validate all schemas
    Then schema "item" should have errors
