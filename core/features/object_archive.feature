Feature: Object archive (soft delete)
  Users can archive objects to hide them from default queries without deleting files.

  Scenario: IsArchived returns true for archived object
    Given a vault is ready
    When I create a "book" object named "archived-book"
    And I archive the object
    Then the object should be archived

  Scenario: IsArchived returns false for non-archived object
    Given a vault is ready
    When I create a "book" object named "normal-book"
    Then the object should not be archived

  Scenario: SetArchived archives an object
    Given a vault is ready
    And a "book" object named "set-archived-book" exists
    When I archive the object
    Then the object should be archived
    And the object frontmatter should contain "archived: true"

  Scenario: SetArchived unarchives an object
    Given a vault is ready
    And a "book" object named "set-unarchive-book" exists
    And I archive the object
    When I unarchive the object
    Then the object should not be archived
    And the object frontmatter should not contain "archived"

  Scenario: Archive a locked object succeeds
    Given a vault is ready
    And a "book" object named "locked-archive-book" exists
    And I lock the object
    When I archive the object
    Then no error should occur
    And the object should be archived

  Scenario: Default query excludes archived objects
    Given a vault is ready
    And a "book" object named "visible-book" exists
    And a "book" object named "hidden-book" exists
    And I archive the "hidden-book" object
    When I query objects with filter "type=book"
    Then the query should return 1 result

  Scenario: Query with include-archived returns all objects
    Given a vault is ready
    And a "book" object named "visible-book2" exists
    And a "book" object named "hidden-book2" exists
    And I archive the "hidden-book2" object
    When I query objects with filter "type=book" including archived
    Then the query should return 2 results

  Scenario: GetObject returns archived objects
    Given a vault is ready
    And a "book" object named "get-archived-book" exists
    And I archive the object
    When I get the object by ID
    Then no error should occur
    And the object should be archived
