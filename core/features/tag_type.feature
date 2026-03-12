Feature: Tag built-in type
  Tags are first-class objects used for categorization.
  The tag type has color and icon properties.

  Scenario: Tag type exists as a built-in type
    Given a vault is ready
    When I create a "tag" object named "go"
    Then the object type should be "tag"
    And the object file should exist on disk

  Scenario: Tag object has color and icon properties
    Given a vault is ready
    When I create a "tag" object named "go"
    Then the object should have property "color" with nil value
    And the object should have property "icon" with nil value

  Scenario: Note type does not have tags property
    Given a vault is ready
    When I create a "note" object named "test-note"
    Then the object should not have property "tags"
