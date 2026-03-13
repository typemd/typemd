Feature: Name template
  Type schemas can define a name template for auto-generating object names at creation time.

  Scenario: Type with name template generates name automatically
    Given a vault is ready
    And a type schema "journal" with name template "日記 {{ date:YYYY-MM-DD }}"
    When I create a "journal" object with no name
    Then no error should occur
    And the object property "name" should contain "日記"

  Scenario: Explicit name overrides template
    Given a vault is ready
    And a type schema "journal" with name template "日記 {{ date:YYYY-MM-DD }}"
    When I create a "journal" object named "我的日記"
    Then no error should occur
    And the object property "name" should be "我的日記"

  Scenario: Object creation fails when no name and no template
    Given a vault is ready
    When I create a "book" object with no name
    Then an error should occur

  Scenario: Schema with name template passes validation
    Given a vault is ready
    And a type schema "journal" with name template "日記 {{ date:YYYY-MM-DD }}"
    When I validate all schemas
    Then schema "journal" should have no errors

  Scenario: Schema with name entry but no template fails validation
    Given a vault is ready
    And a type schema "bad" with name property and type "string"
    When I validate all schemas
    Then schema "bad" should have errors
