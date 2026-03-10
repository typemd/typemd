Feature: Name property
  All objects have a required "name" property that serves as the primary display title.
  The name is stored in YAML frontmatter and accessed via GetName().

  Scenario: GetName returns name property value
    Given a vault is ready
    And a "book" object named "clean-code" exists
    When I set the object name to "Clean Code"
    Then GetName should return "Clean Code"

  Scenario: GetName falls back to DisplayName when name is missing
    Given a vault is ready
    And a "book" object named "clean-code" exists
    When I remove the name property from the object
    Then GetName should return the DisplayName

  Scenario: GetName falls back to DisplayName when name is empty
    Given a vault is ready
    And a "book" object named "clean-code" exists
    When I set the object name to ""
    Then GetName should return the DisplayName

  Scenario: New object has name populated from slug
    Given a vault is ready
    When I create a "book" object named "golang-in-action"
    Then the object property "name" should be "golang-in-action"

  Scenario: Sync adds name to existing object without one
    Given a vault is ready
    And a "book" object named "old-book" exists
    When I remove the name property from the object
    And I save the object
    And I sync the index
    Then the synced object should have name matching its DisplayName

  Scenario: Sync preserves existing name
    Given a vault is ready
    And a "book" object named "my-book" exists
    When I set the object name to "My Awesome Book"
    And I save the object
    And I sync the index
    Then the synced object should have name "My Awesome Book"
