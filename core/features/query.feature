Feature: Query and search
  Objects can be queried by type and property filters, or searched via full-text search.

  Scenario: Query objects by type
    Given a vault is ready
    And a "book" object named "book1" exists
    And a "book" object named "book2" exists
    And a "person" object named "alice" exists
    When I query objects with filter "type=book"
    Then the query should return 2 results
    And all results should have type "book"

  Scenario: Query objects by property
    Given a vault is ready
    And a "book" object named "book1" exists with property "status" set to "reading"
    And a "book" object named "book2" exists with property "status" set to "done"
    When I query objects with filter "type=book status=reading"
    Then the query should return 1 result

  Scenario: Query with empty filter returns all objects
    Given a vault is ready
    And a "book" object named "book1" exists
    And a "person" object named "alice" exists
    When I query objects with filter ""
    Then the query should return 2 results

  Scenario: Search objects by filename
    Given a vault is ready
    And a "book" object named "concurrency-in-go" exists
    And a "book" object named "clean-code" exists
    When I search objects for "concurrency"
    Then the search should return 1 result

  Scenario: Search objects by body content
    Given a vault is ready
    And a "book" object named "mybook" exists with body "This book covers goroutines and channels."
    And a "book" object named "other" exists
    When I search objects for "goroutines"
    Then the search should return 1 result
