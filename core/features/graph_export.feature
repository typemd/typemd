Feature: Graph export
  The vault can export its object graph (relations and wiki-links) in DOT format
  for visualization with external tools like Graphviz.

  Scenario: Empty vault produces empty digraph
    Given a vault is ready
    When I export the graph
    Then no error should occur
    And the DOT output should contain "digraph vault {"
    And the DOT output should contain "}"
    And the DOT output should have 0 nodes
    And the DOT output should have 0 edges

  Scenario: Single object with no relations produces one node
    Given a vault is ready
    And I create a "book" object named "clean-code"
    When I export the graph
    Then no error should occur
    And the DOT output should have 1 node
    And the DOT output should have 0 edges

  Scenario: Relation produces a labeled edge
    Given a vault is ready with relation schemas
    And I create a "book" object named "clean-code"
    And I create a "person" object named "bob"
    And I link "clean-code" to "bob" via "author"
    When I export the graph
    Then no error should occur
    And the DOT output should have 2 nodes
    And the DOT output should have an edge labeled "author"

  Scenario: Wiki-link produces an edge labeled wikilink
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And a "note" object named "my-notes" exists
    And "my-notes" body contains a wiki-link to "clean-code"
    And I sync the index
    When I export the graph
    Then no error should occur
    And the DOT output should have an edge labeled "wikilink"

  Scenario: Type filter limits graph to specified types
    Given a vault is ready with relation schemas
    And I create a "book" object named "clean-code"
    And I create a "person" object named "bob"
    And I link "clean-code" to "bob" via "author"
    When I export the graph with type filter "book"
    Then no error should occur
    And the DOT output should have 1 node
    And the DOT output should have 0 edges

  Scenario: No-relations flag excludes relation edges
    Given a vault is ready with relation schemas
    And I create a "book" object named "clean-code"
    And I create a "person" object named "bob"
    And I link "clean-code" to "bob" via "author"
    When I export the graph without relations
    Then no error should occur
    And the DOT output should have 2 nodes
    And the DOT output should have 0 edges

  Scenario: No-wikilinks flag excludes wiki-link edges
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And a "note" object named "my-notes" exists
    And "my-notes" body contains a wiki-link to "clean-code"
    And I sync the index
    When I export the graph without wikilinks
    Then no error should occur
    And the DOT output should have 0 edges

  Scenario: Unresolved wiki-link is skipped
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And "clean-code" body contains a wiki-link to "person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj"
    And I sync the index
    When I export the graph
    Then no error should occur
    And the DOT output should have 0 edges

  Scenario: Bidirectional relation produces one edge per stored direction
    Given a vault is ready with relation schemas
    And I create a "book" object named "clean-code"
    And I create a "person" object named "bob"
    And I link "clean-code" to "bob" via "author"
    When I export the graph
    Then the DOT output should have an edge labeled "author"
    And the DOT output should have an edge labeled "books"
