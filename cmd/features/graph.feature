Feature: tmd graph command
  The graph command exports the object relation graph in DOT format.

  Background:
    Given a vault is ready

  Scenario: Graph of empty vault
    When I run graph
    Then the command should succeed
    And the output should contain "digraph vault {"

  Scenario: Graph includes objects as nodes
    Given a book object "Clean Code" exists
    When I run graph
    Then the command should succeed
    And the output should contain "Clean Code"

  Scenario: Graph with type filter
    Given a book object "Clean Code" exists
    When I run graph with type "book"
    Then the command should succeed
    And the output should contain "Clean Code"

  Scenario: Graph with no-relations flag
    Given a book object "Clean Code" exists
    When I run graph with no-relations
    Then the command should succeed
    And the output should contain "digraph vault {"

  Scenario: Graph with no-wikilinks flag
    Given a book object "Clean Code" exists
    When I run graph with no-wikilinks
    Then the command should succeed
    And the output should contain "digraph vault {"
