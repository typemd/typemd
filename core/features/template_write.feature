Feature: Template write operations
  Templates can be saved and deleted through the vault facade,
  enabling programmatic management of object templates.

  Scenario: Save template with properties and body
    Given a vault is ready
    When I save template "review" for type "book" with properties "status: draft" and body "## Review"
    Then no error should occur
    And I load template "review" for type "book"
    And the template property "status" should be "draft"
    And the template body should be "## Review"

  Scenario: Save template with body only
    Given a vault is ready
    When I save template "simple" for type "book" with body "## Notes"
    Then no error should occur
    And I load template "simple" for type "book"
    And the template body should be "## Notes"

  Scenario: Save template creates directory if missing
    Given a vault is ready
    When I save template "first" for type "article" with body "## Article"
    Then no error should occur
    And I list templates for type "article"
    And the template list should contain "first"

  Scenario: Save template overwrites existing
    Given a vault is ready
    And a template "review" for type "book" with body "## Old"
    When I save template "review" for type "book" with body "## New"
    Then no error should occur
    And I load template "review" for type "book"
    And the template body should be "## New"

  Scenario: Delete existing template
    Given a vault is ready
    And a template "review" for type "book" with body "## Review"
    And a template "summary" for type "book" with body "## Summary"
    When I delete template "review" for type "book"
    Then no error should occur
    And I list templates for type "book"
    And the template list should contain "summary"

  Scenario: Delete last template removes empty directory
    Given a vault is ready
    And a template "only" for type "book" with body "## Only"
    When I delete template "only" for type "book"
    Then no error should occur
    And the template directory for type "book" should not exist

  Scenario: Delete nonexistent template returns error
    Given a vault is ready
    When I delete template "nonexistent" for type "book"
    Then an error should occur
