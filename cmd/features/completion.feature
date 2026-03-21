Feature: Shell completion
  As a user
  I want tab completion for object IDs, type names, and relation names
  So that I can use the CLI without memorizing long ULID-suffixed identifiers

  Background:
    Given a vault with a book "Clean Code"

  # Object ID completion — two-stage progressive

  Scenario: Complete type prefix when no slash present
    When I request object ID completions for "b"
    Then the completions should include "book/"

  Scenario: Complete object name after type prefix
    When I request object ID completions for "book/"
    Then the completions should include a book object starting with "book/clean-code-"

  Scenario: No matches for unknown type prefix
    When I request object ID completions for "xyz"
    Then the completions should be empty

  # Type name completion

  Scenario: Complete type name with prefix
    When I request type name completions for "b"
    Then the completions should include "book"

  Scenario: Complete type name lists all types with empty prefix
    When I request type name completions for ""
    Then the completions should include "book"
    And the completions should include "tag"
    And the completions should include "page"

  # Relation name completion

  Scenario: Complete relation name from source object schema
    Given the book type has a relation "series" to "book"
    When I request relation name completions for the created book with prefix ""
    Then the completions should include "series"

  Scenario: No relation completion for invalid source object
    When I request relation name completions for "invalid/nonexistent" with prefix ""
    Then the completions should be empty
