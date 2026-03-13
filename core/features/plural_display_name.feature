Feature: Plural display name
  Type schemas can have an optional plural field for grammatically correct
  collection display names (e.g., "books" instead of "book" in group headers).

  Scenario: Type schema with plural defined
    Given a vault is ready
    And a type schema "book" with plural "books"
    When I load type "book"
    Then the loaded schema plural should be "books"
    And the loaded schema PluralName should be "books"

  Scenario: Type schema without plural defined
    Given a vault is ready
    And a type schema "note" without plural
    When I load type "note"
    Then the loaded schema plural should be ""
    And the loaded schema PluralName should be "note"

  Scenario: Built-in tag type has plural
    Given a vault is ready
    When I load type "tag"
    Then the loaded schema plural should be "tags"
    And the loaded schema PluralName should be "tags"
