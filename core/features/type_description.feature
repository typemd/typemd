Feature: Type and property description
  Type schemas and their properties can have optional descriptions for documentation.

  Scenario: Type schema with description
    Given a vault is ready
    And a type schema "presentation" with description "Slide decks and presentation materials"
    When I load type "presentation"
    Then the loaded schema description should be "Slide decks and presentation materials"

  Scenario: Type schema without description
    Given a vault is ready
    And a type schema "item" with a "title" string property
    When I load type "item"
    Then the loaded schema description should be ""

  Scenario: Property with description
    Given a vault is ready
    And a type schema "presentation" with property "speaker" having description "The person who gave this presentation"
    When I load type "presentation"
    Then the loaded property "speaker" description should be "The person who gave this presentation"

  Scenario: Property without description
    Given a vault is ready
    And a type schema "item" with a "title" string property
    When I load type "item"
    Then the loaded property "title" description should be ""
