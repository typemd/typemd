Feature: Wiki-links and backlinks
  Objects can reference each other using [[type/name-ulid]] syntax in their body.
  The system tracks these links and their backlinks in the database.

  Scenario: Wiki-links are parsed and stored on sync
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And a "person" object named "robert-martin" exists
    And "clean-code" body contains a wiki-link to "robert-martin"
    When I sync the index
    Then "clean-code" should have 1 wiki-link
    And the wiki-link target should be "robert-martin"
    And "robert-martin" should have 1 backlink from "clean-code"

  Scenario: Broken wiki-link has empty resolved ID
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And "clean-code" body contains a wiki-link to "person/nobody-01jjjjjjjjjjjjjjjjjjjjjjjj"
    When I sync the index
    Then "clean-code" should have 1 wiki-link
    And the wiki-link should have an empty resolved ID

  Scenario: Wiki-links are updated on re-sync
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And a "person" object named "alice" exists
    And a "person" object named "bob" exists
    And "clean-code" body contains a wiki-link to "alice"
    And I sync the index
    When I change "clean-code" wiki-link to "bob"
    And I sync the index
    Then "clean-code" wiki-link should point to "bob"
    And "alice" should have 0 backlinks

  Scenario: Wiki-links support display text
    Given a vault is ready with note schemas
    And a "book" object named "clean-code" exists
    And a "person" object named "robert-martin" exists
    And "clean-code" body contains a wiki-link to "robert-martin" with display text "Uncle Bob"
    When I sync the index
    Then "clean-code" should have 1 wiki-link
    And the wiki-link display text should be "Uncle Bob"

  Scenario: Backlinks from multiple sources
    Given a vault is ready with note schemas
    And a "note" object named "target" exists
    And a "note" object named "alpha" exists
    And a "note" object named "beta" exists
    And "alpha" body contains a wiki-link to "target"
    And "beta" body contains a wiki-link to "target"
    When I sync the index
    Then "target" should have 2 backlinks
