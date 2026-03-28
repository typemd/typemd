Feature: Query fallback when index is unavailable
  When the SQLite index is unavailable, query operations fall back to
  filesystem scanning with in-memory filtering and sorting.

  Background:
    Given a vault with objects and a broken index

  Scenario: Query with type filter in fallback mode
    When I query with fallback filter "type=book"
    Then the fallback query should return 2 results
    And all fallback results should have type "book"

  Scenario: Query with property filter in fallback mode
    When I query with fallback filter "type=book status=reading"
    Then the fallback query should return 1 result

  Scenario: Query with no filters returns all objects
    When I query with fallback filter ""
    Then the fallback query should return 3 results

  Scenario: Query with sort in fallback mode
    When I query with fallback filter "type=book" sorted by "name" "asc"
    Then the fallback query should return 2 results
    And the first fallback result name should be "Alpha Book"

  Scenario: Search by name in fallback mode
    When I search with fallback for "Alpha"
    Then the fallback search should return 1 result

  Scenario: Search by body in fallback mode
    When I search with fallback for "interesting content"
    Then the fallback search should return 1 result

  Scenario: Search with no match in fallback mode
    When I search with fallback for "nonexistent"
    Then the fallback search should return 0 results

  Scenario: Search is case-insensitive in fallback mode
    When I search with fallback for "alpha"
    Then the fallback search should return 1 result

  Scenario: VaultStats works in fallback mode
    When I request vault stats with fallback
    Then the fallback vault stats total should be 3

  Scenario: TypeStats works in fallback mode
    When I request type stats for "book" with fallback
    Then the fallback type stats count should be 2
