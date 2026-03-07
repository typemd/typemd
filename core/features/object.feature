Feature: Object management
  Objects are Markdown files with YAML frontmatter, identified by type/name-ULID.

  Scenario: Create a new object
    Given a vault is ready
    When I create a "book" object named "golang-in-action"
    Then the object filename should start with "golang-in-action-"
    And the object filename should have a 26-character ULID suffix
    And the object type should be "book"
    And the object file should exist on disk

  Scenario: Same name produces different ULIDs
    Given a vault is ready
    When I create a "book" object named "test"
    And I create another "book" object named "test"
    Then the two objects should have different IDs

  Scenario: Creating object with unknown type fails
    Given a vault is ready
    When I create a "nonexistent" object named "test"
    Then an error should occur

  Scenario: Get object by ID
    Given a vault is ready
    And a "book" object named "golang-in-action" exists
    When I get the object by its ID
    Then the retrieved object should match the created one

  Scenario: Set and persist a property
    Given a vault is ready
    And a "book" object named "test" exists
    When I set property "title" to "Go in Action" on the object
    Then the object property "title" should be "Go in Action"

  Scenario: Property validation rejects wrong type
    Given a vault is ready
    And a "book" object named "test" exists
    When I set property "rating" to "not-a-number" on the object
    Then an error should occur

  Scenario: Save object persists to file and database
    Given a vault is ready
    And a "book" object named "test-book" exists
    When I update the object body to "New body content"
    And I update the object title to "Updated Title"
    And I save the object
    Then the object file should contain "Updated Title"
    And the object file should contain "New body content"
    And getting the object by ID should return body "New body content"
