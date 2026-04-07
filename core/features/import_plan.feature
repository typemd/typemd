Feature: Import plan
  Generate a conversion plan that maps source files to vault objects,
  determines import order by dependency, and detects existing objects.

  Scenario: Generate a plan with dependency ordering
    Given a vault is ready
    And a classification list:
      | source_path | type_name | name       |
      | tag-go.md   | tag       | go         |
      | book-one.md | book      | Clean Code |
    When I generate an import plan
    Then no error should occur
    And the plan should have 0 new types
    And the plan order should place "tag" objects before "book" objects

  Scenario: Plan detects new type schemas needed
    Given a vault is ready
    And a classification list:
      | source_path  | type_name | name           |
      | recipe-01.md | recipe    | Pasta Carbonara|
    When I generate an import plan
    Then no error should occur
    And the plan should have 1 new types

  Scenario: Plan detects existing objects as conflicts
    Given a vault is ready
    And a "book" object named "clean-code" exists
    And a classification list:
      | source_path   | type_name | name       |
      | clean-code.md | book      | clean-code |
    When I generate an import plan
    Then no error should occur
    And the plan object 0 should have conflict "skip"

  Scenario: Plan with circular dependencies
    Given a vault is ready
    And a classification list with circular dependencies
    When I generate an import plan
    Then no error should occur
    And the plan should include all objects in the order
