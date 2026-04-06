Feature: Property access via GetProperty
  Objects expose both stored (frontmatter) and non-stored (derived/computed)
  properties through a unified GetProperty API. Derived properties like
  object_type are inferred from structure; computed properties like links
  are resolved from content.

  Scenario: GetProperty resolves object_type for a book
    Given a vault is ready
    And a "book" object named "get-prop-book" exists
    When I get property "object_type" on the object
    Then the property value should be "book"
    And the property should exist

  Scenario: GetProperty resolves object_type for a person
    Given a vault is ready
    And a "person" object named "get-prop-person" exists
    When I get property "object_type" on the object
    Then the property value should be "person"
    And the property should exist

  Scenario: GetProperty resolves a stored property
    Given a vault is ready
    And a "book" object named "get-prop-stored" exists
    And I set property "title" to "Go in Action" on the object
    When I get property "title" on the object
    Then the property value should be "Go in Action"
    And the property should exist

  Scenario: GetProperty returns false for missing property
    Given a vault is ready
    And a "book" object named "get-prop-missing" exists
    When I get property "nonexistent" on the object
    Then the property should not exist
