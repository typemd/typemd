Feature: Type statistics
  Users can view per-property aggregate statistics for a specific type.

  Scenario: Number property aggregation
    Given a vault is ready
    And a "book" object named "book1" exists with typed property "rating" set to "3"
    And a "book" object named "book2" exists with typed property "rating" set to "4"
    And a "book" object named "book3" exists with typed property "rating" set to "5"
    When I request type stats for "book"
    Then the type stats should show count 3
    And the type stats property "rating" should have type "number"
    And the type stats property "rating" number avg should be 4.0
    And the type stats property "rating" number min should be 3.0
    And the type stats property "rating" number max should be 5.0
    And the type stats property "rating" number sum should be 12.0
    And the type stats property "rating" should show filled 3

  Scenario: Select property distribution
    Given a vault is ready
    And a "book" object named "book1" exists with property "status" set to "reading"
    And a "book" object named "book2" exists with property "status" set to "done"
    And a "book" object named "book3" exists with property "status" set to "done"
    When I request type stats for "book"
    Then the type stats property "status" should have type "select"
    And the type stats property "status" select "reading" should have count 1
    And the type stats property "status" select "done" should have count 2

  Scenario: Checkbox property ratio
    Given a vault is ready
    And a type "task" with a "completed" checkbox property exists
    And a "task" object named "task1" exists with typed property "completed" set to "true"
    And a "task" object named "task2" exists with typed property "completed" set to "true"
    And a "task" object named "task3" exists with typed property "completed" set to "false"
    When I request type stats for "task"
    Then the type stats property "completed" should have type "checkbox"
    And the type stats property "completed" checkbox true count should be 2
    And the type stats property "completed" checkbox false count should be 1

  Scenario: Date property range
    Given a vault is ready
    And a type "event" with a "date" date property exists
    And a "event" object named "event1" exists with property "date" set to "2024-01-15"
    And a "event" object named "event2" exists with property "date" set to "2024-06-20"
    And a "event" object named "event3" exists with property "date" set to "2024-03-10"
    When I request type stats for "event"
    Then the type stats property "date" should have type "date"
    And the type stats property "date" date earliest should be "2024-01-15"
    And the type stats property "date" date latest should be "2024-06-20"

  Scenario: Relation property count
    Given a vault is ready
    And a type "book" with an "author" relation to "person" exists
    And a "person" object named "alice" exists
    And a "person" object named "bob" exists
    And a "book" object named "book1" exists
    And I link "book1" to "alice" via "author"
    And a "book" object named "book2" exists
    And I link "book2" to "bob" via "author"
    When I request type stats for "book"
    Then the type stats property "author" should have type "relation"
    And the type stats property "author" relation count should be 2

  Scenario: Property with no values filled
    Given a vault is ready
    And a "book" object named "book1" exists
    And a "book" object named "book2" exists
    When I request type stats for "book"
    Then the type stats property "rating" should show filled 0

  Scenario: Non-existent type returns error
    Given a vault is ready
    When I request type stats for "nonexistent"
    Then an error should occur
