Feature: Query fallback when index is unavailable
  When the SQLite index is unavailable, query and search operations
  fall back to filesystem scanning instead of failing.

  Background:
    Given a vault with objects and a broken index

  Scenario: Query falls back to filesystem with filter and sort
    When I query with fallback filter "type=book" sorted by "name" "asc"
    Then the fallback query should return 2 results
    And all fallback results should have type "book"
    And the first fallback result name should be "Alpha Book"

  Scenario: Search falls back to substring matching
    When I search with fallback for "interesting content"
    Then the fallback search should return 1 result

  Scenario: Search fallback is case-insensitive
    When I search with fallback for "alpha"
    Then the fallback search should return 1 result
