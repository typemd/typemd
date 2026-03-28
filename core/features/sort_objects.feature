Feature: In-memory object sorting
  Objects can be sorted by property values in memory when SQLite index is unavailable.
  SortObjects is a pure function that operates on []*Object without requiring a vault.

  Scenario: Sort by string property ascending
    Given in-memory objects:
      | type | name      | status  |
      | book | Charlie   | draft   |
      | book | Alpha     | active  |
      | book | Bravo     | reading |
    When I sort objects by "status" "asc"
    Then the sorted object names should be "Alpha,Charlie,Bravo"

  Scenario: Sort by string property descending
    Given in-memory objects:
      | type | name    | status  |
      | book | Alpha   | active  |
      | book | Bravo   | reading |
      | book | Charlie | draft   |
    When I sort objects by "status" "desc"
    Then the sorted object names should be "Bravo,Charlie,Alpha"

  Scenario: Sort by number property ascending
    Given in-memory objects:
      | type | name  | rating |
      | book | Low   | 2      |
      | book | High  | 9      |
      | book | Mid   | 5      |
    When I sort objects by "rating" "asc"
    Then the sorted object names should be "Low,Mid,High"

  Scenario: Sort by type property
    Given in-memory objects:
      | type    | name   |
      | book    | BookA  |
      | article | ArtA   |
      | book    | BookB  |
    When I sort objects by "type" "asc"
    Then the sorted object names should be "ArtA,BookA,BookB"

  Scenario: Sort with missing property values
    Given in-memory objects:
      | type | name    | rating |
      | book | Rated   | 5      |
      | book | Unrated |        |
      | book | Also    | 3      |
    When I sort objects by "rating" "asc"
    Then the object "Unrated" should be last in sorted results

  Scenario: Stable sort preserves original order for equal values
    Given in-memory objects:
      | type | name   | status |
      | book | First  | active |
      | book | Second | active |
      | book | Third  | active |
    When I sort objects by "status" "asc"
    Then the sorted object names should be "First,Second,Third"
