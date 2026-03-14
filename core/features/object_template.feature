Feature: Object templates
  Types can define templates that provide default frontmatter and body content
  when creating new objects.

  Scenario: List templates for type with multiple templates
    Given a vault is ready
    And a template "review" for type "book" with body "## Review"
    And a template "summary" for type "book" with body "## Summary"
    When I list templates for type "book"
    Then the template list should contain "review, summary"

  Scenario: List templates for type with one template
    Given a vault is ready
    And a template "default" for type "book" with body "## Notes"
    When I list templates for type "book"
    Then the template list should contain "default"

  Scenario: List templates for type with no templates directory
    Given a vault is ready
    When I list templates for type "book"
    Then the template list should be empty

  Scenario: List templates for type with empty templates directory
    Given a vault is ready
    And an empty templates directory for type "book"
    When I list templates for type "book"
    Then the template list should be empty

  Scenario: Load template with frontmatter and body
    Given a vault is ready
    And a template "review" for type "book" with frontmatter "status: draft" and body "## Notes"
    When I load template "review" for type "book"
    Then the template property "status" should be "draft"
    And the template body should be "## Notes"

  Scenario: Load template with body only
    Given a vault is ready
    And a template "simple" for type "book" with body "## My Book"
    When I load template "simple" for type "book"
    Then the template body should be "## My Book"

  Scenario: Load template with frontmatter only
    Given a vault is ready
    And a template "preset" for type "book" with frontmatter "status: reading" and body ""
    When I load template "preset" for type "book"
    Then the template property "status" should be "reading"

  Scenario: Load nonexistent template returns error
    Given a vault is ready
    When I load template "nonexistent" for type "book"
    Then the template load should fail

  Scenario: Create object with template body and properties
    Given a vault is ready
    And a template "review" for type "book" with frontmatter "status: draft" and body "## Review Notes"
    When I create a "book" object named "my-book" with template "review"
    Then the object property "status" should be "draft"
    And the object body should be "## Review Notes"

  Scenario: Template overrides schema default
    Given a vault is ready
    And a template "review" for type "book" with frontmatter "status: draft" and body ""
    When I create a "book" object named "my-book" with template "review"
    Then the object property "status" should be "draft"

  Scenario: Create object without template preserves current behavior
    Given a vault is ready
    When I create a "book" object named "my-book"
    Then the object body should be ""

  Scenario: Template with mutable system property description
    Given a vault is ready
    And a template "review" for type "book" with frontmatter "description: A book review" and body ""
    When I create a "book" object named "my-book" with template "review"
    Then the object property "description" should be "A book review"

  Scenario: Template with immutable system property created_at is ignored
    Given a vault is ready
    And a template "review" for type "book" with frontmatter "created_at: 2020-01-01T00:00:00Z" and body ""
    When I create a "book" object named "my-book" with template "review"
    Then the object "created_at" should be recent

  Scenario: Template with unknown property is ignored
    Given a vault is ready
    And a template "review" for type "book" with frontmatter "unknown_prop: value" and body ""
    When I create a "book" object named "my-book" with template "review"
    Then the object should not have property "unknown_prop"

  Scenario: Nonexistent template name in create returns error
    Given a vault is ready
    When I create a "book" object named "my-book" with template "nonexistent"
    Then the object creation should fail
