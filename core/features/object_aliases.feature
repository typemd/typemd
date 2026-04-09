Feature: Object aliases system property
  Objects can declare alternative names via the `aliases` stored system property.
  Aliases are written to frontmatter after `tags` and are indexed for name lookup.

  Scenario: Object with aliases serializes them in frontmatter
    Given a vault is ready
    And a "book" object named "golang-in-action" exists
    When I set the aliases of "golang-in-action" to "Go 語言" and "Golang"
    And I save the object "golang-in-action"
    Then the file of "golang-in-action" should contain "aliases:"

  Scenario: Object without aliases omits the aliases field
    Given a vault is ready
    And a "book" object named "golang-in-action" exists
    When I save the object "golang-in-action"
    Then the file of "golang-in-action" should not contain "aliases:"

  Scenario: User can set aliases on create via frontmatter
    Given a vault is ready
    And a "book" object named "golang-in-action" exists
    When I set the aliases of "golang-in-action" to "Go 語言" and "Golang"
    And I save the object "golang-in-action"
    Then the object "golang-in-action" should have alias "Go 語言"
    And the object "golang-in-action" should have alias "Golang"

  Scenario: Type schema with aliases property is rejected
    Given a vault is ready
    When I define a type schema with a property named "aliases"
    Then schema validation should report a reserved system property error for "aliases"

  Scenario: Object is found by alias via wiki-link resolution
    Given a vault is ready with note schemas
    And a "book" object named "golang-in-action" exists
    When I set the aliases of "golang-in-action" to "Go 語言"
    And I save the object "golang-in-action"
    And a "note" object named "my-notes" exists
    And "my-notes" body contains a shorthand wiki-link "book/go-語言"
    And I sync the index
    Then "my-notes" should have 1 wiki-link
    And the wiki-link should resolve to "golang-in-action"

