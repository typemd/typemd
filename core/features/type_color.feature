Feature: Type schema color
  Type schemas can have an optional color for visual theming in TUI and Web UI.

  Scenario: Type schema with preset color
    Given a vault is ready
    And a type schema "item" with color "green"
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Type schema with hex color
    Given a vault is ready
    And a type schema "item" with color "#FF5733"
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Type schema with 3-digit hex color
    Given a vault is ready
    And a type schema "item" with color "#F53"
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Type schema without color
    Given a vault is ready
    And a type schema "item" with a "title" string property
    When I validate all schemas
    Then schema "item" should have no errors

  Scenario: Invalid color name rejected
    Given a vault is ready
    And a type schema "item" with color "magenta"
    When I validate all schemas
    Then schema "item" should have errors

  Scenario: Invalid hex color rejected
    Given a vault is ready
    And a type schema "item" with color "#GGGGGG"
    When I validate all schemas
    Then schema "item" should have errors

  Scenario: Hex without hash rejected
    Given a vault is ready
    And a type schema "item" with color "FF5733"
    When I validate all schemas
    Then schema "item" should have errors
