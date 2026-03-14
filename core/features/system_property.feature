Feature: System property registry
  typemd maintains a registry of system-managed properties that are
  automatically present on all objects regardless of type schema.

  Scenario: Registry contains all system properties in order
    Then the system property registry should contain "name, description, created_at, updated_at, tags"

  Scenario: IsSystemProperty recognizes system properties
    Then "name" should be a system property
    And "description" should be a system property
    And "created_at" should be a system property
    And "updated_at" should be a system property

  Scenario: IsSystemProperty rejects non-system properties
    Then "title" should not be a system property
    And "author" should not be a system property

  Scenario: Schema validation rejects created_at property
    Given a vault is ready
    And a type schema "bad" with a system property "created_at"
    When I validate all schemas
    Then schema "bad" should have errors

  Scenario: Schema validation rejects updated_at property
    Given a vault is ready
    And a type schema "bad" with a system property "updated_at"
    When I validate all schemas
    Then schema "bad" should have errors

  Scenario: Schema validation rejects description property
    Given a vault is ready
    And a type schema "bad" with a system property "description"
    When I validate all schemas
    Then schema "bad" should have errors

  Scenario: Shared property validation rejects description
    Given a vault is ready
    And a shared properties file with a system property "description"
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Shared property validation rejects created_at
    Given a vault is ready
    And a shared properties file with a system property "created_at"
    When I validate all schemas
    Then shared properties should have errors

  Scenario: New object does not have description
    Given a vault is ready
    When I create a "book" object named "no-desc-book"
    Then the object should not have property "description"

  Scenario: Object with description preserves it
    Given a vault is ready
    And a "book" object named "desc-book" exists
    When I set property "description" to "A great book" on the object
    Then the object property "description" should be "A great book"

  Scenario: SyncIndex preserves description
    Given a vault is ready
    And a raw object file with description exists
    When I sync the index
    Then the indexed properties for the object should contain "description"

  Scenario: Sync does not add description to existing objects
    Given a vault is ready
    And a raw object file without timestamps exists
    When I sync the index
    Then the raw object file should not have description added

  Scenario: New object has created_at and updated_at timestamps
    Given a vault is ready
    When I create a "book" object named "test-book"
    Then the object should have a "created_at" timestamp
    And the object should have an "updated_at" timestamp

  Scenario: created_at is not modified on save
    Given a vault is ready
    And a "book" object named "test-book" exists
    When I save the object
    Then the object "created_at" should not have changed

  Scenario: SaveObject updates updated_at
    Given a vault is ready
    And a "book" object named "test-book" exists
    When I save the object
    Then the object "updated_at" should be recent

  Scenario: SetProperty updates updated_at
    Given a vault is ready
    And a "book" object named "test-book" exists
    When I set property "title" to "Test" on the object
    Then the object "updated_at" should be recent

  Scenario: Frontmatter orders system properties first
    Given a vault is ready
    When I create a "book" object named "ordered-book"
    Then the frontmatter should have system properties before schema properties

  Scenario: Frontmatter orders description between name and created_at
    Given a vault is ready
    And a "book" object named "ordered-desc-book" exists
    When I set property "description" to "A test description" on the object
    Then the frontmatter should have "description" before "created_at"
    And the frontmatter should have "name" before "description"

  Scenario: SyncIndex preserves system properties
    Given a vault is ready
    And a "book" object named "sync-book" exists
    When I sync the index
    Then the indexed properties for the object should contain "created_at"
    And the indexed properties for the object should contain "updated_at"

  Scenario: Existing object without timestamps loads successfully
    Given a vault is ready
    And a raw object file without timestamps exists
    When I sync the index
    Then the raw object file should not have timestamps added

  Scenario: Schema validation rejects tags property
    Given a vault is ready
    And a type schema "bad" with a system property "tags"
    When I validate all schemas
    Then schema "bad" should have errors

  Scenario: Shared property validation rejects tags
    Given a vault is ready
    And a shared properties file with a system property "tags"
    When I validate all schemas
    Then shared properties should have errors

  Scenario: Immutable system properties are identified
    Then "created_at" should be an immutable system property
    And "updated_at" should be an immutable system property

  Scenario: Mutable system properties are not immutable
    Then "name" should not be an immutable system property
    And "description" should not be an immutable system property
    And "tags" should not be an immutable system property

  Scenario: Non-system properties are not immutable
    Then "title" should not be an immutable system property
