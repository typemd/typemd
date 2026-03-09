Feature: Property filtering during sync
  As a tmd user
  I want sync to only index schema-defined properties
  So that the type schema is the single source of truth for managed properties

  Background:
    Given a vault is ready

  Scenario: Sync filters out undefined properties
    Given a type schema "book" with properties "title,status,rating"
    And a raw object file "book/clean-code.md" with properties:
      | key    | value      |
      | title  | Clean Code |
      | status | reading    |
      | mood   | happy      |
    When I sync the index
    Then the indexed properties for "book/clean-code" should contain "title"
    And the indexed properties for "book/clean-code" should contain "status"
    And the indexed properties for "book/clean-code" should not contain "mood"

  Scenario: Object with only undefined properties stores empty object
    Given a type schema "book" with properties "title,status,rating"
    And a raw object file "book/mystery.md" with properties:
      | key   | value |
      | mood  | happy |
      | color | blue  |
    When I sync the index
    Then the indexed properties for "book/mystery" should be empty

  Scenario: Object type without schema retains all properties
    And a raw object file "recipe/pasta.md" with properties:
      | key  | value |
      | name | Pasta |
      | time | 30min |
    When I sync the index
    Then the indexed properties for "recipe/pasta" should contain "name"
    And the indexed properties for "recipe/pasta" should contain "time"

  Scenario: Frontmatter file is not modified during sync
    Given a type schema "book" with properties "title,status,rating"
    And a raw object file "book/clean-code.md" with properties:
      | key    | value      |
      | title  | Clean Code |
      | mood   | happy      |
    When I sync the index
    Then the file "book/clean-code.md" should still contain "mood" in frontmatter
