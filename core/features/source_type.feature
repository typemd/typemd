Feature: Source built-in type
  Sources track raw materials (books, articles, videos) that have been
  ingested into the vault. The source type exists in every vault by default.

  Scenario: Source type exists as a built-in type
    Given a vault is ready
    When I create a "source" object named "golang-blog-post"
    Then the object type should be "source"
    And the object file should exist on disk

  Scenario: Source type loads without custom schema
    Given a vault is ready
    When I load type "source"
    Then no error should occur
    And the loaded schema should have emoji "📥"
    And the loaded schema plural should be "sources"
    And the loaded schema should have unique false

  Scenario: Source type has provenance properties
    Given a vault is ready
    When I load type "source"
    Then no error should occur
    And the loaded type should have 3 properties
    And the loaded property "url" should have type "string"
    And the loaded property "url" should have emoji "🔗"
    And the loaded property "author" should have type "string"
    And the loaded property "author" should have emoji "✍️"
    And the loaded property "ingested_at" should have type "date"
    And the loaded property "ingested_at" should have emoji "📅"

  Scenario: Deleting source type is rejected
    Given a vault is ready
    When I delete type "source"
    Then an error should occur
    And the error message should contain "cannot delete built-in type"

  Scenario: Source type appears in type listing
    Given a vault is ready
    When I list all types
    Then the type list should contain "source"

  Scenario: Custom source schema overrides built-in
    Given a vault is ready
    And a custom "source" type schema with emoji "📦"
    When I load type "source"
    Then no error should occur
    And the loaded schema should have emoji "📦"
