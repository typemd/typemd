Feature: Frontmatter parsing and writing
  Objects are stored as Markdown files with YAML frontmatter.
  The frontmatter contains object properties.

  Scenario: Write frontmatter with properties
    Given properties with "title" set to "Go in Action"
    And properties with "rating" set to "4.5"
    When I write frontmatter with no body
    Then the output should start with "---"
    And the output should contain "title: Go in Action"

  Scenario: Write frontmatter with body
    Given properties with "title" set to "Test"
    When I write frontmatter with body "Hello world"
    Then the output should contain "Hello world"

  Scenario: Write frontmatter with empty properties
    Given empty properties
    When I write frontmatter with no body
    Then the output should equal "---\n---\n"

  Scenario: Parse frontmatter from markdown
    Given markdown content:
      """
      ---
      title: Go in Action
      rating: 4.5
      ---

      Some body here.
      """
    When I parse the frontmatter
    Then the parsed property "title" should be "Go in Action"
    And the parsed body should be "Some body here."
