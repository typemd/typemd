Feature: Page built-in type
  Pages are general-purpose content containers for free-form writing.
  The page type exists in every vault by default.

  Scenario: Page type exists as a built-in type
    Given a vault is ready
    When I create a "page" object named "hello-world"
    Then the object type should be "page"
    And the object file should exist on disk

  Scenario: Page type loads without custom schema
    Given a vault is ready
    When I load type "page"
    Then no error should occur
    And the loaded schema should have emoji "📄"
    And the loaded schema plural should be "pages"
    And the loaded schema should have unique false

  Scenario: Page type has no custom properties
    Given a vault is ready
    When I load type "page"
    Then no error should occur
    And the loaded type should have 0 properties

  Scenario: Deleting page type is rejected
    Given a vault is ready
    When I delete type "page"
    Then an error should occur
    And the error message should contain "cannot delete built-in type"

  Scenario: Page type appears in type listing
    Given a vault is ready
    When I list all types
    Then the type list should contain "page"
    And the type list should contain "tag"

  Scenario: Custom page schema overrides built-in
    Given a vault is ready
    And a custom "page" type schema with emoji "📝"
    When I load type "page"
    Then no error should occur
    And the loaded schema should have emoji "📝"
