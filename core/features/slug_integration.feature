Feature: Object creation with slug conversion
  As a user
  I want to create objects with natural-language names
  So that filenames are valid slugs while display names preserve my input

  Scenario: Natural-language name is slugified for filename
    Given a vault is ready
    When I create a "book" object named "Some Great Thought"
    Then no error should occur
    And the current object ID should contain "some-great-thought"
    And the current object name property should be "Some Great Thought"

  Scenario: Pre-slugified name passes through unchanged
    Given a vault is ready
    When I create a "book" object named "already-slugified"
    Then no error should occur
    And the current object ID should contain "already-slugified"
    And the current object name property should be "already-slugified"

  Scenario: Name with special characters is slugified
    Given a vault is ready
    When I create a "book" object named "What's Next?"
    Then no error should occur
    And the current object ID should contain "whats-next"
    And the current object name property should be "What's Next?"
