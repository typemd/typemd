Feature: Tags system property
  Every object has a tags system property that holds references to tag objects.

  Scenario: Tags is recognized as a system property
    Then "tags" should be a system property

  Scenario: Tags appears in system property names
    Then the system property registry should contain "name, description, created_at, updated_at, tags"

  Scenario: Frontmatter orders tags after updated_at
    Given a vault is ready
    And a "book" object named "ordered-tags-book" exists
    When I set tags on the object to a tag reference
    Then the frontmatter should have "updated_at" before "tags"
    And the frontmatter should have "tags" before "title"
