Feature: Object archive and unarchive commands
  Users can archive and unarchive objects via CLI.

  Scenario: Archive an object
    Given a vault with a "book" object named "archive-test"
    When I run tmd object archive on the object
    Then the output should contain "Archived"
    And the object should be archived

  Scenario: Archive an already archived object
    Given a vault with a "book" object named "already-archived"
    And the object is archived
    When I run tmd object archive on the object
    Then the output should contain "already archived"

  Scenario: Unarchive an archived object
    Given a vault with a "book" object named "unarchive-test"
    And the object is archived
    When I run tmd object unarchive on the object
    Then the output should contain "Unarchived"
    And the object should not be archived

  Scenario: Unarchive a non-archived object
    Given a vault with a "book" object named "not-archived"
    When I run tmd object unarchive on the object
    Then the output should contain "not archived"

  Scenario: List excludes archived objects by default
    Given a vault with a "book" object named "visible-book"
    And a book object "hidden-book" exists
    And the "hidden-book" object is archived
    When I run object list
    Then the output should contain "visible-book"
    And the output should not contain "hidden-book"

  Scenario: List with include-archived shows all objects
    Given a vault with a "book" object named "visible-book2"
    And a book object "hidden-book2" exists
    And the "hidden-book2" object is archived
    When I run object list with include-archived
    Then the output should contain "visible-book2"
    And the output should contain "hidden-book2"
