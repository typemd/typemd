Feature: Import execute
  Execute a confirmed import plan: create types, create objects in
  dependency order, handle conflicts, and report results.

  Scenario: Execute a plan creating objects
    Given a vault is ready
    And an import plan with objects:
      | source_path | type_name | name        |
      | note-01.md  | page      | First Note  |
      | note-02.md  | page      | Second Note |
    When I execute the import plan
    Then no error should occur
    And the import report should show 2 created
    And the import report should show 0 skipped
    And the import report should show 0 failed

  Scenario: Execute a plan with new types
    Given a vault is ready
    And an import plan with new type "recipe" and objects:
      | source_path  | type_name | name            |
      | recipe-01.md | recipe    | Pasta Carbonara |
    When I execute the import plan
    Then no error should occur
    And the import report should show 1 types created
    And the import report should show 1 created

  Scenario: Skip conflicting objects
    Given a vault is ready
    And an import plan with a skipped object:
      | source_path   | type_name | name  | conflict |
      | old-note.md   | page      | old   | skip     |
    When I execute the import plan
    Then no error should occur
    And the import report should show 0 created
    And the import report should show 1 skipped

  Scenario: Report includes suggestions for failures
    Given a vault is ready
    And an import plan with objects for a nonexistent type:
      | source_path | type_name  | name  |
      | bad.md      | nonexistent| bad   |
    When I execute the import plan
    Then no error should occur
    And the import report should show 1 failed
    And the import report should suggest reviewing failed files
